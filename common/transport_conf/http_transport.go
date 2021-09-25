package transport_conf

import (
	"deblasis.net/space-traffic-control/common/errs"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

func GetCommonHTTPServerOptions(logger log.Logger) []httptransport.ServerOption {

	return []httptransport.ServerOption{
		httptransport.ServerBefore(errs.ErrorContainerInjectorHTTP),
		httptransport.ServerErrorHandler(transport.ErrorHandlerFunc(
			errs.TranportErrorHandler(log.With(logger, "component", "TranportErrorHandler")),
		)),
		httptransport.ServerErrorEncoder(errs.EncodeErrorHTTP),
	}

}
