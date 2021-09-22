package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	RegisterShipEndpoint endpoint.Endpoint
	GetAllShipsEndpoint  endpoint.Endpoint

	RegisterStationEndpoint endpoint.Endpoint
	GetAllStationsEndpoint  endpoint.Endpoint

	GetNextAvailableDockingStationEndpoint endpoint.Endpoint
	RegisterShipLandingEndpoint            endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var registerShipEndpoint endpoint.Endpoint
	{
		registerShipEndpoint = MakeRegisterShipEndpoint(s)
		registerShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerShipEndpoint)
		registerShipEndpoint = opentracing.TraceServer(otTracer, "RegisterShip")(registerShipEndpoint)
		if zipkinTracer != nil {
			registerShipEndpoint = zipkin.TraceEndpoint(zipkinTracer, "RegisterShip")(registerShipEndpoint)
		}
		registerShipEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterShip"))(registerShipEndpoint)
		registerShipEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "RegisterShip"))(registerShipEndpoint)
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

	var registerStationEndpoint endpoint.Endpoint
	{
		registerStationEndpoint = MakeRegisterStationEndpoint(s)
		registerStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerStationEndpoint)
		registerStationEndpoint = opentracing.TraceServer(otTracer, "RegisterStation")(registerStationEndpoint)
		if zipkinTracer != nil {
			registerStationEndpoint = zipkin.TraceEndpoint(zipkinTracer, "RegisterStation")(registerStationEndpoint)
		}
		registerStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterStation"))(registerStationEndpoint)
		registerStationEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "RegisterStation"))(registerStationEndpoint)
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
		landShipToDockEndpoint = MakeRegisterShipLandingEndpoint(s)
		landShipToDockEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landShipToDockEndpoint)
		landShipToDockEndpoint = opentracing.TraceServer(otTracer, "RegisterShipLanding")(landShipToDockEndpoint)
		if zipkinTracer != nil {
			landShipToDockEndpoint = zipkin.TraceEndpoint(zipkinTracer, "RegisterShipLanding")(landShipToDockEndpoint)
		}
		landShipToDockEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RegisterShipLanding"))(landShipToDockEndpoint)
		landShipToDockEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "RegisterShipLanding"))(landShipToDockEndpoint)
	}

	return EndpointSet{
		StatusEndpoint: healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),

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
