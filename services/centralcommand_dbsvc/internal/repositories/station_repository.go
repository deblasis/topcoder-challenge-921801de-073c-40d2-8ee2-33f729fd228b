package repositories

import (
	"context"

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

func (u stationRepository) GetById(ctx context.Context, id string) (*model.Station, error) {
	//TODO use validate
	if id == "" {
		level.Debug(u.logger).Log("err", errs.ErrBadRequest)
		return nil, errs.ErrBadRequest
	}

	var ret model.Station
	err := u.Db.WithContext(ctx).Model(&ret).
		Relation("Docks").
		Where("id = ?", id).
		Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("msg", "no rows")
		return nil, nil
	}

	if err != nil {
		return nil, errs.ErrCannotSelectEntity
	}

	return &ret, nil
}

func (u stationRepository) Create(ctx context.Context, station model.Station) (*model.Station, error) {

	err := u.Db.RunInTransaction(ctx, func(t *pg.Tx) error {

		result, err := t.Exec("insert into stations (id, capacity) VALUES (?,?)", station.Id, station.Capacity)
		if err != nil {
			level.Debug(u.logger).Log("err", err)
			return errs.ErrCannotInsertEntity
		}

		if result != nil {
			if result.RowsAffected() == 0 {
				level.Debug(u.logger).Log("err", err)
				return errs.ErrCannotInsertEntity
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
				level.Debug(u.logger).Log("err", err)
				return errs.ErrCannotInsertEntity
			}
		}
		return nil
	})

	return &station, err
}

func (u stationRepository) GetAll(ctx context.Context) ([]model.Station, error) {
	var stations []model.Station

	err := u.Db.WithContext(ctx).
		Model(&stations).
		Relation("Docks").
		Select()
	if err != nil {
		level.Debug(u.logger).Log("err", err)
		return nil, errs.ErrCannotSelectEntities
	}

	return stations, nil
}
