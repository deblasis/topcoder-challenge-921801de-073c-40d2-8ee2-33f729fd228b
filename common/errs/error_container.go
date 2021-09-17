package errs

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type errorContainer struct {
	Transport error
	Domain    error
}

type errorContainerKeyType struct{}

var errorContainerKey errorContainerKeyType

func GetErrorContainer(ctx context.Context) *errorContainer {
	v := ctx.Value(errorContainerKey)
	if v == nil {
		panic("no error container set")
	}
	c, ok := v.(*errorContainer)
	if !ok {
		panic("invalid error container")
	}
	return c
}

func ErrorContainerInjectorHTTP(ctx context.Context, r *http.Request) context.Context {
	return injectErrorContainer(ctx)
}

func ErrorContainerInjectorGRPC(ctx context.Context, m metadata.MD) context.Context {
	return injectErrorContainer(ctx)
}

func injectErrorContainer(ctx context.Context) context.Context {
	return context.WithValue(ctx, errorContainerKey, &errorContainer{}) // new empty errorContainer with each request
}
