package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
)

type EndpointSet struct {
	StatusEndpoint endpoint.Endpoint

	RequestLandingEndpoint endpoint.Endpoint
	LandingEndpoint        endpoint.Endpoint

	logger log.Logger
}

func NewEndpointSet(s service.ShippingStationService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EndpointSet {

	var requestLandingEndpoint endpoint.Endpoint
	{
		requestLandingEndpoint = MakeRequestLandingEndpoint(s)
		requestLandingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(requestLandingEndpoint)
		requestLandingEndpoint = opentracing.TraceServer(otTracer, "RequestLanding")(requestLandingEndpoint)
		if zipkinTracer != nil {
			requestLandingEndpoint = zipkin.TraceEndpoint(zipkinTracer, "RequestLanding")(requestLandingEndpoint)
		}
		requestLandingEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "RequestLanding"))(requestLandingEndpoint)
		requestLandingEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "RequestLanding"))(requestLandingEndpoint)
	}

	var landingEndpoint endpoint.Endpoint
	{
		landingEndpoint = MakeLandingEndpoint(s)
		landingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(landingEndpoint)
		landingEndpoint = opentracing.TraceServer(otTracer, "Landing")(landingEndpoint)
		if zipkinTracer != nil {
			landingEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Landing")(landingEndpoint)
		}
		landingEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "Landing"))(landingEndpoint)
		landingEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "Landing"))(landingEndpoint)
	}

	return EndpointSet{
		StatusEndpoint: healthcheck.MakeStatusEndpoint(logger, duration, otTracer, zipkinTracer),

		RequestLandingEndpoint: requestLandingEndpoint,
		LandingEndpoint:        landingEndpoint,

		logger: logger,
	}
}

func MakeRequestLandingEndpoint(s service.ShippingStationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.RequestLandingRequest)
		resp, err := s.RequestLanding(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}

func MakeLandingEndpoint(s service.ShippingStationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.LandingRequest)
		resp, err := s.Landing(ctx, req)
		if resp != nil && !errs.IsNil(resp.Failed()) {
			errs.GetErrorContainer(ctx).Domain = resp.Error
		}
		return resp, err
	}
}
