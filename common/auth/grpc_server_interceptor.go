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
	"deblasis.net/space-traffic-control/common/errs"
	gk "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServerInterceptor struct {
	logger log.Logger

	jwtHandler *JwtHandler
	aclRules   map[string]ACLRule
}

func NewAuthServerInterceptor(logger log.Logger, jwtHandler *JwtHandler, aclRules map[string]ACLRule) *AuthServerInterceptor {
	return &AuthServerInterceptor{
		logger: logger,

		jwtHandler: jwtHandler,
		aclRules:   aclRules,
	}
}

func (interceptor *AuthServerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		//heathcheck, it's ok anyway, let's not flood logs and stuff
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}
		interceptor.logger.Log("server_unaryinterceptor", info.FullMethod)

		ctx, err := interceptor.checkAuth(ctx, info.FullMethod, req, log.With(interceptor.logger, "component", "checkAuth"))
		if err != nil {
			errs.InjectGrpcErrorStatusCode(ctx, err, http.StatusUnauthorized)
			return nil, err
		}

		//this is necessary for Go-Kit go work, it provides execution metadata to the framework
		ctx = context.WithValue(ctx, gk.ContextKeyRequestMethod, info.FullMethod)

		return handler(ctx, req)
	}
}

func (interceptor *AuthServerInterceptor) checkAuth(ctx context.Context, method string, req interface{}, logger log.Logger) (context.Context, error) {

	reqType := fmt.Sprintf("%T", req)
	aclRule, ok := interceptor.aclRules[method]
	level.Debug(logger).Log("method", method, "reqType", reqType)

	if !ok {
		level.Debug(logger).Log("msg", "everyone can access")
		// everyone can access
		return ctx, nil
	}

	acl := aclRule(req, log.With(interceptor.logger, "component", "aclRule"))
	if acl.allGood {
		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	//I don't want to have the token in logs, just flagging its presence is fine for debugging purposes
	level.Debug(logger).Log("accessToken_len", len(accessToken))

	if len(strings.TrimLeft(accessToken, "Bearer ")) < len(accessToken) {
		accessToken = strings.Split(accessToken, "Bearer ")[1]
	}

	token, err := interceptor.jwtHandler.VerifyToken(accessToken)
	if err != nil {
		return ctx, status.Errorf(acl.statusCodeOnFailure, "access token is invalid: %v", err)
	}

	aclRuleErr := acl.tokenValidator(token)
	if aclRuleErr == nil {
		level.Debug(logger).Log("aclRuleErr", "SUCCESS")
		claims, ok := token.Claims.(*STCClaims)
		if ok {
			level.Debug(logger).Log("msg", "setting creds into context")
			ctx = context.WithValue(ctx, common.ContextKeyUserId, claims.UserId)
			ctx = context.WithValue(ctx, common.ContextKeyUserRole, claims.Role)
		}
		return ctx, nil
	}
	return ctx, status.Error(acl.statusCodeOnFailure, "you are not allowed to do that")
}

type ACLRule func(req interface{}, logger log.Logger) ACLDescriptor

type ACLDescriptor struct {
	allGood             bool
	tokenValidator      TokenValidator
	statusCodeOnFailure codes.Code
}

func NewAllGoodAclDescriptor() ACLDescriptor {
	return ACLDescriptor{allGood: true}
}
func NewMustCheckTokenDescriptor(tokenValidator TokenValidator) ACLDescriptor {
	return ACLDescriptor{allGood: false, tokenValidator: tokenValidator, statusCodeOnFailure: codes.Unauthenticated}
}

func NewMustCheckTokenDescriptorWithCustomStatusCode(tokenValidator TokenValidator, statusCodeOnFailure codes.Code) ACLDescriptor {
	return ACLDescriptor{allGood: false, tokenValidator: tokenValidator, statusCodeOnFailure: statusCodeOnFailure}
}
