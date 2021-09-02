package endpoints

import (
	"context"
	"reflect"
	"strings"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.UserManager, logger log.Logger) EndpointSet {
	return EndpointSet{
		StatusEndpoint:            makeStatusEndpoint(s),
		GetUserByUsernameEndpoint: makeGetUserByUsernameEndpoint(s),
		CreateUserEndpoint:        makeCreateUserEndpointEndpoint(s),
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
		req := request.(model.User)

		var err error
		err = validate.Struct(req)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return -1, errors.Wrap(validationErrors, "Validation failed")
		}

		id, err := s.CreateUser(ctx, &req)
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
