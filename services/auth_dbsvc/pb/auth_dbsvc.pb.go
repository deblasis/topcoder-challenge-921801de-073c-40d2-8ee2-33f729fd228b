// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: auth_dbsvc.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type User_Role int32

const (
	User_UNSPECIFIED User_Role = 0
	User_SHIP        User_Role = 1
	User_STATION     User_Role = 2
	User_COMMAND     User_Role = 3
)

// Enum value maps for User_Role.
var (
	User_Role_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "SHIP",
		2: "STATION",
		3: "COMMAND",
	}
	User_Role_value = map[string]int32{
		"UNSPECIFIED": 0,
		"SHIP":        1,
		"STATION":     2,
		"COMMAND":     3,
	}
)

func (x User_Role) Enum() *User_Role {
	p := new(User_Role)
	*p = x
	return p
}

func (x User_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (User_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_auth_dbsvc_proto_enumTypes[0].Descriptor()
}

func (User_Role) Type() protoreflect.EnumType {
	return &file_auth_dbsvc_proto_enumTypes[0]
}

func (x User_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use User_Role.Descriptor instead.
func (User_Role) EnumDescriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{0, 0}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username string    `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password string    `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Role     User_Role `protobuf:"varint,4,opt,name=role,proto3,enum=User_Role" json:"role,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *User) GetRole() User_Role {
	if x != nil {
		return x.Role
	}
	return User_UNSPECIFIED
}

type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type CreateUserReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CreateUserReply) Reset() {
	*x = CreateUserReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserReply) ProtoMessage() {}

func (x *CreateUserReply) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserReply.ProtoReflect.Descriptor instead.
func (*CreateUserReply) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUserReply) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CreateUserReply) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetUserByUsernameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *GetUserByUsernameRequest) Reset() {
	*x = GetUserByUsernameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByUsernameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByUsernameRequest) ProtoMessage() {}

func (x *GetUserByUsernameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByUsernameRequest.ProtoReflect.Descriptor instead.
func (*GetUserByUsernameRequest) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserByUsernameRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type UserReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User  *User  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UserReply) Reset() {
	*x = UserReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserReply) ProtoMessage() {}

func (x *UserReply) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserReply.ProtoReflect.Descriptor instead.
func (*UserReply) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{4}
}

func (x *UserReply) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UserReply) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type ServiceStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServiceStatusRequest) Reset() {
	*x = ServiceStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStatusRequest) ProtoMessage() {}

func (x *ServiceStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStatusRequest.ProtoReflect.Descriptor instead.
func (*ServiceStatusRequest) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{5}
}

type ServiceStatusReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Err  string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *ServiceStatusReply) Reset() {
	*x = ServiceStatusReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_dbsvc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStatusReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStatusReply) ProtoMessage() {}

func (x *ServiceStatusReply) ProtoReflect() protoreflect.Message {
	mi := &file_auth_dbsvc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStatusReply.ProtoReflect.Descriptor instead.
func (*ServiceStatusReply) Descriptor() ([]byte, []int) {
	return file_auth_dbsvc_proto_rawDescGZIP(), []int{6}
}

func (x *ServiceStatusReply) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ServiceStatusReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_auth_dbsvc_proto protoreflect.FileDescriptor

var file_auth_dbsvc_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x64, 0x62, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0a, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x22, 0x3b, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x53, 0x48, 0x49, 0x50, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x54, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x03,
	0x22, 0x2e, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x22, 0x37, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x36, 0x0a, 0x18, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x3c, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x19,
	0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22,
	0x16, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x0a, 0x12, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x65, 0x72, 0x72, 0x32, 0xbe, 0x01, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x44, 0x42, 0x53, 0x76,
	0x63, 0x12, 0x3d, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x15, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x12, 0x34, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x42, 0x18, 0x5a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x64, 0x62, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_dbsvc_proto_rawDescOnce sync.Once
	file_auth_dbsvc_proto_rawDescData = file_auth_dbsvc_proto_rawDesc
)

func file_auth_dbsvc_proto_rawDescGZIP() []byte {
	file_auth_dbsvc_proto_rawDescOnce.Do(func() {
		file_auth_dbsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_dbsvc_proto_rawDescData)
	})
	return file_auth_dbsvc_proto_rawDescData
}

var file_auth_dbsvc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_auth_dbsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_auth_dbsvc_proto_goTypes = []interface{}{
	(User_Role)(0),                   // 0: User.Role
	(*User)(nil),                     // 1: User
	(*CreateUserRequest)(nil),        // 2: CreateUserRequest
	(*CreateUserReply)(nil),          // 3: CreateUserReply
	(*GetUserByUsernameRequest)(nil), // 4: GetUserByUsernameRequest
	(*UserReply)(nil),                // 5: UserReply
	(*ServiceStatusRequest)(nil),     // 6: ServiceStatusRequest
	(*ServiceStatusReply)(nil),       // 7: ServiceStatusReply
}
var file_auth_dbsvc_proto_depIdxs = []int32{
	0, // 0: User.role:type_name -> User.Role
	1, // 1: CreateUserRequest.user:type_name -> User
	1, // 2: UserReply.user:type_name -> User
	6, // 3: AuthDBSvc.ServiceStatus:input_type -> ServiceStatusRequest
	2, // 4: AuthDBSvc.CreateUser:input_type -> CreateUserRequest
	4, // 5: AuthDBSvc.GetUserByUsername:input_type -> GetUserByUsernameRequest
	7, // 6: AuthDBSvc.ServiceStatus:output_type -> ServiceStatusReply
	3, // 7: AuthDBSvc.CreateUser:output_type -> CreateUserReply
	5, // 8: AuthDBSvc.GetUserByUsername:output_type -> UserReply
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_auth_dbsvc_proto_init() }
func file_auth_dbsvc_proto_init() {
	if File_auth_dbsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_dbsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_auth_dbsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserRequest); i {
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
		file_auth_dbsvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserReply); i {
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
		file_auth_dbsvc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserByUsernameRequest); i {
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
		file_auth_dbsvc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserReply); i {
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
		file_auth_dbsvc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStatusRequest); i {
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
		file_auth_dbsvc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStatusReply); i {
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
			RawDescriptor: file_auth_dbsvc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_dbsvc_proto_goTypes,
		DependencyIndexes: file_auth_dbsvc_proto_depIdxs,
		EnumInfos:         file_auth_dbsvc_proto_enumTypes,
		MessageInfos:      file_auth_dbsvc_proto_msgTypes,
	}.Build()
	File_auth_dbsvc_proto = out.File
	file_auth_dbsvc_proto_rawDesc = nil
	file_auth_dbsvc_proto_goTypes = nil
	file_auth_dbsvc_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthDBSvcClient is the client API for AuthDBSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthDBSvcClient interface {
	ServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusReply, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error)
	GetUserByUsername(ctx context.Context, in *GetUserByUsernameRequest, opts ...grpc.CallOption) (*UserReply, error)
}

type authDBSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthDBSvcClient(cc grpc.ClientConnInterface) AuthDBSvcClient {
	return &authDBSvcClient{cc}
}

func (c *authDBSvcClient) ServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusReply, error) {
	out := new(ServiceStatusReply)
	err := c.cc.Invoke(ctx, "/AuthDBSvc/ServiceStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authDBSvcClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserReply, error) {
	out := new(CreateUserReply)
	err := c.cc.Invoke(ctx, "/AuthDBSvc/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authDBSvcClient) GetUserByUsername(ctx context.Context, in *GetUserByUsernameRequest, opts ...grpc.CallOption) (*UserReply, error) {
	out := new(UserReply)
	err := c.cc.Invoke(ctx, "/AuthDBSvc/GetUserByUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthDBSvcServer is the server API for AuthDBSvc service.
type AuthDBSvcServer interface {
	ServiceStatus(context.Context, *ServiceStatusRequest) (*ServiceStatusReply, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error)
	GetUserByUsername(context.Context, *GetUserByUsernameRequest) (*UserReply, error)
}

// UnimplementedAuthDBSvcServer can be embedded to have forward compatible implementations.
type UnimplementedAuthDBSvcServer struct {
}

func (*UnimplementedAuthDBSvcServer) ServiceStatus(context.Context, *ServiceStatusRequest) (*ServiceStatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceStatus not implemented")
}
func (*UnimplementedAuthDBSvcServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (*UnimplementedAuthDBSvcServer) GetUserByUsername(context.Context, *GetUserByUsernameRequest) (*UserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByUsername not implemented")
}

func RegisterAuthDBSvcServer(s *grpc.Server, srv AuthDBSvcServer) {
	s.RegisterService(&_AuthDBSvc_serviceDesc, srv)
}

func _AuthDBSvc_ServiceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthDBSvcServer).ServiceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthDBSvc/ServiceStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthDBSvcServer).ServiceStatus(ctx, req.(*ServiceStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthDBSvc_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthDBSvcServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthDBSvc/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthDBSvcServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthDBSvc_GetUserByUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthDBSvcServer).GetUserByUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AuthDBSvc/GetUserByUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthDBSvcServer).GetUserByUsername(ctx, req.(*GetUserByUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthDBSvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "AuthDBSvc",
	HandlerType: (*AuthDBSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ServiceStatus",
			Handler:    _AuthDBSvc_ServiceStatus_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _AuthDBSvc_CreateUser_Handler,
		},
		{
			MethodName: "GetUserByUsername",
			Handler:    _AuthDBSvc_GetUserByUsername_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth_dbsvc.proto",
}
