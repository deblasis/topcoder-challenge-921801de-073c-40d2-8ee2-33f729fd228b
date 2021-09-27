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

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-pg/pg/v10"
)

type auxRepository struct {
	Db     *pg.DB
	logger log.Logger
}

func NewAuxRepository(db *pg.DB, logger log.Logger) AuxRepository {
	return &auxRepository{
		Db:     db,
		logger: logger,
	}
}

func (u auxRepository) Cleanup(ctx context.Context) (err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "Cleanup", "err", err)
		}
	}()

	ret, err := u.Db.WithContext(ctx).Exec("truncate table docked_ships, stations, docks, ships; ")
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot cleanup", err)
		return err
	}
	level.Debug(u.logger).Log("method", "Cleanup", "ret", ret.RowsAffected())

	return nil
}
