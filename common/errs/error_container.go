//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
