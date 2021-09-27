//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
