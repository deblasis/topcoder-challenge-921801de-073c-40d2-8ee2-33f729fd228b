package repositories

import (
	"context"

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

func (u shipRepository) GetById(ctx context.Context, id string) (*model.Ship, error) {
	//TODO use validate
	if id == "" {
		level.Debug(u.logger).Log("err", errs.ErrBadRequest)
		return nil, errs.ErrBadRequest
	}

	var ret model.Ship
	err := u.Db.WithContext(ctx).Model(&ret).
		Where("id = ?", id).Select()

	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("msg", "no rows")
		return nil, nil
	}
	return &ret, errs.ErrCannotSelectEntity
}

func (u shipRepository) Create(ctx context.Context, ship model.Ship) (*model.Ship, error) {
	result, err := u.Db.WithContext(ctx).Model(&ship).
		ExcludeColumn("status").
		Returning("id").Insert(&ship.Id)
	if err != nil {
		level.Debug(u.logger).Log("err", err)
		return nil, errs.ErrCannotInsertEntity
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			level.Debug(u.logger).Log("err", err)
			return nil, errs.ErrCannotInsertEntity
		}
	}
	return &ship, nil
}

func (u shipRepository) GetAll(ctx context.Context) ([]model.Ship, error) {
	var ret []model.Ship

	err := u.Db.WithContext(ctx).Model(&ret).Select()
	if err != nil {
		level.Debug(u.logger).Log("err", err)
		return nil, errs.ErrCannotSelectEntities
	}

	return ret, nil
}
