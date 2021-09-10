package endpoints

import (
	"context"
	"reflect"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-playground/validator"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	RegisterShipEndpoint endpoint.Endpoint
	GetAllShipsEndpoint  endpoint.Endpoint

	RegisterStationEndpoint endpoint.Endpoint
	GetAllStationsEndpoint  endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.CentralCommandService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var registerShipEndpoint endpoint.Endpoint
	{
		registerShipEndpoint = MakeRegisterShipEndpoint(s)
		registerShipEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(registerShipEndpoint)
		registerShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerShipEndpoint)
		registerShipEndpoint = opentracing.TraceServer(otTracer, "Signup")(registerShipEndpoint)
		if zipkinTracer != nil {
			registerShipEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Signup")(registerShipEndpoint)
		}
		registerShipEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Signup"))(registerShipEndpoint)
		registerShipEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Signup"))(registerShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = MakeGetAllShipsEndpoint(s)
		getAllShipsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(getAllShipsEndpoint)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllShipsEndpoint)
		getAllShipsEndpoint = opentracing.TraceServer(otTracer, "Login")(getAllShipsEndpoint)
		if zipkinTracer != nil {
			getAllShipsEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Login")(getAllShipsEndpoint)
		}
		getAllShipsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Login"))(getAllShipsEndpoint)
		getAllShipsEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Login"))(getAllShipsEndpoint)
	}

	var registerStationEndpoint endpoint.Endpoint
	{
		registerStationEndpoint = MakeRegisterStationEndpoint(s)
		registerStationEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(registerStationEndpoint)
		registerStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerStationEndpoint)
		registerStationEndpoint = opentracing.TraceServer(otTracer, "Signup")(registerStationEndpoint)
		if zipkinTracer != nil {
			registerStationEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Signup")(registerStationEndpoint)
		}
		registerStationEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Signup"))(registerStationEndpoint)
		registerStationEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Signup"))(registerStationEndpoint)
	}

	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = MakeGetAllStationsEndpoint(s)
		getAllStationsEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(getAllStationsEndpoint)
		getAllStationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getAllStationsEndpoint)
		getAllStationsEndpoint = opentracing.TraceServer(otTracer, "Login")(getAllStationsEndpoint)
		if zipkinTracer != nil {
			getAllStationsEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Login")(getAllStationsEndpoint)
		}
		getAllStationsEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Login"))(getAllStationsEndpoint)
		getAllStationsEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Login"))(getAllStationsEndpoint)
	}

	return EndpointSet{
		StatusEndpoint: healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),

		RegisterShipEndpoint: registerShipEndpoint,
		GetAllShipsEndpoint:  getAllShipsEndpoint,

		RegisterStationEndpoint: registerStationEndpoint,
		GetAllStationsEndpoint:  getAllStationsEndpoint,

		logger: logger,
	}
}

func MakeRegisterShipEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp *pb.RegisterShipResponse
			err  error
		)

		req := request.(*pb.RegisterShipRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.RegisterShip(ctx, req)
		return resp, err
	}
}

func MakeGetAllShipsEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp *pb.GetAllShipsResponse
			err  error
		)

		req := request.(*pb.GetAllShipsRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.GetAllShips(ctx, req)
		return resp, err
	}
}

func MakeRegisterStationEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp *pb.RegisterStationResponse
			err  error
		)

		req := request.(*pb.RegisterStationRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.RegisterStation(ctx, req)
		return resp, err
	}
}

func MakeGetAllStationsEndpoint(s service.CentralCommandService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp *pb.GetAllStationsResponse
			err  error
		)

		req := request.(*pb.GetAllStationsRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.GetAllStations(ctx, req)
		return resp, err
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
