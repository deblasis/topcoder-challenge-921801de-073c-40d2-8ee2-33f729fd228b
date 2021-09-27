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
package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/cache"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	common_v1 "deblasis.net/space-traffic-control/gen/proto/go/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ServiceName      = "deblasis-v1-AuthService"
	ShortServiceName = "authsvc"
	Namespace        = "deblasis"
	Tags             = []string{}

	GrpcServerPort = 9082 //TODO config
)

type AuthService interface {
	Signup(ctx context.Context, request *pb.SignupRequest) (*pb.SignupResponse, error)
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type authService struct {
	logger             log.Logger
	validate           *validator.Validate
	db_svc_endpointset dbe.EndpointSet
	jwtConfig          config.JWTConfig
	tokensCache        cache.TokensCache
	jwtHandler         *auth.JwtHandler
}

func NewAuthService(logger log.Logger, jwtConfig config.JWTConfig, db_svc_endpointset dbe.EndpointSet, tokensCache cache.TokensCache, jwtHandler *auth.JwtHandler) AuthService {
	return &authService{
		logger:             logger,
		validate:           common.GetValidator(),
		db_svc_endpointset: db_svc_endpointset,
		jwtConfig:          jwtConfig,
		tokensCache:        tokensCache,
		jwtHandler:         jwtHandler,
	}
}

func (s *authService) Signup(ctx context.Context, request *pb.SignupRequest) (resp *pb.SignupResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "Signup",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "Signup", "err", err)

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.SignupResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	level.Debug(s.logger).Log("signup_attempt", request.Username)

	ret, err := s.db_svc_endpointset.CreateUser(ctx, converters.SignupRequestToDBDto(request))
	if err != nil {
		return nil, err
	}

	if !errs.IsNil(ret.Failed()) {
		return &pb.SignupResponse{
			Error: errs.ToProtoV1(ret.Failed()),
		}, nil
	}

	jwtTokenClaims, expiresAt, err := s.jwtHandler.NewJWTToken(*ret.Id, request.Username, request.Role, "http://"+ServiceName) //TODO cfg

	resp = &pb.SignupResponse{
		Token: &pb.Token{
			Token:     jwtTokenClaims.Token,
			ExpiresAt: expiresAt,
		},
		Error: errs.ToProtoV1(err),
	}

	err = s.authUserSession(ctx, jwtTokenClaims.Claims)
	if err != nil {
		resp.Error = errs.ToProtoV1(err)
	}

	return resp, nil
}

func (s *authService) Login(ctx context.Context, request *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "Login", "err", err)
		}
	}()
	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.LoginResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	level.Debug(s.logger).Log("login_attempt", request.Username)

	ret, err := s.db_svc_endpointset.GetUserByUsername(ctx, &dtos.GetUserByUsernameRequest{
		Username: request.Username,
	})
	if err != nil || ret.User == nil {
		return unauthorized("unknown username")
	}

	if !errs.IsNil(ret.Failed()) {
		return &pb.LoginResponse{
			Error: errs.ToProtoV1(ret.Failed()),
		}, nil
	}
	user := ret.User

	bytesHashed := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(bytesHashed, []byte(request.Password+auth.PWDSALT))
	if err != nil {
		return unauthorized("wrong password")
	}

	jwtTokenClaims, expiresAt, err := s.jwtHandler.NewJWTToken(uuid.MustParse(user.Id), user.Username, user.Role, "http://"+ServiceName) //TODO cfg

	resp = &pb.LoginResponse{
		Token: &pb.Token{
			Token:     jwtTokenClaims.Token,
			ExpiresAt: expiresAt,
		}, Error: errs.ToProtoV1(err),
	}

	err = s.authUserSession(ctx, jwtTokenClaims.Claims)
	if err != nil {
		resp.Error = errs.ToProtoV1(err)
	}

	return resp, nil
}

func (s *authService) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (resp *pb.CheckTokenResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "CheckToken", "err", err)
		}
	}()
	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.CheckTokenResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	claims, err := s.jwtHandler.ExtractClaims(request.Token)
	if err != nil {
		return &pb.CheckTokenResponse{
			TokenPayload: nil,
			Error:        errs.ToProtoV1(err),
		}, nil
	}
	authorizedUser, err := s.tokensCache.Get(ctx, claims.Id)
	if err != nil {
		err = errs.NewError(http.StatusUnauthorized, "token not authorized", err)
		return &pb.CheckTokenResponse{
			TokenPayload: nil,
			Error:        errs.ToProtoV1(err),
		}, nil
	}
	user := authorizedUser.(dtos.User)
	return &pb.CheckTokenResponse{
		TokenPayload: &pb.TokenPayload{
			TokenId:  claims.Id,
			UserId:   user.Id,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

func unauthorized(reason string) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Error: &common_v1.Error{
			Code:    http.StatusUnauthorized,
			Message: fmt.Sprintf("unauthorized: %v", reason),
		},
	}, nil
}

func (s *authService) authUserSession(ctx context.Context, claims auth.STCClaims) error {

	expires := time.Unix(claims.ExpiresAt, 0)
	now := time.Now()

	err := s.tokensCache.Set(ctx, claims.Id, dtos.User{
		Id:       claims.UserId,
		Username: claims.Username,
		Role:     claims.Role,
	}, expires.Sub(now))

	if err != nil {
		return err
	}
	return nil

}
