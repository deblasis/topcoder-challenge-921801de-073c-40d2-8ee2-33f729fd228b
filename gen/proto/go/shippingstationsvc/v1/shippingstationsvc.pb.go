// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: shippingstationsvc/v1/shippingstationsvc.proto

package shippingstationsvc_v1

import (
	v1 "deblasis.net/space-traffic-control/gen/proto/go/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
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

type RequestLandingResponse_Command int32

const (
	RequestLandingResponse_WAIT RequestLandingResponse_Command = 0
	RequestLandingResponse_LAND RequestLandingResponse_Command = 1
)

// Enum value maps for RequestLandingResponse_Command.
var (
	RequestLandingResponse_Command_name = map[int32]string{
		0: "WAIT",
		1: "LAND",
	}
	RequestLandingResponse_Command_value = map[string]int32{
		"WAIT": 0,
		"LAND": 1,
	}
)

func (x RequestLandingResponse_Command) Enum() *RequestLandingResponse_Command {
	p := new(RequestLandingResponse_Command)
	*p = x
	return p
}

func (x RequestLandingResponse_Command) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestLandingResponse_Command) Descriptor() protoreflect.EnumDescriptor {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_enumTypes[0].Descriptor()
}

func (RequestLandingResponse_Command) Type() protoreflect.EnumType {
	return &file_shippingstationsvc_v1_shippingstationsvc_proto_enumTypes[0]
}

func (x RequestLandingResponse_Command) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestLandingResponse_Command.Descriptor instead.
func (RequestLandingResponse_Command) EnumDescriptor() ([]byte, []int) {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP(), []int{1, 0}
}

type RequestLandingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//@gotags: validate:"uuid4"
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" validate:"uuid4"`
	//@gotags: validate:"required"
	Time int64 `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty" validate:"required"`
}

func (x *RequestLandingRequest) Reset() {
	*x = RequestLandingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestLandingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestLandingRequest) ProtoMessage() {}

func (x *RequestLandingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestLandingRequest.ProtoReflect.Descriptor instead.
func (*RequestLandingRequest) Descriptor() ([]byte, []int) {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP(), []int{0}
}

func (x *RequestLandingRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RequestLandingRequest) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type RequestLandingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command RequestLandingResponse_Command `protobuf:"varint,1,opt,name=command,proto3,enum=deblasis.v1.RequestLandingResponse_Command" json:"command,omitempty"`
	// Types that are assignable to DockingStationIdOrDuration:
	//	*RequestLandingResponse_DockingStationId
	//	*RequestLandingResponse_Duration
	DockingStationIdOrDuration isRequestLandingResponse_DockingStationIdOrDuration `protobuf_oneof:"docking_station_id_or_duration"`
	//@gotags: model:"-"
	Error *v1.Error `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty" model:"-"`
}

func (x *RequestLandingResponse) Reset() {
	*x = RequestLandingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestLandingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestLandingResponse) ProtoMessage() {}

func (x *RequestLandingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestLandingResponse.ProtoReflect.Descriptor instead.
func (*RequestLandingResponse) Descriptor() ([]byte, []int) {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP(), []int{1}
}

func (x *RequestLandingResponse) GetCommand() RequestLandingResponse_Command {
	if x != nil {
		return x.Command
	}
	return RequestLandingResponse_WAIT
}

func (m *RequestLandingResponse) GetDockingStationIdOrDuration() isRequestLandingResponse_DockingStationIdOrDuration {
	if m != nil {
		return m.DockingStationIdOrDuration
	}
	return nil
}

func (x *RequestLandingResponse) GetDockingStationId() string {
	if x, ok := x.GetDockingStationIdOrDuration().(*RequestLandingResponse_DockingStationId); ok {
		return x.DockingStationId
	}
	return ""
}

func (x *RequestLandingResponse) GetDuration() int64 {
	if x, ok := x.GetDockingStationIdOrDuration().(*RequestLandingResponse_Duration); ok {
		return x.Duration
	}
	return 0
}

func (x *RequestLandingResponse) GetError() *v1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type isRequestLandingResponse_DockingStationIdOrDuration interface {
	isRequestLandingResponse_DockingStationIdOrDuration()
}

type RequestLandingResponse_DockingStationId struct {
	DockingStationId string `protobuf:"bytes,2,opt,name=docking_station_id,json=dockingStationId,proto3,oneof"`
}

type RequestLandingResponse_Duration struct {
	Duration int64 `protobuf:"varint,3,opt,name=duration,proto3,oneof"`
}

func (*RequestLandingResponse_DockingStationId) isRequestLandingResponse_DockingStationIdOrDuration() {
}

func (*RequestLandingResponse_Duration) isRequestLandingResponse_DockingStationIdOrDuration() {}

type LandingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//@gotags: validate:"uuid4,required"
	ShipId string `protobuf:"bytes,1,opt,name=ship_id,json=shipId,proto3" json:"ship_id,omitempty" validate:"uuid4,required"`
	//@gotags: validate:"uuid4,required"
	DockId string `protobuf:"bytes,2,opt,name=dock_id,json=dockId,proto3" json:"dock_id,omitempty" validate:"uuid4,required"`
	//@gotags: validate:"required"
	Time int64 `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty" validate:"required"`
}

func (x *LandingRequest) Reset() {
	*x = LandingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LandingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LandingRequest) ProtoMessage() {}

func (x *LandingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LandingRequest.ProtoReflect.Descriptor instead.
func (*LandingRequest) Descriptor() ([]byte, []int) {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP(), []int{2}
}

func (x *LandingRequest) GetShipId() string {
	if x != nil {
		return x.ShipId
	}
	return ""
}

func (x *LandingRequest) GetDockId() string {
	if x != nil {
		return x.DockId
	}
	return ""
}

func (x *LandingRequest) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type LandingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//@gotags: model:"-"
	Error *v1.Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty" model:"-"`
}

