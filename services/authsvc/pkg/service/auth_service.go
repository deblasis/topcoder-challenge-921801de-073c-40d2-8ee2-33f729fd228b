package service

import (
	"context"
	"time"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/cache"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/converters"

	//"deblasis.net/space-traffic-control/services/authsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (s *authService) Signup(ctx context.Context, request *pb.SignupRequest) (*pb.SignupResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.SignupResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := s.db_svc_endpointset.CreateUser(ctx, converters.SignupRequestToDBDto(request))
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.SignupResponse{
			Error: ret.Failed().Error(),
		}, nil
	}

	jwtTokenClaims, expiresAt, err := s.jwtHandler.NewJWTToken(*ret.Id, request.Username, request.Role, "http://"+ServiceName) //TODO cfg

	resp := &pb.SignupResponse{
		Token: &pb.Token{
			Token:     jwtTokenClaims.Token,
			ExpiresAt: expiresAt,
		},
		Error: errs.Err2str(err),
	}

	err = s.authUserSession(ctx, jwtTokenClaims.Claims)
	if err != nil {
		resp.Error = err.Error()
	}

	return resp, nil
}

func (s *authService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.LoginResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	getUserResponse, err := s.db_svc_endpointset.GetUserByUsername(ctx, &dtos.GetUserByUsernameRequest{
		Username: request.Username,
	})
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return unauthorized()
	}
	user := getUserResponse.User

	bytesHashed := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(bytesHashed, []byte(request.Password+auth.PWDSALT))
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return unauthorized()
	}

	jwtTokenClaims, expiresAt, err := s.jwtHandler.NewJWTToken(uuid.MustParse(user.Id), user.Username, user.Role, "http://"+ServiceName) //TODO cfg

	resp := &pb.LoginResponse{
		Token: &pb.Token{
			Token:     jwtTokenClaims.Token,
			ExpiresAt: expiresAt,
		}, Error: errs.Err2str(err),
	}

	err = s.authUserSession(ctx, jwtTokenClaims.Claims)
	if err != nil {
		resp.Error = err.Error()
	}

	return resp, err
}

func (s *authService) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.CheckTokenResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	claims, err := s.jwtHandler.ExtractClaims(request.Token)
	if err != nil {
		return &pb.CheckTokenResponse{
			TokenPayload: nil,
			Error:        errs.Err2str(err),
		}, nil
	}
	authorizedUser, err := s.tokensCache.Get(ctx, claims.Id)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return &pb.CheckTokenResponse{
			TokenPayload: nil,
			Error:        errs.ErrUnauthorized.Error(),
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

func unauthorized() (*pb.LoginResponse, error) {
	return &pb.LoginResponse{Error: errs.ErrUnauthorized.Error()}, nil
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
