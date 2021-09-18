package transport

import (
	"context"
	"strings"
	"time"

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/service"
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

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.CentralCommandDBService {
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

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateShip",
			encodeGRPCCreateShipRequest,
			decodeGRPCCreateShipResponse,
			pb.CreateShipResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		createShipEndpoint = opentracing.TraceClient(otTracer, "CreateShip")(createShipEndpoint)
		createShipEndpoint = limiter(createShipEndpoint)
		createShipEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateShip",
			Timeout: 30 * time.Second,
		}))(createShipEndpoint)
	}

	var getAllShipsEndpoint endpoint.Endpoint
	{
		getAllShipsEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
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

	var createStationEndpoint endpoint.Endpoint
	{
		createStationEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateStation",
			encodeGRPCCreateStationRequest,
			decodeGRPCCreateStationResponse,
			pb.CreateStationResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		createStationEndpoint = opentracing.TraceClient(otTracer, "CreateStation")(createStationEndpoint)
		createStationEndpoint = limiter(createStationEndpoint)
		createStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateStation",
			Timeout: 30 * time.Second,
		}))(createStationEndpoint)
	}

	var getAllStationsEndpoint endpoint.Endpoint
	{
		getAllStationsEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
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

	var getNextAvailableDockingStationEndpoint endpoint.Endpoint
	{
		getNextAvailableDockingStationEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetNextAvailableDockingStation",
			encodeGRPCGetNextAvailableDockingStationRequest,
			decodeGRPCGetNextAvailableDockingStationResponse,
			pb.GetNextAvailableDockingStationResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		getNextAvailableDockingStationEndpoint = opentracing.TraceClient(otTracer, "GetNextAvailableDockingStation")(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = limiter(getNextAvailableDockingStationEndpoint)
		getNextAvailableDockingStationEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetNextAvailableDockingStation",
			Timeout: 30 * time.Second,
		}))(getNextAvailableDockingStationEndpoint)
	}

	return endpoints.EndpointSet{
		CreateShipEndpoint:  createShipEndpoint,
		GetAllShipsEndpoint: getAllShipsEndpoint,

		CreateStationEndpoint:  createStationEndpoint,
		GetAllStationsEndpoint: getAllStationsEndpoint,

		GetNextAvailableDockingStationEndpoint: getNextAvailableDockingStationEndpoint,
	}
}

func encodeGRPCCreateShipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateShipRequest)
	return converters.CreateShipRequestToProto(req), nil
}
func decodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateShipResponse)
	return converters.ProtoCreateShipResponseToDto(response), nil
}

func encodeGRPCGetAllShipsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetAllShipsRequest)
	return converters.GetAllShipsRequestToProto(req), nil
}
func decodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	return converters.ProtoGetAllShipsResponseToDto(response), nil
}

func encodeGRPCCreateStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateStationRequest)
	return converters.CreateStationRequestToProto(req), nil
}

func decodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateStationResponse)
	return converters.ProtoCreateStationResponseToDto(response), nil
}

func encodeGRPCGetAllStationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetAllStationsRequest)
	return converters.GetAllStationsRequestToProto(req), nil
}
func decodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	return converters.ProtoGetAllStationsResponseToDto(response), nil
}

func encodeGRPCGetNextAvailableDockingStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.GetNextAvailableDockingStationRequest)
	return converters.GetNextAvailableDockingStationRequestToProto(req), nil
}
func decodeGRPCGetNextAvailableDockingStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetNextAvailableDockingStationResponse)
	return converters.ProtoGetNextAvailableDockingStationResponseToDto(response), nil
}
