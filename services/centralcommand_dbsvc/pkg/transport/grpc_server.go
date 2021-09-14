package transport

import (
	"context"
	"fmt"

	"deblasis.net/space-traffic-control/common/errors"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
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
	fmt.Printf("CreateShipRequest.ship %+v <\n", req.Ship)
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
	return nil, errors.Str2err("cannot unmarshal GetAllShipsRequest")
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
	return nil, errors.Str2err("cannot unmarshal GetAllStationsRequest")
}
func encodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetAllStationsResponse)
	return converters.GetAllStationsResponseToProto(response), nil
}
