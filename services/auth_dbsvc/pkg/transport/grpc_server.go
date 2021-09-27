//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/auth_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedAuthDBServiceServer
	createUser        grpctransport.Handler
	getUserByUsername grpctransport.Handler
	getUserById       grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.AuthDBServiceServer {
	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{
		createUser: grpctransport.NewServer(
			e.CreateUserEndpoint,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
			options...,
		),
		getUserByUsername: grpctransport.NewServer(
			e.GetUserByUsernameEndpoint,
			decodeGRPCGetUserByUsernameRequest,
			encodeGRPCGetUserByUsernameResponse,
			options...,
		),
		getUserById: grpctransport.NewServer(
			e.GetUserByIdEndpoint,
			decodeGRPCGetUserByIdRequest,
			encodeGRPCGetUserByIdResponse,
			options...,
		),
	}

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

func (g *grpcServer) GetUserById(ctx context.Context, r *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	_, rep, err := g.getUserById.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetUserByIdResponse), nil
}
func decodeGRPCCreateUserRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateUserRequest)
	return &dtos.CreateUserRequest{
		Id:       req.User.Id,
		Username: req.User.Username,
		Password: req.User.Password,
		//TODO centralize
		// Role: strings.Title(strings.ToLower(strings.TrimLeft(req.User.Role.String(), "ROLE_"))),
		Role: req.User.Role,
	}, nil
}
func encodeGRPCCreateUserResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.CreateUserResponse)

	id := ""
	if response.Id != nil {
		id = response.Id.String()
	}

	return &pb.CreateUserResponse{
		Id:    id,
		Error: errs.ToProtoV1(response.Error),
	}, nil
}

func decodeGRPCGetUserByUsernameRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserByUsernameRequest)
	return &dtos.GetUserByUsernameRequest{
		Username: req.Username,
	}, nil
}
func decodeGRPCGetUserByIdRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserByIdRequest)
	return &dtos.GetUserByIdRequest{
		Id: req.Id,
	}, nil
}

func encodeGRPCGetUserByIdResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetUserResponse)

	var user *pb.User
	if response.User != nil {
		user = &pb.User{
			Id:       response.User.Id,
			Username: response.User.Username,
			Password: response.User.Password,
			Role:     response.User.Role,
		}
	}

	return &pb.GetUserByIdResponse{
		User:  user,
		Error: errs.ToProtoV1(response.Error),
	}, nil

}

func encodeGRPCGetUserByUsernameResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*dtos.GetUserResponse)

	var user *pb.User
	if response.User != nil {
		user = &pb.User{
			Id:       response.User.Id,
			Username: response.User.Username,
			Password: response.User.Password,
			Role:     response.User.Role,
		}
	}

	return &pb.GetUserByUsernameResponse{
		User:  user,
		Error: errs.ToProtoV1(response.Error),
	}, nil

}
