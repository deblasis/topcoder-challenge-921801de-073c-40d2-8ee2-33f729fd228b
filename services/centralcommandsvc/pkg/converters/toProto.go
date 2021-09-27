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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func DBDtoCreateStationResponseToProto(src dtos.CreateStationResponse) *pb.RegisterStationResponse {
	if !errs.IsNil(src.Error) {
		return &pb.RegisterStationResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.RegisterStationResponse{
		Station: &pb.Station{},
	}
	if errs := m.Copy(ret.Station, src.Station); len(errs) > 0 {
		panic(errs[0])
	}
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
		errs := m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((*dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		errs := m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		errs := m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		errs := m.Copy(ret, v)

		//it's ignored so we map it manually
		ret.Status = v.Status

		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func DBDtoGetAllStationsResponseToProto(src dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	if !errs.IsNil(src.Error) {
		return &pb.GetAllStationsResponse{Error: errs.ToProtoV1(src.Error)}
	}
	ret := &pb.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
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
