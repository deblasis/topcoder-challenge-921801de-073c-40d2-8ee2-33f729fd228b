package service

import (
	"context"
	"fmt"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/errors"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	ServiceName = "authsvc.v1.AuthService"
	Namespace   = "deblasis"
	Tags        = []string{}
)

type AuthService interface {
	Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error)
	Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error)
}

type authService struct {
	logger             log.Logger
	validate           *validator.Validate
	db_svc_endpointset dbe.EndpointSet
	jwtConfig          config.JWTConfig
}

func NewAuthService(logger log.Logger, jwtConfig config.JWTConfig, db_svc_endpointset dbe.EndpointSet) AuthService {
	return &authService{
		logger:             logger,
		validate:           validator.New(),
		db_svc_endpointset: db_svc_endpointset,
		jwtConfig:          jwtConfig,
	}
}

func (s *authService) Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error) {
	level.Info(s.logger).Log("handling request", "Signup")
	defer level.Info(s.logger).Log("handled request", "Signup")
	var resp pb.SignupResponse

	req := dtos.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
		Role:     converters.ProtoToDTORole(request.Role),
	}
	_, err := s.db_svc_endpointset.CreateUserEndpoint(ctx, req)
	if err != nil {
		return resp, err
	}

	token, expiresAt, err := ca.NewJWTToken(s.jwtConfig, request.Username, request.Role.String(), ServiceName)

	resp = pb.SignupResponse{
		Token: &pb.Token{
			Token:     token,
			ExpiresAt: expiresAt,
		},
		Error: errors.Err2str(err),
	}

	return resp, nil
}

func (s *authService) Login(ctx context.Context, request pb.LoginRequest) (pb.LoginResponse, error) {
	level.Info(s.logger).Log("handling request", "Login")
	defer level.Info(s.logger).Log("handled request", "Login")
	var resp pb.LoginResponse

	user, err := s.db_svc_endpointset.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return unauthorized(err)
	}

	bytesHashed := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(bytesHashed, []byte(request.Password+ca.PWDSALT))
	if err != nil {
		return unauthorized(nil)
	}

	token, expiresAt, err := ca.NewJWTToken(s.jwtConfig, user.Username, user.Role, ServiceName)

	resp = pb.LoginResponse{
		Token: &pb.Token{
			Token:     token,
			ExpiresAt: expiresAt,
		}, Error: errors.Err2str(err),
	}

	return resp, err
}

func unauthorized(err error) (pb.LoginResponse, error) {
	//TODO refactor
	return pb.LoginResponse{
		Error: fmt.Sprintf("Unauthorized: %v", err),
	}, nil
}
