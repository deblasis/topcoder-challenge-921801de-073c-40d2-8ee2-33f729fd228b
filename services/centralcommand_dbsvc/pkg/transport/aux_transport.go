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

	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

type auxGrpcServer struct {
	pb.CentralCommandDBAuxServiceServer
	cleanup grpctransport.Handler
}

func NewAuxGrpcServer(e endpoints.AuxEndpointSet, l log.Logger) pb.CentralCommandDBAuxServiceServer {
	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &auxGrpcServer{
		cleanup: grpctransport.NewServer(
			e.CleanupEndpoint,
			decodeGRPCCleanupRequest,
			encodeGRPCCleanupResponse,
			options...,
		),
	}
}

func (g *auxGrpcServer) Cleanup(ctx context.Context, r *pb.CleanupRequest) (*pb.CleanupResponse, error) {
	_, rep, err := g.cleanup.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CleanupResponse), nil
}
func encodeGRPCCleanupResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CleanupResponse)
	return response, nil
}

func decodeGRPCCleanupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CleanupRequest)
	return req, nil
}

//Client
func NewAuxGrpcClient(conn *grpc.ClientConn, logger log.Logger) service.CentralCommandDBAuxService {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var options []grpctransport.ClientOption

	var cleanupEndpoint endpoint.Endpoint
	{
		cleanupEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.AuxServiceName, "-", ".", -1),
			"Cleanup",
			encodeGRPCCleanupRequest,
			decodeGRPCCleanupResponse,
			pb.CleanupResponse{},
			options...,
		).Endpoint()

		cleanupEndpoint = limiter(cleanupEndpoint)
		cleanupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Cleanup",
			Timeout: 30 * time.Second,
		}))(cleanupEndpoint)
	}

	return endpoints.AuxEndpointSet{
		CleanupEndpoint: cleanupEndpoint,
	}
}

func encodeGRPCCleanupRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CleanupRequest)
	return req, nil
}
func decodeGRPCCleanupResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CleanupResponse)
	return response, nil
}
