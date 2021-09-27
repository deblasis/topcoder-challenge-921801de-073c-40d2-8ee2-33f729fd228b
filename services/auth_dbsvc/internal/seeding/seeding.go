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
package seeding

import (
	"context"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

func SeedDB(cfg config.Config) {

	ctx := context.Background()
	logger := cfg.Logger
	// prepare pg connection
	pgClient := db.NewPostgresClientFromConfig(cfg)
	dbConn := pgClient.GetConnection()
	defer dbConn.Close()

	//Check if seeding table exists

	mustSeedQuery := `SELECT 'ImustSeedAndDestroyThisTable' where to_regclass('seeding_tmp') is not null`
	res, err := dbConn.WithContext(ctx).Exec(mustSeedQuery)
	if err != nil {
		level.Debug(logger).Log("seeding_err", err)
	}
	mustSeed := res.RowsReturned()
	level.Debug(logger).Log("seeding", mustSeed)
	if mustSeed == 0 {
		//nothing to do, bye
		return
	}

	//seed and destroy
	users := make([]seedModel, 0)
	err = dbConn.WithContext(ctx).Model(&users).Column("role", "username", "password").Select()
	if err != nil {
		level.Error(logger).Log(
			"seeding_err", err,
		)
	}

	for _, u := range users {

		hashedPwd, err := ca.HashPwd(u.Password)
		if err != nil {
			level.Error(logger).Log(
				"username", u.Username,
				"seeding_err", err,
			)
		}

		user := &model.User{
			Id:       uuid.NewString(),
			Username: u.Username,
			Password: hashedPwd,
			Role:     u.Role,
		}
		//insert into users and hash
		_, err = dbConn.Model(user).Insert()
		if err != nil {
			level.Error(logger).Log(
				"username", u.Username,
				"seeding_err", err,
			)
		}

	}
	//burn after reading
	level.Debug(logger).Log("msg", "seeding completed, systems operational")
	dbConn.Exec("drop table if exists seeding_tmp")
}

type seedModel struct {
	tableName struct{} `pg:"seeding_tmp"`

	Role     string `db:"role"`
	Username string `db:"username"`
	Password string `db:"password"`
}
