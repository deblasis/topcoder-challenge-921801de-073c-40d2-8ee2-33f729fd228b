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
