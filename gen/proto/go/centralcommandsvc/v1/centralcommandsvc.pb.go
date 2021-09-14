// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: centralcommandsvc/v1/centralcommandsvc.proto

package centralcommandsvc_v1

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

type Dock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	StationId       string  `protobuf:"bytes,2,opt,name=station_id,json=stationId,proto3" json:"station_id,omitempty"`
	NumDockingPorts int64   `protobuf:"varint,3,opt,name=num_docking_ports,json=numDockingPorts,proto3" json:"num_docking_ports,omitempty"`
	Occupied        int64   `protobuf:"varint,4,opt,name=occupied,proto3" json:"occupied,omitempty"`
	Weight          float32 `protobuf:"fixed32,5,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *Dock) Reset() {
	*x = Dock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dock) ProtoMessage() {}

func (x *Dock) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dock.ProtoReflect.Descriptor instead.
func (*Dock) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{0}
}

func (x *Dock) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Dock) GetStationId() string {
	if x != nil {
		return x.StationId
	}
	return ""
}

func (x *Dock) GetNumDockingPorts() int64 {
	if x != nil {
		return x.NumDockingPorts
	}
	return 0
}

func (x *Dock) GetOccupied() int64 {
	if x != nil {
		return x.Occupied
	}
	return 0
}

func (x *Dock) GetWeight() float32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type Station struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Capacity     float32 `protobuf:"fixed32,2,opt,name=capacity,proto3" json:"capacity,omitempty"`
	UsedCapacity float32 `protobuf:"fixed32,3,opt,name=used_capacity,json=usedCapacity,proto3" json:"used_capacity,omitempty"`
	Docks        []*Dock `protobuf:"bytes,4,rep,name=docks,proto3" json:"docks,omitempty"`
}

func (x *Station) Reset() {
	*x = Station{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Station) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Station) ProtoMessage() {}

func (x *Station) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Station.ProtoReflect.Descriptor instead.
func (*Station) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{1}
}

func (x *Station) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Station) GetCapacity() float32 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

func (x *Station) GetUsedCapacity() float32 {
	if x != nil {
		return x.UsedCapacity
	}
	return 0
}

func (x *Station) GetDocks() []*Dock {
	if x != nil {
		return x.Docks
	}
	return nil
}

type Ship struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	//@gotags: model:"-"
	Status string  `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty" model:"-"`
	Weight float32 `protobuf:"fixed32,3,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *Ship) Reset() {
	*x = Ship{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ship) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ship) ProtoMessage() {}

func (x *Ship) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ship.ProtoReflect.Descriptor instead.
func (*Ship) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{2}
}

func (x *Ship) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Ship) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Ship) GetWeight() float32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type RegisterStationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Station *Station `protobuf:"bytes,1,opt,name=station,proto3" json:"station,omitempty"`
}

func (x *RegisterStationRequest) Reset() {
	*x = RegisterStationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterStationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterStationRequest) ProtoMessage() {}

func (x *RegisterStationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterStationRequest.ProtoReflect.Descriptor instead.
func (*RegisterStationRequest) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterStationRequest) GetStation() *Station {
	if x != nil {
		return x.Station
	}
	return nil
}

type RegisterStationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Station *Station `protobuf:"bytes,1,opt,name=station,proto3" json:"station,omitempty"`
	Error   string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *RegisterStationResponse) Reset() {
	*x = RegisterStationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterStationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterStationResponse) ProtoMessage() {}

func (x *RegisterStationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterStationResponse.ProtoReflect.Descriptor instead.
func (*RegisterStationResponse) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterStationResponse) GetStation() *Station {
	if x != nil {
		return x.Station
	}
	return nil
}

func (x *RegisterStationResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type RegisterShipRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ship *Ship `protobuf:"bytes,1,opt,name=ship,proto3" json:"ship,omitempty"`
}

func (x *RegisterShipRequest) Reset() {
	*x = RegisterShipRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterShipRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterShipRequest) ProtoMessage() {}

func (x *RegisterShipRequest) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterShipRequest.ProtoReflect.Descriptor instead.
func (*RegisterShipRequest) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterShipRequest) GetShip() *Ship {
	if x != nil {
		return x.Ship
	}
	return nil
}

type RegisterShipResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ship  *Ship  `protobuf:"bytes,1,opt,name=ship,proto3" json:"ship,omitempty"`
	Error string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *RegisterShipResponse) Reset() {
	*x = RegisterShipResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterShipResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterShipResponse) ProtoMessage() {}

func (x *RegisterShipResponse) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterShipResponse.ProtoReflect.Descriptor instead.
func (*RegisterShipResponse) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{6}
}

func (x *RegisterShipResponse) GetShip() *Ship {
	if x != nil {
		return x.Ship
	}
	return nil
}

func (x *RegisterShipResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetAllShipsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllShipsRequest) Reset() {
	*x = GetAllShipsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllShipsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllShipsRequest) ProtoMessage() {}

func (x *GetAllShipsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllShipsRequest.ProtoReflect.Descriptor instead.
func (*GetAllShipsRequest) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{7}
}

type GetAllShipsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ships []*Ship `protobuf:"bytes,1,rep,name=ships,proto3" json:"ships,omitempty"`
	Error string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetAllShipsResponse) Reset() {
	*x = GetAllShipsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllShipsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllShipsResponse) ProtoMessage() {}

func (x *GetAllShipsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllShipsResponse.ProtoReflect.Descriptor instead.
func (*GetAllShipsResponse) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{8}
}

func (x *GetAllShipsResponse) GetShips() []*Ship {
	if x != nil {
		return x.Ships
	}
	return nil
}

func (x *GetAllShipsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetAllStationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllStationsRequest) Reset() {
	*x = GetAllStationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllStationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllStationsRequest) ProtoMessage() {}

func (x *GetAllStationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllStationsRequest.ProtoReflect.Descriptor instead.
func (*GetAllStationsRequest) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{9}
}

type GetAllStationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stations []*Station `protobuf:"bytes,1,rep,name=stations,proto3" json:"stations,omitempty"`
	Error    string     `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetAllStationsResponse) Reset() {
	*x = GetAllStationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllStationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllStationsResponse) ProtoMessage() {}

