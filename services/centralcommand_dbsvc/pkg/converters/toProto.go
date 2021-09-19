package converters

import (
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
	ret := &pb.GetAllShipsResponse{}

	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
func GetAllStationsRequestToProto(src *dtos.GetAllStationsRequest) *pb.GetAllStationsRequest {
	ret := &pb.GetAllStationsRequest{}
	if *src == (dtos.GetAllStationsRequest{}) {
		return ret
	}
	if errs := m.Copy(ret, &src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func GetAllStationsResponseToProto(src *dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	ret := &pb.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func GetNextAvailableDockingStationResponseToProto(src *dtos.GetNextAvailableDockingStationResponse) *pb.GetNextAvailableDockingStationResponse {
	ret := &pb.GetNextAvailableDockingStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func GetNextAvailableDockingStationRequestToProto(src *dtos.GetNextAvailableDockingStationRequest) *pb.GetNextAvailableDockingStationRequest {
	ret := &pb.GetNextAvailableDockingStationRequest{}
	if errs := m.Copy(ret, &src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
func LandShipToDockRequestToProto(src *dtos.LandShipToDockRequest) *pb.LandShipToDockRequest {
	ret := &pb.LandShipToDockRequest{}
	if errs := m.Copy(ret, &src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}

func LandShipToDockResponseToProto(src *dtos.LandShipToDockResponse) *pb.LandShipToDockResponse {
	ret := &pb.LandShipToDockResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}
