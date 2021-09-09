package converters

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func ProtoRegisterShipResponseToDto(src pb.RegisterShipResponse) dtos.RegisterShipResponse {
	ret := dtos.RegisterShipResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ProtoRegisterStationRequestToDto(src pb.RegisterStationRequest) dtos.RegisterStationRequest {
	ret := &dtos.RegisterStationRequest{}
	if errs := m.Copy(ret, src.Station); len(errs) > 0 {
		panic(errs[0])
	}
	return *ret
}

func ProtoRegisterStationResponseToDto(src pb.RegisterStationResponse) *dtos.RegisterStationResponse {
	ret := &dtos.RegisterStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ProtoGetAllStationsResponseToDto(src pb.GetAllStationsResponse) *dtos.GetAllStationsResponse {
	ret := &dtos.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ProtoGetAllShipsResponseToDto(src pb.GetAllShipsResponse) *dtos.GetAllShipsResponse {
	ret := &dtos.GetAllShipsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}
