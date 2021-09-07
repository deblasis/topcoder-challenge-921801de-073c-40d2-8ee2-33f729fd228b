package transport

import (
	"context"
	"errors"
	"strings"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/service"
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

func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.AuthService {
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
			decodeGRPCServiceStatusReply,
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

	var signupEndpoint endpoint.Endpoint
	{
		signupEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"Signup",
			encodeGRPCSignupRequest,
			decodeGRPCSignupReply,
			pb.SignupResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		signupEndpoint = opentracing.TraceClient(otTracer, "Signup")(signupEndpoint)
		signupEndpoint = limiter(signupEndpoint)
		signupEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Signup",
			Timeout: 30 * time.Second,
		}))(signupEndpoint)
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			service.ServiceName,
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCLoginReply,
			pb.LoginResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()

		loginEndpoint = opentracing.TraceClient(otTracer, "Login")(loginEndpoint)
		loginEndpoint = limiter(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Login",
			Timeout: 30 * time.Second,
		}))(loginEndpoint)
	}

	return endpoints.EndpointSet{
		StatusEndpoint: statusEndpoint,
		SignupEndpoint: signupEndpoint,
		LoginEndpoint:  loginEndpoint,
	}

}

func encodeGRPCServiceStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	return &pb.ServiceStatusRequest{}, nil
}

func decodeGRPCServiceStatusReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ServiceStatusResponse)
	return healthcheck.ServiceStatusResponse{Code: reply.Code, Err: reply.Err}, nil
}

func encodeGRPCSignupRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.SignupRequest)

	//TODO: centralise
	roleId := pb.SignupRequest_Role_value[strings.ToUpper("ROLE_"+req.Role)]
	if roleId <= 0 {
		return nil, errors.New("cannot unmarshal role")
	}

	return &pb.SignupRequest{
		Username: req.Username,
		Password: req.Password,
		Role:     pb.SignupRequest_Role(roleId),
	}, nil
}

func decodeGRPCSignupReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SignupResponse)
	return dtos.SignupResponse{
		Token: dtos.Token{
			Token:     reply.Token.Token,
			ExpiresAt: reply.Token.ExpiresAt,
		},
		Err: reply.Error,
	}, nil
}

func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(dtos.LoginRequest)
	return &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func decodeGRPCLoginReply(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.LoginResponse)
	return dtos.LoginResponse{
		Token: dtos.Token{
			Token:     reply.Token.Token,
			ExpiresAt: reply.Token.ExpiresAt,
		},
		Err: reply.Error,
	}, nil
}
