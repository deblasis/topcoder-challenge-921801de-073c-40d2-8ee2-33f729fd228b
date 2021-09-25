package healthcheck

import (
	"context"
	"net/http"
	"time"

	"deblasis.net/space-traffic-control/common/middlewares"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

func makeStatusEndpoint(logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := serviceStatus(ctx, logger)
		if err != nil {
			return ServiceStatusResponse{Code: code, Error: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Error: ""}, nil
	}
}

func MakeStatusEndpoint(logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) endpoint.Endpoint {
	statusEndpoint := makeStatusEndpoint(logger)

	statusEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(statusEndpoint)
	statusEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(statusEndpoint)
	statusEndpoint = opentracing.TraceServer(otTracer, "ServiceStatus")(statusEndpoint)
	if zipkinTracer != nil {
		statusEndpoint = zipkin.TraceEndpoint(zipkinTracer, "ServiceStatus")(statusEndpoint)
	}
	statusEndpoint = middlewares.LoggingMiddleware(log.With(logger, "method", "ServiceStatus"))(statusEndpoint)
	statusEndpoint = middlewares.InstrumentingMiddleware(duration.With("method", "ServiceStatus"))(statusEndpoint)
	return statusEndpoint
}

type ServiceStatusRequest struct{}
type ServiceStatusResponse struct {
	Code  int64  `json:"code"`
	Error string `json:"error,omitempty"`
}

func serviceStatus(ctx context.Context, logger log.Logger) (int64, error) {
	level.Debug(logger).Log("handlingrequest", "ServiceStatus")
	defer level.Debug(logger).Log("handledrequest", "ServiceStatus")
	return http.StatusOK, nil
}

func DecodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req ServiceStatusRequest
	return req, nil
}
