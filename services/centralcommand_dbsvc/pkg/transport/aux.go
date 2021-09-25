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
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
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
