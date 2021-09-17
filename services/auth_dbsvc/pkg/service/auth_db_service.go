package service

import (
	"context"

	"deblasis.net/space-traffic-control/common"
	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "deblasis-state-v1-AuthDBService"
	Namespace   = "stc"
	Tags        = []string{}
)

type AuthDBService interface {
	GetUserByUsername(context.Context, *dtos.GetUserByUsernameRequest) (*dtos.GetUserResponse, error)
	GetUserById(context.Context, *dtos.GetUserByIdRequest) (*dtos.GetUserResponse, error)
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
		validate:   common.GetValidator(),
	}
}

func (u *authDbService) GetUserByUsername(ctx context.Context, request *dtos.GetUserByUsernameRequest) (*dtos.GetUserResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.GetUserResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	user, err := u.repository.GetUserByUsername(ctx, request.Username)
	return &dtos.GetUserResponse{
		User:  user,
		Error: errs.Err2str(err),
	}, nil
}

func (u *authDbService) GetUserById(ctx context.Context, request *dtos.GetUserByIdRequest) (*dtos.GetUserResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.GetUserResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	user, err := u.repository.GetUserById(ctx, request.Id)
	return &dtos.GetUserResponse{
		User:  user,
		Error: errs.Err2str(err),
	}, nil
}

func (u *authDbService) CreateUser(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.CreateUserResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	existing, err := u.repository.GetUserByUsername(ctx, request.Username)
	if existing != nil {
		return &dtos.CreateUserResponse{
			Error: errs.ErrCannotInsertAlreadyExistingEntity.Error(),
		}, nil
	}

	hashedPassword, err := ca.HashPwd(request.Password)
	if err != nil {
		return nil, err
	}
	request.Password = hashedPassword
	user := model.User(*request)

	id, err := u.repository.CreateUser(ctx, &user)

	return &dtos.CreateUserResponse{
		Id:    id,
		Error: errs.Err2str(err),
	}, nil
}
