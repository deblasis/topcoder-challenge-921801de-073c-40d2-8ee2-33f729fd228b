package transport

import (
	"context"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
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

	var createShipEndpoint endpoint.Endpoint
	{
		createShipEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
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

	var createStationEndpoint endpoint.Endpoint
	{
		createStationEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
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
		StatusEndpoint: statusEndpoint,

		CreateShipEndpoint:  createShipEndpoint,
		GetAllShipsEndpoint: getAllShipsEndpoint,

		CreateStationEndpoint:  createStationEndpoint,
		GetAllStationsEndpoint: getAllStationsEndpoint,
	}
}

func encodeGRPCServiceStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &pb.ServiceStatusRequest{}, nil
}

func decodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.ServiceStatusResponse)
	return healthcheck.ServiceStatusResponse{Code: response.Code, Err: response.Err}, nil
}

func encodeGRPCCreateShipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.CreateShipRequest)

	//TODO: centralise
	// roleId := pb.User_Role_value[strings.ToUpper("ROLE_"+req.Role)]
	// if roleId <= 0 {
	// 	return nil, errors.New("cannot unmarshal role")
	// }
	return &pb.CreateShipRequest{
		Ship: &pb.Ship{
			Weight: float32(req.Weight),
		},
	}, nil
}
func decodeGRPCCreateShipResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateShipResponse)
	return dtos.CreateShipResponse{
		Ship: &model.Ship{
			ID: response.Ship.Id,
			//TODO converter
			Status: response.Ship.Status.String(),
			Weight: float32(response.Ship.Weight),
		},
		Err: response.Error,
	}, nil
}

func encodeGRPCGetAllShipsRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &pb.GetAllShipsRequest{}, nil

}
func decodeGRPCGetAllShipsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllShipsResponse)
	return dtos.GetAllShipsResponse{
		Ships: []*model.Ship{},
		Err:   response.Error,
	}, nil
}

func encodeGRPCCreateStationRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.CreateStationRequest)

	//TODO: centralise
	// roleId := pb.User_Role_value[strings.ToUpper("ROLE_"+req.Role)]
	// if roleId <= 0 {
	// 	return nil, errors.New("cannot unmarshal role")
	// }
	return &pb.CreateStationRequest{
		Station: &pb.Station{
			Capacity: req.Capacity,
			Docks:    modelToProtoDocks(req.Docks),
		},
	}, nil
}
func decodeGRPCCreateStationResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateStationResponse)
	return dtos.CreateStationResponse{
		Station: &model.Station{
			Capacity: response.Station.Capacity,
			Docks:    protoToModelDocks(response.Station.Docks),
		},
		Err: response.Error,
	}, nil
}

func encodeGRPCGetAllStationsRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &pb.GetAllStationsRequest{}, nil

}
func decodeGRPCGetAllStationsResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetAllStationsResponse)
	return dtos.GetAllStationsResponse{
		Stations: []*model.Station{},
		Err:      response.Error,
	}, nil
}