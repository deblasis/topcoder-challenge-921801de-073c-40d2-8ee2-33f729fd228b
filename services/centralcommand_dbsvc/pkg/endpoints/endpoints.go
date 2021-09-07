package endpoints

import (
	"context"
	"reflect"
	"strings"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	CreateShipEndpoint  endpoint.Endpoint
	GetAllShipsEndpoint endpoint.Endpoint

	CreateStationEndpoint  endpoint.Endpoint
	GetAllStationsEndpoint endpoint.Endpoint

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

	return EndpointSet{
		StatusEndpoint:         healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),
		CreateShipEndpoint:     createShipEndpoint,
		GetAllShipsEndpoint:    getAllShipsEndpoint,
		CreateStationEndpoint:  createStationEndpoint,
		GetAllStationsEndpoint: getAllStationsEndpoint,
		logger:                 logger,
	}
}

func MakeCreateShipEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dtos.CreateShipRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		ret, err := s.CreateShip(ctx, model.Ship{
			Weight: req.Weight,
		})

		if err != nil {
			return dtos.CreateShipResponse{
				Ship: ret,
				Err:  err.Error(),
			}, err
		}
		return dtos.CreateShipResponse{
			Ship: ret,
			Err:  "",
		}, nil
	}
}

func MakeGetAllShipsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dtos.GetAllShipsRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		ret, err := s.GetAllShips(ctx)

		if err != nil {
			return dtos.GetAllShipsResponse{
				Ships: ret,
				Err:   err.Error(),
			}, err
		}
		return dtos.GetAllShipsResponse{
			Ships: ret,
			Err:   "",
		}, nil
	}
}

func MakeCreateStationEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dtos.CreateStationRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		station := model.Station{
			Capacity: req.Capacity,
			Docks:    []*model.Dock{},
		}
		for _, d := range req.Docks {
			station.Docks = append(station.Docks, &model.Dock{
				NumDockingPorts: d.NumDockingPorts,
			})
		}

		ret, err := s.CreateStation(ctx, station)

		if err != nil {
			return dtos.CreateStationResponse{
				Station: nil,
				Err:     err.Error(),
			}, err
		}
		return dtos.CreateStationResponse{
			Station: ret,
			Err:     "",
		}, nil
	}
}

func MakeGetAllStationsEndpoint(s service.CentralCommandDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dtos.GetAllStationsRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		ret, err := s.GetAllStations(ctx)

		if err != nil {
			return dtos.GetAllStationsResponse{
				Stations: ret,
				Err:      err.Error(),
			}, err
		}
		return dtos.GetAllStationsResponse{
			Stations: ret,
			Err:      "",
		}, nil
	}
}

//TODO see singleton init
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		switch field.Kind() {
		case reflect.String:
			return len(strings.TrimSpace(field.String())) > 0
		case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
			return field.Len() > 0
		case reflect.Ptr, reflect.Interface, reflect.Func:
			return !field.IsNil()
		default:
			return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
		}
	})
}