package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	pb.CentralCommandServiceServer
	registerShip grpctransport.Handler
	getAllShips  grpctransport.Handler

	registerStation                grpctransport.Handler
	getAllStations                 grpctransport.Handler
	getNextAvailableDockingStation grpctransport.Handler
	registerShipLanding            grpctransport.Handler
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
		registerShipLanding: grpctransport.NewServer(
			e.RegisterShipLandingEndpoint,
			decodeGRPCRegisterShipLandingRequest,
			encodeGRPCRegisterShipLandingResponse,
			options...,
		),
	}
}

func (g *grpcServer) RegisterShip(ctx context.Context, r *pb.RegisterShipRequest) (*httpbody.HttpBody, error) {
	_, _, err := g.registerShip.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	//resp := rep.(*pb.RegisterShipResponse)
	// if errors.Is(resp.Error, errs.ErrCannotInsertAlreadyExistingEntity) {
	// 	//this will trigger the error handlers so we can alter body and header
	// 	return nil, resp.Error
	// }

	return &httpbody.HttpBody{
		ContentType: "application/json",
	}, nil

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

	json := serializeRegisterStationResponse(resp)

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
	json := serializeGetAllStationsResponse(resp)

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

func (g *grpcServer) RegisterShipLanding(ctx context.Context, r *pb.RegisterShipLandingRequest) (*pb.RegisterShipLandingResponse, error) {
	_, rep, err := g.getAllStations.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.RegisterShipLandingResponse)
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
	resp := grpcResponse.(*pb.RegisterShipResponse)

	header := metadata.Pairs(
		"x-no-content", "true",
	)
	//TODO refactor
	if !errs.IsNil(resp.Failed()) {
		header.Set("x-http-code", fmt.Sprintf("%v", resp.Error.Code))
	}
	grpc.SendHeader(ctx, header)

	//return converters.RegisterShipResponseToProto(*response), nil
	return resp, nil
}

func decodeGRPCGetAllShipsRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetAllShipsRequest)
	return req, nil
}
func encodeGRPCGetAllShipsResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	resp := grpcResponse.(*pb.GetAllShipsResponse)
	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}
	//return converters.GetAllShipsResponseToProto(*response), nil
	return resp, nil
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
	resp := grpcResponse.(*pb.RegisterStationResponse)

	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.RegisterStationResponseToProto(*response), nil
	return resp, nil
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
	resp := grpcResponse.(*pb.GetAllStationsResponse)

	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.GetAllStationsResponseToProto(*response), nil
	return resp, nil
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
	resp := grpcResponse.(*pb.GetNextAvailableDockingStationResponse)

	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.GetNextAvailableDockingStationResponseToProto(*response), nil
	return resp, nil
}

func decodeGRPCRegisterShipLandingRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterShipLandingRequest)
	return req, nil
}
func encodeGRPCRegisterShipLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.LandShipToDockResponse)
	return response, nil
}

func serializeRegisterStationResponse(resp *pb.RegisterStationResponse) []byte {
	type dock struct {
		Id              string `json:"id"`
		NumDockingPorts int64  `json:"numDockingPorts"`
	}
	var station map[string]interface{}
	if resp.Station != nil {
		station = map[string]interface{}{
			"id":    resp.Station.Id,
			"docks": []dock{},
		}

		for _, d := range resp.Station.Docks {
			station["docks"] = append(station["docks"].([]dock), dock{
				Id:              d.Id,
				NumDockingPorts: d.NumDockingPorts,
			})
		}
	}
	json, _ := json.Marshal(station)
	return json
}

func serializeGetAllStationsResponse(resp *pb.GetAllStationsResponse) []byte {
	type dock struct {
		Id              string   `json:"id"`
		NumDockingPorts int64    `json:"numDockingPorts"`
		Occupied        *int64   `json:"occupied"`
		Weight          *float32 `json:"weight"`
	}
	type station struct {
		Id           string   `json:"id"`
		Capacity     int64    `json:"capacity"`
		UsedCapacity *float32 `json:"used_capacity"`
		Docks        []*dock  `json:"docks"`
	}
	stations := make([]station, 0)
	if resp.Stations != nil {
		for _, s := range resp.Stations {

			station := station{
				Id:           s.Id,
				Capacity:     int64(s.Capacity),
				UsedCapacity: &s.UsedCapacity,
				Docks:        []*dock{},
			}
			for _, d := range s.Docks {
				station.Docks = append(station.Docks, &dock{
					Id:              d.Id,
					NumDockingPorts: d.NumDockingPorts,
					Occupied:        &d.Occupied,
					Weight:          &d.Weight,
				})
			}

			stations = append(stations, station)
		}
	}
	json, _ := json.Marshal(stations)
	return json
}
