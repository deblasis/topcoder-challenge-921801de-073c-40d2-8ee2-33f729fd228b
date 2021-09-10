package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"github.com/go-kit/kit/endpoint"
)

//ServiceStatus(ctx context.Context) (int64, error)
func (s EndpointSet) ServiceStatus(ctx context.Context) (int64, error) {
	resp, err := s.StatusEndpoint(ctx, healthcheck.ServiceStatusRequest{})
	if err != nil {
		return 0, err
	}
	response := resp.(healthcheck.ServiceStatusResponse)
	return response.Code, errors.Str2err(response.Err)
}

// RegisterShip(ctx context.Context, request pb.RegisterShipRequest) (pb.RegisterShipResponse, error)
func (s EndpointSet) RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error) {
	resp, err := s.RegisterShipEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RegisterShipResponse)
	return response, errors.Str2err(response.Error)
}

// GetAllShips(ctx context.Context, request pb.GetAllShipsRequest) (pb.GetAllShipsResponse, error)
func (s EndpointSet) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	resp, err := s.GetAllShipsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.GetAllShipsResponse)
	return response, errors.Str2err(response.Error)
}

// RegisterStation(ctx context.Context, request pb.RegisterStationRequest) (pb.RegisterStationResponse, error)
func (s EndpointSet) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
	resp, err := s.RegisterStationEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RegisterStationResponse)
	return response, errors.Str2err(response.Error)
}

// GetAllStations(ctx context.Context, request pb.GetAllStationsRequest) (pb.GetAllStationsResponse, error)
func (s EndpointSet) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	resp, err := s.GetAllStationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.GetAllStationsResponse)
	return response, errors.Str2err(response.Error)
}

var (
	_ endpoint.Failer = pb.RegisterShipResponse{}
	_ endpoint.Failer = pb.GetAllShipsResponse{}
	_ endpoint.Failer = pb.RegisterStationResponse{}
	_ endpoint.Failer = pb.GetAllStationsResponse{}
)
