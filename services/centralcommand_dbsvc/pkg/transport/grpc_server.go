package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.CentralCommandDBServiceServer
	serviceStatus grpctransport.Handler

	createShip  grpctransport.Handler
	getAllShips grpctransport.Handler

	createStation  grpctransport.Handler
	getAllStations grpctransport.Handler

	getNextAvailableDockingStation grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandDBServiceServer {

	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{
		createShip: grpctransport.NewServer(
			e.CreateShipEndpoint,
			decodeGRPCCreateShipRequest,
			encodeGRPCCreateShipResponse,
			options...,
		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
			options...,
		),
		createStation: grpctransport.NewServer(
			e.CreateStationEndpoint,
			decodeGRPCCreateStationRequest,
			encodeGRPCCreateStationResponse,
			options...,
		),
		getAllStations: grpctransport.NewServer(
			e.GetAllStationsEndpoint,
			decodeGRPCGetAllStationsRequest,
			encodeGRPCGetAllStationsResponse,
			options...,
		),
		getNextAvailableDockingStation: grpctransport.NewServer(
			e.GetNextAvailableDockingStationEndpoint,
			decodeGRPCGetNextAvailableDockingStationRequest,
			encodeGRPCGetNextAvailableDockingStationResponse,
			options...,
		),
	}
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

func (g *grpcServer) GetNextAvailableDockingStation(ctx context.Context, r *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	_, rep, err := g.getNextAvailableDockingStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetNextAvailableDockingStationResponse), nil
}

func decodeGRPCCreateShipRequest(c context.Context, grpcReq interface{}) (interface{}, error) {

	req := grpcReq.(*pb.CreateShipRequest)
	return converters.ProtoCreateShipRequestToDto(req), nil

	// Id:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.CreateShipResponse)
	return converters.CreateShipResponseToProto(response), nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	if req != nil {
		return &dtos.GetAllShipsRequest{}, nil
	}
	return nil, errs.ErrBadRequest
}
func encodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetAllShipsResponse)
	return converters.GetAllShipsResponseToProto(response), nil
}

func decodeGRPCCreateStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateStationRequest)
	return converters.ProtoCreateStationRequestToDto(req), nil

	// Id:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.CreateStationResponse)
	return converters.CreateStationResponseToProto(response), nil
}

func decodeGRPCGetAllStationsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllStationsRequest)
	if req != nil {
		return &dtos.GetAllStationsRequest{}, nil
	}
	return nil, errs.ErrBadRequest
}
func encodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetAllStationsResponse)
	return converters.GetAllStationsResponseToProto(response), nil
}

func decodeGRPCGetNextAvailableDockingStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetNextAvailableDockingStationRequest)
	return converters.ProtoGetNextAvailableDockingStationRequestToDto(req), nil
}
func encodeGRPCGetNextAvailableDockingStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetNextAvailableDockingStationResponse)
	return converters.GetNextAvailableDockingStationResponseToProto(response), nil
}
