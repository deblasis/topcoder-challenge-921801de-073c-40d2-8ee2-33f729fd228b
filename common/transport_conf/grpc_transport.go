package transport_conf

import (
	"deblasis.net/space-traffic-control/common/errs"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func GetCommonGRPCServerOptions(logger log.Logger) []grpctransport.ServerOption {
	return []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerBefore(errs.ErrorContainerInjectorGRPC),
		grpctransport.ServerErrorHandler(transport.ErrorHandlerFunc(errs.ErrorHandlerGRPC)),
	}
}
