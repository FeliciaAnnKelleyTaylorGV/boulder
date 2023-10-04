package ratelimits

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jmhodges/clock"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Allowed is used for rate limit metrics, it's the value of the 'decision'
	// label when a request was allowed.
	Allowed = "allowed"

	// Denied is used for rate limit metrics, it's the value of the 'decision'
	// label when a request was denied.
	Denied = "denied"
)

// ErrInvalidCost indicates that the cost specified was <= 0.
var ErrInvalidCost = fmt.Errorf("invalid cost, must be > 0")

// ErrInvalidCostForCheck indicates that the check cost specified was < 0.
var ErrInvalidCostForCheck = fmt.Errorf("invalid check cost, must be >= 0")

// ErrInvalidCostOverLimit indicates that the cost specified was > limit.Burst.
var ErrInvalidCostOverLimit = fmt.Errorf("invalid cost, must be <= limit.Burst")

// errLimitDisabled indicates that the limit name specified is valid but is not
// currently configured.
var errLimitDisabled = errors.New("limit disabled")

// disabledLimitDecision is an "allowed" *Decision that should be returned when
// a checked limit is found to be disabled.
var disabledLimitDecision = &Decision{true, 0, 0, 0, time.Time{}}

// BatchEntry is used to batch requests to the limiter.
type BatchEntry struct {
	Name Name
	Id   string
	Cost int64
}

// bucketKey returns the key used to store the bucket in the source.
func (b *BatchEntry) bucketKey() string {
	return fmt.Sprintf("%s:%s", nameToEnumString(b.Name), b.Id)
}

// Batch is a slice of *BatchEntry used to batch requests to the limiter.
type Batch []*BatchEntry

// Limiter provides a high-level interface for rate limiting requests by
// utilizing a leaky bucket-style approach.
type Limiter struct {
	// defaults stores default limits by 'name'.
	defaults limits

	// overrides stores override limits by 'name:id'.
	overrides limits

	// source is used to store buckets. It must be safe for concurrent use.
	source source
	clk    clock.Clock

	spendLatency       *prometheus.HistogramVec
	overrideUsageGauge *prometheus.GaugeVec
}

// NewLimiter returns a new *Limiter. The provided source must be safe for
// concurrent use. The defaults and overrides paths are expected to be paths to
// YAML files that contain the default and override limits, respectively. The
// overrides file is optional, all other arguments are required.
func NewLimiter(clk clock.Clock, source source, defaults, overrides string, stats prometheus.Registerer) (*Limiter, error) {
	limiter := &Limiter{source: source, clk: clk}

	var err error
	limiter.defaults, err = loadAndParseDefaultLimits(defaults)
	if err != nil {
		return nil, err
	}

	limiter.spendLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ratelimits_spend_latency",
		Help: fmt.Sprintf("Latency of ratelimit checks labeled by limit=[name] and decision=[%s|%s], in seconds", Allowed, Denied),
		// Exponential buckets ranging from 0.0005s to 3s.
		Buckets: prometheus.ExponentialBuckets(0.0005, 3, 8),
	}, []string{"limit", "decision"})
	stats.MustRegister(limiter.spendLatency)

	if overrides == "" {
		// No overrides specified, initialize an empty map.
		limiter.overrides = make(limits)
		return limiter, nil
	}

	limiter.overrides, err = loadAndParseOverrideLimits(overrides)
	if err != nil {
		return nil, err
	}

	limiter.overrideUsageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ratelimits_override_usage",
		Help: "Proportion of override limit used, by limit name and client id.",
	}, []string{"limit", "client_id"})
	stats.MustRegister(limiter.overrideUsageGauge)

	return limiter, nil
}

type Decision struct {
	// Allowed is true if the bucket possessed enough capacity to allow the
	// request given the cost.
	Allowed bool

	// Remaining is the number of requests the client is allowed to make before
	// they're rate limited.
	Remaining int64

	// RetryIn is the duration the client MUST wait before they're allowed to
	// make a request.
	RetryIn time.Duration

	// ResetIn is the duration the bucket will take to refill to its maximum
	// capacity, assuming no further requests are made.
	ResetIn time.Duration

	// newTAT indicates the time at which the bucket will be full. It is the
	// theoretical arrival time (TAT) of next request. It must be no more than
	// (burst * (period / count)) in the future at any single point in time.
	newTAT time.Time
}

