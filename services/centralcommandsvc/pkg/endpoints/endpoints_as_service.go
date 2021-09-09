package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
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

// RegisterShip(ctx context.Context, request dtos.RegisterShipRequest) (dtos.RegisterShipResponse, error)
func (s EndpointSet) RegisterShip(ctx context.Context, request dtos.RegisterShipRequest) (*dtos.RegisterShipResponse, error) {
	resp, err := s.RegisterShipEndpoint(ctx, dtos.RegisterShipRequest{
		Weight: request.Weight,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.RegisterShipResponse)
	return response, errors.Str2err(response.Err)
}

// GetAllShips(ctx context.Context, request dtos.GetAllShipsRequest) (dtos.GetAllShipsResponse, error)
func (s EndpointSet) GetAllShips(ctx context.Context, request dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error) {
	resp, err := s.GetAllShipsEndpoint(ctx, dtos.GetAllShipsRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetAllShipsResponse)
	return response, errors.Str2err(response.Err)
}

// RegisterStation(ctx context.Context, request dtos.RegisterStationRequest) (dtos.RegisterStationResponse, error)
func (s EndpointSet) RegisterStation(ctx context.Context, request dtos.RegisterStationRequest) (*dtos.RegisterStationResponse, error) {
	resp, err := s.RegisterStationEndpoint(ctx, dtos.RegisterStationRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.RegisterStationResponse)
	return response, errors.Str2err(response.Err)
}

// GetAllStations(ctx context.Context, request dtos.GetAllStationsRequest) (dtos.GetAllStationsResponse, error)
func (s EndpointSet) GetAllStations(ctx context.Context, request dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error) {
	resp, err := s.GetAllStationsEndpoint(ctx, dtos.GetAllStationsRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetAllStationsResponse)
	return response, errors.Str2err(response.Err)
}

var (
	_ endpoint.Failer = dtos.RegisterShipResponse{}
	_ endpoint.Failer = dtos.GetAllShipsResponse{}
	_ endpoint.Failer = dtos.RegisterStationResponse{}
	_ endpoint.Failer = dtos.GetAllStationsResponse{}
)
