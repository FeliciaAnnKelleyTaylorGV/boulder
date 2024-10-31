// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.20.1
// source: va.proto

package proto

import (
	proto "github.com/letsencrypt/boulder/core/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IsCAAValidRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// NOTE: Domain may be a name with a wildcard prefix (e.g. `*.example.com`)
	Domain           string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	ValidationMethod string `protobuf:"bytes,2,opt,name=validationMethod,proto3" json:"validationMethod,omitempty"`
	AccountURIID     int64  `protobuf:"varint,3,opt,name=accountURIID,proto3" json:"accountURIID,omitempty"`
	AuthzID          string `protobuf:"bytes,4,opt,name=authzID,proto3" json:"authzID,omitempty"`
}

func (x *IsCAAValidRequest) Reset() {
	*x = IsCAAValidRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_va_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsCAAValidRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsCAAValidRequest) ProtoMessage() {}

func (x *IsCAAValidRequest) ProtoReflect() protoreflect.Message {
	mi := &file_va_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsCAAValidRequest.ProtoReflect.Descriptor instead.
func (*IsCAAValidRequest) Descriptor() ([]byte, []int) {
	return file_va_proto_rawDescGZIP(), []int{0}
}

func (x *IsCAAValidRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *IsCAAValidRequest) GetValidationMethod() string {
	if x != nil {
		return x.ValidationMethod
	}
	return ""
}

func (x *IsCAAValidRequest) GetAccountURIID() int64 {
	if x != nil {
		return x.AccountURIID
	}
	return 0
}

func (x *IsCAAValidRequest) GetAuthzID() string {
	if x != nil {
		return x.AuthzID
	}
	return ""
}

// If CAA is valid for the requested domain, the problem will be empty
type IsCAAValidResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Problem *proto.ProblemDetails `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
}

func (x *IsCAAValidResponse) Reset() {
	*x = IsCAAValidResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_va_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsCAAValidResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsCAAValidResponse) ProtoMessage() {}

func (x *IsCAAValidResponse) ProtoReflect() protoreflect.Message {
	mi := &file_va_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsCAAValidResponse.ProtoReflect.Descriptor instead.
func (*IsCAAValidResponse) Descriptor() ([]byte, []int) {
	return file_va_proto_rawDescGZIP(), []int{1}
}

func (x *IsCAAValidResponse) GetProblem() *proto.ProblemDetails {
	if x != nil {
		return x.Problem
	}
	return nil
}

type AuthzMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegID int64  `protobuf:"varint,2,opt,name=regID,proto3" json:"regID,omitempty"`
}

func (x *AuthzMeta) Reset() {
	*x = AuthzMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_va_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthzMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthzMeta) ProtoMessage() {}

func (x *AuthzMeta) ProtoReflect() protoreflect.Message {
	mi := &file_va_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthzMeta.ProtoReflect.Descriptor instead.
func (*AuthzMeta) Descriptor() ([]byte, []int) {
	return file_va_proto_rawDescGZIP(), []int{2}
}

func (x *AuthzMeta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AuthzMeta) GetRegID() int64 {
	if x != nil {
		return x.RegID
	}
	return 0
}

type PerformValidationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DnsName                  string           `protobuf:"bytes,1,opt,name=dnsName,proto3" json:"dnsName,omitempty"`
	Challenge                *proto.Challenge `protobuf:"bytes,2,opt,name=challenge,proto3" json:"challenge,omitempty"`
	Authz                    *AuthzMeta       `protobuf:"bytes,3,opt,name=authz,proto3" json:"authz,omitempty"`
	ExpectedKeyAuthorization string           `protobuf:"bytes,4,opt,name=expectedKeyAuthorization,proto3" json:"expectedKeyAuthorization,omitempty"`
}

func (x *PerformValidationRequest) Reset() {
	*x = PerformValidationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_va_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PerformValidationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PerformValidationRequest) ProtoMessage() {}

func (x *PerformValidationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_va_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PerformValidationRequest.ProtoReflect.Descriptor instead.
func (*PerformValidationRequest) Descriptor() ([]byte, []int) {
	return file_va_proto_rawDescGZIP(), []int{3}
}

func (x *PerformValidationRequest) GetDnsName() string {
	if x != nil {
		return x.DnsName
	}
	return ""
}

func (x *PerformValidationRequest) GetChallenge() *proto.Challenge {
	if x != nil {
		return x.Challenge
	}
	return nil
}

func (x *PerformValidationRequest) GetAuthz() *AuthzMeta {
	if x != nil {
		return x.Authz
	}
	return nil
}

func (x *PerformValidationRequest) GetExpectedKeyAuthorization() string {
	if x != nil {
		return x.ExpectedKeyAuthorization
	}
	return ""
}

type ValidationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records     []*proto.ValidationRecord `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
	Problems    *proto.ProblemDetails     `protobuf:"bytes,2,opt,name=problems,proto3" json:"problems,omitempty"`
	Perspective string                    `protobuf:"bytes,3,opt,name=perspective,proto3" json:"perspective,omitempty"`
	Rir         string                    `protobuf:"bytes,4,opt,name=rir,proto3" json:"rir,omitempty"`
}

func (x *ValidationResult) Reset() {
	*x = ValidationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_va_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationResult) ProtoMessage() {}

