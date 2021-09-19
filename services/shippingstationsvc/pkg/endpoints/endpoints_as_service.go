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
