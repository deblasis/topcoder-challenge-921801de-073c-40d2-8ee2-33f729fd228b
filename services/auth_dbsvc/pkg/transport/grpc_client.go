package transport

import (
	"context"
	"strings"
	"time"

	pb "deblasis.net/space-traffic-control/gen/proto/go/auth_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.AuthDBService {
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

	var getUserByUsernameEndpoint endpoint.Endpoint
	{
		getUserByUsernameEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"GetUserByUsername",
			encodeGRPCGetUserByUsernameRequest,
			decodeGRPCGetUserByUsernameResponse,
			pb.GetUserByUsernameResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		getUserByUsernameEndpoint = opentracing.TraceClient(otTracer, "GetUserByUsername")(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = limiter(getUserByUsernameEndpoint)
		getUserByUsernameEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetUserByUsername",
			Timeout: 30 * time.Second,
		}))(getUserByUsernameEndpoint)
	}
	var createUserEndpoint endpoint.Endpoint
	{
		createUserEndpoint = grpctransport.NewClient(
			conn,
			strings.Replace(service.ServiceName, "-", ".", -1),
			"CreateUser",
			encodeGRPCCreateUserRequest,
			decodeGRPCCreateUserResponse,
			pb.CreateUserResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		createUserEndpoint = opentracing.TraceClient(otTracer, "CreateUser")(createUserEndpoint)
		createUserEndpoint = limiter(createUserEndpoint)
		createUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateUser",
			Timeout: 30 * time.Second,
		}))(createUserEndpoint)
	}

	return endpoints.EndpointSet{
		GetUserByUsernameEndpoint: getUserByUsernameEndpoint,
		CreateUserEndpoint:        createUserEndpoint,
	}

}

func encodeGRPCGetUserByUsernameRequest(_ context.Context, request interface{}) (interface{}, error) {
	//TODO converters
	req := request.(*dtos.GetUserByUsernameRequest)
	return &pb.GetUserByUsernameRequest{
		Username: req.Username,
	}, nil
}

func decodeGRPCGetUserByUsernameResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.GetUserByUsernameResponse)

	var user *model.User
	if response.User != nil {
		user = &model.User{
			Id:       response.User.Id,
			Username: response.User.Username,
			Password: response.User.Password,
			//Role:     strings.Title(strings.ToLower(strings.TrimLeft(response.User.Role.String(), "ROLE_"))),
			Role: response.User.Role,
		}
	}

	return &dtos.GetUserResponse{
		User:  user,
		Error: response.Error,
	}, nil
}

func encodeGRPCCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*dtos.CreateUserRequest)

	// //TODO: centralise
	// roleId := pb.User_Role_value[strings.ToUpper("ROLE_"+req.Role)]
	// if roleId <= 0 {
	// 	return nil, errors.New("cannot unmarshal role")
	// }
	return &pb.CreateUserRequest{
		User: &pb.User{
			Id:       req.Id,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,
		},
	}, nil
}
func decodeGRPCCreateUserResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*pb.CreateUserResponse)

	id, _ := uuid.Parse(response.Id)

	return &dtos.CreateUserResponse{
		Id:    &id,
		Error: response.Error,
	}, nil
}
