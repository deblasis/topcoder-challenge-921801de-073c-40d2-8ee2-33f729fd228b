package transport

import (
	"context"
	"errors"
	"strings"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pb"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	serviceStatus     grpctransport.Handler
	createUser        grpctransport.Handler
	getUserByUsername grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.AuthDBSvcServer {
	return &grpcServer{
		serviceStatus: grpctransport.NewServer(
			e.StatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusReply,
		),
		createUser: grpctransport.NewServer(
			e.CreateUserEndpoint,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserReply,
		),
		getUserByUsername: grpctransport.NewServer(
			e.GetUserByUsernameEndpoint,
			decodeGRPCGetUserByUsernameRequest,
			encodeGRPCGetUserByUsernameReply,
		),
	}

}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceStatusReply), nil
}

func (g *grpcServer) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	_, rep, err := g.createUser.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateUserReply), nil
}
func (g *grpcServer) GetUserByUsername(ctx context.Context, r *pb.GetUserByUsernameRequest) (*pb.UserReply, error) {
	_, rep, err := g.getUserByUsername.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserReply), nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (interface{}, error) {
	var req model.ServiceStatusRequest
	return req, nil
}
func encodeGRPCServiceStatusReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(model.ServiceStatusReply)
	return &pb.ServiceStatusReply{Code: reply.Code, Err: reply.Err}, nil
}

func decodeGRPCCreateUserRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateUserRequest)
	return model.CreateUserRequest{
		ID:       req.User.Id,
		Username: req.User.Username,
		Password: req.User.Password,
		Role:     strings.Title(strings.ToLower(req.User.Role.String())),
	}, nil
}
func encodeGRPCCreateUserReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(model.CreateUserReply)
	return &pb.CreateUserReply{
		Id:    reply.Id,
		Error: reply.Err,
	}, nil
}

func decodeGRPCGetUserByUsernameRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserByUsernameRequest)
	return model.GetUserByUsernameRequest{
		Username: req.Username,
	}, nil
}

func encodeGRPCGetUserByUsernameReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(model.GetUserByUsernameReply)

	roleId := pb.User_Role_value[strings.ToUpper(reply.User.Role)]
	if roleId <= 0 {
		return nil, errors.New("cannot unmarshal role")
	}
	return &pb.UserReply{
		User: &pb.User{
			Id:       reply.User.ID,
			Username: reply.User.Username,
			Password: reply.User.Password,
			Role:     pb.User_Role(roleId),
		},
		Error: "",
	}, nil
}
