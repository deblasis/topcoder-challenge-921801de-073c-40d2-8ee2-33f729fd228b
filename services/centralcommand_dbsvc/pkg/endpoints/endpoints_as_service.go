package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// ServiceStatus(ctx context.Context) (int64, error)
func (s EndpointSet) ServiceStatus(ctx context.Context) (int64, error) {
	resp, err := s.StatusEndpoint(ctx, healthcheck.ServiceStatusRequest{})
	if err != nil {
		return 0, err
	}
	response := resp.(healthcheck.ServiceStatusResponse)
	return response.Code, errors.Str2err(response.Err)
}

// CreateShip(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) CreateShip(ctx context.Context, ship model.Ship) (*model.Ship, error) {
	resp, err := s.CreateShipEndpoint(ctx, dtos.CreateShipRequest{
		Weight: ship.Weight,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(dtos.CreateShipResponse)
	return response.Ship, errors.Str2err(response.Err)
}

// GetAllShips(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) GetAllShips(ctx context.Context) ([]*model.Ship, error) {
	resp, err := s.GetAllShipsEndpoint(ctx, dtos.GetAllShipsRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(dtos.GetAllShipsResponse)
	return response.Ships, errors.Str2err(response.Err)
}

// CreateStation(ctx context.Context, station *model.Station) (string, error)
func (s EndpointSet) CreateStation(ctx context.Context, station model.Station) (*model.Station, error) {
	resp, err := s.CreateStationEndpoint(ctx, dtos.CreateStationRequest{
		Capacity: station.Capacity,
		Docks:    station.Docks,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(dtos.CreateStationResponse)
	return response.Station, errors.Str2err(response.Err)
}

// GetAllStations(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) GetAllStations(ctx context.Context) ([]*model.Station, error) {
	resp, err := s.GetAllStationsEndpoint(ctx, dtos.GetAllStationsRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(dtos.GetAllStationsResponse)
	return response.Stations, errors.Str2err(response.Err)
}

var (
	_ endpoint.Failer = dtos.CreateStationResponse{}
	_ endpoint.Failer = dtos.CreateShipResponse{}
)
