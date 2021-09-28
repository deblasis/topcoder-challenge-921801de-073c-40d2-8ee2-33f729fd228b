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
package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	pb.CentralCommandServiceServer
	registerShip grpctransport.Handler
	getAllShips  grpctransport.Handler

	registerStation                grpctransport.Handler
	getAllStations                 grpctransport.Handler
	getNextAvailableDockingStation grpctransport.Handler
	registerShipLanding            grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandServiceServer {

	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{

		registerShip: grpctransport.NewServer(
			e.RegisterShipEndpoint,
			decodeGRPCRegisterShipRequest,
			encodeGRPCRegisterShipResponse,
			options...,
		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
			options...,
		),
		registerStation: grpctransport.NewServer(
			e.RegisterStationEndpoint,
			decodeGRPCRegisterStationRequest,
			encodeGRPCRegisterStationResponse,
			options...,
		),
		getAllStations: grpctransport.NewServer(
			e.GetAllStationsEndpoint,
			decodeGRPCGetAllStationsRequest,
			encodeGRPCGetAllStationsResponse,
			options...,
		),
		getNextAvailableDockingStation: grpctransport.NewServer(
			e.GetNextAvailableDockingStationEndpoint,
			decodeGRPCGetNextAvailableDockingStationRequest,
			encodeGRPCGetNextAvailableDockingStationResponse,
			options...,
		),
		registerShipLanding: grpctransport.NewServer(
			e.RegisterShipLandingEndpoint,
			decodeGRPCRegisterShipLandingRequest,
			encodeGRPCRegisterShipLandingResponse,
			options...,
		),
	}
}

func (g *grpcServer) RegisterShip(ctx context.Context, r *pb.RegisterShipRequest) (*httpbody.HttpBody, error) { //TODO should this return empty?
	_, rep, err := g.registerShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	header := metadata.Pairs(
		"x-no-content", "true",
	)
	grpc.SendHeader(ctx, header)
	return &httpbody.HttpBody{
		ContentType: "application/json",
	}, nil

}
func (g *grpcServer) GetAllShips(ctx context.Context, r *pb.GetAllShipsRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.getAllShips.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	resp := rep.(*pb.GetAllShipsResponse)
	json := serializeGetAllShipsResponse(resp)
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func (g *grpcServer) RegisterStation(ctx context.Context, r *pb.RegisterStationRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.registerStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	resp := rep.(*pb.RegisterStationResponse)
	json := serializeRegisterStationResponse(resp)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil
}

func (g *grpcServer) GetAllStations(ctx context.Context, r *pb.GetAllStationsRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	resp := rep.(*pb.GetAllStationsResponse)
	json := serializeGetAllStationsResponse(resp)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func (g *grpcServer) GetNextAvailableDockingStation(ctx context.Context, r *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	_, rep, err := g.getNextAvailableDockingStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	resp := rep.(*pb.GetNextAvailableDockingStationResponse)
	return resp, nil
}

func (g *grpcServer) RegisterShipLanding(ctx context.Context, r *pb.RegisterShipLandingRequest) (*pb.RegisterShipLandingResponse, error) {
	_, rep, err := g.registerShipLanding.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	resp := rep.(*pb.RegisterShipLandingResponse)
	return resp, nil
}

func decodeGRPCRegisterShipRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.RegisterShipRequest), nil
}
func encodeGRPCRegisterShipResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.RegisterShipResponse), nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.GetAllShipsRequest), nil
}
func encodeGRPCGetAllShipsResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.GetAllShipsResponse), nil
}

func decodeGRPCRegisterStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.RegisterStationRequest), nil
}
func encodeGRPCRegisterStationResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.RegisterStationResponse), nil
}

func decodeGRPCGetAllStationsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.GetAllStationsRequest), nil
}
func encodeGRPCGetAllStationsResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.GetAllStationsResponse), nil
}

func decodeGRPCGetNextAvailableDockingStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.GetNextAvailableDockingStationRequest), nil
}
func encodeGRPCGetNextAvailableDockingStationResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.GetNextAvailableDockingStationResponse), nil
}

func decodeGRPCRegisterShipLandingRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.RegisterShipLandingRequest), nil
}
func encodeGRPCRegisterShipLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.RegisterShipLandingResponse), nil
}

func serializeRegisterStationResponse(resp *pb.RegisterStationResponse) []byte {
	if !errs.IsNil(resp.Failed()) {
		return []byte(fmt.Sprintf(`{"error":"%v"}`, resp.Failed()))
	}
	type dock struct {
		Id              string `json:"id"`
		NumDockingPorts int64  `json:"numDockingPorts"`
	}
	var station map[string]interface{}
	if resp.Station != nil {
		station = map[string]interface{}{
			"id":    resp.Station.Id,
			"docks": []dock{},
		}

		for _, d := range resp.Station.Docks {
			station["docks"] = append(station["docks"].([]dock), dock{
				Id:              d.Id,
				NumDockingPorts: d.NumDockingPorts,
			})
		}
	}
	json, _ := json.Marshal(station)
	return json
}

func serializeGetAllStationsResponse(resp *pb.GetAllStationsResponse) []byte {
	if !errs.IsNil(resp.Failed()) {
		return []byte(fmt.Sprintf(`{"error":"%v"}`, resp.Failed()))
	}
	type dock struct {
		Id              string   `json:"id"`
		NumDockingPorts int64    `json:"numDockingPorts"`
		Occupied        *int64   `json:"occupied"`
		Weight          *float32 `json:"weight"`
	}
	type station struct {
		Id           string   `json:"id"`
		Capacity     int64    `json:"capacity"`
		UsedCapacity *float32 `json:"usedCapacity"`
		Docks        []*dock  `json:"docks"`
	}
	stations := make([]station, 0)
	if resp.Stations != nil {
		for _, s := range resp.Stations {

			station := station{
				Id:           s.Id,
				Capacity:     int64(s.Capacity),
				UsedCapacity: &s.UsedCapacity,
				Docks:        []*dock{},
			}
			for _, d := range s.Docks {
				station.Docks = append(station.Docks, &dock{
					Id:              d.Id,
					NumDockingPorts: d.NumDockingPorts,
					Occupied:        &d.Occupied,
					Weight:          &d.Weight,
				})
			}

			stations = append(stations, station)
		}
	}
	json, _ := json.Marshal(stations)
	return json
}

func serializeGetAllShipsResponse(resp *pb.GetAllShipsResponse) []byte {
	if !errs.IsNil(resp.Failed()) {
		return []byte(fmt.Sprintf(`{"error":"%v"}`, resp.Failed()))
	}
	type ship struct {
		Id     string  `json:"id"`
		Status string  `json:"status"`
		Weight float32 `json:"weight"`
	}

	ships := make([]ship, 0)

	if resp.Ships != nil {
		for _, s := range resp.Ships {
			ship := ship{
				Id:     s.Id,
				Status: s.Status,
				Weight: s.Weight,
			}
			ships = append(ships, ship)
		}
	}
	json, _ := json.Marshal(ships)
	return json
}
