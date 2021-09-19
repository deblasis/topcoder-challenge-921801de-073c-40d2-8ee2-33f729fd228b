package transport

import (
	"context"
	"strings"
	"time"

	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.ShippingStationService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []grpctransport.ClientOption

	if zipkinTracer != nil {
		// Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing client can be instantiated
		// without an operation name and fed to each Go kit client as ClientOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path.
		//
		// In this example, we demonstrace a global tracing client.
		options = append(options, zipkin.GRPCClientTrace(zipkinTracer))

	}

	var requestLandingEndpoint endpoint.Endpoint
	{
		requestLandingEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"RequestLanding",
			encodeGRPCRequestLandingRequest,
			decodeGRPCRequestLandingResponse,
			pb.RequestLandingResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		requestLandingEndpoint = opentracing.TraceClient(otTracer, "RequestLanding")(requestLandingEndpoint)
		requestLandingEndpoint = limiter(requestLandingEndpoint)
		requestLandingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RequestLanding",
			Timeout: 30 * time.Second,
		}))(requestLandingEndpoint)
	}
	var landingEndpoint endpoint.Endpoint
	{
		landingEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"Landing",
			encodeGRPCLandingRequest,
			decodeGRPCLandingResponse,
			pb.LandingResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		landingEndpoint = opentracing.TraceClient(otTracer, "Landing")(landingEndpoint)
		landingEndpoint = limiter(landingEndpoint)
		landingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Landing",
			Timeout: 30 * time.Second,
		}))(landingEndpoint)
	}

	return endpoints.EndpointSet{
		RequestLandingEndpoint: requestLandingEndpoint,
		LandingEndpoint:        landingEndpoint,
	}

}

func encodeGRPCRequestLandingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RequestLandingRequest)
	//return converters.RequestLandingRequestToProto(req), nil
	return req, nil
}

func decodeGRPCRequestLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RequestLandingResponse)
	//return converters.ProtoRequestLandingResponseToDto(*response), nil
	return response, nil
}

func encodeGRPCLandingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LandingRequest)
	//return converters.LandingRequestToProto(req), nil
	return req, nil
}

func decodeGRPCLandingResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.LandingResponse)
	//return converters.ProtoLandingResponseToDto(*response), nil
	return response, nil
}
