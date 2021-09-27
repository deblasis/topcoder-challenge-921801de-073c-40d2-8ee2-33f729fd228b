// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>  
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE. 
//
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
