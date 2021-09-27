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
package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common"
	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
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

func (u *authDbService) GetUserByUsername(ctx context.Context, request *dtos.GetUserByUsernameRequest) (resp *dtos.GetUserResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetUserByUsername", "err", err)
		}
	}()

	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.GetUserResponse{
			Error: err,
		}, nil
	}

	user, err := u.repository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot get user", err)
	}
	return &dtos.GetUserResponse{
		User:  user,
		Error: err,
	}, nil
}

func (u *authDbService) GetUserById(ctx context.Context, request *dtos.GetUserByIdRequest) (resp *dtos.GetUserResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetUserById", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.GetUserResponse{
			Error: err,
		}, nil
	}

	user, err := u.repository.GetUserById(ctx, request.Id)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot get user", err)
	}
	return &dtos.GetUserResponse{
		User:  user,
		Error: err,
	}, nil
}

func (u *authDbService) CreateUser(ctx context.Context, request *dtos.CreateUserRequest) (resp *dtos.CreateUserResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.CreateUserResponse{
			Error: err,
		}, nil
	}

	existing, err := u.repository.GetUserByUsername(ctx, request.Username)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot check user existence", err)
		return &dtos.CreateUserResponse{
			Error: err,
		}, nil
	}
	if existing != nil {
		err = errs.NewError(http.StatusBadRequest, "username already taken", errs.ErrCannotInsertAlreadyExistingEntity)
		return &dtos.CreateUserResponse{
			Error: err,
		}, nil
	}

	hashedPassword, herr := ca.HashPwd(request.Password)
	if herr != nil {
		err = errs.NewError(http.StatusInternalServerError, "unable to hash password", herr)
		return nil, err
	}
	request.Password = hashedPassword
	user := model.User(*request)

	id, err := u.repository.CreateUser(ctx, &user)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot create user", err)
	}
	return &dtos.CreateUserResponse{
		Id:    id,
		Error: err,
	}, nil
}
