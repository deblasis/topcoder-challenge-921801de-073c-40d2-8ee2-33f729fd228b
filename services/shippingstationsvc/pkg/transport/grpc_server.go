package transport

import (
	"context"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
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

func (g *grpcServer) RequestLanding(ctx context.Context, r *pb.RequestLandingRequest) (*pb.RequestLandingResponse, error) {
	_, rep, err := g.requestLanding.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.RequestLandingResponse)
	return resp, nil
}

func (g *grpcServer) Landing(ctx context.Context, r *pb.LandingRequest) (*pb.LandingResponse, error) {
	_, rep, err := g.landing.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	resp := rep.(*pb.LandingResponse)
	return resp, nil
}

func decodeGRPCRequestLandingRequest(c context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RequestLandingRequest)
	return req, nil
}
func encodeGRPCRequestLandingResponse(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RequestLandingResponse)

	//TODO refactor
	if f, ok := grpcResponse.(endpoint.Failer); ok && f.Failed() != nil {

		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(f.Failed())),
			"x-stc-error", f.Failed().Error(),
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
	if response.Failed() != nil {
		errs.GetErrorContainer(ctx).Domain = errs.Str2err(response.Error)
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", errs.Err2code(errs.Str2err(response.Error))),
			"x-stc-error", response.Failed().Error(),
		)
		grpc.SendHeader(ctx, header)
	}

	//return converters.LandingResponseToProto(*response), nil
	return response, nil
}
