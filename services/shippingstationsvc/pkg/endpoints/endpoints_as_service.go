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
