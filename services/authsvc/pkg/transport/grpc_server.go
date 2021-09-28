// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
package transport

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/genproto/googleapis/api/httpbody"
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

func (g *grpcServer) Signup(ctx context.Context, r *pb.SignupRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.signup.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	return rep.(*pb.SignupResponse), nil
}

func decodeGRPCSignupRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.SignupRequest), nil
}
func encodeGRPCSignupResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	resp := grpcResponse.(*pb.SignupResponse)
	return resp, nil
}

func (g *grpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)

	return rep.(*pb.LoginResponse), nil
}
func decodeGRPCLoginRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.LoginRequest), nil
}

func encodeGRPCLoginResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.LoginResponse), nil
}

func (g *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	errs.InjectGrpcStatusCode(ctx, rep)
	return rep.(*pb.CheckTokenResponse), nil
}
func decodeGRPCCheckTokenRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.CheckTokenRequest), nil
}

func encodeGRPCCheckTokenResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	return grpcResponse.(*pb.CheckTokenResponse), nil
}
