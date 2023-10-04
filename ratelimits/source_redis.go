package ratelimits

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jmhodges/clock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// Compile-time check that RedisSource implements the source interface.
var _ source = (*RedisSource)(nil)

// RedisSource is a ratelimits source backed by sharded Redis.
type RedisSource struct {
	client  *redis.Ring
	clk     clock.Clock
	latency *prometheus.HistogramVec
}

// NewRedisSource returns a new Redis backed source using the provided
// *redis.Ring client.
func NewRedisSource(client *redis.Ring, clk clock.Clock, stats prometheus.Registerer) *RedisSource {
	latency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "ratelimits_latency",
			Help: "Histogram of Redis call latencies labeled by call=[set|get|delete|ping] and result=[success|error]",
			// Exponential buckets ranging from 0.0005s to 3s.
			Buckets: prometheus.ExponentialBucketsRange(0.0005, 3, 8),
		},
		[]string{"call", "result"},
	)
	stats.MustRegister(latency)

	return &RedisSource{
		client:  client,
		clk:     clk,
		latency: latency,
	}
}

// resultForError returns a string representing the result of the operation
// based on the provided error.
func resultForError(err error) string {
	if errors.Is(redis.Nil, err) {
		// Bucket key does not exist.
		return "notFound"
	} else if errors.Is(err, context.DeadlineExceeded) {
		// Client read or write deadline exceeded.
		return "deadlineExceeded"
	} else if errors.Is(err, context.Canceled) {
		// Caller canceled the operation.
		return "canceled"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		// Dialer timed out connecting to Redis.
		return "timeout"
	}
	var redisErr redis.Error
	if errors.Is(err, redisErr) {
		// An internal error was returned by the Redis server.
		return "redisError"
	}
	return "failed"
}

// Set stores the TAT at the specified bucketKey ('name:id'). It returns an
// error if the operation failed and nil otherwise. If the bucketKey does not
// exist, it will be created.
func (r *RedisSource) Set(ctx context.Context, bucketKey string, tat time.Time) error {
	start := r.clk.Now()

	err := r.client.Set(ctx, bucketKey, tat.UnixNano(), 0).Err()
	if err != nil {
		r.latency.With(prometheus.Labels{"call": "set", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return err
	}

	r.latency.With(prometheus.Labels{"call": "set", "result": "success"}).Observe(time.Since(start).Seconds())
	return nil
}

// BatchSet stores the bucket keys and corresponding TATs using Redis' MSET
// command. It returns an error if the operation failed and nil otherwise. If a
// key for a given bucket does not exist, it will be created.
func (r *RedisSource) BatchSet(ctx context.Context, buckets map[string]time.Time) error {
	start := r.clk.Now()

	// Flatten the map into a slice of alternating keys and values.
	bucketsSlice := make([]interface{}, 0, len(buckets)*2)
	for k, v := range buckets {
		bucketsNano := v.UTC().UnixNano()
		bucketsSlice = append(bucketsSlice, k, bucketsNano)
	}

	err := r.client.MSet(ctx, bucketsSlice).Err()
	if err != nil {
		r.latency.With(prometheus.Labels{"call": "mset", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return err
	}

	r.latency.With(prometheus.Labels{"call": "mset", "result": "success"}).Observe(time.Since(start).Seconds())
	return nil
}

// Get retrieves the TAT at the specified bucketKey ('name:id'). It returns the
// TAT and nil if the operation succeeded, or an error if the operation failed.
// If the bucketKey does not exist, it returns ErrBucketNotFound.
func (r *RedisSource) Get(ctx context.Context, bucketKey string) (time.Time, error) {
	start := r.clk.Now()

	tatNano, err := r.client.Get(ctx, bucketKey).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// Bucket key does not exist.
			r.latency.With(prometheus.Labels{"call": "get", "result": "notFound"}).Observe(time.Since(start).Seconds())
			return time.Time{}, ErrBucketNotFound
		}
		r.latency.With(prometheus.Labels{"call": "get", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return time.Time{}, err
	}

	r.latency.With(prometheus.Labels{"call": "get", "result": "success"}).Observe(time.Since(start).Seconds())
	return time.Unix(0, tatNano).UTC(), nil
}

// BatchGet retrieves the TATs at the specified bucketKeys ('name:id'). It
// returns the TATs and nil if the operation succeeded, or an error if the
// operation failed. If a bucketKey does not exist, it will be returned in the
// map with a zero value.
func (r *RedisSource) BatchGet(ctx context.Context, bucketKeys []string) (map[string]time.Time, error) {
	start := r.clk.Now()

	results, err := r.client.MGet(ctx, bucketKeys...).Result()
	if err != nil {
		r.latency.With(prometheus.Labels{"call": "mget", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return nil, err
	}

	buckets := make(map[string]time.Time, len(bucketKeys))
	for i, iface := range results {
		if iface == nil {
			// Bucket key does not exist. A zero value timestamp is used to
			// indicate that the requested bucket does not exist.
			buckets[bucketKeys[i]] = time.Time{}
			continue
		}
		tatRaw, ok := iface.(string)
		if !ok {
			// This should never happen, Set() and BatchSet() should never
			// store a TATs using a different Redis type.
			return nil, fmt.Errorf("invalid type for bucketKey %q: %T", bucketKeys[i], iface)
		}
		tatNano, err := strconv.ParseInt(tatRaw, 10, 64)
		if err != nil {
			// This should never happen, Set() and BatchSet() should never
			// store a TATs as anything other than a 64-bit integer.
			return nil, fmt.Errorf("invalid value for bucketKey %q: %s", bucketKeys[i], tatRaw)
		}
		buckets[bucketKeys[i]] = time.Unix(0, tatNano).UTC()
	}

	r.latency.With(prometheus.Labels{"call": "mget", "result": "success"}).Observe(time.Since(start).Seconds())
	return buckets, nil
}

// Delete deletes the TAT at the specified bucketKey ('name:id'). It returns an
// error if the operation failed and nil otherwise. A nil return value does not
// indicate that the bucketKey existed.
func (r *RedisSource) Delete(ctx context.Context, bucketKey string) error {
	start := r.clk.Now()

	err := r.client.Del(ctx, bucketKey).Err()
	if err != nil {
		r.latency.With(prometheus.Labels{"call": "delete", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return err
	}

	r.latency.With(prometheus.Labels{"call": "delete", "result": "success"}).Observe(time.Since(start).Seconds())
	return nil
}

// Ping checks that each shard of the *redis.Ring is reachable using the PING
// command. It returns an error if any shard is unreachable and nil otherwise.
func (r *RedisSource) Ping(ctx context.Context) error {
	start := r.clk.Now()

	err := r.client.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		r.latency.With(prometheus.Labels{"call": "ping", "result": resultForError(err)}).Observe(time.Since(start).Seconds())
		return err
	}
	r.latency.With(prometheus.Labels{"call": "ping", "result": "success"}).Observe(time.Since(start).Seconds())
	return nil
}
