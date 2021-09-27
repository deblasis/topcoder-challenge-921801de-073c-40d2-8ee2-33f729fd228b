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
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func StationToDto(src *model.Station) *dtos.Station {
	ret := &dtos.Station{}
	m.Copy(ret, src)
	return ret
}

func StationsToDto(src []model.Station) []dtos.Station {
	ret := []dtos.Station{}
	for _, x := range src {
		ret = append(ret, *StationToDto(&x))
	}
	return ret
}

func DockToDto(src model.Dock) dtos.Dock {
	ret := &dtos.Dock{}
	m.Copy(ret, src)
	return *ret
}

func DocksToDto(src []model.Dock) []dtos.Dock {
	var ret []dtos.Dock
	for _, x := range src {
		ret = append(ret, DockToDto(x))
	}
	return ret
}

func ShipToDto(src *model.Ship) *dtos.Ship {
	ret := &dtos.Ship{}
	m.Copy(ret, src)
	return ret
}

func ShipsToDto(src []model.Ship) []dtos.Ship {
	ret := []dtos.Ship{}
	for _, x := range src {
		ret = append(ret, *ShipToDto(&x))
	}
	return ret
}

func NextAvailableDockingStationToDto(src *model.NextAvailableDockingStation) *dtos.NextAvailableDockingStation {
	if src == nil {
		return nil
	}
	ret := &dtos.NextAvailableDockingStation{}
	m.Copy(ret, src)
	return ret
}

func ProtoCreateStationRequestToDto(src *pb.CreateStationRequest) *dtos.CreateStationRequest {
	ret := &dtos.CreateStationRequest{}
	srcStuct := src.Station
	if srcStuct == nil {
		return ret
	}

	m.Copy(ret, srcStuct)
	return ret
}

func ProtoGetNextAvailableDockingStationRequestToDto(src *pb.GetNextAvailableDockingStationRequest) *dtos.GetNextAvailableDockingStationRequest {
	return &dtos.GetNextAvailableDockingStationRequest{ShipId: src.ShipId}
}

func ProtoCreateShipRequestToDto(src *pb.CreateShipRequest) *dtos.CreateShipRequest {
	ret := &dtos.CreateShipRequest{}
	srcStuct := src.Ship
	if srcStuct == nil {
		return ret
	}
	m.Copy(ret, srcStuct)
	return ret
}
func ProtoCreateShipResponseToDto(src *pb.CreateShipResponse) *dtos.CreateShipResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.CreateShipResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.CreateShipResponse{
		Ship: &dtos.Ship{},
	}
	m.Copy(ret.Ship, src.Ship)

	return ret
}

func ProtoCreateStationResponseToDto(src *pb.CreateStationResponse) *dtos.CreateStationResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.CreateStationResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.CreateStationResponse{
		Station: &dtos.Station{},
	}
	m.Copy(ret.Station, src.Station)

	return ret
}

func ProtoGetAllStationsResponseToDto(src *pb.GetAllStationsResponse) *dtos.GetAllStationsResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.GetAllStationsResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.GetAllStationsResponse{
		Stations: []dtos.Station{},
	}
	if src.Stations == nil {
		return ret
	}
	m.Copy(ret, src)
	return ret
}

func ProtoGetAllShipsResponseToDto(src *pb.GetAllShipsResponse) *dtos.GetAllShipsResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.GetAllShipsResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.GetAllShipsResponse{
		Ships: []dtos.Ship{},
	}
	if src.Ships == nil {
		return ret
	}
	m.Copy(ret, src)
	return ret
}

func ProtoGetNextAvailableDockingStationResponseToDto(src *pb.GetNextAvailableDockingStationResponse) *dtos.GetNextAvailableDockingStationResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.GetNextAvailableDockingStationResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.GetNextAvailableDockingStationResponse{}
	m.Copy(ret, src)
	return ret
}

func ProtoLandShipToDockResponseToDto(src *pb.LandShipToDockResponse) *dtos.LandShipToDockResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.LandShipToDockResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.LandShipToDockResponse{}
	m.Copy(ret, src)
	return ret
}
func ProtoLandShipToDockRequestToDto(src *pb.LandShipToDockRequest) *dtos.LandShipToDockRequest {
	ret := &dtos.LandShipToDockRequest{}
	m.Copy(ret, src)
	return ret
}
