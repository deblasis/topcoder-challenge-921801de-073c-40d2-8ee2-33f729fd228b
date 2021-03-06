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

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// CreateShip(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) CreateShip(ctx context.Context, request *dtos.CreateShipRequest) (*dtos.CreateShipResponse, error) {

	resp, err := s.CreateShipEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.CreateShipResponse)
	return response, nil
}

// GetAllShips(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) GetAllShips(ctx context.Context, request *dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error) {
	resp, err := s.GetAllShipsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetAllShipsResponse)
	return response, nil
}

// CreateStation(ctx context.Context, station *model.Station) (string, error)
func (s EndpointSet) CreateStation(ctx context.Context, request *dtos.CreateStationRequest) (*dtos.CreateStationResponse, error) {
	resp, err := s.CreateStationEndpoint(ctx, request)

	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.CreateStationResponse)
	return response, nil
}

// GetAllStations(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) GetAllStations(ctx context.Context, request *dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error) {
	resp, err := s.GetAllStationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetAllStationsResponse)
	return response, nil
}

// GetNextAvailableDockingStation(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) GetNextAvailableDockingStation(ctx context.Context, request *dtos.GetNextAvailableDockingStationRequest) (*dtos.GetNextAvailableDockingStationResponse, error) {
	resp, err := s.GetNextAvailableDockingStationEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetNextAvailableDockingStationResponse)
	return response, nil
}

// LandShipToDock(ctx context.Context, ship *model.Ship) (int64, error)
func (s EndpointSet) LandShipToDock(ctx context.Context, request *dtos.LandShipToDockRequest) (*dtos.LandShipToDockResponse, error) {
	resp, err := s.LandShipToDockEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.LandShipToDockResponse)
	return response, nil
}

var (
	_ endpoint.Failer = dtos.CreateStationResponse{}
	_ endpoint.Failer = dtos.CreateShipResponse{}
	_ endpoint.Failer = dtos.GetAllShipsResponse{}
	_ endpoint.Failer = dtos.GetAllStationsResponse{}
	_ endpoint.Failer = dtos.GetNextAvailableDockingStationResponse{}
	_ endpoint.Failer = dtos.LandShipToDockResponse{}
)
