// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package centralcommandsvc_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CentralCommandServiceClient is the client API for CentralCommandService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CentralCommandServiceClient interface {
	RegisterStation(ctx context.Context, in *RegisterStationRequest, opts ...grpc.CallOption) (*RegisterStationResponse, error)
	RegisterShip(ctx context.Context, in *RegisterShipRequest, opts ...grpc.CallOption) (*RegisterShipResponse, error)
	GetAllShips(ctx context.Context, in *GetAllShipsRequest, opts ...grpc.CallOption) (*GetAllShipsResponse, error)
	GetAllStations(ctx context.Context, in *GetAllStationsRequest, opts ...grpc.CallOption) (*GetAllStationsResponse, error)
}

type centralCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCentralCommandServiceClient(cc grpc.ClientConnInterface) CentralCommandServiceClient {
	return &centralCommandServiceClient{cc}
}

func (c *centralCommandServiceClient) RegisterStation(ctx context.Context, in *RegisterStationRequest, opts ...grpc.CallOption) (*RegisterStationResponse, error) {
	out := new(RegisterStationResponse)
	err := c.cc.Invoke(ctx, "/centralcommandsvc.v1.CentralCommandService/RegisterStation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) RegisterShip(ctx context.Context, in *RegisterShipRequest, opts ...grpc.CallOption) (*RegisterShipResponse, error) {
	out := new(RegisterShipResponse)
	err := c.cc.Invoke(ctx, "/centralcommandsvc.v1.CentralCommandService/RegisterShip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) GetAllShips(ctx context.Context, in *GetAllShipsRequest, opts ...grpc.CallOption) (*GetAllShipsResponse, error) {
	out := new(GetAllShipsResponse)
	err := c.cc.Invoke(ctx, "/centralcommandsvc.v1.CentralCommandService/GetAllShips", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) GetAllStations(ctx context.Context, in *GetAllStationsRequest, opts ...grpc.CallOption) (*GetAllStationsResponse, error) {
	out := new(GetAllStationsResponse)
	err := c.cc.Invoke(ctx, "/centralcommandsvc.v1.CentralCommandService/GetAllStations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CentralCommandServiceServer is the server API for CentralCommandService service.
// All implementations should embed UnimplementedCentralCommandServiceServer
// for forward compatibility
type CentralCommandServiceServer interface {
	RegisterStation(context.Context, *RegisterStationRequest) (*RegisterStationResponse, error)
	RegisterShip(context.Context, *RegisterShipRequest) (*RegisterShipResponse, error)
	GetAllShips(context.Context, *GetAllShipsRequest) (*GetAllShipsResponse, error)
	GetAllStations(context.Context, *GetAllStationsRequest) (*GetAllStationsResponse, error)
}

// UnimplementedCentralCommandServiceServer should be embedded to have forward compatible implementations.
type UnimplementedCentralCommandServiceServer struct {
}

func (UnimplementedCentralCommandServiceServer) RegisterStation(context.Context, *RegisterStationRequest) (*RegisterStationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterStation not implemented")
}
func (UnimplementedCentralCommandServiceServer) RegisterShip(context.Context, *RegisterShipRequest) (*RegisterShipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterShip not implemented")
}
func (UnimplementedCentralCommandServiceServer) GetAllShips(context.Context, *GetAllShipsRequest) (*GetAllShipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllShips not implemented")
}
func (UnimplementedCentralCommandServiceServer) GetAllStations(context.Context, *GetAllStationsRequest) (*GetAllStationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllStations not implemented")
}

// UnsafeCentralCommandServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CentralCommandServiceServer will
// result in compilation errors.
type UnsafeCentralCommandServiceServer interface {
	mustEmbedUnimplementedCentralCommandServiceServer()
}

func RegisterCentralCommandServiceServer(s grpc.ServiceRegistrar, srv CentralCommandServiceServer) {
	s.RegisterService(&CentralCommandService_ServiceDesc, srv)
}

func _CentralCommandService_RegisterStation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterStationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).RegisterStation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/centralcommandsvc.v1.CentralCommandService/RegisterStation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).RegisterStation(ctx, req.(*RegisterStationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralCommandService_RegisterShip_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterShipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).RegisterShip(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/centralcommandsvc.v1.CentralCommandService/RegisterShip",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).RegisterShip(ctx, req.(*RegisterShipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralCommandService_GetAllShips_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllShipsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).GetAllShips(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/centralcommandsvc.v1.CentralCommandService/GetAllShips",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).GetAllShips(ctx, req.(*GetAllShipsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralCommandService_GetAllStations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllStationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).GetAllStations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/centralcommandsvc.v1.CentralCommandService/GetAllStations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).GetAllStations(ctx, req.(*GetAllStationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CentralCommandService_ServiceDesc is the grpc.ServiceDesc for CentralCommandService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CentralCommandService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "centralcommandsvc.v1.CentralCommandService",
	HandlerType: (*CentralCommandServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterStation",
			Handler:    _CentralCommandService_RegisterStation_Handler,
		},
		{
			MethodName: "RegisterShip",
			Handler:    _CentralCommandService_RegisterShip_Handler,
		},
		{
			MethodName: "GetAllShips",
			Handler:    _CentralCommandService_GetAllShips_Handler,
		},
		{
			MethodName: "GetAllStations",
			Handler:    _CentralCommandService_GetAllStations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "centralcommandsvc/v1/centralcommandsvc.proto",
}
