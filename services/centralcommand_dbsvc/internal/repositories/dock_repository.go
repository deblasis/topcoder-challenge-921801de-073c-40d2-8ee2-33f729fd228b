package repositories

import (
	"context"
	"fmt"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type dockRepository struct {
	Db     *pg.DB
	logger log.Logger
}

func NewDockRepository(db *pg.DB, logger log.Logger) DockRepository {
	return &dockRepository{
		Db:     db,
		logger: logger,
	}
}

func (u dockRepository) GetById(ctx context.Context, id string) (*model.Dock, error) {
	//TODO use validate
	if id == "" {
		err := errors.New("id is empty")
		level.Debug(u.logger).Log("err", err)
		return nil, err
	}

	var ret model.Dock
	err := u.Db.WithContext(ctx).Model(&ret).
		Where("id = ?", id).Select()

	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return nil, nil
	}
	return &ret, err
}

func (u dockRepository) Create(ctx context.Context, dock model.Dock) (*model.Dock, error) {
	dock.Id = uuid.New().String()
	result, err := u.Db.WithContext(ctx).Model(dock).Returning("id").Insert(&dock.Id)
	if err != nil {
		err = errors.Wrapf(err, "Failed to insert dock %v", dock)
		level.Debug(u.logger).Log("err", err)
		return nil, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errors.New("Failed to insert, affected is 0")
			level.Debug(u.logger).Log("err", err)
			return nil, err
		}
	}
	return &dock, nil
}

func (u dockRepository) GetNextAvailableDockingStation(ctx context.Context, shipId uuid.UUID) (*model.NextAvailableDockingStation, error) {

	var nextAvail model.NextAvailableDockingStation
	_, err := u.Db.WithContext(ctx).Model(&nextAvail).
		QueryOne(&nextAvail, fmt.Sprintf("select * from %v(?)", model.GetNextAvailableDockingStationForShipFunctionName), shipId)

	if err != nil {
		level.Debug(u.logger).Log("err", err)
		return nil, errs.ErrCannotSelectEntities
	}

	return &nextAvail, nil
}
