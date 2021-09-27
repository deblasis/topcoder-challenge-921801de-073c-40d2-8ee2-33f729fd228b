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
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"deblasis.net/space-traffic-control/common"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/status"
)

type HttpAuthProvider struct {
	logger     log.Logger
	jwtHandler *JwtHandler
}

func NewHttpAuthProvider(logger log.Logger, jwtHandler *JwtHandler) *HttpAuthProvider {
	return &HttpAuthProvider{
		logger:     logger,
		jwtHandler: jwtHandler,
	}
}

func (a *HttpAuthProvider) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := a.checkAuth(r.Context(), r.RequestURI, w, r, log.With(a.logger, "component", "checkAuth"))
		if err != nil {
			level.Error(a.logger).Log("method", "Handler", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (interceptor *HttpAuthProvider) checkAuth(ctx context.Context, method string, w http.ResponseWriter, req *http.Request, logger log.Logger) (context.Context, error) {

	reqType := fmt.Sprintf("%T", req)

	level.Debug(logger).Log("method", method, "reqType", reqType)

	authHeader := req.Header.Get("Authorization")
	if req.Header.Get("Authorization") == "" {
		return ctx, nil
	}

	ah := strings.Split(authHeader, "Bearer ")
	if len(ah) != 2 {
		level.Error(logger).Log("method", method, "msg", "Malformed token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Malformed Token"))
	} else {
		jwtToken := ah[1]
		token, err := interceptor.jwtHandler.VerifyToken(jwtToken)
		if err != nil {
			return ctx, status.Errorf(http.StatusUnauthorized, "access token is invalid: %v", err)
		}
		claims, ok := token.Claims.(*STCClaims)
		if ok {
			level.Debug(logger).Log("msg", "setting creds into context")
			ctx = context.WithValue(ctx, common.ContextKeyUserId, claims.UserId)
			ctx = context.WithValue(ctx, common.ContextKeyUserRole, claims.Role)
		}
	}
	return ctx, nil
}