func (x *GetAllStationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllStationsResponse.ProtoReflect.Descriptor instead.
func (*GetAllStationsResponse) Descriptor() ([]byte, []int) {
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP(), []int{10}
}

func (x *GetAllStationsResponse) GetStations() []*Station {
	if x != nil {
		return x.Stations
	}
	return nil
}

func (x *GetAllStationsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_centralcommandsvc_v1_centralcommandsvc_proto protoreflect.FileDescriptor

var file_centralcommandsvc_v1_centralcommandsvc_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x73, 0x76, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14,
	0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76,
	0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x04, 0x44, 0x6f, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x11, 0x6e,
	0x75, 0x6d, 0x5f, 0x64, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x0f, 0x6e, 0x75,
	0x6d, 0x44, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x69, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x6f, 0x63, 0x63, 0x75, 0x70, 0x69, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x22, 0x98, 0x01, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a,
	0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x42,
	0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12,
	0x23, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x64, 0x43, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x12, 0x36, 0x0a, 0x05, 0x64, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x42,
	0x04, 0xe2, 0x41, 0x01, 0x02, 0x52, 0x05, 0x64, 0x6f, 0x63, 0x6b, 0x73, 0x22, 0x4c, 0x0a, 0x04,
	0x53, 0x68, 0x69, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x06,
	0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x42, 0x04, 0xe2, 0x41,
	0x01, 0x02, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x51, 0x0a, 0x16, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x68, 0x0a,
	0x17, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x45, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x53, 0x68, 0x69, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e,
	0x0a, 0x04, 0x73, 0x68, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63,
	0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x69, 0x70, 0x52, 0x04, 0x73, 0x68, 0x69, 0x70, 0x22, 0x5c,
	0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x68, 0x69, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x73, 0x68, 0x69, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x69, 0x70,
	0x52, 0x04, 0x73, 0x68, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x14, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x68, 0x69, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x5d, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x68, 0x69, 0x70,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x68, 0x69,
	0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72,
	0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x68, 0x69, 0x70, 0x52, 0x05, 0x73, 0x68, 0x69, 0x70, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x69, 0x0a, 0x16, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xcd, 0x04, 0x0a, 0x15, 0x43, 0x65, 0x6e, 0x74, 0x72, 0x61,
	0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x9a, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x22, 0x19, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x3a, 0x07, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x8b, 0x01, 0x0a,
	0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x68, 0x69, 0x70, 0x12, 0x29, 0x2e,
	0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x68, 0x69,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72,
	0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x22, 0x16, 0x2f, 0x63,
	0x65, 0x6e, 0x74, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x69, 0x70, 0x2f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x3a, 0x04, 0x73, 0x68, 0x69, 0x70, 0x12, 0x7d, 0x0a, 0x0b, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x68, 0x69, 0x70, 0x73, 0x12, 0x28, 0x2e, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x68, 0x69, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x53, 0x68, 0x69, 0x70, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x19,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x68, 0x69, 0x70, 0x2f, 0x61, 0x6c, 0x6c, 0x12, 0x89, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2b, 0x2e, 0x63,
	0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x63, 0x65, 0x6e, 0x74,
	0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12,
	0x14, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x61, 0x6c, 0x6c, 0x42, 0x61, 0x5a, 0x5f, 0x64, 0x65, 0x62, 0x6c, 0x61, 0x73, 0x69,
	0x73, 0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2d, 0x74, 0x72, 0x61, 0x66,
	0x66, 0x69, 0x63, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x76, 0x31, 0x3b, 0x63, 0x65, 0x6e, 0x74, 0x72, 0x61, 0x6c, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x73, 0x76, 0x63, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescOnce sync.Once
	file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescData = file_centralcommandsvc_v1_centralcommandsvc_proto_rawDesc
)

func file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescGZIP() []byte {
	file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescOnce.Do(func() {
		file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescData = protoimpl.X.CompressGZIP(file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescData)
	})
	return file_centralcommandsvc_v1_centralcommandsvc_proto_rawDescData
}

var file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_centralcommandsvc_v1_centralcommandsvc_proto_goTypes = []interface{}{
	(*Dock)(nil),                    // 0: centralcommandsvc.v1.Dock
	(*Station)(nil),                 // 1: centralcommandsvc.v1.Station
	(*Ship)(nil),                    // 2: centralcommandsvc.v1.Ship
	(*RegisterStationRequest)(nil),  // 3: centralcommandsvc.v1.RegisterStationRequest
	(*RegisterStationResponse)(nil), // 4: centralcommandsvc.v1.RegisterStationResponse
	(*RegisterShipRequest)(nil),     // 5: centralcommandsvc.v1.RegisterShipRequest
	(*RegisterShipResponse)(nil),    // 6: centralcommandsvc.v1.RegisterShipResponse
	(*GetAllShipsRequest)(nil),      // 7: centralcommandsvc.v1.GetAllShipsRequest
	(*GetAllShipsResponse)(nil),     // 8: centralcommandsvc.v1.GetAllShipsResponse
	(*GetAllStationsRequest)(nil),   // 9: centralcommandsvc.v1.GetAllStationsRequest
	(*GetAllStationsResponse)(nil),  // 10: centralcommandsvc.v1.GetAllStationsResponse
}
var file_centralcommandsvc_v1_centralcommandsvc_proto_depIdxs = []int32{
	0,  // 0: centralcommandsvc.v1.Station.docks:type_name -> centralcommandsvc.v1.Dock
	1,  // 1: centralcommandsvc.v1.RegisterStationRequest.station:type_name -> centralcommandsvc.v1.Station
	1,  // 2: centralcommandsvc.v1.RegisterStationResponse.station:type_name -> centralcommandsvc.v1.Station
	2,  // 3: centralcommandsvc.v1.RegisterShipRequest.ship:type_name -> centralcommandsvc.v1.Ship
	2,  // 4: centralcommandsvc.v1.RegisterShipResponse.ship:type_name -> centralcommandsvc.v1.Ship
	2,  // 5: centralcommandsvc.v1.GetAllShipsResponse.ships:type_name -> centralcommandsvc.v1.Ship
	1,  // 6: centralcommandsvc.v1.GetAllStationsResponse.stations:type_name -> centralcommandsvc.v1.Station
	3,  // 7: centralcommandsvc.v1.CentralCommandService.RegisterStation:input_type -> centralcommandsvc.v1.RegisterStationRequest
	5,  // 8: centralcommandsvc.v1.CentralCommandService.RegisterShip:input_type -> centralcommandsvc.v1.RegisterShipRequest
	7,  // 9: centralcommandsvc.v1.CentralCommandService.GetAllShips:input_type -> centralcommandsvc.v1.GetAllShipsRequest
	9,  // 10: centralcommandsvc.v1.CentralCommandService.GetAllStations:input_type -> centralcommandsvc.v1.GetAllStationsRequest
	4,  // 11: centralcommandsvc.v1.CentralCommandService.RegisterStation:output_type -> centralcommandsvc.v1.RegisterStationResponse
	6,  // 12: centralcommandsvc.v1.CentralCommandService.RegisterShip:output_type -> centralcommandsvc.v1.RegisterShipResponse
	8,  // 13: centralcommandsvc.v1.CentralCommandService.GetAllShips:output_type -> centralcommandsvc.v1.GetAllShipsResponse
	10, // 14: centralcommandsvc.v1.CentralCommandService.GetAllStations:output_type -> centralcommandsvc.v1.GetAllStationsResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_centralcommandsvc_v1_centralcommandsvc_proto_init() }
func file_centralcommandsvc_v1_centralcommandsvc_proto_init() {
	if File_centralcommandsvc_v1_centralcommandsvc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dock); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Station); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ship); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterStationRequest); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterStationResponse); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterShipRequest); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterShipResponse); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllShipsRequest); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllShipsResponse); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllStationsRequest); i {
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
		file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllStationsResponse); i {
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
			RawDescriptor: file_centralcommandsvc_v1_centralcommandsvc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_centralcommandsvc_v1_centralcommandsvc_proto_goTypes,
		DependencyIndexes: file_centralcommandsvc_v1_centralcommandsvc_proto_depIdxs,
		MessageInfos:      file_centralcommandsvc_v1_centralcommandsvc_proto_msgTypes,
	}.Build()
	File_centralcommandsvc_v1_centralcommandsvc_proto = out.File
	file_centralcommandsvc_v1_centralcommandsvc_proto_rawDesc = nil
	file_centralcommandsvc_v1_centralcommandsvc_proto_goTypes = nil
	file_centralcommandsvc_v1_centralcommandsvc_proto_depIdxs = nil
}
