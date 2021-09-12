package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/auth_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	serviceStatus     grpctransport.Handler
	createUser        grpctransport.Handler
	getUserByUsername grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.AuthDBServiceServer {
	return &grpcServer{
		serviceStatus: grpctransport.NewServer(
			e.StatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
		),
		createUser: grpctransport.NewServer(
			e.CreateUserEndpoint,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
		),
		getUserByUsername: grpctransport.NewServer(
			e.GetUserByUsernameEndpoint,
			decodeGRPCGetUserByUsernameRequest,
			encodeGRPCGetUserByUsernameResponse,
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

func (g *grpcServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, rep, err := g.createUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserResponse), nil
}
func (g *grpcServer) GetUserByUsername(ctx context.Context, r *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	_, rep, err := g.getUserByUsername.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserByUsernameResponse), nil
}

func decodeGRPCCreateUserRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateUserRequest)
	return dtos.CreateUserRequest{
		Id:       req.User.Id,
		Username: req.User.Username,
		Password: req.User.Password,
		//TODO centralize
		// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),
		Role: req.User.Role,
	}, nil
}
func encodeGRPCCreateUserResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.CreateUserResponse)
	return &pb.CreateUserResponse{
		Id:    response.Id,
		Error: response.Err,
	}, nil
}

func decodeGRPCGetUserByUsernameRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserByUsernameRequest)
	return dtos.GetUserByUsernameRequest{
		Username: req.Username,
	}, nil
}

func encodeGRPCGetUserByUsernameResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(dtos.GetUserByUsernameResponse)

	//TODO: centralise
	// roleId := pb.User_Role_value[strings.ToUpper("ROLE_"+response.User.Role)]
	// if roleId <= 0 {
	// 	return nil, errors.New("cannot unmarshal role")
	// }
	return &pb.GetUserByUsernameResponse{
		User: &pb.User{
			Id:       response.User.Id,
			Username: response.User.Username,
			Password: response.User.Password,
			Role:     response.User.Role,
		},
		Error: response.Err,
	}, nil
}
