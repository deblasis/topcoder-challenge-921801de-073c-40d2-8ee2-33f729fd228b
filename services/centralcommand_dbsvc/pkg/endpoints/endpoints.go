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

	"github.com/sony/gobreaker"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type EndpointSet struct {
	CreateShipEndpoint  endpoint.Endpoint
	GetAllShipsEndpoint endpoint.Endpoint

	CreateStationEndpoint  endpoint.Endpoint
	GetAllStationsEndpoint endpoint.Endpoint

	GetNextAvailableDockingStationEndpoint endpoint.Endpoint
	LandShipToDockEndpoint                 endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandDBService, logger log.Logger) EndpointSet {

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = MakeCreateShipEndpoint(s)
		createShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createShipEndpoint)
		createShipEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateShip"))(createShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = MakeGetAllShipsEndpoint(s)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllShipsEndpoint)
		getAllShipsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetAllShips"))(getAllShipsEndpoint)
	}

	var createStationEndpoint endpoint.Endpoint
	{
		createStationEndpoint = MakeCreateStationEndpoint(s)
		createStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createStationEndpoint)
		createStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateStation"))(createStationEndpoint)
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
		landShipToDockEndpoint = MakeLandShipToDockEndpoint(s)
		landShipToDockEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landShipToDockEndpoint)
		landShipToDockEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "LandShipToDock"))(landShipToDockEndpoint)
	}
	return EndpointSet{
		CreateShipEndpoint:                     createShipEndpoint,
		GetAllShipsEndpoint:                    getAllShipsEndpoint,
		CreateStationEndpoint:                  createStationEndpoint,
		GetAllStationsEndpoint:                 getAllStationsEndpoint,
		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
		LandShipToDockEndpoint:                 landShipToDockEndpoint,
		logger:                                 logger,
	}
}

func MakeCreateShipEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateShipRequest)
		resp, err := s.CreateShip(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetAllShipsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetAllShipsRequest)
		resp, err := s.GetAllShips(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeCreateStationEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateStationRequest)
		resp, err := s.CreateStation(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetAllStationsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetAllStationsRequest)
		resp, err := s.GetAllStations(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetNextAvailableDockingStationEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetNextAvailableDockingStationRequest)
		resp, err := s.GetNextAvailableDockingStation(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeLandShipToDockEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.LandShipToDockRequest)
		resp, err := s.LandShipToDock(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
