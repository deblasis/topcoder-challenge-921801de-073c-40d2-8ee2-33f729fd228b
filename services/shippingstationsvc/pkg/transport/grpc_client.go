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

	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) service.ShippingStationService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var requestLandingEndpoint endpoint.Endpoint
	{
		requestLandingEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"RequestLanding",
			encodeGRPCRequestLandingRequest,
			decodeGRPCRequestLandingResponse,
			pb.RequestLandingResponse{},
		).Endpoint()

		requestLandingEndpoint = limiter(requestLandingEndpoint)
		requestLandingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RequestLanding",
			Timeout: 30 * time.Second,
		}))(requestLandingEndpoint)
	}
	var landingEndpoint endpoint.Endpoint
	{
		landingEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"Landing",
			encodeGRPCLandingRequest,
			decodeGRPCLandingResponse,
			pb.LandingResponse{},
		).Endpoint()

		landingEndpoint = limiter(landingEndpoint)
		landingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Landing",
			Timeout: 30 * time.Second,
		}))(landingEndpoint)
	}

	return endpoints.EndpointSet{
		RequestLandingEndpoint: requestLandingEndpoint,
		LandingEndpoint:        landingEndpoint,
	}

}

func encodeGRPCRequestLandingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RequestLandingRequest)
	return req, nil
}

func decodeGRPCRequestLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RequestLandingResponse)
	return response, nil
}

func encodeGRPCLandingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LandingRequest)
	return req, nil
}

func decodeGRPCLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.LandingResponse)
	return response, nil
}
