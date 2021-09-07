package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	serviceStatus grpctransport.Handler

	createShip  grpctransport.Handler
	getAllShips grpctransport.Handler

	createStation  grpctransport.Handler
	getAllStations grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandDBServiceServer {
	return &grpcServer{
		serviceStatus: grpctransport.NewServer(
			e.StatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
		),
		createShip: grpctransport.NewServer(
			e.CreateShipEndpoint,
			decodeGRPCCreateShipRequest,
			encodeGRPCCreateShipResponse,
		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
		),
		createStation: grpctransport.NewServer(
			e.CreateStationEndpoint,
			decodeGRPCCreateStationRequest,
			encodeGRPCCreateStationResponse,
		),
		getAllStations: grpctransport.NewServer(
			e.GetAllStationsEndpoint,
			decodeGRPCGetAllStationsRequest,
			encodeGRPCGetAllStationsResponse,
		),
	}

}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusResponse, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceStatusResponse), nil
}
func decodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (interface{}, error) {
	var req healthcheck.ServiceStatusRequest
	return req, nil
}
func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	Response := grpcResponse.(healthcheck.ServiceStatusResponse)
	return &pb.ServiceStatusResponse{Code: Response.Code, Err: Response.Err}, nil
}

func (g *grpcServer) CreateShip(ctx context.Context, r *pb.CreateShipRequest) (*pb.CreateShipResponse, error) {
	_, rep, err := g.createShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateShipResponse), nil
}
func (g *grpcServer) GetAllShips(ctx context.Context, r *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	_, rep, err := g.getAllShips.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllShipsResponse), nil
}

func (g *grpcServer) CreateStation(ctx context.Context, r *pb.CreateStationRequest) (*pb.CreateStationResponse, error) {
	_, rep, err := g.createStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateStationResponse), nil
}
func (g *grpcServer) GetAllStations(ctx context.Context, r *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllStationsResponse), nil
}

func decodeGRPCCreateShipRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateShipRequest)
	return dtos.CreateShipRequest{
		Weight: req.Ship.Weight,
	}, nil

	// ID:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.CreateShipResponse)
	return &pb.CreateShipResponse{
		Ship: &pb.Ship{
			Id:     response.Ship.ID,
			Status: pb.Ship_Status(0), //TODO
			Weight: response.Ship.Weight,
		},
		Error: response.Err,
	}, nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	if req != nil {
		return dtos.GetAllShipsRequest{}, nil
	}
	return nil, errors.Str2err("cannot unmarshal GetAllShipsRequest")
}
func encodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.GetAllShipsResponse)
	return &pb.GetAllShipsResponse{
		Ships: modelToProtoShips(response.Ships),
	}, nil
}

func decodeGRPCCreateStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateStationRequest)
	return dtos.CreateStationRequest{
		Capacity: req.Station.Capacity,
		Docks:    protoToModelDocks(req.Station.Docks), //TODO converte
	}, nil

	// ID:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.CreateStationResponse)
	return &pb.CreateStationResponse{
		Station: modelToProtoStation(response.Station),
		Error:   response.Err,
	}, nil
}
func decodeGRPCGetAllStationsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllStationsRequest)
	if req != nil {
		return dtos.GetAllStationsRequest{}, nil
	}
	return nil, errors.Str2err("cannot unmarshal GetAllStationsRequest")
}
func encodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.GetAllStationsResponse)
	return &pb.GetAllStationsResponse{
		Stations: modelToProtoStations(response.Stations),
	}, nil
}

//TODO refactor into converters across the board

func protoToModelDocks(src []*pb.Dock) []*model.Dock {
	var ret []*model.Dock
	for _, x := range src {
		ret = append(ret, protoToModelDock(x))
	}
	return ret
}

func protoToModelDock(src *pb.Dock) *model.Dock {
	return &model.Dock{
		ID:              src.Id,
		StationId:       src.StationId,
		NumDockingPorts: src.NumDockingPorts,
		Occupied:        src.Occupied,
		Weight:          src.Weight,
	}
}

func protoToModelStation(src *pb.Station) *model.Station {
	return &model.Station{
		ID:           src.Id,
		Capacity:     src.Capacity,
		UsedCapacity: src.UsedCapacity,
		Docks:        protoToModelDocks(src.Docks),
	}
}

func modelToProtoDocks(src []*model.Dock) []*pb.Dock {
	var ret []*pb.Dock
	for _, x := range src {
		ret = append(ret, modelToProtoDock(x))
	}
	return ret
}

func modelToProtoDock(src *model.Dock) *pb.Dock {
	return &pb.Dock{
		Id:              src.ID,
		StationId:       src.StationId,
		NumDockingPorts: src.NumDockingPorts,
		Occupied:        src.Occupied,
		Weight:          src.Weight,
	}
}

func modelToProtoStation(src *model.Station) *pb.Station {
	return &pb.Station{
		Id:           src.ID,
		Capacity:     src.Capacity,
		UsedCapacity: src.UsedCapacity,
		Docks:        modelToProtoDocks(src.Docks),
	}
}

func modelToProtoStations(src []*model.Station) []*pb.Station {
	var ret []*pb.Station
	for _, x := range src {
		ret = append(ret, modelToProtoStation(x))
	}
	return ret
}

func modelToProtoShip(src *model.Ship) *pb.Ship {
	return &pb.Ship{
		Id:     src.ID,
		Status: modelToProtoShipStatus(src.Status),
		Weight: src.Weight,
	}
}

func modelToProtoShips(src []*model.Ship) []*pb.Ship {
	var ret []*pb.Ship
	for _, x := range src {
		ret = append(ret, modelToProtoShip(x))
	}
	return ret
}

func modelToProtoShipStatus(src string) pb.Ship_Status {
	switch src {
	case "in-flight": //TODO const
		return pb.Ship_STATUS_INFLIGHT
	case "docked": //TODO const
		return pb.Ship_STATUS_DOCKED
	default:
		return pb.Ship_STATUS_UNSPECIFIED
	}
}

func protoToModelShip(src *pb.Ship) *model.Ship {
	return &model.Ship{
		ID:     src.Id,
		Status: protoToModelShipStatus(src.Status),
		Weight: src.Weight,
	}
}

func protoToModelShips(src []*pb.Ship) []*model.Ship {
	var ret []*model.Ship
	for _, x := range src {
		ret = append(ret, protoToModelShip(x))
	}
	return ret
}

func protoToModelShipStatus(src pb.Ship_Status) string {
	switch src {
	case pb.Ship_STATUS_INFLIGHT:
		return "in-flight" //TODO const
	case pb.Ship_STATUS_DOCKED:
		return "docked" //TODO const
	default:
		return "" //TODO const
	}
}

// func decodeGRPCGetUserByUsernameRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
// 	req := grpcReq.(*pb.GetUserByUsernameRequest)
// 	return dtos.GetUserByUsernameRequest{
// 		Username: req.Username,
// 	}, nil
// }

// func encodeGRPCGetUserByUsernameResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
// 	response := grpcResponse.(dtos.GetUserByUsernameResponse)

// 	//TODO: centralise
// 	roleId := pb.User_Role_value[strings.ToUpper("ROLE_"+response.User.Role)]
// 	if roleId <= 0 {
// 		return nil, errors.New("cannot unmarshal role")
// 	}
// 	return &pb.GetUserByUsernameResponse{
// 		User: &pb.User{
// 			Id:       response.User.ID,
// 			Username: response.User.Username,
// 			Password: response.User.Password,
// 			Role:     pb.User_Role(roleId),
// 		},
// 		Error: response.Err,
// 	}, nil
// }
