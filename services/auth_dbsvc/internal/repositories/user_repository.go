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
package repositories

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type userRepository struct {
	Db     *pg.DB
	logger log.Logger
}

func NewUserRepository(db *pg.DB, logger log.Logger) UserRepository {
	return &userRepository{
		Db:     db,
		logger: logger,
	}
}

func (u userRepository) GetUserByUsername(ctx context.Context, username string) (resp *model.User, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetUserByUsername", "err", err)
		}
	}()
	if username == "" {
		err = errs.NewError(http.StatusBadRequest, "empty username", errs.ErrBadRequest)
		return nil, err
	}

	var user model.User
	err = u.Db.WithContext(ctx).Model(&user).Where("username = ?", username).Select()
	level.Debug(u.logger).Log("method", "GetUserByUsername", "msg", "no rows")
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}

func (u userRepository) GetUserById(ctx context.Context, id string) (resp *model.User, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetUserById", "err", err)
		}
	}()
	if id == "" {
		err = errs.NewError(http.StatusBadRequest, "empty id", errs.ErrBadRequest)
		return nil, err
	}

	var user model.User
	err = u.Db.WithContext(ctx).Model(&user).Where("id = ?", id).Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("method", "GetUserById", "msg", "no rows")
		return nil, nil
	}
	return &user, nil
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) (resp *uuid.UUID, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	if user == nil {
		err = errs.NewError(http.StatusBadRequest, "empty user", errs.ErrBadRequest)
		return nil, err
	}

	id := uuid.New()
	user.Id = id.String()

	result, err := u.Db.WithContext(ctx).Model(user).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			err = errs.NewError(http.StatusBadRequest, "user already exists", pgErr)
			return nil, err
		} else {
			err = errs.NewError(http.StatusInternalServerError, "cannot insert user", err)
			return nil, err
		}
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errs.NewError(http.StatusInternalServerError, "cannot insert user", errs.ErrCannotInsertEntity)
			return nil, err
		}
	}

	return &id, nil
}
