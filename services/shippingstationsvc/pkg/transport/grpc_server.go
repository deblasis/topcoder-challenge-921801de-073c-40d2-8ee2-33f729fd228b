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
	"encoding/json"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	//TODO check for consistency
	pb.UnimplementedShippingStationServiceServer

	requestLanding grpctransport.Handler
	landing        grpctransport.Handler
}

func NewGRPCServer(e endpoints.EndpointSet, l log.Logger) pb.ShippingStationServiceServer {

	options := transport_conf.GetCommonGRPCServerOptions(l)

	return &grpcServer{

		requestLanding: grpctransport.NewServer(
			e.RequestLandingEndpoint,
			decodeGRPCRequestLandingRequest,
			encodeGRPCRequestLandingResponse,
			options...,
		),
		landing: grpctransport.NewServer(
			e.LandingEndpoint,
			decodeGRPCLandingRequest,
			encodeGRPCLandingResponse,
			options...,
		),
	}
}

func (g *grpcServer) RequestLanding(ctx context.Context, r *pb.RequestLandingRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.requestLanding.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.RequestLandingResponse)

	json := serializeRequestLandingResponse(resp)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func (g *grpcServer) Landing(ctx context.Context, r *pb.LandingRequest) (*httpbody.HttpBody, error) {
	_, rep, err := g.landing.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.LandingResponse)
	json := serializeLandingResponse(resp)

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        []byte(json),
	}, nil

}

func decodeGRPCRequestLandingRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RequestLandingRequest)
	return req, nil
}
func encodeGRPCRequestLandingResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RequestLandingResponse)

	//TODO refactor
	if !errs.IsNil(response.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", response.Error.Code),
			"x-no-content", "true",
		)
		grpc.SendHeader(ctx, header)
	}

	return response, nil
}

func decodeGRPCLandingRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LandingRequest)
	return req, nil
}
func encodeGRPCLandingResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {

	response := grpcResponse.(*pb.LandingResponse)
	//TODO: refactor
	if !errs.IsNil(response.Failed()) {
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", response.Error.Code),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.LandingResponseToProto(*response), nil
	return response, nil
}

func serializeRequestLandingResponse(resp *pb.RequestLandingResponse) []byte {
	type landResponse struct {
		Command        string `json:"command"`
		DockingStation string `json:"dockingStation"`
	}
	type waitResponse struct {
		Command  string `json:"command"`
		Duration int    `json:"duration"`
		Error    string `json:"error,omitempty"`
	}
	var ret []byte
	if resp.Command == pb.RequestLandingResponse_LAND || resp.GetDuration() < 0 {
		ret, _ = json.Marshal(&landResponse{
			Command:        "land",
			DockingStation: resp.GetDockingStationId(),
		})
	} else if errs.IsNil(resp.Failed()) {
		ret, _ = json.Marshal(&waitResponse{
			Command:  "wait",
			Duration: int(resp.GetDuration()),
		})
	} else {
		ret, _ = json.Marshal(&waitResponse{
			Command: "wait",
			Error:   resp.Failed().Error(),
		})
	}
	return ret

}

func serializeLandingResponse(resp *pb.LandingResponse) []byte {
	var ret []byte
	if !errs.IsNil(resp.Failed()) {
		return []byte(fmt.Sprintf(`{"error":"%v"}`, resp.Failed()))
	} else {
		type okResponse struct {
			Message string `json:"message"`
		}
		ret, _ = json.Marshal(&okResponse{
			Message: "Landed successfully",
		})
	}
	return ret
}
