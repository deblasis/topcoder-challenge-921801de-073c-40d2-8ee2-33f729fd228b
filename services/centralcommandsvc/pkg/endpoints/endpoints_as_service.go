//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package endpoints

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
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
func (s EndpointSet) RegisterShipLanding(ctx context.Context, request *pb.RegisterShipLandingRequest) (*pb.RegisterShipLandingResponse, error) {
	resp, err := s.RegisterShipLandingEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pb.RegisterShipLandingResponse)
	return response, nil
}

var (
	_ endpoint.Failer = pb.RegisterShipResponse{}
	_ endpoint.Failer = pb.GetAllShipsResponse{}
	_ endpoint.Failer = pb.RegisterStationResponse{}
	_ endpoint.Failer = pb.GetAllStationsResponse{}
	_ endpoint.Failer = pb.GetNextAvailableDockingStationResponse{}
	_ endpoint.Failer = pb.RegisterShipLandingResponse{}
)
