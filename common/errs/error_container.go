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
