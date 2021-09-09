package converters

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func StationToDto(src *model.Station) *dtos.Station {
	ret := &dtos.Station{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ShipsToDto(src []model.Ship) []dtos.Ship {
	ret := []dtos.Ship{}
	for _, x := range src {
		ret = append(ret, *ShipToDto(&x))
	}
	return ret
}

func ProtoCreateShipResponseToDto(src pb.CreateShipResponse) dtos.CreateShipResponse {
	ret := dtos.CreateShipResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ProtoCreateStationRequestToDto(src pb.CreateStationRequest) dtos.CreateStationRequest {
	ret := &dtos.CreateStationRequest{}
	if errs := m.Copy(ret, src.Station); len(errs) > 0 {
		panic(errs[0])
	}
	return *ret
}

func ProtoCreateStationResponseToDto(src pb.CreateStationResponse) *dtos.CreateStationResponse {
	ret := &dtos.CreateStationResponse{}
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
