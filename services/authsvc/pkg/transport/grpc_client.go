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
package transport

import (
	"context"
	"strings"
	"time"

	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) service.AuthService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var signupEndpoint endpoint.Endpoint
	{
		signupEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"Signup",
			encodeGRPCSignupRequest,
			decodeGRPCSignupReply,
			pb.SignupResponse{},
		).Endpoint()

		signupEndpoint = limiter(signupEndpoint)
		signupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Signup",
			Timeout: 30 * time.Second,
		}))(signupEndpoint)
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCLoginReply,
			pb.LoginResponse{},
		).Endpoint()

		loginEndpoint = limiter(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Login",
			Timeout: 30 * time.Second,
		}))(loginEndpoint)
	}
	var checkTokenEndpoint endpoint.Endpoint
	{
		checkTokenEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CheckToken",
			encodeGRPCCheckTokenRequest,
			decodeGRPCCheckTokenReply,
			pb.CheckTokenResponse{},
		).Endpoint()

		checkTokenEndpoint = limiter(checkTokenEndpoint)
		checkTokenEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Login",
			Timeout: 30 * time.Second,
		}))(checkTokenEndpoint)
	}

	return endpoints.EndpointSet{
		SignupEndpoint:     signupEndpoint,
		LoginEndpoint:      loginEndpoint,
		CheckTokenEndpoint: checkTokenEndpoint,
	}

}

func encodeGRPCSignupRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SignupRequest)
	return req, nil
}

func decodeGRPCSignupReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SignupResponse)
	return reply, nil
}

func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginRequest)
	return req, nil
}

func decodeGRPCLoginReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.LoginResponse)
	return reply, nil
}

func encodeGRPCCheckTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CheckTokenRequest)
	return req, nil
}

func decodeGRPCCheckTokenReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CheckTokenResponse)
	return reply, nil
}
