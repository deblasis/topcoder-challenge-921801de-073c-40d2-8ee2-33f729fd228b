package seeding

import (
	ca "deblasis.net/space-traffic-control/common/auth"
	"github.com/go-kit/kit/log/level"

	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
)

func SeedDB(cfg config.Config) {

	logger := cfg.Logger
	// prepare pg connection
	pgClient := db.NewPostgresClientFromConfig(cfg)
	dbConn := pgClient.GetConnection()
	defer dbConn.Close()

	//Check if seeding table exists

	mustSeedQuery := `SELECT 'ImustSeedAndDestroyThisTable' where to_regclass('seeding_tmp') is not null`
	res, err := dbConn.Exec(mustSeedQuery)
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
	err = dbConn.Model(&users).Column("role", "username", "password").Select()
	if err != nil {
		panic(err)
	}
	for _, u := range users {

		hashedPwd, err := ca.HashPwd(u.Password)
		if err != nil {
			panic(err)
		}
		user := &model.User{
			Username: u.Username,
			Password: hashedPwd,
			Role:     u.Role,
		}
		//insert into users and hash
		_, err = dbConn.Model(user).Insert()
		if err != nil {
			panic(err)
		}

		//burn after reading
		dbConn.Exec("drop table seeding_tmp")
		level.Debug(logger).Log("msg", "seeding completed, systems operational")
	}

}

type seedModel struct {
	tableName struct{} `pg:"seeding_tmp"`

	Role     string `db:"role"`
	Username string `db:"username"`
	Password string `db:"password"`
}
