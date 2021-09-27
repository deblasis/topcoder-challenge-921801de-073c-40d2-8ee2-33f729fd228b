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
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.AuthService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []grpctransport.ClientOption

	if zipkinTracer != nil {
		// Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing client can be instantiated
		// without an operation name and fed to each Go kit client as ClientOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path.
		//
		// In this example, we demonstrace a global tracing client.
		options = append(options, zipkin.GRPCClientTrace(zipkinTracer))

	}

	var signupEndpoint endpoint.Endpoint
	{
		signupEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"Signup",
			encodeGRPCSignupRequest,
			decodeGRPCSignupReply,
			pb.SignupResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		signupEndpoint = opentracing.TraceClient(otTracer, "Signup")(signupEndpoint)
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
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		loginEndpoint = opentracing.TraceClient(otTracer, "Login")(loginEndpoint)
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
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		checkTokenEndpoint = opentracing.TraceClient(otTracer, "Login")(checkTokenEndpoint)
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
