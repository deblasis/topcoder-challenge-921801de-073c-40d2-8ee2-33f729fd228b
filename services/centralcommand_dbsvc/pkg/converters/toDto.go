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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
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
		//panic(errs[0])
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
		//panic(errs[0])
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

func NextAvailableDockingStationToDto(src *model.NextAvailableDockingStation) *dtos.NextAvailableDockingStation {
	if src == nil {
		return nil
	}
	ret := &dtos.NextAvailableDockingStation{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func ProtoCreateStationRequestToDto(src *pb.CreateStationRequest) *dtos.CreateStationRequest {
	ret := &dtos.CreateStationRequest{}
	srcStuct := src.Station
	if srcStuct == nil {
		return ret
	}

	if errs := m.Copy(ret, srcStuct); len(errs) > 0 {
		//panic(errs[0])
	}
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
	if errs := m.Copy(ret, srcStuct); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}
func ProtoCreateShipResponseToDto(src *pb.CreateShipResponse) *dtos.CreateShipResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.CreateShipResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.CreateShipResponse{
		Ship: &dtos.Ship{},
	}
	if errs := m.Copy(ret.Ship, src.Ship); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func ProtoCreateStationResponseToDto(src *pb.CreateStationResponse) *dtos.CreateStationResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.CreateStationResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.CreateStationResponse{
		Station: &dtos.Station{},
	}
	if errs := m.Copy(ret.Station, src.Station); len(errs) > 0 {
		//panic(errs[0])
	}
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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
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
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func ProtoGetNextAvailableDockingStationResponseToDto(src *pb.GetNextAvailableDockingStationResponse) *dtos.GetNextAvailableDockingStationResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.GetNextAvailableDockingStationResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.GetNextAvailableDockingStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func ProtoLandShipToDockResponseToDto(src *pb.LandShipToDockResponse) *dtos.LandShipToDockResponse {
	if !errs.IsNil(src.Error) {
		return &dtos.LandShipToDockResponse{Error: errs.FromProtoV1(src.Error)}
	}
	ret := &dtos.LandShipToDockResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
func ProtoLandShipToDockRequestToDto(src *pb.LandShipToDockRequest) *dtos.LandShipToDockRequest {
	ret := &dtos.LandShipToDockRequest{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