// Check returns a *Decision that indicates whether there's enough capacity to
// allow the request, given the cost, for the specified limit Name and client
// id. However, it DOES NOT deduct the cost of the request from the bucket's
// capacity. Hence, the returned *Decision represents the hypothetical state of
// the bucket if the cost WERE to be deducted. The returned *Decision will
// always include the number of remaining requests in the bucket, the required
// wait time before the client can make another request, and the time until the
// bucket refills to its maximum capacity (resets). If no bucket exists for the
// given limit Name and client id, a new one will be created WITHOUT the
// request's cost deducted from its initial capacity.
func (l *Limiter) Check(ctx context.Context, name Name, id string, cost int64) (*Decision, error) {
	if cost < 0 {
		return nil, ErrInvalidCostForCheck
	}

	limit, err := l.getLimit(name, id)
	if err != nil {
		if errors.Is(err, errLimitDisabled) {
			return disabledLimitDecision, nil
		}
		return nil, err
	}

	if cost > limit.Burst {
		return nil, ErrInvalidCostOverLimit
	}

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	tat, err := l.source.Get(ctx, bucketKey(name, id))
	if err != nil {
		if !errors.Is(err, ErrBucketNotFound) {
			return nil, err
		}
		// First request from this client. The cost is not deducted from the
		// initial capacity because this is only a check.
		d, err := l.initialize(ctx, limit, name, id, 0)
		if err != nil {
			return nil, err
		}
		return maybeSpend(l.clk, limit, d.newTAT, cost), nil
	}
	return maybeSpend(l.clk, limit, tat, cost), nil
}

// Spend returns a *Decision that indicates if enough capacity was available to
// process the request, given the cost, for the specified limit Name and client
// id. If capacity existed, the cost of the request HAS been deducted from the
// bucket's capacity, otherwise no cost was deducted. The returned *Decision
// will always include the number of remaining requests in the bucket, the
// required wait time before the client can make another request, and the time
// until the bucket refills to its maximum capacity (resets). If no bucket
// exists for the given limit Name and client id, a new one will be created WITH
// the request's cost deducted from its initial capacity.
func (l *Limiter) Spend(ctx context.Context, name Name, id string, cost int64) (*Decision, error) {
	if cost <= 0 {
		return nil, ErrInvalidCost
	}

	limit, err := l.getLimit(name, id)
	if err != nil {
		if errors.Is(err, errLimitDisabled) {
			return disabledLimitDecision, nil
		}
		return nil, err
	}

	if cost > limit.Burst {
		return nil, ErrInvalidCostOverLimit
	}

	start := l.clk.Now()
	status := Denied
	defer func() {
		l.spendLatency.WithLabelValues(name.String(), status).Observe(l.clk.Since(start).Seconds())
	}()

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	tat, err := l.source.Get(ctx, bucketKey(name, id))
	if err != nil {
		if errors.Is(err, ErrBucketNotFound) {
			// First request from this client.
			d, err := l.initialize(ctx, limit, name, id, cost)
			if err != nil {
				return nil, err
			}
			if d.Allowed {
				status = Allowed
			}
			return d, nil
		}
		return nil, err
	}

	d := maybeSpend(l.clk, limit, tat, cost)

	if limit.isOverride {
		// Calculate the current utilization of the override limit for the
		// specified client id.
		utilization := float64(limit.Burst-d.Remaining) / float64(limit.Burst)
		l.overrideUsageGauge.WithLabelValues(name.String(), id).Set(utilization)
	}

	if !d.Allowed {
		return d, nil
	}

	err = l.source.Set(ctx, bucketKey(name, id), d.newTAT)
	if err != nil {
		return nil, err
	}
	status = Allowed
	return d, nil
}

