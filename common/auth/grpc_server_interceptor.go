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
package auth

import (
	"context"
	"fmt"
	"strings"

	"deblasis.net/space-traffic-control/common"
	"github.com/go-kit/kit/log"
	gk "github.com/go-kit/kit/transport/grpc"
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
