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
		if !errs.IsNil(err) {
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
		if !errs.IsNil(err) {
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
		if !errs.IsNil(err) {
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
