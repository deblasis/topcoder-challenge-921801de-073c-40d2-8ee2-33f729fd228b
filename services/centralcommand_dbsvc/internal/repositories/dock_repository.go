package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
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

func (u dockRepository) GetById(ctx context.Context, id string) (resp *model.Dock, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	if id == "" {
		err := errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return nil, err
	}
	var ret model.Dock
	err = u.Db.WithContext(ctx).Model(&ret).
		Where("id = ?", id).Select()

	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("method", "GetById", "msg", "no rows")
		return nil, nil
	}
	return &ret, nil
}

func (u dockRepository) Create(ctx context.Context, dock model.Dock) (resp *model.Dock, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	dock.Id = uuid.New().String()
	result, err := u.Db.WithContext(ctx).Model(dock).Returning("id").Insert(&dock.Id)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "failed to insert dock", err)
		return nil, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errs.NewError(http.StatusInternalServerError, "failed to insert dock", errs.ErrCannotInsertEntity)
			return nil, err
		}
	}
	return &dock, nil
}

func (u dockRepository) GetNextAvailableDockingStation(ctx context.Context, shipId uuid.UUID) (resp *model.NextAvailableDockingStation, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	var nextAvail model.NextAvailableDockingStation
	_, err = u.Db.WithContext(ctx).Model(&nextAvail).
		QueryOne(&nextAvail, fmt.Sprintf("select * from %v(?)", model.GetNextAvailableDockingStationForShipFunctionName), shipId)

	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot determine next available docking station", err)
		return nil, err
	}

	return &nextAvail, nil
}

func (u dockRepository) LandShipToDock(ctx context.Context, shipId uuid.UUID, dockId uuid.UUID, duration int64) (resp *model.DockedShip, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "CreateUser", "err", err)
		}
	}()
	dockedShip := &model.DockedShip{
		DockId:       dockId.String(),
		ShipId:       shipId.String(),
		DockedSince:  time.Now().UTC(),
		DockDuration: int64(duration),
	}

	result, err := u.Db.WithContext(ctx).Model(dockedShip).Insert()
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "failed to insert docked ship", err)
		return nil, err
	}

	if result != nil {
		if result.RowsAffected() == 0 {
			err = errs.NewError(http.StatusInternalServerError, "failed to insert docked ship", errs.ErrCannotInsertEntity)
			return nil, err
		}
	}

	return dockedShip, nil
}