// BatchSpend checks each bucket in the batch to see if there's enough capacity
// to allow the request, given the cost. If capacity exists, the cost of the
// request HAS been deducted from the bucket's capacity, otherwise no cost was
// deducted. If one or more buckets in the batch do not exist, they will be
// created WITH the request's cost deducted from the initial capacity. If one or
// more buckets in the batch are disabled, they will be ignored. The returned
// *Decision will always include the number of remaining requests in the bucket,
// the required wait time before the client can make another request, and the
// time until the bucket refills to its maximum capacity (resets). This is
// achieved by taking the minimum of the Remaining values for each bucket in the
// batch, the maximum of the RetryIn and ResetIn values.
func (l *Limiter) BatchSpend(ctx context.Context, batch Batch) (*Decision, error) {
	if len(batch) == 0 {
		return nil, ErrInvalidCost
	}

	bucketKeys := make([]string, 0, len(batch))
	for _, entry := range batch {
		if entry.Cost <= 0 {
			return nil, ErrInvalidCost
		}
		bucketKeys = append(bucketKeys, entry.bucketKey())
	}

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	tats, err := l.source.BatchGet(ctx, bucketKeys)
	if err != nil {
		return nil, err
	}

	var minRemaining int64 = math.MaxInt64
	var maxRetryIn time.Duration
	var maxResetIn time.Duration
	newTATs := make(map[string]time.Time)
	allowed := true

	// Assign nowTAT outside of the loop to avoid clock skew.
	nowTAT := l.clk.Now()

	for _, entry := range batch {
		bucketKey := entry.bucketKey()
		tat, exists := tats[bucketKey]
		limit, err := l.getLimit(entry.Name, entry.Id)
		if err != nil {
			if errors.Is(err, errLimitDisabled) {
				// Ignore disabled limit.
				continue
			}
			return nil, err
		}

		if !exists || tat.IsZero() {
			// First request from this client. A TAT of "now" is equivalent to a
			// full bucket.
			tat = nowTAT
		}

		// Spend the cost and update the consolidated decision.
		d := maybeSpend(l.clk, limit, tat, entry.Cost)
		if d.Allowed {
			newTATs[bucketKey] = d.newTAT
		}

		// All spend decisions must be allowed for the batch to be considered
		// allowed.
		allowed = allowed && d.Allowed
		minRemaining = min(minRemaining, d.Remaining)
		maxRetryIn = max(maxRetryIn, d.RetryIn)
		maxResetIn = max(maxResetIn, d.ResetIn)
	}

	// Conditionally, spend the batch.
	if len(newTATs) > 0 && allowed {
		err = l.source.BatchSet(ctx, newTATs)
		if err != nil {
			return nil, err
		}
	}

	// Consolidated decision for the batch.
	return &Decision{
		Allowed:   allowed,
		Remaining: minRemaining,
		RetryIn:   maxRetryIn,
		ResetIn:   maxResetIn,
	}, nil
}

// Refund attempts to refund the cost to the bucket identified by limit name and
// client id. The returned *Decision indicates whether the refund was successful
// or not. If the refund was successful, the cost of the request was added back
// to the bucket's capacity. If the refund is not possible (i.e., the bucket is
// already full or the refund amount is invalid), no cost is refunded.
//
// Note: The amount refunded cannot cause the bucket to exceed its maximum
// capacity. However, partial refunds are allowed and are considered successful.
// For instance, if a bucket has a maximum capacity of 10 and currently has 5
// requests remaining, a refund request of 7 will result in the bucket reaching
// its maximum capacity of 10, not 12.
func (l *Limiter) Refund(ctx context.Context, name Name, id string, cost int64) (*Decision, error) {
	if cost <= 0 {
		return nil, ErrInvalidCost
	}

	limit, err := l.getLimit(name, id)
	if err != nil {
		if errors.Is(err, errLimitDisabled) {
			return disabledLimitDecision, nil
		}
		return nil, err
	}

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	tat, err := l.source.Get(ctx, bucketKey(name, id))
	if err != nil {
		return nil, err
	}
	d := maybeRefund(l.clk, limit, tat, cost)
	if !d.Allowed {
		// The bucket is already at maximum capacity.
		return d, nil
	}
	return d, l.source.Set(ctx, bucketKey(name, id), d.newTAT)
}

