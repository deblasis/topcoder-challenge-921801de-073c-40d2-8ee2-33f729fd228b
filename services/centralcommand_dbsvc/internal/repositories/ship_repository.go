package repositories

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
)

type shipRepository struct {
	Db     *pg.DB
	logger log.Logger
}

func NewShipRepository(db *pg.DB, logger log.Logger) ShipRepository {
	return &shipRepository{
		Db:     db,
		logger: logger,
	}
}

func (u shipRepository) GetById(ctx context.Context, id string) (resp *model.Ship, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "GetById", "err", err)
		}
	}()
	if id == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return nil, err
	}

	var ret model.Ship
	err = u.Db.WithContext(ctx).Model(&ret).
		Where("id = ?", id).Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("method", "GetById", "msg", "no rows")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ret, err
}

func (u shipRepository) Create(ctx context.Context, ship model.Ship) (resp *model.Ship, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "Create", "err", err)
		}
	}()
	result, err := u.Db.WithContext(ctx).Model(&ship).
		ExcludeColumn("status").
		Returning("id").Insert(&ship.Id)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot insert ship", err)
		return nil, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errs.NewError(http.StatusInternalServerError, "cannot insert ship", errs.ErrCannotInsertEntity)
			return nil, err
		}
	}
	return &ship, nil
}

func (u shipRepository) GetAll(ctx context.Context) (resp []model.Ship, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "GetAll", "err", err)
		}
	}()
	var ret []model.Ship

	err = u.Db.WithContext(ctx).Model(&ret).Select()
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot select ships", err)
		return nil, err
	}

	return ret, nil
}
