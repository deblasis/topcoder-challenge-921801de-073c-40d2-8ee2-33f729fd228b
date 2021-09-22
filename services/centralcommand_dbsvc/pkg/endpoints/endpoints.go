package endpoints

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	CreateShipEndpoint  endpoint.Endpoint
	GetAllShipsEndpoint endpoint.Endpoint

	CreateStationEndpoint  endpoint.Endpoint
	GetAllStationsEndpoint endpoint.Endpoint

	GetNextAvailableDockingStationEndpoint endpoint.Endpoint
	LandShipToDockEndpoint                 endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandDBService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = MakeCreateShipEndpoint(s)
		createShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createShipEndpoint)
		createShipEndpoint = opentracing.TraceServer(otTracer, "CreateShip")(createShipEndpoint)
		if zipkinTracer != nil {
			createShipEndpoint = zipkin.TraceEndpoint(zipkinTracer, "CreateShip")(createShipEndpoint)
		}
		createShipEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateShip"))(createShipEndpoint)
		createShipEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "CreateShip"))(createShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = MakeGetAllShipsEndpoint(s)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllShipsEndpoint)
		getAllShipsEndpoint = opentracing.TraceServer(otTracer, "GetAllShips")(getAllShipsEndpoint)
		if zipkinTracer != nil {
			getAllShipsEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetAllShips")(getAllShipsEndpoint)
		}
		getAllShipsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetAllShips"))(getAllShipsEndpoint)
		getAllShipsEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetAllShips"))(getAllShipsEndpoint)
	}

	var createStationEndpoint endpoint.Endpoint
	{
		createStationEndpoint = MakeCreateStationEndpoint(s)
		createStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createStationEndpoint)
		createStationEndpoint = opentracing.TraceServer(otTracer, "CreateStation")(createStationEndpoint)
		if zipkinTracer != nil {
			createStationEndpoint = zipkin.TraceEndpoint(zipkinTracer, "CreateStation")(createStationEndpoint)
		}
		createStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateStation"))(createStationEndpoint)
		createStationEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "CreateStation"))(createStationEndpoint)
	}

	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = MakeGetAllStationsEndpoint(s)
		getAllStationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllStationsEndpoint)
		getAllStationsEndpoint = opentracing.TraceServer(otTracer, "GetAllStations")(getAllStationsEndpoint)
		if zipkinTracer != nil {
			getAllStationsEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetAllStations")(getAllStationsEndpoint)
		}
		getAllStationsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetAllStations"))(getAllStationsEndpoint)
		getAllStationsEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetAllStations"))(getAllStationsEndpoint)
	}

	var getNextAvailableDockingStationEndpoint endpoint.Endpoint
	{
		getNextAvailableDockingStationEndpoint = MakeGetNextAvailableDockingStationEndpoint(s)
		getNextAvailableDockingStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = opentracing.TraceServer(otTracer, "GetNextAvailableDockingStation")(getNextAvailableDockingStationEndpoint)
		if zipkinTracer != nil {
			getNextAvailableDockingStationEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetNextAvailableDockingStation")(getNextAvailableDockingStationEndpoint)
		}
		getNextAvailableDockingStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetNextAvailableDockingStation"))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetNextAvailableDockingStation"))(getNextAvailableDockingStationEndpoint)
	}

	var landShipToDockEndpoint endpoint.Endpoint
	{
		landShipToDockEndpoint = MakeLandShipToDockEndpoint(s)
		landShipToDockEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landShipToDockEndpoint)
		landShipToDockEndpoint = opentracing.TraceServer(otTracer, "LandShipToDock")(landShipToDockEndpoint)
		if zipkinTracer != nil {
			landShipToDockEndpoint = zipkin.TraceEndpoint(zipkinTracer, "LandShipToDock")(landShipToDockEndpoint)
		}
		landShipToDockEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "LandShipToDock"))(landShipToDockEndpoint)
		landShipToDockEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "LandShipToDock"))(landShipToDockEndpoint)
	}
	return EndpointSet{
		StatusEndpoint:                         healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),
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
