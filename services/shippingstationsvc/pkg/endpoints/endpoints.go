// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/sony/gobreaker"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	RequestLandingEndpoint endpoint.Endpoint
	LandingEndpoint        endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.ShippingStationService, logger log.Logger) EndpointSet {

	var requestLandingEndpoint endpoint.Endpoint
	{
		requestLandingEndpoint = MakeRequestLandingEndpoint(s)
		requestLandingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(requestLandingEndpoint)
		requestLandingEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RequestLanding"))(requestLandingEndpoint)
	}

	var landingEndpoint endpoint.Endpoint
	{
		landingEndpoint = MakeLandingEndpoint(s)
		landingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landingEndpoint)
		landingEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Landing"))(landingEndpoint)
	}

	return EndpointSet{
		StatusEndpoint: healthcheck.MakeStatusEndpoint(logger),

		RequestLandingEndpoint: requestLandingEndpoint,
		LandingEndpoint:        landingEndpoint,

		logger: logger,
	}
}

func MakeRequestLandingEndpoint(s service.ShippingStationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.RequestLandingRequest)
		resp, err := s.RequestLanding(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeLandingEndpoint(s service.ShippingStationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.LandingRequest)
		resp, err := s.Landing(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
