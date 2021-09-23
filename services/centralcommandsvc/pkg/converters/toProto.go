package converters

import (
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
