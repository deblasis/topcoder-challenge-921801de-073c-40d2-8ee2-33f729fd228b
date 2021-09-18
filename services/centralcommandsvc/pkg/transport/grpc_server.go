package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcServer struct {
	pb.CentralCommandServiceServer
	registerShip grpctransport.Handler
	getAllShips  grpctransport.Handler

	registerStation                grpctransport.Handler
	getAllStations                 grpctransport.Handler
	getNextAvailableDockingStation grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.CentralCommandServiceServer {

	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{

		registerShip: grpctransport.NewServer(
			e.RegisterShipEndpoint,
			decodeGRPCRegisterShipRequest,
			encodeGRPCRegisterShipResponse,
			options...,
		// grpctransport.ServerBefore(func(c context.Context, m metadata.MD) context.Context {
		// 	level.Info(l).Log("ServerBefore", "RegisterShip",
		// 		"userId", c.Value(common.ContextKeyUserId).(string),
		// 		"role", c.Value(common.ContextKeyUserRole).(string),
		// 	)
		// 	level.Info(l).Log("ServerBefore", "RegisterShip",
		// 		"userId-md", m["x-stc-user-id"],
		// 		"role-md", m["x-stc-user-role"],
		// 	)
		// 	// level.Info(l).Log("ServerBefore", "RegisterShip",
		// 	// 	"userId-md[0]", m["x-stc-user-id"][0],
		// 	// 	"role-md[0]", m["x-stc-user-role"][0],
		// 	// )

		// 	return c
		// }),
		// grpctransport.ServerAfter(func(ctx context.Context, header, trailer *metadata.MD) context.Context {
		// 	var statusCode int
		// 	statusCode, ok := ctx.Value(common.ContextKeyReturnCode).(int)
		// 	if ok {
		// 		header.Append("x-http-code", fmt.Sprintf("%v", statusCode))
		// 	}
		// 	return ctx
		// }),

		),
		getAllShips: grpctransport.NewServer(
			e.GetAllShipsEndpoint,
			decodeGRPCGetAllShipsRequest,
			encodeGRPCGetAllShipsResponse,
			options...,
		),
		registerStation: grpctransport.NewServer(
			e.RegisterStationEndpoint,
			decodeGRPCRegisterStationRequest,
			encodeGRPCRegisterStationResponse,
			options...,
		// grpctransport.ServerErrorHandler(transport.ErrorHandlerFunc(func(ctx context.Context, err error) {

		// 	level.Info(l).Log("ServerErrorHandler", "RegisterStation",
		// 		"err", err,
		// 	)
		// 	if err != nil {
		// 		st, ok := status.FromError(err)
		// 		if ok && st.Code() == codes.AlreadyExists {
		// 			header := metadata.Pairs("x-http-code", "400")
		// 			grpc.SendHeader(ctx, header)
		// 			level.Info(l).Log("ServerErrorHandler", "RegisterStation",
		// 				"grpcheader", "x-http-code",
		// 				"msg", "sent",
		// 			)
		// 		}
		// 	}
		// })),
		// grpctransport.ServerAfter(func(ctx context.Context, header, trailer *metadata.MD) context.Context {
		// 	level.Info(l).Log("ServerAfter", "RegisterStation",
		// 		"statusCode_from_ctx", ctx.Value(common.ContextKeyReturnCode),
		// 	)

		// 	//panic(ctx.Value(common.ContextKeyReturnCode))
		// 	// var statusCode int
		// 	// statusCode, ok := ctx.Value(common.ContextKeyReturnCode).(int)
		// 	// if ok {
		// 	// 	level.Info(l).Log("ServerAfter", "RegisterStations",
		// 	// 		"statusCode", statusCode,
		// 	// 		"msg", "appending to header md as x-http-code",
		// 	// 	)
		// 	// 	header.Append("x-http-code", fmt.Sprintf("%v", statusCode))
		// 	// }

		// 	return ctx
		// }),
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

func (g *grpcServer) RegisterShip(ctx context.Context, r *pb.RegisterShipRequest) (*emptypb.Empty, error) {
	_, rep, err := g.registerShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	resp := rep.(*pb.RegisterShipResponse)
	if resp.Error == errs.ErrCannotInsertAlreadyExistingEntity.Error() {
		//this will trigger the error handlers so we can alter body and header
		return nil, errs.ErrCannotInsertAlreadyExistingEntity
	}

	return &emptypb.Empty{}, nil

}
func (g *grpcServer) GetAllShips(ctx context.Context, r *pb.GetAllShipsRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.getAllShips.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.GetAllShipsResponse)
	json, _ := json.Marshal(resp.Ships)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func (g *grpcServer) RegisterStation(ctx context.Context, r *pb.RegisterStationRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.registerStation.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	resp := rep.(*pb.RegisterStationResponse)
	json, _ := json.Marshal(resp)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil
}
func (g *grpcServer) GetAllStations(ctx context.Context, r *pb.GetAllStationsRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.GetAllStationsResponse)
	json, _ := json.Marshal(resp.Stations)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func (g *grpcServer) GetNextAvailableDockingStation(ctx context.Context, r *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.GetNextAvailableDockingStationResponse)
	return resp, nil
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
func encodeGRPCRegisterShipResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterShipResponse)

	//TODO refactor
	if f, ok := grpcResponse.(endpoint.Failer); ok && f.Failed() != nil {

		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(f.Failed())),
			"x-stc-error", f.Failed().Error(),
			"x-no-content", "true",
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.RegisterShipResponseToProto(*response), nil
	return response, nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	return req, nil
}
func encodeGRPCGetAllShipsResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	response := grpcResponse.(*pb.GetAllShipsResponse)
	//TODO: refactor
	if response.Failed() != nil {
		errs.GetErrorContainer(ctx).Domain = errs.Str2err(response.Error)
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(errs.Str2err(response.Error))),
			"x-stc-error", response.Failed().Error(),
		)
		grpc.SendHeader(ctx, header)
	}

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
func encodeGRPCRegisterStationResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	// if f, ok := grpcResponse.(endpoint.Failer); ok && f.Failed() != nil {
	// 	return errorEncoder(ctx, f.Failed(), grpcResponse.(*pb.RegisterStationResponse))
	// }
	response := grpcResponse.(*pb.RegisterStationResponse)

	//TODO refactor
	if f, ok := grpcResponse.(endpoint.Failer); ok && f.Failed() != nil {

		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(f.Failed())),
			"x-stc-error", f.Failed().Error(),
		)
		grpc.SendHeader(ctx, header)
	}

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
func encodeGRPCGetAllStationsResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)

	//TODO: refactor
	if response.Failed() != nil {
		errs.GetErrorContainer(ctx).Domain = errs.Str2err(response.Error)
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(errs.Str2err(response.Error))),
			"x-stc-error", response.Failed().Error(),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.GetAllStationsResponseToProto(*response), nil
	return response, nil
}

func decodeGRPCGetNextAvailableDockingStationRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetNextAvailableDockingStationRequest)
	// if req != nil {
	// 	return pb.GetNextAvailableDockingStationRequest{}, nil
	// }
	// return nil, errors.Str2err("cannot unmarshal GetNextAvailableDockingStationRequest")
	return req, nil
}
func encodeGRPCGetNextAvailableDockingStationResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetNextAvailableDockingStationResponse)

	//TODO: refactor
	if response.Failed() != nil {
		errs.GetErrorContainer(ctx).Domain = errs.Str2err(response.Error)
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(errs.Str2err(response.Error))),
			"x-stc-error", response.Failed().Error(),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.GetNextAvailableDockingStationResponseToProto(*response), nil
	return response, nil
}
