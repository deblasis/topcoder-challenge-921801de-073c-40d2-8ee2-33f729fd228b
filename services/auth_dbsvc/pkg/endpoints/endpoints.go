//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package endpoints

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
)

type EndpointSet struct {
	StatusEndpoint            endpoint.Endpoint
	GetUserByUsernameEndpoint endpoint.Endpoint
	GetUserByIdEndpoint       endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
	logger                    log.Logger
}

func NewEndpointSet(s service.AuthDBService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var getUserByUsernameEndpoint endpoint.Endpoint
	{
		getUserByUsernameEndpoint = MakeGetUserByUsernameEndpoint(s)
		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = opentracing.TraceServer(otTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
		if zipkinTracer != nil {
			getUserByUsernameEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
		}
		getUserByUsernameEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetUserByUsername"))(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetUserByUsername"))(getUserByUsernameEndpoint)
	}

	var getUserByIdEndpoint endpoint.Endpoint
	{
		getUserByIdEndpoint = MakeGetUserByIdEndpoint(s)
		getUserByIdEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserByIdEndpoint)
		getUserByIdEndpoint = opentracing.TraceServer(otTracer, "GetUserById")(getUserByIdEndpoint)
		if zipkinTracer != nil {
			getUserByIdEndpoint = zipkin.TraceEndpoint(zipkinTracer, "GetUserById")(getUserByIdEndpoint)
		}
		getUserByIdEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "GetUserById"))(getUserByIdEndpoint)
		getUserByIdEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "GetUserById"))(getUserByIdEndpoint)
	}

	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = MakeCreateUserEndpoint(s)
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
		GetUserByIdEndpoint:       getUserByIdEndpoint,
		CreateUserEndpoint:        createUserEndpoint,
		logger:                    logger,
	}
}

func MakeGetUserByIdEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetUserByIdRequest)
		resp, err := s.GetUserById(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeGetUserByUsernameEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.GetUserByUsernameRequest)
		resp, err := s.GetUserByUsername(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeCreateUserEndpoint(s service.AuthDBService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dtos.CreateUserRequest)
		resp, err := s.CreateUser(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
