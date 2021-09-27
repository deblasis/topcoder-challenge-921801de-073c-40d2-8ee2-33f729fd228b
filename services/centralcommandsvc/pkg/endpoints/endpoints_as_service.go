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
