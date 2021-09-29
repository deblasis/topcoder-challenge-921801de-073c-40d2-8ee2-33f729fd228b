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
package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
)

type EndpointSet struct {
	RegisterShipEndpoint endpoint.Endpoint
	GetAllShipsEndpoint  endpoint.Endpoint

	RegisterStationEndpoint endpoint.Endpoint
	GetAllStationsEndpoint  endpoint.Endpoint

	GetNextAvailableDockingStationEndpoint endpoint.Endpoint
	RegisterShipLandingEndpoint            endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandService, logger log.Logger) EndpointSet {

	var registerShipEndpoint endpoint.Endpoint
	{
		registerShipEndpoint = MakeRegisterShipEndpoint(s)
		registerShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerShipEndpoint)
		registerShipEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterShip"))(registerShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = MakeGetAllShipsEndpoint(s)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllShipsEndpoint)
		getAllShipsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetAllShips"))(getAllShipsEndpoint)
	}

	var registerStationEndpoint endpoint.Endpoint
	{
		registerStationEndpoint = MakeRegisterStationEndpoint(s)
		registerStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerStationEndpoint)
		registerStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterStation"))(registerStationEndpoint)
	}

	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = MakeGetAllStationsEndpoint(s)
		getAllStationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllStationsEndpoint)
		getAllStationsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetAllStations"))(getAllStationsEndpoint)
	}

	var getNextAvailableDockingStationEndpoint endpoint.Endpoint
	{
		getNextAvailableDockingStationEndpoint = MakeGetNextAvailableDockingStationEndpoint(s)
		getNextAvailableDockingStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetNextAvailableDockingStation"))(getNextAvailableDockingStationEndpoint)
	}

	var landShipToDockEndpoint endpoint.Endpoint
	{
		landShipToDockEndpoint = MakeRegisterShipLandingEndpoint(s)
		landShipToDockEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landShipToDockEndpoint)
		landShipToDockEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterShipLanding"))(landShipToDockEndpoint)
	}

	return EndpointSet{
		RegisterShipEndpoint: registerShipEndpoint,
		GetAllShipsEndpoint:  getAllShipsEndpoint,

		RegisterStationEndpoint: registerStationEndpoint,
		GetAllStationsEndpoint:  getAllStationsEndpoint,

		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
		RegisterShipLandingEndpoint:            landShipToDockEndpoint,

		logger: logger,
	}
}

func MakeRegisterShipEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.RegisterShipRequest)
		resp, err := s.RegisterShip(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetAllShipsEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetAllShipsRequest)
		resp, err := s.GetAllShips(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeRegisterStationEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.RegisterStationRequest)
		resp, err := s.RegisterStation(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetAllStationsEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetAllStationsRequest)
		resp, err := s.GetAllStations(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetNextAvailableDockingStationEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.GetNextAvailableDockingStationRequest)
		resp, err := s.GetNextAvailableDockingStation(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
func MakeRegisterShipLandingEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.RegisterShipLandingRequest)
		resp, err := s.RegisterShipLanding(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
