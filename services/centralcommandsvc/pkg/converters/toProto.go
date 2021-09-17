package converters

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func DBDtoCreateShipResponseToProto(src dtos.CreateShipResponse) *pb.RegisterShipResponse {
	ret := &pb.RegisterShipResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func DBDtoCreateStationResponseToProto(src dtos.CreateStationResponse) *pb.RegisterStationResponse {
	ret := &pb.RegisterStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func DBDtoGetAllShipsResponseToProto(src dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
	ret := &pb.GetAllShipsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func DBDtoGetAllStationsResponseToProto(src dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	ret := &pb.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
