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
package converters

import (
	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func CreateStationRequestToProto(src *dtos.CreateStationRequest) *pb.CreateStationRequest {
	ret := &pb.CreateStationRequest{
		Station: &pb.Station{},
	}
	m.Copy(ret.Station, src)

	return ret
}

func CreateStationResponseToProto(src *dtos.CreateStationResponse) *pb.CreateStationResponse {
	ret := &pb.CreateStationResponse{}
	m.Copy(ret, src)
	ret.Error = errs.ToProtoV1(src.Error)
	return ret
}
func CreateShipRequestToProto(src *dtos.CreateShipRequest) *pb.CreateShipRequest {
	ret := &pb.CreateShipRequest{
		Ship: &pb.Ship{},
	}
	m.Copy(ret.Ship, src)

	return ret
}

func CreateShipResponseToProto(src *dtos.CreateShipResponse) *pb.CreateShipResponse {
	ret := &pb.CreateShipResponse{}
	m.Copy(ret, src)
	ret.Error = errs.ToProtoV1(src.Error)
	return ret
}
func GetAllShipsRequestToProto(src *dtos.GetAllShipsRequest) *pb.GetAllShipsRequest {
	ret := &pb.GetAllShipsRequest{}
	if *src == (dtos.GetAllShipsRequest{}) {
		return ret
	}
	m.Copy(ret, &src)
	return ret
}

func GetAllShipsResponseToProto(src *dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllShipsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllShipsResponse{}
	m.Copy(ret, src)
	return ret
}
func GetAllStationsRequestToProto(src *dtos.GetAllStationsRequest) *pb.GetAllStationsRequest {
	ret := &pb.GetAllStationsRequest{ShipId: src.ShipId}
	if *src == (dtos.GetAllStationsRequest{}) {
		return ret
	}
	m.Copy(ret, &src)
	return ret
}

func GetAllStationsResponseToProto(src *dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllStationsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllStationsResponse{}
	m.Copy(ret, src)

	return ret
}

func GetNextAvailableDockingStationResponseToProto(src *dtos.GetNextAvailableDockingStationResponse) *pb.GetNextAvailableDockingStationResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetNextAvailableDockingStationResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetNextAvailableDockingStationResponse{}
	m.Copy(ret, src)
	return ret
}

func GetNextAvailableDockingStationRequestToProto(src *dtos.GetNextAvailableDockingStationRequest) *pb.GetNextAvailableDockingStationRequest {
	return &pb.GetNextAvailableDockingStationRequest{ShipId: src.ShipId}
}
func LandShipToDockRequestToProto(src *dtos.LandShipToDockRequest) *pb.LandShipToDockRequest {
	ret := &pb.LandShipToDockRequest{}
	m.Copy(ret, src)
	return ret
}

func LandShipToDockResponseToProto(src *dtos.LandShipToDockResponse) *pb.LandShipToDockResponse {
	if !errs.IsNil(src.Error) {
		return &pb.LandShipToDockResponse{Error: errs.ToProtoV1(src.Error)}
	}
	return &pb.LandShipToDockResponse{}
}
