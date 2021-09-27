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
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	pb.UnimplementedAuthServiceServer
	signup     grpctransport.Handler
	login      grpctransport.Handler
	checkToken grpctransport.Handler
}

func NewGRPCServer(l log.Logger, e endpoints.EndpointSet) pb.AuthServiceServer {
	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{
		signup: grpctransport.NewServer(
			e.SignupEndpoint,
			decodeGRPCSignupRequest,
			encodeGRPCSignupResponse,
			options...,
		),
		login: grpctransport.NewServer(
			e.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
			options...,
		),
		checkToken: grpctransport.NewServer(
			e.CheckTokenEndpoint,
			decodeGRPCCheckTokenRequest,
			encodeGRPCCheckTokenResponse,
			options...,
		),
	}
}

func (g *grpcServer) Signup(ctx context.Context, r *pb.SignupRequest) (*pb.SignupResponse, error) {
	_, rep, err := g.signup.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SignupResponse), nil
}

func decodeGRPCSignupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.SignupRequest), nil
}
func encodeGRPCSignupResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	resp := grpcResponse.(*pb.SignupResponse)
	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	return resp, nil
}

func (g *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginResponse), nil
}
func decodeGRPCLoginRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.LoginRequest), nil
}

func encodeGRPCLoginResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	resp := grpcResponse.(*pb.LoginResponse)
	//TODO: refactor
	if !errs.IsNil(resp.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", resp.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	return resp, nil
}

func (g *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CheckTokenResponse), nil
}
func decodeGRPCCheckTokenRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.CheckTokenRequest), nil
}

func encodeGRPCCheckTokenResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.CheckTokenResponse), nil
}
