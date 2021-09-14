package transport

import (
	"context"

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	registerShip grpctransport.Handler
	getAllShips  grpctransport.Handler

	registerStation grpctransport.Handler
	getAllStations  grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandServiceServer {
	return &grpcServer{

		registerShip: grpctransport.NewServer(
			e.RegisterShipEndpoint,
			decodeGRPCRegisterShipRequest,
			encodeGRPCRegisterShipResponse,
		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
		),
		registerStation: grpctransport.NewServer(
			e.RegisterStationEndpoint,
			decodeGRPCRegisterStationRequest,
			encodeGRPCRegisterStationResponse,
		),
		getAllStations: grpctransport.NewServer(
			e.GetAllStationsEndpoint,
			decodeGRPCGetAllStationsRequest,
			encodeGRPCGetAllStationsResponse,
		),
	}
}

func (g *grpcServer) RegisterShip(ctx context.Context, r *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error) {
	_, rep, err := g.registerShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RegisterShipResponse), nil
}
func (g *grpcServer) GetAllShips(ctx context.Context, r *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	_, rep, err := g.getAllShips.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllShipsResponse), nil
}

func (g *grpcServer) RegisterStation(ctx context.Context, r *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
	_, rep, err := g.registerStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.RegisterStationResponse), nil
}
func (g *grpcServer) GetAllStations(ctx context.Context, r *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetAllStationsResponse), nil
}

func decodeGRPCRegisterShipRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterShipRequest)
	return req, nil

	// Id:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCRegisterShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterShipResponse)
	//return converters.RegisterShipResponseToProto(*response), nil
	return response, nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	return req, nil
}
func encodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	//return converters.GetAllShipsResponseToProto(*response), nil
	return response, nil
}

func decodeGRPCRegisterStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterStationRequest)
	//return converters.ProtoRegisterStationRequestToDto(*req), nil
	return req, nil

	// Id:       req.User.Id,
	// Username: req.User.Username,
	// Password: req.User.Password,
	// //TODO centralize
	// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),

}
func encodeGRPCRegisterStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterStationResponse)
	//return converters.RegisterStationResponseToProto(*response), nil
	return response, nil
}

func decodeGRPCGetAllStationsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllStationsRequest)
	// if req != nil {
	// 	return pb.GetAllStationsRequest{}, nil
	// }
	// return nil, errors.Str2err("cannot unmarshal GetAllStationsRequest")
	return req, nil
}
func encodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	//return converters.GetAllStationsResponseToProto(*response), nil
	return response, nil
}
