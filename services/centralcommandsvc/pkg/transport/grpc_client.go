package transport

import (
	"context"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/service"
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

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.CentralCommandService {
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
	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"ServiceStatus",
			encodeGRPCServiceStatusRequest,
			decodeGRPCServiceStatusResponse,
			pb.ServiceStatusResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		statusEndpoint = opentracing.TraceClient(otTracer, "ServiceStatus")(statusEndpoint)
		statusEndpoint = limiter(statusEndpoint)
		statusEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "ServiceStatus",
			Timeout: 30 * time.Second,
		}))(statusEndpoint)
	}

	var registerShipEndpoint endpoint.Endpoint
	{
		registerShipEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"RegisterShip",
			encodeGRPCRegisterShipRequest,
			decodeGRPCRegisterShipResponse,
			pb.RegisterShipResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		registerShipEndpoint = opentracing.TraceClient(otTracer, "RegisterShip")(registerShipEndpoint)
		registerShipEndpoint = limiter(registerShipEndpoint)
		registerShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RegisterShip",
			Timeout: 30 * time.Second,
		}))(registerShipEndpoint)
	}
	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"GetAllShips",
			encodeGRPCGetAllShipsRequest,
			decodeGRPCGetAllShipsResponse,
			pb.GetAllShipsResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		getAllShipsEndpoint = opentracing.TraceClient(otTracer, "GetAllShips")(getAllShipsEndpoint)
		getAllShipsEndpoint = limiter(getAllShipsEndpoint)
		getAllShipsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetAllShips",
			Timeout: 30 * time.Second,
		}))(getAllShipsEndpoint)
	}

	var registerStationEndpoint endpoint.Endpoint
	{
		registerStationEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"RegisterStation",
			encodeGRPCRegisterStationRequest,
			decodeGRPCRegisterStationResponse,
			pb.RegisterStationResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		registerStationEndpoint = opentracing.TraceClient(otTracer, "RegisterStation")(registerStationEndpoint)
		registerStationEndpoint = limiter(registerStationEndpoint)
		registerStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RegisterStation",
			Timeout: 30 * time.Second,
		}))(registerStationEndpoint)
	}
	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"GetAllStations",
			encodeGRPCGetAllStationsRequest,
			decodeGRPCGetAllStationsResponse,
			pb.GetAllStationsResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		getAllStationsEndpoint = opentracing.TraceClient(otTracer, "GetAllStations")(getAllStationsEndpoint)
		getAllStationsEndpoint = limiter(getAllStationsEndpoint)
		getAllStationsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetAllStations",
			Timeout: 30 * time.Second,
		}))(getAllStationsEndpoint)
	}

	return endpoints.EndpointSet{
		StatusEndpoint:       statusEndpoint,
		RegisterShipEndpoint: registerShipEndpoint,
		GetAllShipsEndpoint:  getAllShipsEndpoint,

		RegisterStationEndpoint: registerStationEndpoint,
		GetAllStationsEndpoint:  getAllStationsEndpoint,
	}

}

func encodeGRPCServiceStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &pb.ServiceStatusRequest{}, nil
}

func decodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	reply := grpcResponse.(*pb.ServiceStatusResponse)
	return healthcheck.ServiceStatusResponse{Code: reply.Code, Err: reply.Err}, nil
}

func encodeGRPCRegisterShipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.RegisterShipRequest)
	return converters.RegisterShipRequestToProto(req), nil
}

func decodeGRPCRegisterShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterShipResponse)
	return converters.ProtoRegisterShipResponseToDto(*response), nil
}

func encodeGRPCGetAllShipsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.GetAllShipsRequest)
	return converters.GetAllShipsRequestToProto(req), nil
}

func decodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	return converters.ProtoGetAllShipsResponseToDto(*response), nil
}

//

func encodeGRPCRegisterStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.RegisterStationRequest)
	return converters.RegisterStationRequestToProto(req), nil
}

func decodeGRPCRegisterStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.RegisterStationResponse)
	return converters.ProtoRegisterStationResponseToDto(*response), nil
}

func encodeGRPCGetAllStationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.GetAllStationsRequest)
	return converters.GetAllStationsRequestToProto(req), nil
}

func decodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	return converters.ProtoGetAllStationsResponseToDto(*response), nil
}