// BatchRefund attempts to refund quota to the specified buckets in the batch.
// If a refund was successful for a bucket, the cost is added back to its
// capacity. If a refund is not possible for a bucket (e.g., already full,
// invalid amount), no cost is refunded for that bucket. The consolidated
// *Decision returned indicates whether at least 1 refund was successful
// (Allowed), the minimum remaining capacity across all buckets (Remaining), and
// the maximum RetryIn and ResetIn values across all buckets. If one or more
// buckets in the batch do not exist, they will be ignored.
func (l *Limiter) BatchRefund(ctx context.Context, batch Batch) (*Decision, error) {
	if len(batch) == 0 {
		return nil, ErrInvalidCost
	}

	bucketKeys := make([]string, 0, len(batch))
	for _, entry := range batch {
		if entry.Cost <= 0 {
			return nil, ErrInvalidCost
		}
		bucketKeys = append(bucketKeys, entry.bucketKey())
	}

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	tats, err := l.source.BatchGet(ctx, bucketKeys)
	if err != nil {
		return nil, err
	}

	var minRemaining int64 = math.MaxInt64
	var maxRetryIn time.Duration
	var maxResetIn time.Duration
	var allowed bool
	newTATs := make(map[string]time.Time)

	for _, entry := range batch {
		bucketKey := entry.bucketKey()
		tat, exists := tats[bucketKey]
		limit, err := l.getLimit(entry.Name, entry.Id)
		if err != nil {
			if errors.Is(err, errLimitDisabled) {
				// Ignore disabled limit.
				continue
			}
			return nil, err
		}

		if !exists || tat.IsZero() {
			// If the bucket no longer exists, ignore it. A missing bucket is
			// equivalent to a full bucket.
			continue
		}

		// Refund the cost and update the consolidated decision.
		d := maybeRefund(l.clk, limit, tat, entry.Cost)
		if d.Allowed {
			newTATs[bucketKey] = d.newTAT
		}

		// At least one refund must be allowed for the batch to be considered
		// allowed.
		allowed = allowed || d.Allowed
		minRemaining = min(minRemaining, d.Remaining)
		maxRetryIn = max(maxRetryIn, d.RetryIn)
		maxResetIn = max(maxResetIn, d.ResetIn)
	}

	// Conditionally, refund the batch.
	if len(newTATs) > 0 {
		err = l.source.BatchSet(ctx, newTATs)
		if err != nil {
			return nil, err
		}
	}

	// Consolidated decision for the batch.
	return &Decision{
		Allowed:   allowed,
		Remaining: minRemaining,
		RetryIn:   maxRetryIn,
		ResetIn:   maxResetIn,
	}, nil
}

// Reset resets the specified bucket.
func (l *Limiter) Reset(ctx context.Context, name Name, id string) error {
	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	return l.source.Delete(ctx, bucketKey(name, id))
}

// initialize creates a new bucket, specified by limit name and id, with the
// cost of the request factored into the initial state.
func (l *Limiter) initialize(ctx context.Context, rl limit, name Name, id string, cost int64) (*Decision, error) {
	d := maybeSpend(l.clk, rl, l.clk.Now(), cost)

	// Remove cancellation from the request context so that transactions are not
	// interrupted by a client disconnect.
	ctx = context.WithoutCancel(ctx)
	err := l.source.Set(ctx, bucketKey(name, id), d.newTAT)
	if err != nil {
		return nil, err
	}
	return d, nil

}

// GetLimit returns the limit for the specified by name and id, name is
// required, id is optional. If id is left unspecified, the default limit for
// the limit specified by name is returned. If no default limit exists for the
// specified name, ErrLimitDisabled is returned.
func (l *Limiter) getLimit(name Name, id string) (limit, error) {
	if !name.isValid() {
		// This should never happen. Callers should only be specifying the limit
		// Name enums defined in this package.
		return limit{}, fmt.Errorf("specified name enum %q, is invalid", name)
	}
	if id != "" {
		// Check for override.
		ol, ok := l.overrides[bucketKey(name, id)]
		if ok {
			return ol, nil
		}
	}
	dl, ok := l.defaults[nameToEnumString(name)]
	if ok {
		return dl, nil
	}
	return limit{}, errLimitDisabled
}
