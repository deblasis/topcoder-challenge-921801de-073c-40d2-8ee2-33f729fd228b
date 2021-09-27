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

	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"github.com/go-kit/kit/endpoint"
)

// RequestLanding(ctx context.Context, request pb.RequestLandingRequest) (pb.RequestLandingResponse, error)
func (s EndpointSet) RequestLanding(ctx context.Context, request *pb.RequestLandingRequest) (*pb.RequestLandingResponse, error) {
	resp, err := s.RequestLandingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RequestLandingResponse)
	return response, nil
}

// Landing(ctx context.Context, request pb.LandingRequest) (pb.LandingResponse, error)
func (s EndpointSet) Landing(ctx context.Context, request *pb.LandingRequest) (*pb.LandingResponse, error) {
	resp, err := s.LandingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.LandingResponse)
	return response, nil
}

var (
	_ endpoint.Failer = pb.RequestLandingResponse{}
	_ endpoint.Failer = pb.LandingResponse{}
)
