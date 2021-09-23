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
		v := in.Interface().(dtos.Ship)
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
