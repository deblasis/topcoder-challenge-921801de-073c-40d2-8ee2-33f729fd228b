// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: authsvc/v1/authsvc.proto

package authsvc_v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//@gotags: validate:"required,notblank"
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty" validate:"required,notblank"`
	//@gotags: validate:"required,notblank"
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" validate:"required,notblank"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *LoginResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type SignupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//@gotags: validate:"required,notblank"
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty" validate:"required,notblank"`
	//@gotags: validate:"required,notblank"
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" validate:"required,notblank"`
	//@gotags: validate:"required,oneof=Ship Station Command"
	Role string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty" validate:"required,oneof=Ship Station Command"`
}

func (x *SignupRequest) Reset() {
	*x = SignupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupRequest) ProtoMessage() {}

func (x *SignupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupRequest.ProtoReflect.Descriptor instead.
func (*SignupRequest) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{2}
}

func (x *SignupRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SignupRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *SignupRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type SignupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token *Token `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SignupResponse) Reset() {
	*x = SignupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupResponse) ProtoMessage() {}

func (x *SignupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupResponse.ProtoReflect.Descriptor instead.
func (*SignupResponse) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{3}
}

func (x *SignupResponse) GetToken() *Token {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *SignupResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type CheckTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CheckTokenRequest) Reset() {
	*x = CheckTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenRequest) ProtoMessage() {}

func (x *CheckTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenRequest.ProtoReflect.Descriptor instead.
func (*CheckTokenRequest) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{4}
}

func (x *CheckTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type CheckTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenPayload *TokenPayload `protobuf:"bytes,1,opt,name=token_payload,json=tokenPayload,proto3" json:"token_payload,omitempty"`
	Error        string        `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CheckTokenResponse) Reset() {
	*x = CheckTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenResponse) ProtoMessage() {}

func (x *CheckTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenResponse.ProtoReflect.Descriptor instead.
func (*CheckTokenResponse) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{5}
}

func (x *CheckTokenResponse) GetTokenPayload() *TokenPayload {
	if x != nil {
		return x.TokenPayload
	}
	return nil
}

func (x *CheckTokenResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token     string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresAt int64  `protobuf:"varint,2,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{6}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Token) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

type TokenPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenId  string `protobuf:"bytes,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
	UserId   int64  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Role     string `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *TokenPayload) Reset() {
	*x = TokenPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authsvc_v1_authsvc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenPayload) ProtoMessage() {}

func (x *TokenPayload) ProtoReflect() protoreflect.Message {
	mi := &file_authsvc_v1_authsvc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenPayload.ProtoReflect.Descriptor instead.
func (*TokenPayload) Descriptor() ([]byte, []int) {
	return file_authsvc_v1_authsvc_proto_rawDescGZIP(), []int{7}
}

func (x *TokenPayload) GetTokenId() string {
	if x != nil {
		return x.TokenId
	}
	return ""
}

func (x *TokenPayload) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *TokenPayload) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *TokenPayload) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

var File_authsvc_v1_authsvc_proto protoreflect.FileDescriptor

var file_authsvc_v1_authsvc_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x75, 0x74, 0x68, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x64, 0x65, 0x62, 0x6c,
	0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x52, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x02,
	0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x4f, 0x0a, 0x0d, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x65, 0x62,
	0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x6d, 0x0a, 0x0d, 0x53,
	0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04,
	0xe2, 0x41, 0x01, 0x02, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x18, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04,
	0xe2, 0x41, 0x01, 0x02, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x50, 0x0a, 0x0e, 0x53, 0x69,
	0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x64, 0x65,
	0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x2f, 0x0a, 0x11,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6a, 0x0a,
	0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x65, 0x62,
	0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x0c, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3c, 0x0a, 0x05, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x73, 0x41, 0x74, 0x22, 0x72, 0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x32, 0xc1, 0x01, 0x0a, 0x0b,
	0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x05, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x12, 0x19, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x3a, 0x01, 0x2a, 0x12, 0x5a, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x1a, 0x2e,
	0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x65, 0x62, 0x6c,
	0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x22, 0x0c,
	0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x3a, 0x01, 0x2a, 0x42,
	0x4d, 0x5a, 0x4b, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x6e, 0x65, 0x74, 0x2f,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x2d, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2d, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x73, 0x76, 0x63, 0x5f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authsvc_v1_authsvc_proto_rawDescOnce sync.Once
	file_authsvc_v1_authsvc_proto_rawDescData = file_authsvc_v1_authsvc_proto_rawDesc
)

func file_authsvc_v1_authsvc_proto_rawDescGZIP() []byte {
	file_authsvc_v1_authsvc_proto_rawDescOnce.Do(func() {
		file_authsvc_v1_authsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_authsvc_v1_authsvc_proto_rawDescData)
	})
	return file_authsvc_v1_authsvc_proto_rawDescData
}

var file_authsvc_v1_authsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_authsvc_v1_authsvc_proto_goTypes = []interface{}{
	(*LoginRequest)(nil),       // 0: deblasis.v1.LoginRequest
	(*LoginResponse)(nil),      // 1: deblasis.v1.LoginResponse
	(*SignupRequest)(nil),      // 2: deblasis.v1.SignupRequest
	(*SignupResponse)(nil),     // 3: deblasis.v1.SignupResponse
	(*CheckTokenRequest)(nil),  // 4: deblasis.v1.CheckTokenRequest
	(*CheckTokenResponse)(nil), // 5: deblasis.v1.CheckTokenResponse
	(*Token)(nil),              // 6: deblasis.v1.Token
	(*TokenPayload)(nil),       // 7: deblasis.v1.TokenPayload
}
var file_authsvc_v1_authsvc_proto_depIdxs = []int32{
	6, // 0: deblasis.v1.LoginResponse.token:type_name -> deblasis.v1.Token
	6, // 1: deblasis.v1.SignupResponse.token:type_name -> deblasis.v1.Token
	7, // 2: deblasis.v1.CheckTokenResponse.token_payload:type_name -> deblasis.v1.TokenPayload
	0, // 3: deblasis.v1.AuthService.Login:input_type -> deblasis.v1.LoginRequest
	2, // 4: deblasis.v1.AuthService.Signup:input_type -> deblasis.v1.SignupRequest
	1, // 5: deblasis.v1.AuthService.Login:output_type -> deblasis.v1.LoginResponse
	3, // 6: deblasis.v1.AuthService.Signup:output_type -> deblasis.v1.SignupResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_authsvc_v1_authsvc_proto_init() }
func file_authsvc_v1_authsvc_proto_init() {
	if File_authsvc_v1_authsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_authsvc_v1_authsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRequest); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResponse); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignupRequest); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignupResponse); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenRequest); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenResponse); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
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
		file_authsvc_v1_authsvc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenPayload); i {
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
			RawDescriptor: file_authsvc_v1_authsvc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_authsvc_v1_authsvc_proto_goTypes,
		DependencyIndexes: file_authsvc_v1_authsvc_proto_depIdxs,
		MessageInfos:      file_authsvc_v1_authsvc_proto_msgTypes,
	}.Build()
	File_authsvc_v1_authsvc_proto = out.File
	file_authsvc_v1_authsvc_proto_rawDesc = nil
	file_authsvc_v1_authsvc_proto_goTypes = nil
	file_authsvc_v1_authsvc_proto_depIdxs = nil
}
