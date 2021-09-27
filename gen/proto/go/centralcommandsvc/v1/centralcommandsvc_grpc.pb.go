// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
//
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package centralcommandsvc_v1

import (
	context "context"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
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
	RegisterStation(ctx context.Context, in *RegisterStationRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	RegisterShip(ctx context.Context, in *RegisterShipRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	GetAllShips(ctx context.Context, in *GetAllShipsRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	GetAllStations(ctx context.Context, in *GetAllStationsRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error)
	GetNextAvailableDockingStation(ctx context.Context, in *GetNextAvailableDockingStationRequest, opts ...grpc.CallOption) (*GetNextAvailableDockingStationResponse, error)
	RegisterShipLanding(ctx context.Context, in *RegisterShipLandingRequest, opts ...grpc.CallOption) (*RegisterShipLandingResponse, error)
}

type centralCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCentralCommandServiceClient(cc grpc.ClientConnInterface) CentralCommandServiceClient {
	return &centralCommandServiceClient{cc}
}

func (c *centralCommandServiceClient) RegisterStation(ctx context.Context, in *RegisterStationRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/RegisterStation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) RegisterShip(ctx context.Context, in *RegisterShipRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/RegisterShip", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) GetAllShips(ctx context.Context, in *GetAllShipsRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/GetAllShips", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) GetAllStations(ctx context.Context, in *GetAllStationsRequest, opts ...grpc.CallOption) (*httpbody.HttpBody, error) {
	out := new(httpbody.HttpBody)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/GetAllStations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) GetNextAvailableDockingStation(ctx context.Context, in *GetNextAvailableDockingStationRequest, opts ...grpc.CallOption) (*GetNextAvailableDockingStationResponse, error) {
	out := new(GetNextAvailableDockingStationResponse)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/GetNextAvailableDockingStation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *centralCommandServiceClient) RegisterShipLanding(ctx context.Context, in *RegisterShipLandingRequest, opts ...grpc.CallOption) (*RegisterShipLandingResponse, error) {
	out := new(RegisterShipLandingResponse)
	err := c.cc.Invoke(ctx, "/deblasis.v1.CentralCommandService/RegisterShipLanding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CentralCommandServiceServer is the server API for CentralCommandService service.
// All implementations must embed UnimplementedCentralCommandServiceServer
// for forward compatibility
type CentralCommandServiceServer interface {
	RegisterStation(context.Context, *RegisterStationRequest) (*httpbody.HttpBody, error)
	RegisterShip(context.Context, *RegisterShipRequest) (*httpbody.HttpBody, error)
	GetAllShips(context.Context, *GetAllShipsRequest) (*httpbody.HttpBody, error)
	GetAllStations(context.Context, *GetAllStationsRequest) (*httpbody.HttpBody, error)
	GetNextAvailableDockingStation(context.Context, *GetNextAvailableDockingStationRequest) (*GetNextAvailableDockingStationResponse, error)
	RegisterShipLanding(context.Context, *RegisterShipLandingRequest) (*RegisterShipLandingResponse, error)
	mustEmbedUnimplementedCentralCommandServiceServer()
}

// UnimplementedCentralCommandServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCentralCommandServiceServer struct {
}

func (UnimplementedCentralCommandServiceServer) RegisterStation(context.Context, *RegisterStationRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterStation not implemented")
}
func (UnimplementedCentralCommandServiceServer) RegisterShip(context.Context, *RegisterShipRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterShip not implemented")
}
func (UnimplementedCentralCommandServiceServer) GetAllShips(context.Context, *GetAllShipsRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllShips not implemented")
}
func (UnimplementedCentralCommandServiceServer) GetAllStations(context.Context, *GetAllStationsRequest) (*httpbody.HttpBody, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllStations not implemented")
}
func (UnimplementedCentralCommandServiceServer) GetNextAvailableDockingStation(context.Context, *GetNextAvailableDockingStationRequest) (*GetNextAvailableDockingStationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextAvailableDockingStation not implemented")
}
func (UnimplementedCentralCommandServiceServer) RegisterShipLanding(context.Context, *RegisterShipLandingRequest) (*RegisterShipLandingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterShipLanding not implemented")
}
func (UnimplementedCentralCommandServiceServer) mustEmbedUnimplementedCentralCommandServiceServer() {}

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
		FullMethod: "/deblasis.v1.CentralCommandService/RegisterStation",
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
		FullMethod: "/deblasis.v1.CentralCommandService/RegisterShip",
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
		FullMethod: "/deblasis.v1.CentralCommandService/GetAllShips",
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
		FullMethod: "/deblasis.v1.CentralCommandService/GetAllStations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).GetAllStations(ctx, req.(*GetAllStationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralCommandService_GetNextAvailableDockingStation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNextAvailableDockingStationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).GetNextAvailableDockingStation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deblasis.v1.CentralCommandService/GetNextAvailableDockingStation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).GetNextAvailableDockingStation(ctx, req.(*GetNextAvailableDockingStationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CentralCommandService_RegisterShipLanding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterShipLandingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CentralCommandServiceServer).RegisterShipLanding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deblasis.v1.CentralCommandService/RegisterShipLanding",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CentralCommandServiceServer).RegisterShipLanding(ctx, req.(*RegisterShipLandingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CentralCommandService_ServiceDesc is the grpc.ServiceDesc for CentralCommandService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CentralCommandService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "deblasis.v1.CentralCommandService",
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
		{
			MethodName: "GetNextAvailableDockingStation",
			Handler:    _CentralCommandService_GetNextAvailableDockingStation_Handler,
		},
		{
			MethodName: "RegisterShipLanding",
			Handler:    _CentralCommandService_RegisterShipLanding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "centralcommandsvc/v1/centralcommandsvc.proto",
}
