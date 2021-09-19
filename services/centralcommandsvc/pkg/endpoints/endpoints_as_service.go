package endpoints

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// RegisterShip(ctx context.Context, request pb.RegisterShipRequest) (pb.RegisterShipResponse, error)
func (s EndpointSet) RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error) {
	resp, err := s.RegisterShipEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RegisterShipResponse)
	return response, nil
}

// GetAllShips(ctx context.Context, request pb.GetAllShipsRequest) (pb.GetAllShipsResponse, error)
func (s EndpointSet) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	resp, err := s.GetAllShipsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.GetAllShipsResponse)
	return response, nil
}

// RegisterStation(ctx context.Context, request pb.RegisterStationRequest) (pb.RegisterStationResponse, error)
func (s EndpointSet) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
	resp, err := s.RegisterStationEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RegisterStationResponse)
	return response, nil
}

// GetAllStations(ctx context.Context, request pb.GetAllStationsRequest) (pb.GetAllStationsResponse, error)
func (s EndpointSet) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	resp, err := s.GetAllStationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.GetAllStationsResponse)
	return response, nil
}

// GetNextAvailableDockingStation(ctx context.Context, request pb.GetNextAvailableDockingStationRequest) (pb.GetNextAvailableDockingStationResponse, error)
func (s EndpointSet) GetNextAvailableDockingStation(ctx context.Context, request *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	resp, err := s.GetNextAvailableDockingStationEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.GetNextAvailableDockingStationResponse)
	return response, nil
}

// LandShipToDock(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) LandShipToDock(ctx context.Context, request *pb.LandShipToDockRequest) (*pb.LandShipToDockResponse, error) {
	resp, err := s.LandShipToDockEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.LandShipToDockResponse)
	return response, nil
}

var (
	_ endpoint.Failer = pb.RegisterShipResponse{}
	_ endpoint.Failer = pb.GetAllShipsResponse{}
	_ endpoint.Failer = pb.RegisterStationResponse{}
	_ endpoint.Failer = pb.GetAllStationsResponse{}
	_ endpoint.Failer = pb.GetNextAvailableDockingStationResponse{}
	_ endpoint.Failer = dtos.LandShipToDockResponse{}
)
