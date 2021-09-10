package converters

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func DBDtoCreateShipResponseToProto(src dtos.CreateShipResponse) *pb.RegisterShipResponse {
	ret := &pb.RegisterShipResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

func DBDtoCreateStationResponseToProto(src dtos.CreateStationResponse) *pb.RegisterStationResponse {
	ret := &pb.RegisterStationResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

// func GetAllShipsRequestToDBDto(src pb.GetAllShipsRequest) dtos.GetAllShipsRequest {
// 	ret := &dtos.GetAllShipsRequest{}
// 	if src == (pb.GetAllShipsRequest{}) {
// 		return dtos.GetAllShipsRequest{}
// 	}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }
func DBDtoGetAllShipsResponseToProto(src dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
	ret := &pb.GetAllShipsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

// func GetAllStationsRequestToDBDto(src pb.GetAllStationsRequest) dtos.GetAllStationsRequest {
// 	ret := &dtos.GetAllStationsRequest{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }
func DBDtoGetAllStationsResponseToProto(src dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
	ret := &pb.GetAllStationsResponse{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		panic(errs[0])
	}
	return ret
}

// import (
// 	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
// 	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
// 	m "gopkg.in/jeevatkm/go-model.v1"
// )

// func RegisterStationRequestToProto(src dtos.RegisterStationRequest) pb.RegisterStationRequest {
// 	ret := &pb.RegisterStationRequest{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }

// func RegisterStationResponseToProto(src dtos.RegisterStationResponse) *pb.RegisterStationResponse {
// 	ret := &pb.RegisterStationResponse{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return ret
// }
// func RegisterShipRequestToProto(src dtos.RegisterShipRequest) pb.RegisterShipRequest {
// 	ret := &pb.RegisterShipRequest{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }

// func RegisterShipResponseToProto(src dtos.RegisterShipResponse) *pb.RegisterShipResponse {
// 	ret := &pb.RegisterShipResponse{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return ret
// }
// func GetAllShipsRequestToProto(src dtos.GetAllShipsRequest) pb.GetAllShipsRequest {
// 	ret := &pb.GetAllShipsRequest{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }

// func GetAllShipsResponseToProto(src dtos.GetAllShipsResponse) *pb.GetAllShipsResponse {
// 	ret := &pb.GetAllShipsResponse{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return ret
// }
// func GetAllStationsRequestToProto(src dtos.GetAllStationsRequest) pb.GetAllStationsRequest {
// 	ret := &pb.GetAllStationsRequest{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return *ret
// }

// func GetAllStationsResponseToProto(src dtos.GetAllStationsResponse) *pb.GetAllStationsResponse {
// 	ret := &pb.GetAllStationsResponse{}
// 	if errs := m.Copy(ret, src); len(errs) > 0 {
// 		panic(errs[0])
// 	}
// 	return ret
// }
