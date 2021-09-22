package repositories

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type stationRepository struct {
	Db     *pg.DB
	logger log.Logger
}

func NewStationRepository(db *pg.DB, logger log.Logger) StationRepository {
	return &stationRepository{
		Db:     db,
		logger: logger,
	}
}

func (u stationRepository) GetById(ctx context.Context, id string) (resp *model.Station, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "GetById", "err", err)
		}
	}()

	if id == "" {
		err = errs.NewError(http.StatusBadRequest, "empty id", errs.ErrBadRequest)
		return nil, err
	}

	var ret model.Station
	err = u.Db.WithContext(ctx).Model(&ret).
		Relation("Docks").
		Where("id = ?", id).
		Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("method", "GetById", "msg", "no rows")
		return nil, nil
	}

	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot get station", errs.ErrCannotSelectEntity)
		return nil, err
	}

	return &ret, nil
}

func (u stationRepository) Create(ctx context.Context, station model.Station) (resp *model.Station, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "Create", "err", err)
		}
	}()
	err = u.Db.RunInTransaction(ctx, func(t *pg.Tx) error {

		result, err := t.Exec("insert into stations (id, capacity) VALUES (?,?)", station.Id, station.Capacity)
		if err != nil {
			err = errs.NewError(http.StatusInternalServerError, "cannot insert station", err)
			return err
		}

		if result != nil {
			if result.RowsAffected() == 0 {
				err = errs.NewError(http.StatusInternalServerError, "cannot insert station", errs.ErrCannotInsertEntity)
				return err
			}
		}

		//insert docks
		for _, dock := range station.Docks {
			dock.Id = uuid.NewString()
			dock.StationId = station.Id

			_, err = t.Model(dock).
				ExcludeColumn("occupied", "weight").
				Returning("id").Insert(&dock.Id)
			if err != nil {
				err = errs.NewError(http.StatusInternalServerError, "cannot insert dock", err)
				return err
			}
		}
		return nil
	})

	return &station, err
}

func (u stationRepository) GetAll(ctx context.Context) (resp []model.Station, err error) {
	defer func() {
		if err != nil {
			level.Debug(u.logger).Log("method", "GetAll", "err", err)
		}
	}()
	var stations []model.Station

	err = u.Db.WithContext(ctx).
		Model(&stations).
		Relation("Docks").
		Select()
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot select docks", err)
		return nil, err
	}

	return stations, nil
}
