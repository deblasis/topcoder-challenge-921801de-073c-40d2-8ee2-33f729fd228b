// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
package repositories

import (
	"context"
	"fmt"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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
		if !errs.IsNil(err) {
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
		if !errs.IsNil(err) {
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
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetAll", "err", err)
		}
	}()
	var stations []model.Station

	err = u.Db.WithContext(ctx).
		Model(&stations).
		Relation("Docks").
		Select()
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot select stations", err)
		return nil, err
	}

	return stations, nil
}

func (u stationRepository) GetAvailableForShip(ctx context.Context, shipId uuid.UUID) (resp []model.Station, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetAvailableForShip", "err", err)
		}
	}()

	stations := make([]model.Station, 0)
	stationMap := make(map[string]*model.Station)

	var availableStations []model.AvailableStationsForShip
	_, err = u.Db.WithContext(ctx).Model(&availableStations).
		Query(&availableStations, fmt.Sprintf("select * from %v(?)", model.GetAvailableStationsForShipFunctionName), shipId)

	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot select stations", err)
		return nil, err
	}

	for _, a := range availableStations {
		if _, ok := stationMap[a.StationId]; !ok {
			stationMap[a.StationId] = &model.Station{
				Id:           a.StationId,
				Capacity:     a.Capacity,
				UsedCapacity: a.UsedCapacity,
				Docks:        []*model.Dock{},
			}
		}
		stationMap[a.StationId].Docks = append(stationMap[a.StationId].Docks, &model.Dock{
			Id:              a.DockId,
			StationId:       a.StationId,
			NumDockingPorts: a.NumDockingPorts,
			Occupied:        a.Occupied,
			Weight:          a.Weight,
		})
	}

	for _, s := range stationMap {
		stations = append(stations, *s)
	}

	return stations, nil
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
