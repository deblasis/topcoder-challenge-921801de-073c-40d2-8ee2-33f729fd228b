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
	"strings"
	"time"

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) service.CentralCommandDBService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateShip",
			encodeGRPCCreateShipRequest,
			decodeGRPCCreateShipResponse,
			pb.CreateShipResponse{},
		).Endpoint()

		createShipEndpoint = limiter(createShipEndpoint)
		createShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateShip",
			Timeout: 30 * time.Second,
		}))(createShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetAllShips",
			encodeGRPCGetAllShipsRequest,
			decodeGRPCGetAllShipsResponse,
			pb.GetAllShipsResponse{},
		).Endpoint()

		getAllShipsEndpoint = limiter(getAllShipsEndpoint)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetAllShips",
			Timeout: 30 * time.Second,
		}))(getAllShipsEndpoint)
	}

	var createStationEndpoint endpoint.Endpoint
	{
		createStationEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateStation",
			encodeGRPCCreateStationRequest,
			decodeGRPCCreateStationResponse,
			pb.CreateStationResponse{},
		).Endpoint()

		createStationEndpoint = limiter(createStationEndpoint)
		createStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateStation",
			Timeout: 30 * time.Second,
		}))(createStationEndpoint)
	}

	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetAllStations",
			encodeGRPCGetAllStationsRequest,
			decodeGRPCGetAllStationsResponse,
			pb.GetAllStationsResponse{},
		).Endpoint()

		getAllStationsEndpoint = limiter(getAllStationsEndpoint)
		getAllStationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetAllStations",
			Timeout: 30 * time.Second,
		}))(getAllStationsEndpoint)
	}

	var getNextAvailableDockingStationEndpoint endpoint.Endpoint
	{
		getNextAvailableDockingStationEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetNextAvailableDockingStation",
			encodeGRPCGetNextAvailableDockingStationRequest,
			decodeGRPCGetNextAvailableDockingStationResponse,
			pb.GetNextAvailableDockingStationResponse{},
		).Endpoint()

		getNextAvailableDockingStationEndpoint = limiter(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetNextAvailableDockingStation",
			Timeout: 30 * time.Second,
		}))(getNextAvailableDockingStationEndpoint)
	}

	var landShipToDockEndpoint endpoint.Endpoint
	{
		landShipToDockEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"LandShipToDock",
			encodeGRPCLandShipToDockRequest,
			decodeGRPCLandShipToDockResponse,
			pb.LandShipToDockResponse{},
		).Endpoint()

		landShipToDockEndpoint = limiter(landShipToDockEndpoint)
		landShipToDockEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "LandShipToDock",
			Timeout: 30 * time.Second,
		}))(landShipToDockEndpoint)
	}

	return endpoints.EndpointSet{
		CreateShipEndpoint:  createShipEndpoint,
		GetAllShipsEndpoint: getAllShipsEndpoint,

		CreateStationEndpoint:  createStationEndpoint,
		GetAllStationsEndpoint: getAllStationsEndpoint,

		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
		LandShipToDockEndpoint:                 landShipToDockEndpoint,
	}
}

func encodeGRPCCreateShipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateShipRequest)
	return converters.CreateShipRequestToProto(req), nil
}
func decodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateShipResponse)
	return converters.ProtoCreateShipResponseToDto(response), nil
}

func encodeGRPCGetAllShipsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetAllShipsRequest)
	return converters.GetAllShipsRequestToProto(req), nil
}
func decodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	return converters.ProtoGetAllShipsResponseToDto(response), nil
}

func encodeGRPCCreateStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateStationRequest)
	return converters.CreateStationRequestToProto(req), nil
}

func decodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateStationResponse)
	return converters.ProtoCreateStationResponseToDto(response), nil
}

func encodeGRPCGetAllStationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetAllStationsRequest)
	return converters.GetAllStationsRequestToProto(req), nil
}
func decodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	return converters.ProtoGetAllStationsResponseToDto(response), nil
}

func encodeGRPCGetNextAvailableDockingStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetNextAvailableDockingStationRequest)
	return converters.GetNextAvailableDockingStationRequestToProto(req), nil
}
func decodeGRPCGetNextAvailableDockingStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetNextAvailableDockingStationResponse)
	return converters.ProtoGetNextAvailableDockingStationResponseToDto(response), nil
}

func encodeGRPCLandShipToDockRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.LandShipToDockRequest)
	return converters.LandShipToDockRequestToProto(req), nil
}
func decodeGRPCLandShipToDockResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.LandShipToDockResponse)
	return converters.ProtoLandShipToDockResponseToDto(response), nil
}
