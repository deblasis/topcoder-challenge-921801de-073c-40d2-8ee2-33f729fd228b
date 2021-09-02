package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/repositories"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type UserManager interface {
	ServiceStatus(ctx context.Context) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
}

type userManager struct {
	repository repositories.UserRepository
	logger     log.Logger
	validate   *validator.Validate
}

func NewUserManager(repository repositories.UserRepository, logger log.Logger) UserManager {
	return &userManager{
		repository: repository,
		logger:     logger,
		validate:   validator.New(),
	}
}

func (u *userManager) ServiceStatus(ctx context.Context) (int64, error) {
	level.Info(u.logger).Log("handling request", "ServiceStatus")
	defer level.Info(u.logger).Log("handled request", "ServiceStatus")
	return http.StatusOK, nil
}

func (u *userManager) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user, err := u.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return model.User{}, errors.Wrap(err, "Failed to get user ")
	}
	return user, nil
}

func (u *userManager) CreateUser(ctx context.Context, user *model.User) (int64, error) {

	err := u.validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return -1, errors.Wrap(validationErrors, "Failed to create user")
	}

	id, err := u.repository.CreateUser(ctx, user)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to create user")
	}
	return id, nil
}
