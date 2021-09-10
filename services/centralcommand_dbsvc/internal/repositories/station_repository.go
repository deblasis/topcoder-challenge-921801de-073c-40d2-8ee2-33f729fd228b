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
		err := errors.New("id is empty")
		level.Debug(u.logger).Log(err)
		return nil, err
	}

	var ret model.Station
	err := u.Db.WithContext(ctx).Model(&ret).
		Relation("Docks").
		Where("id = ?", id).
		Select()
	if err == pg.ErrNoRows {
		level.Debug(u.logger).Log("no rows")
		return nil, nil
	}

	//u.Db.WithContext(ctx).Model(&ret.Docks).Where("station_id = ?", id).Select()

	return &ret, err
}

func (u stationRepository) Create(ctx context.Context, station model.Station) (*model.Station, error) {

	err := u.Db.RunInTransaction(ctx, func(t *pg.Tx) error {
		station.Id = uuid.NewString()
		result, err := t.Exec("insert into stations (id, capacity) VALUES (?,?)", station.Id, station.Capacity)
		if err != nil {
			err = errors.Wrapf(err, "Failed to insert station %v", station)
			level.Debug(u.logger).Log(err)
			return err
		}

		if result != nil {
			if result.RowsAffected() == 0 {
				err = errors.New("Failed to insert, affected is 0")
				level.Debug(u.logger).Log(err)
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
				err = errors.Wrapf(err, "Failed to insert dock %v", station)
				level.Debug(u.logger).Log(err)
				return err
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
		err = errors.Wrapf(err, "Failed to select stations")
		level.Debug(u.logger).Log(err)
		return nil, err
	}

	return stations, nil
}
