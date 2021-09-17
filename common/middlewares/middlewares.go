package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"deblasis.net/space-traffic-control/common/errs"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			error_container := errs.GetErrorContainer(ctx)

			defer func(begin time.Time) {
				logger.Log(
					"transport_error", error_container.Transport,
					"domain_error", error_container.Domain,
					"took", time.Since(begin),
				)
			}(time.Now())
			return next(ctx, request)

		}
	}
}

func JsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/vnd.deblasis.spacetrafficcontrol-v1+/json; charset=utf-8")
		next.ServeHTTP(rw, r)
	})
}
