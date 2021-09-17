package transport_conf

import (
	"deblasis.net/space-traffic-control/common/errs"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func GetCommonHTTPServerOptions(logger log.Logger) []httptransport.ServerOption {

	return []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerBefore(errs.ErrorContainerInjectorHTTP),
		httptransport.ServerErrorEncoder(errs.EncodeErrorHTTP),
	}

}
