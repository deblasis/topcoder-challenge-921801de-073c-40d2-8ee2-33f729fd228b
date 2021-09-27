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

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"github.com/go-kit/kit/endpoint"
)

// Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error)
func (s EndpointSet) Signup(ctx context.Context, request *pb.SignupRequest) (*pb.SignupResponse, error) {
	var ret *pb.SignupResponse
	resp, err := s.SignupEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.SignupResponse)
	return response, nil
}

// Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error)
func (s EndpointSet) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	var ret *pb.LoginResponse
	resp, err := s.LoginEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.LoginResponse)
	return response, nil
	//return response, errors.Str2err(response.Error)
}

func (s EndpointSet) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var ret *pb.CheckTokenResponse
	resp, err := s.CheckTokenEndpoint(ctx, request)
	if err != nil {
		return ret, err
	}
	response := resp.(*pb.CheckTokenResponse)
	return response, nil
	//return response, errors.Str2err(response.Error)
}

var (
	_ endpoint.Failer = pb.SignupResponse{}
	_ endpoint.Failer = pb.LoginResponse{}
	_ endpoint.Failer = pb.CheckTokenResponse{}
)
