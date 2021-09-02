package endpoints

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.UserManager, logger log.Logger) EndpointSet {

	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = makeStatusEndpoint(s)
		statusEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(statusEndpoint)
		statusEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(statusEndpoint)
		//statusEndpoint = opentracing.TraceServer(otTracer, "ServiceStatus")(statusEndpoint)
		// if zipkinTracer != nil {
		// 	statusEndpoint = zipkin.TraceEndpoint(zipkinTracer, "ServiceStatus")(statusEndpoint)
		// }
		statusEndpoint = LoggingMiddleware(log.With(logger, "method", "ServiceStatus"))(statusEndpoint)
		//statusEndpoint = InstrumentingMiddleware(duration.With("method", "ServiceStatus"))(statusEndpoint)

	}

	return EndpointSet{
		StatusEndpoint:            statusEndpoint,
		GetUserByUsernameEndpoint: LoggingMiddleware(log.With(logger, "method", "GetUserByUsername"))(makeGetUserByUsernameEndpoint(s)),
		CreateUserEndpoint:        LoggingMiddleware(log.With(logger, "method", "CreateUser"))(makeCreateUserEndpointEndpoint(s)),
		logger:                    logger,
	}
}

func makeStatusEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(model.ServiceStatusRequest)
		code, err := s.ServiceStatus(ctx)
		if err != nil {
			return model.ServiceStatusReply{Code: code, Err: err.Error()}, nil
		}
		return model.ServiceStatusReply{Code: code, Err: ""}, nil
	}
}

func makeGetUserByUsernameEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserByUsernameRequest)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		user, err := s.GetUserByUsername(ctx, req.Username)
		if err != nil {
			return model.GetUserByUsernameReply{
				User: user,
				Err:  err.Error(),
			}, nil
		}
		return model.GetUserByUsernameReply{
			User: user,
			Err:  "",
		}, nil
	}
}

func makeCreateUserEndpointEndpoint(s service.UserManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)

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
			return model.CreateUserReply{
				Id:  -1,
				Err: err.Error(),
			}, err
		}
		return model.CreateUserReply{
			Id:  id,
			Err: "",
		}, nil
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

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
