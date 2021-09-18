package endpoints

import (
	"context"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
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

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandDBService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = MakeCreateShipEndpoint(s)
		createShipEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(createShipEndpoint)
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
		getAllShipsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(getAllShipsEndpoint)
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
		createStationEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(createStationEndpoint)
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
		getAllStationsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(getAllStationsEndpoint)
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
		getNextAvailableDockingStationEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = opentracing.TraceServer(otTracer, "GetNextAvailableDockingStation")(getNextAvailableDockingStationEndpoint)
		if zipkinTracer != nil {
			getNextAvailableDockingStationEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetNextAvailableDockingStation")(getNextAvailableDockingStationEndpoint)
		}
		getNextAvailableDockingStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetNextAvailableDockingStation"))(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetNextAvailableDockingStation"))(getNextAvailableDockingStationEndpoint)
	}

	return EndpointSet{
		StatusEndpoint:                         healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),
		CreateShipEndpoint:                     createShipEndpoint,
		GetAllShipsEndpoint:                    getAllShipsEndpoint,
		CreateStationEndpoint:                  createStationEndpoint,
		GetAllStationsEndpoint:                 getAllStationsEndpoint,
		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
		logger:                                 logger,
	}
}

func MakeCreateShipEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateShipRequest)
		return s.CreateShip(ctx, req)
	}
}

func MakeGetAllShipsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetAllShipsRequest)
		return s.GetAllShips(ctx, req)
	}
}

func MakeCreateStationEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateStationRequest)
		return s.CreateStation(ctx, req)
	}
}

func MakeGetAllStationsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetAllStationsRequest)
		return s.GetAllStations(ctx, req)

		// 	if err != nil {
		// 		return dtos.GetAllStationsResponse{
		// 			Err: err.Error(),
		// 		}, err
		// 	}

		// 	stations := make([]dtos.Station, 0)
		// 	for _, x := range ret.Stations {

		// 		station := &dtos.Station{}
		// 		errs := m.Copy(station, x)
		// 		if len(errs) > 0 {
		// 			return nil, errors.Wrap(errs[0], "Failed to map station")
		// 		}
		// 		stations = append(stations, *station)
		// 	}

		// 	return dtos.GetAllStationsResponse{
		// 		Stations: stations,
		// 	}, nil
		// }

	}
}

func MakeGetNextAvailableDockingStationEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetNextAvailableDockingStationRequest)
		return s.GetNextAvailableDockingStation(ctx, req)
	}
}
