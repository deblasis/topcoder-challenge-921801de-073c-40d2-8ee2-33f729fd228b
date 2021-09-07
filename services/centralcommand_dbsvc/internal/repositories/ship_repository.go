package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		err := errors.New("id is empty")
		level.Debug(u.logger).Log(err)
		return nil, err
	}

	var ret model.Ship
	err := u.Db.WithContext(ctx).Model(&ret).
		Where("id = ?", id).Select()

	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return nil, nil
	}
	return &ret, err
}

func (u shipRepository) Create(ctx context.Context, ship model.Ship) (*model.Ship, error) {
	ship.ID = uuid.New().String()
	result, err := u.Db.WithContext(ctx).Model(&ship).
		ExcludeColumn("status").
		Returning("id").Insert(&ship.ID)
	if err != nil {
		err = errors.Wrapf(err, "Failed to insert ship %v", ship)
		level.Debug(u.logger).Log(err)
		return nil, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errors.New("Failed to insert, affected is 0")
			level.Debug(u.logger).Log(err)
			return nil, err
		}
	}
	return &ship, nil
}

func (u shipRepository) GetAll(ctx context.Context) ([]*model.Ship, error) {
	var ret []*model.Ship

	err := u.Db.WithContext(ctx).Model(&ret).Select()
	if err != nil {
		err = errors.Wrapf(err, "Failed to select ships")
		level.Debug(u.logger).Log(err)
		return nil, err
	}

	return ret, nil
}
