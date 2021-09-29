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
package transport

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type grpcServer struct {
	pb.UnimplementedCentralCommandDBServiceServer

	createShip  grpctransport.Handler
	getAllShips grpctransport.Handler

	createStation  grpctransport.Handler
	getAllStations grpctransport.Handler

	getNextAvailableDockingStation grpctransport.Handler
	landShipToDock                 grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandDBServiceServer {

	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{
		createShip: grpctransport.NewServer(
			e.CreateShipEndpoint,
			decodeGRPCCreateShipRequest,
			encodeGRPCCreateShipResponse,
			options...,
		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
			options...,
		),
		createStation: grpctransport.NewServer(
			e.CreateStationEndpoint,
			decodeGRPCCreateStationRequest,
			encodeGRPCCreateStationResponse,
			options...,
		),
		getAllStations: grpctransport.NewServer(
			e.GetAllStationsEndpoint,
			decodeGRPCGetAllStationsRequest,
			encodeGRPCGetAllStationsResponse,
			options...,
		),
		getNextAvailableDockingStation: grpctransport.NewServer(
			e.GetNextAvailableDockingStationEndpoint,
			decodeGRPCGetNextAvailableDockingStationRequest,
			encodeGRPCGetNextAvailableDockingStationResponse,
			options...,
		),
		landShipToDock: grpctransport.NewServer(
			e.LandShipToDockEndpoint,
			decodeGRPCLandShipToDockRequest,
			encodeGRPCLandShipToDockResponse,
			options...,
		),
	}
}

func (g *grpcServer) CreateShip(ctx context.Context, r *pb.CreateShipRequest) (*pb.CreateShipResponse, error) {
	_, rep, err := g.createShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateShipResponse), nil
}
func (g *grpcServer) GetAllShips(ctx context.Context, r *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	_, rep, err := g.getAllShips.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllShipsResponse), nil
}

func (g *grpcServer) CreateStation(ctx context.Context, r *pb.CreateStationRequest) (*pb.CreateStationResponse, error) {
	_, rep, err := g.createStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateStationResponse), nil
}
func (g *grpcServer) GetAllStations(ctx context.Context, r *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllStationsResponse), nil
}

func (g *grpcServer) GetNextAvailableDockingStation(ctx context.Context, r *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	_, rep, err := g.getNextAvailableDockingStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetNextAvailableDockingStationResponse), nil
}
func (g *grpcServer) LandShipToDock(ctx context.Context, r *pb.LandShipToDockRequest) (*pb.LandShipToDockResponse, error) {
	_, rep, err := g.landShipToDock.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LandShipToDockResponse), nil
}

func decodeGRPCCreateShipRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateShipRequest)
	return converters.ProtoCreateShipRequestToDto(req), nil
}
func encodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.CreateShipResponse)
	return converters.CreateShipResponseToProto(response), nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	if req != nil {
		return &dtos.GetAllShipsRequest{}, nil
	}
	return nil, errs.NewError(http.StatusInternalServerError, "cannot decode request", errs.ErrException)
}
func encodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetAllShipsResponse)
	return converters.GetAllShipsResponseToProto(response), nil
}

func decodeGRPCCreateStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateStationRequest)
	return converters.ProtoCreateStationRequestToDto(req), nil
}
func encodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.CreateStationResponse)
	return converters.CreateStationResponseToProto(response), nil
}

func decodeGRPCGetAllStationsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllStationsRequest)
	if req != nil {
		return &dtos.GetAllStationsRequest{ShipId: req.ShipId}, nil
	}
	return nil, errs.NewError(http.StatusInternalServerError, "cannot decode request", errs.ErrException)
}
func encodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetAllStationsResponse)
	return converters.GetAllStationsResponseToProto(response), nil
}

func decodeGRPCGetNextAvailableDockingStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetNextAvailableDockingStationRequest)
	return converters.ProtoGetNextAvailableDockingStationRequestToDto(req), nil
}
func encodeGRPCGetNextAvailableDockingStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetNextAvailableDockingStationResponse)
	return converters.GetNextAvailableDockingStationResponseToProto(response), nil
}

func decodeGRPCLandShipToDockRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LandShipToDockRequest)
	return converters.ProtoLandShipToDockRequestToDto(req), nil
}
func encodeGRPCLandShipToDockResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.LandShipToDockResponse)
	return converters.LandShipToDockResponseToProto(response), nil
}
