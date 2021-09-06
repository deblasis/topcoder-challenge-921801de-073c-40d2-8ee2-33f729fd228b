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
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.UserManager, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var getUserByUsernameEndpoint endpoint.Endpoint
	{
		getUserByUsernameEndpoint = MakeGetUserByUsernameEndpoint(s, logger)
		getUserByUsernameEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = opentracing.TraceServer(otTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
		if zipkinTracer != nil {
			getUserByUsernameEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
		}
		getUserByUsernameEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetUserByUsername"))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetUserByUsername"))(getUserByUsernameEndpoint)
	}

	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = MakeCreateUserEndpoint(s, logger)
		createUserEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 30))(createUserEndpoint)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createUserEndpoint)
		createUserEndpoint = opentracing.TraceServer(otTracer, "CreateUser")(createUserEndpoint)
		if zipkinTracer != nil {
			createUserEndpoint = zipkin.TraceEndpoint(zipkinTracer, "CreateUser")(createUserEndpoint)
		}
		createUserEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "CreateUser"))(createUserEndpoint)
		createUserEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "CreateUser"))(createUserEndpoint)
	}

	return EndpointSet{
		StatusEndpoint:            healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),
		GetUserByUsernameEndpoint: getUserByUsernameEndpoint,
		CreateUserEndpoint:        createUserEndpoint,
		logger:                    logger,
	}
}

func MakeGetUserByUsernameEndpoint(s service.UserManager, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		level.Info(logger).Log("handling request", "MakeGetUserByUsernameEndpoint")
		defer level.Info(logger).Log("handled request", "MakeGetUserByUsernameEndpoint")

		req := request.(dtos.GetUserByUsernameRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		user, err := s.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return dtos.GetUserByUsernameResponse{
				User: user,
				Err:  err.Error(),
			}, nil
		}
		return dtos.GetUserByUsernameResponse{
			User: user,
			Err:  "",
		}, nil
	}
}

func MakeCreateUserEndpoint(s service.UserManager, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		level.Info(logger).Log("handling request", "MakeCreateUserEndpoint")
		defer level.Info(logger).Log("handled request", "MakeCreateUserEndpoint")

		req := request.(dtos.CreateUserRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		id, err := s.CreateUser(ctx, &model.User{
			ID:       req.ID,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,
		})

		if err != nil {
			return dtos.CreateUserResponse{
				Id:  -1,
				Err: err.Error(),
			}, err
		}
		return dtos.CreateUserResponse{
			Id:  id,
			Err: "",
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
