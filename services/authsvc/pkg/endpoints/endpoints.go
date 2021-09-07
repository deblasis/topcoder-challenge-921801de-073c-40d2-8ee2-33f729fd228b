package endpoints

import (
	"context"
	"reflect"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
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
	SignupEndpoint endpoint.Endpoint
	LoginEndpoint  endpoint.Endpoint
	logger         log.Logger
}

func NewEndpointSet(s service.AuthService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var signupEndpoint endpoint.Endpoint
	{
		signupEndpoint = MakeSignupEndpoint(s)
		signupEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(signupEndpoint)
		signupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(signupEndpoint)
		signupEndpoint = opentracing.TraceServer(otTracer, "Signup")(signupEndpoint)
		if zipkinTracer != nil {
			signupEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Signup")(signupEndpoint)
		}
		signupEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Signup"))(signupEndpoint)
		signupEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Signup"))(signupEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(s)
		loginEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(loginEndpoint)
		loginEndpoint = opentracing.TraceServer(otTracer, "Login")(loginEndpoint)
		if zipkinTracer != nil {
			loginEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Login")(loginEndpoint)
		}
		loginEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
		loginEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Login"))(loginEndpoint)
	}

	return EndpointSet{
		StatusEndpoint: healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),
		SignupEndpoint: signupEndpoint,
		LoginEndpoint:  loginEndpoint,
		logger:         logger,
	}
}

func MakeSignupEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp dtos.SignupResponse
			err  error
		)

		req := request.(dtos.SignupRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.Signup(ctx, req)
		return resp, err
	}
}

func MakeLoginEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			resp dtos.LoginResponse
			err  error
		)

		req := request.(dtos.LoginRequest)

		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return resp, errors.Wrap(validationErrors, "Validation failed")
		}

		resp, err = s.Login(ctx, req)
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
