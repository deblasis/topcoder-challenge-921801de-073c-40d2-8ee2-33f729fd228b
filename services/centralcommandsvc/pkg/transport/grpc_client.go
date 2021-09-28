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

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) service.CentralCommandService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var registerShipEndpoint endpoint.Endpoint
	{
		registerShipEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"RegisterShip",
			encodeGRPCRegisterShipRequest,
			decodeGRPCRegisterShipResponse,
			pb.RegisterShipResponse{},
		).Endpoint()

		registerShipEndpoint = limiter(registerShipEndpoint)
		registerShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RegisterShip",
			Timeout: 30 * time.Second,
		}))(registerShipEndpoint)
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

	var registerStationEndpoint endpoint.Endpoint
	{
		registerStationEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"RegisterStation",
			encodeGRPCRegisterStationRequest,
			decodeGRPCRegisterStationResponse,
			pb.RegisterStationResponse{},
		).Endpoint()

		registerStationEndpoint = limiter(registerStationEndpoint)
		registerStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RegisterStation",
			Timeout: 30 * time.Second,
		}))(registerStationEndpoint)
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

	var registerShipLandingEndpoint endpoint.Endpoint
	{
		registerShipLandingEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"RegisterShipLanding",
			encodeGRPCRegisterShipLandingRequest,
			decodeGRPCRegisterShipLandingResponse,
			pb.RegisterShipLandingResponse{},
		).Endpoint()

		registerShipLandingEndpoint = limiter(registerShipLandingEndpoint)
		registerShipLandingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RegisterShipLanding",
			Timeout: 30 * time.Second,
		}))(registerShipLandingEndpoint)
	}

	return endpoints.EndpointSet{
		RegisterShipEndpoint: registerShipEndpoint,
		GetAllShipsEndpoint:  getAllShipsEndpoint,

		RegisterStationEndpoint: registerStationEndpoint,
		GetAllStationsEndpoint:  getAllStationsEndpoint,

		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
		RegisterShipLandingEndpoint:            registerShipLandingEndpoint,
	}

}

func encodeGRPCRegisterShipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RegisterShipRequest)
	//return converters.RegisterShipRequestToProto(req), nil
	return req, nil
}

func decodeGRPCRegisterShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterShipResponse)
	//return converters.ProtoRegisterShipResponseToDto(*response), nil
	return response, nil
}

func encodeGRPCGetAllShipsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetAllShipsRequest)
	//return converters.GetAllShipsRequestToProto(req), nil
	return req, nil
}

func decodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	//return converters.ProtoGetAllShipsResponseToDto(*response), nil
	return response, nil
}

//

func encodeGRPCRegisterStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RegisterStationRequest)
	//return converters.RegisterStationRequestToProto(req), nil
	return req, nil
}

func decodeGRPCRegisterStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterStationResponse)
	//return converters.ProtoRegisterStationResponseToDto(*response), nil
	return response, nil
}

func encodeGRPCGetAllStationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetAllStationsRequest)
	//return converters.GetAllStationsRequestToProto(req), nil
	return req, nil
}

func decodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	//return converters.ProtoGetAllStationsResponseToDto(*response), nil
	return response, nil
}

func encodeGRPCGetNextAvailableDockingStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetNextAvailableDockingStationRequest)
	//return converters.GetNextAvailableDockingStationRequestToProto(req), nil
	return req, nil
}

func decodeGRPCGetNextAvailableDockingStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetNextAvailableDockingStationResponse)
	//return converters.ProtoGetNextAvailableDockingStationResponseToDto(*response), nil
	return response, nil
}
func encodeGRPCRegisterShipLandingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RegisterShipLandingRequest)
	return req, nil
}
func decodeGRPCRegisterShipLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterShipLandingResponse)
	return response, nil
}
