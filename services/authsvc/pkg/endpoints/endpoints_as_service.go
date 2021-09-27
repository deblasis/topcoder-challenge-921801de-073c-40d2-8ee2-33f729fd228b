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
