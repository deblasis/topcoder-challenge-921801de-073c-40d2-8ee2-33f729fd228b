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
	if errs := m.Copy(ret.Station, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func CreateStationResponseToProto(src *dtos.CreateStationResponse) *pb.CreateStationResponse {
	ret := &pb.CreateStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	ret.Error = errs.ToProtoV1(src.Error)
	return ret
}
func CreateShipRequestToProto(src *dtos.CreateShipRequest) *pb.CreateShipRequest {
	ret := &pb.CreateShipRequest{
		Ship: &pb.Ship{},
	}
	if errs := m.Copy(ret.Ship, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func CreateShipResponseToProto(src *dtos.CreateShipResponse) *pb.CreateShipResponse {
	ret := &pb.CreateShipResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	ret.Error = errs.ToProtoV1(src.Error)
	return ret
}
func GetAllShipsRequestToProto(src *dtos.GetAllShipsRequest) *pb.GetAllShipsRequest {
	ret := &pb.GetAllShipsRequest{}
	if *src == (dtos.GetAllShipsRequest{}) {
		return ret
	}
	if errs := m.Copy(ret, &src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func GetAllShipsResponseToProto(src *dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllShipsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllShipsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}
func GetAllStationsRequestToProto(src *dtos.GetAllStationsRequest) *pb.GetAllStationsRequest {
	ret := &pb.GetAllStationsRequest{ShipId: src.ShipId}
	if *src == (dtos.GetAllStationsRequest{}) {
		return ret
	}
	if errs := m.Copy(ret, &src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func GetAllStationsResponseToProto(src *dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllStationsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}

	return ret
}

func GetNextAvailableDockingStationResponseToProto(src *dtos.GetNextAvailableDockingStationResponse) *pb.GetNextAvailableDockingStationResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetNextAvailableDockingStationResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetNextAvailableDockingStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}

	return ret
}

func GetNextAvailableDockingStationRequestToProto(src *dtos.GetNextAvailableDockingStationRequest) *pb.GetNextAvailableDockingStationRequest {
	return &pb.GetNextAvailableDockingStationRequest{ShipId: src.ShipId}
}
func LandShipToDockRequestToProto(src *dtos.LandShipToDockRequest) *pb.LandShipToDockRequest {
	ret := &pb.LandShipToDockRequest{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func LandShipToDockResponseToProto(src *dtos.LandShipToDockResponse) *pb.LandShipToDockResponse {
	if !errs.IsNil(src.Error) {
		return &pb.LandShipToDockResponse{Error: errs.ToProtoV1(src.Error)}
	}
	return &pb.LandShipToDockResponse{}
}
