package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
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

func (u userRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	if username == "" {
		err := errors.New("Username is empty")
		level.Debug(u.logger).Log(err)
		return model.User{}, err
	}

	var user model.User
	err := u.Db.WithContext(ctx).Model(&user).Where("username = ?", username).Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return model.User{}, nil
	}
	return user, err
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	if user == nil {
		err := errors.New("Input parameter user is nil")
		level.Debug(u.logger).Log(err)
		return -1, err
	}

	result, err := u.Db.WithContext(ctx).Model(user).Returning("id").Insert(&user.Id)
	if err != nil {
		err = errors.Wrapf(err, "Failed to insert user %v", user)
		level.Debug(u.logger).Log(err)
		return -1, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errors.New("Failed to insert, affected is 0")
			level.Debug(u.logger).Log(err)
			return -1, err
		}
	}
	return user.Id, nil
}
