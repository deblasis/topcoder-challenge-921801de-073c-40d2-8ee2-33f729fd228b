package repositories

import (
	"context"

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

func (u userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		level.Debug(u.logger).Log("err", errs.ErrBadRequest)
		return nil, errs.ErrBadRequest
	}

	var user model.User
	err := u.Db.WithContext(ctx).Model(&user).Where("username = ?", username).Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return nil, nil
	}
	return &user, errs.ErrCannotSelectEntity
}

func (u userRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		level.Debug(u.logger).Log("err", errs.ErrBadRequest)
		return nil, errs.ErrBadRequest
	}

	var user model.User
	err := u.Db.WithContext(ctx).Model(&user).Where("id = ?", id).Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return nil, nil
	}
	return &user, errs.ErrCannotSelectEntity
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) (*uuid.UUID, error) {

	if user == nil {
		level.Debug(u.logger).Log("err", errs.ErrBadRequest)
		return nil, errs.ErrBadRequest
	}

	id := uuid.New()
	user.Id = id.String()

	result, err := u.Db.WithContext(ctx).Model(user).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			level.Debug(u.logger).Log("err", pgErr)
			return nil, errs.ErrCannotInsertAlreadyExistingEntity
		} else {
			level.Debug(u.logger).Log("err", err)
			return nil, errs.ErrCannotInsertEntity
		}
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			level.Debug(u.logger).Log("err", err)
			return nil, errs.ErrCannotInsertEntity
		}
	}

	return &id, nil
}
