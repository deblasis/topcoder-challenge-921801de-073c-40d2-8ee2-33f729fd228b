package endpoints

import (
	"context"
	"time"

	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/service"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
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
		requestLandingEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(requestLandingEndpoint)
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
		landingEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(landingEndpoint)
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
		return s.RequestLanding(ctx, req)
	}
}

func MakeLandingEndpoint(s service.ShippingStationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.LandingRequest)
		return s.Landing(ctx, req)
	}
}