func (x *ValidationResult) ProtoReflect() protoreflect.Message {
	mi := &file_va_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationResult.ProtoReflect.Descriptor instead.
func (*ValidationResult) Descriptor() ([]byte, []int) {
	return file_va_proto_rawDescGZIP(), []int{4}
}

func (x *ValidationResult) GetRecords() []*proto.ValidationRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

func (x *ValidationResult) GetProblems() *proto.ProblemDetails {
	if x != nil {
		return x.Problems
	}
	return nil
}

func (x *ValidationResult) GetPerspective() string {
	if x != nil {
		return x.Perspective
	}
	return ""
}

func (x *ValidationResult) GetRir() string {
	if x != nil {
		return x.Rir
	}
	return ""
}

var File_va_proto protoreflect.FileDescriptor

var file_va_proto_rawDesc = []byte{
	0x0a, 0x08, 0x76, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x61, 0x1a, 0x15,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a, 0x11, 0x49, 0x73, 0x43, 0x41, 0x41, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x12, 0x2a, 0x0a, 0x10, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x22, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x52, 0x49, 0x49, 0x44, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x52,
	0x49, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x49, 0x44, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x49, 0x44, 0x22, 0x44, 0x0a,
	0x12, 0x49, 0x73, 0x43, 0x41, 0x41, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x62,
	0x6c, 0x65, 0x6d, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x62,
	0x6c, 0x65, 0x6d, 0x22, 0x31, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x7a, 0x4d, 0x65, 0x74, 0x61,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x67, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x72, 0x65, 0x67, 0x49, 0x44, 0x22, 0xc4, 0x01, 0x0a, 0x18, 0x50, 0x65, 0x72, 0x66, 0x6f,
	0x72, 0x6d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x6e, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6e, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a,
	0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67,
	0x65, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x05,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x76, 0x61,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x7a, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x05, 0x61, 0x75, 0x74, 0x68,
	0x7a, 0x12, 0x3a, 0x0a, 0x18, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4b, 0x65, 0x79,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x18, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4b, 0x65, 0x79,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xaa, 0x01,
	0x0a, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x30, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x12, 0x30, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x72,
	0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x69, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x69, 0x72, 0x32, 0x4f, 0x0a, 0x02, 0x56, 0x41,
	0x12, 0x49, 0x0a, 0x11, 0x50, 0x65, 0x72, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x76, 0x61, 0x2e, 0x50, 0x65, 0x72, 0x66, 0x6f,
	0x72, 0x6d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x76, 0x61, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x32, 0x44, 0x0a, 0x03, 0x43,
	0x41, 0x41, 0x12, 0x3d, 0x0a, 0x0a, 0x49, 0x73, 0x43, 0x41, 0x41, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x12, 0x15, 0x2e, 0x76, 0x61, 0x2e, 0x49, 0x73, 0x43, 0x41, 0x41, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x76, 0x61, 0x2e, 0x49, 0x73, 0x43,
	0x41, 0x41, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6c, 0x65, 0x74, 0x73, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x2f, 0x62, 0x6f, 0x75, 0x6c,
	0x64, 0x65, 0x72, 0x2f, 0x76, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_va_proto_rawDescOnce sync.Once
	file_va_proto_rawDescData = file_va_proto_rawDesc
)

func file_va_proto_rawDescGZIP() []byte {
	file_va_proto_rawDescOnce.Do(func() {
		file_va_proto_rawDescData = protoimpl.X.CompressGZIP(file_va_proto_rawDescData)
	})
	return file_va_proto_rawDescData
}

var file_va_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_va_proto_goTypes = []interface{}{
	(*IsCAAValidRequest)(nil),        // 0: va.IsCAAValidRequest
	(*IsCAAValidResponse)(nil),       // 1: va.IsCAAValidResponse
	(*AuthzMeta)(nil),                // 2: va.AuthzMeta
	(*PerformValidationRequest)(nil), // 3: va.PerformValidationRequest
	(*ValidationResult)(nil),         // 4: va.ValidationResult
	(*proto.ProblemDetails)(nil),     // 5: core.ProblemDetails
	(*proto.Challenge)(nil),          // 6: core.Challenge
	(*proto.ValidationRecord)(nil),   // 7: core.ValidationRecord
}
var file_va_proto_depIdxs = []int32{
	5, // 0: va.IsCAAValidResponse.problem:type_name -> core.ProblemDetails
	6, // 1: va.PerformValidationRequest.challenge:type_name -> core.Challenge
	2, // 2: va.PerformValidationRequest.authz:type_name -> va.AuthzMeta
	7, // 3: va.ValidationResult.records:type_name -> core.ValidationRecord
	5, // 4: va.ValidationResult.problems:type_name -> core.ProblemDetails
	3, // 5: va.VA.PerformValidation:input_type -> va.PerformValidationRequest
	0, // 6: va.CAA.IsCAAValid:input_type -> va.IsCAAValidRequest
	4, // 7: va.VA.PerformValidation:output_type -> va.ValidationResult
	1, // 8: va.CAA.IsCAAValid:output_type -> va.IsCAAValidResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_va_proto_init() }
func file_va_proto_init() {
	if File_va_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_va_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsCAAValidRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_va_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsCAAValidResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_va_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthzMeta); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_va_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PerformValidationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_va_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_va_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_va_proto_goTypes,
		DependencyIndexes: file_va_proto_depIdxs,
		MessageInfos:      file_va_proto_msgTypes,
	}.Build()
	File_va_proto = out.File
	file_va_proto_rawDesc = nil
	file_va_proto_goTypes = nil
	file_va_proto_depIdxs = nil
}
