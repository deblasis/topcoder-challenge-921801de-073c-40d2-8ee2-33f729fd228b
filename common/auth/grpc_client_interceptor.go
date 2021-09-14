package auth

import (
	"context"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthClientInterceptor struct {
	logger log.Logger

	tokenVerifier TokenVerifier //this is implemented by AuthSvc
	token         string
}

func NewAuthClientInterceptor(logger log.Logger, tokenVerifier TokenVerifier) *AuthClientInterceptor {
	return &AuthClientInterceptor{
		logger:        logger,
		tokenVerifier: tokenVerifier,
	}
}
func (interceptor *AuthClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		interceptor.logger.Log("client_unaryinterceptor", method)
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthClientInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.token)
}