func (x *LandingResponse) Reset() {
	*x = LandingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LandingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LandingResponse) ProtoMessage() {}

func (x *LandingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LandingResponse.ProtoReflect.Descriptor instead.
func (*LandingResponse) Descriptor() ([]byte, []int) {
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP(), []int{3}
}

func (x *LandingResponse) GetError() *v1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_shippingstationsvc_v1_shippingstationsvc_proto protoreflect.FileDescriptor

var file_shippingstationsvc_v1_shippingstationsvc_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x0e, 0x76,
	0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x62, 0x6f,
	0x64, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x15, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x9f, 0x02, 0x0a, 0x16, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4c,
	0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45,
	0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x2b, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x2e, 0x0a, 0x12, 0x64, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x10, 0x64, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x1d, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12,
	0x08, 0x0a, 0x04, 0x57, 0x41, 0x49, 0x54, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x41, 0x4e,
	0x44, 0x10, 0x01, 0x42, 0x20, 0x0a, 0x1e, 0x64, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x5f, 0x6f, 0x72, 0x5f, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x56, 0x0a, 0x0e, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x68, 0x69, 0x70, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x68, 0x69, 0x70, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x64, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x42, 0x0a,
	0x0f, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x32, 0xf3, 0x01, 0x0a, 0x16, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x78, 0x0a, 0x0e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x22,
	0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x48, 0x74, 0x74, 0x70, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x2c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26,
	0x22, 0x21, 0x2f, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2d, 0x6c, 0x61, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x3a, 0x01, 0x2a, 0x12, 0x5f, 0x0a, 0x07, 0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x12, 0x1b, 0x2e, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x61, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x48, 0x74, 0x74, 0x70,
	0x42, 0x6f, 0x64, 0x79, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x22, 0x16, 0x2f, 0x73,
	0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x2d, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f,
	0x6c, 0x61, 0x6e, 0x64, 0x3a, 0x01, 0x2a, 0x42, 0x63, 0x5a, 0x61, 0x64, 0x65, 0x62, 0x6c, 0x61,
	0x73, 0x69, 0x73, 0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2d, 0x74, 0x72,
	0x61, 0x66, 0x66, 0x69, 0x63, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x67, 0x65,
	0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x73, 0x68, 0x69, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x76, 0x63, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescOnce sync.Once
	file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescData = file_shippingstationsvc_v1_shippingstationsvc_proto_rawDesc
)

func file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescGZIP() []byte {
	file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescOnce.Do(func() {
		file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescData)
	})
	return file_shippingstationsvc_v1_shippingstationsvc_proto_rawDescData
}

var file_shippingstationsvc_v1_shippingstationsvc_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_shippingstationsvc_v1_shippingstationsvc_proto_goTypes = []interface{}{
	(RequestLandingResponse_Command)(0), // 0: deblasis.v1.RequestLandingResponse.Command
	(*RequestLandingRequest)(nil),       // 1: deblasis.v1.RequestLandingRequest
	(*RequestLandingResponse)(nil),      // 2: deblasis.v1.RequestLandingResponse
	(*LandingRequest)(nil),              // 3: deblasis.v1.LandingRequest
	(*LandingResponse)(nil),             // 4: deblasis.v1.LandingResponse
	(*v1.Error)(nil),                    // 5: deblasis.common.v1.Error
	(*httpbody.HttpBody)(nil),           // 6: google.api.HttpBody
}
var file_shippingstationsvc_v1_shippingstationsvc_proto_depIdxs = []int32{
	0, // 0: deblasis.v1.RequestLandingResponse.command:type_name -> deblasis.v1.RequestLandingResponse.Command
	5, // 1: deblasis.v1.RequestLandingResponse.error:type_name -> deblasis.common.v1.Error
	5, // 2: deblasis.v1.LandingResponse.error:type_name -> deblasis.common.v1.Error
	1, // 3: deblasis.v1.ShippingStationService.RequestLanding:input_type -> deblasis.v1.RequestLandingRequest
	3, // 4: deblasis.v1.ShippingStationService.Landing:input_type -> deblasis.v1.LandingRequest
	6, // 5: deblasis.v1.ShippingStationService.RequestLanding:output_type -> google.api.HttpBody
	6, // 6: deblasis.v1.ShippingStationService.Landing:output_type -> google.api.HttpBody
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_shippingstationsvc_v1_shippingstationsvc_proto_init() }
func file_shippingstationsvc_v1_shippingstationsvc_proto_init() {
	if File_shippingstationsvc_v1_shippingstationsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestLandingRequest); i {
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
		file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestLandingResponse); i {
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
		file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LandingRequest); i {
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
		file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LandingResponse); i {
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
	file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*RequestLandingResponse_DockingStationId)(nil),
		(*RequestLandingResponse_Duration)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shippingstationsvc_v1_shippingstationsvc_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shippingstationsvc_v1_shippingstationsvc_proto_goTypes,
		DependencyIndexes: file_shippingstationsvc_v1_shippingstationsvc_proto_depIdxs,
		EnumInfos:         file_shippingstationsvc_v1_shippingstationsvc_proto_enumTypes,
		MessageInfos:      file_shippingstationsvc_v1_shippingstationsvc_proto_msgTypes,
	}.Build()
	File_shippingstationsvc_v1_shippingstationsvc_proto = out.File
	file_shippingstationsvc_v1_shippingstationsvc_proto_rawDesc = nil
	file_shippingstationsvc_v1_shippingstationsvc_proto_goTypes = nil
	file_shippingstationsvc_v1_shippingstationsvc_proto_depIdxs = nil
}
