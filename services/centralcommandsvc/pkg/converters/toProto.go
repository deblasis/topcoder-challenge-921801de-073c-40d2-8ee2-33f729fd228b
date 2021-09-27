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
	"reflect"

	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func DBDtoCreateShipResponseToProto(src dtos.CreateShipResponse) *pb.RegisterShipResponse {
	if !errs.IsNil(src.Error) {
		return &pb.RegisterShipResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.RegisterShipResponse{}
	m.Copy(ret, src)
	return ret
}

func DBDtoCreateStationResponseToProto(src dtos.CreateStationResponse) *pb.RegisterStationResponse {
	if !errs.IsNil(src.Error) {
		return &pb.RegisterStationResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.RegisterStationResponse{
		Station: &pb.Station{},
	}
	m.Copy(ret.Station, src.Station)
	return ret
}

func DBDtoGetAllShipsResponseToProto(src dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllShipsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllShipsResponse{}
	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((*dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		return reflect.ValueOf(ret), nil
	})

	m.Copy(ret, src)
	return ret
}

func DBDtoGetAllStationsResponseToProto(src dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllStationsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllStationsResponse{}
	m.Copy(ret, src)
	return ret
}

func DBDtoGetNextAvailableDockingStationResponseToProto(src *dtos.GetNextAvailableDockingStationResponse) *pb.GetNextAvailableDockingStationResponse {
	ret := &pb.GetNextAvailableDockingStationResponse{}
	na := src.NextAvailableDockingStation
	if na != nil {
		ret.NextAvailableDockingStation = &pb.NextAvailableDockingStation{
			DockId:     na.DockId,
			StationId:  na.StationId,
			ShipWeight: na.ShipWeight,
		}
		if na.AvailableCapacity != nil {
			ret.NextAvailableDockingStation.AvailableCapacity = *na.AvailableCapacity
		}
		if na.AvailableDocksAtStation != nil {
			ret.NextAvailableDockingStation.AvailableDocksAtStation = *na.AvailableDocksAtStation
		}
		if na.SecondsUntilNextAvailable != nil {
			ret.NextAvailableDockingStation.SecondsUntilNextAvailable = *na.SecondsUntilNextAvailable
		}
	}
	return ret
}
