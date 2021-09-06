package service

import (
	"context"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/errors"
	dbdtos "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/dtos"
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
	Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error)
	Login(ctx context.Context, request dtos.LoginRequest) (dtos.LoginResponse, error)
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

func (s *authService) Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error) {
	level.Info(s.logger).Log("handling request", "Signup")
	defer level.Info(s.logger).Log("handled request", "Signup")
	var resp dtos.SignupResponse

	_, err := s.db_svc_endpointset.CreateUserEndpoint(ctx, dbdtos.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
		Role:     request.Role,
	})
	if err != nil {
		return resp, err
	}

	token, expiresAt, err := ca.NewJWTToken(s.jwtConfig, request.Username, request.Role, ServiceName)

	resp = dtos.SignupResponse{
		Token: dtos.Token{
			Token:     token,
			ExpiresAt: expiresAt,
		},
		Err: errors.Err2str(err),
	}

	return resp, nil
}

func (s *authService) Login(ctx context.Context, request dtos.LoginRequest) (dtos.LoginResponse, error) {
	level.Info(s.logger).Log("handling request", "Login")
	defer level.Info(s.logger).Log("handled request", "Login")
	var resp dtos.LoginResponse

	user, err := s.db_svc_endpointset.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return unauthorized(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password+ca.PWDSALT), bcrypt.DefaultCost+1)
	if err != nil {
		return resp, err
	}

	if user.Password != string(hashedPassword) {
		return unauthorized(nil)
	}

	token, expiresAt, err := ca.NewJWTToken(s.jwtConfig, user.Username, user.Role, ServiceName)

	resp = dtos.LoginResponse{
		Token: dtos.Token{
			Token:     token,
			ExpiresAt: expiresAt,
		}, Err: errors.Err2str(err),
	}

	return resp, err
}

func unauthorized(err error) (dtos.LoginResponse, error) {
	//TODO refactor
	return dtos.LoginResponse{
		Err: "Unauthorized",
	}, err
}

// func (u *userManager) Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error) {
// 	level.Info(u.logger).Log("handling request", "Signup")
// 	defer level.Info(u.logger).Log("handled request", "Signup")
// 	return http.StatusOK, nil
// }

// func (u *userManager) GetUserByUsername(ctx context.Context, username string) (dtos.User, error) {
// 	user, err := u.repository.GetUserByUsername(ctx, username)
// 	if err != nil {
// 		return dtos.User{}, errors.Wrap(err, "Failed to get user ")
// 	}
// 	return user, nil
// }
