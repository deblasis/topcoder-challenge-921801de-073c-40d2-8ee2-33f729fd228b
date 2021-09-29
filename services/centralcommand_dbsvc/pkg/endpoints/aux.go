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
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
)

type AuxEndpointSet struct {
	CleanupEndpoint endpoint.Endpoint
	logger          log.Logger
}

func NewAuxEndpointSet(s service.CentralCommandDBAuxService, logger log.Logger) AuxEndpointSet {

	var cleanupEndpoint endpoint.Endpoint
	{
		cleanupEndpoint = MakeCleanupEndpoint(s)
		cleanupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(cleanupEndpoint)
		cleanupEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Cleanup"))(cleanupEndpoint)
	}

	return AuxEndpointSet{
		CleanupEndpoint: cleanupEndpoint,
		logger:          logger,
	}
}

func MakeCleanupEndpoint(s service.CentralCommandDBAuxService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.CleanupRequest)
		resp, err := s.Cleanup(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

// Cleanup(ctx context.Context, request *pb.CleanupRequest) (*pb.CleanupResponse, error)
func (s AuxEndpointSet) Cleanup(ctx context.Context, request *pb.CleanupRequest) (*pb.CleanupResponse, error) {

	resp, err := s.CleanupEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.CleanupResponse)
	return response, nil
}
