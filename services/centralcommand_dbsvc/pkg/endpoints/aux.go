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

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
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
