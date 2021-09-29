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
package seeding

import (
	"context"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/go-kit/log/level"
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
