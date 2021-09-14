package service

import (
	"context"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "deblasis-v1-AuthDBService"
	Namespace   = "stc"
	Tags        = []string{}
)

type AuthDBService interface {
	GetUserByUsername(context.Context, *dtos.GetUserByUsernameRequest) (*dtos.GetUserByUsernameResponse, error)
	CreateUser(context.Context, *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
}

type authDbService struct {
	repository repositories.UserRepository
	logger     log.Logger
	validate   *validator.Validate
}

func NewAuthDBService(repository repositories.UserRepository, logger log.Logger) AuthDBService {
	return &authDbService{
		repository: repository,
		logger:     logger,
		validate:   validator.New(),
	}
}

func (u *authDbService) GetUserByUsername(ctx context.Context, request *dtos.GetUserByUsernameRequest) (*dtos.GetUserByUsernameResponse, error) {
	level.Info(u.logger).Log("handling request", "GetUserByUsername")
	defer level.Info(u.logger).Log("handled request", "GetUserByUsername")
	user, err := u.repository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user ")
	}
	return &dtos.GetUserByUsernameResponse{
		User: user,
	}, nil
}

func (u *authDbService) CreateUser(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	level.Info(u.logger).Log("handling request", "CreateUser")
	defer level.Info(u.logger).Log("handled request", "CreateUser")
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, errors.Wrap(validationErrors, "Failed to create user")
	}

	hashedPassword, err := ca.HashPwd(request.Password)
	if err != nil {
		return nil, err
	}
	request.Password = hashedPassword
	user := model.User(*request)

	id, err := u.repository.CreateUser(ctx, &user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create user")
	}
	return &dtos.CreateUserResponse{
		Id: id,
	}, nil
}
