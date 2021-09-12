package seeding

import (
	"context"
	"fmt"

	ca "deblasis.net/space-traffic-control/common/auth"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/db"
	. "deblasis.net/space-traffic-control/common/utils"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/go-kit/kit/log/level"
	"github.com/kong/go-kong/kong"
)

func Seed(cfg config.Config) {
	SeedDB(cfg)
	SeedKong(cfg)
}

func SeedKong(cfg config.Config) {
	ctx := context.Background()
	logger := cfg.Logger
	// prepare pg connection
	pgClient := db.NewPostgresClientFromConfig(cfg)
	dbConn := pgClient.GetConnection()
	defer dbConn.Close()
	kongClient, err := kong.NewClient(String(cfg.Kong.BaseUrl), nil)
	if err != nil {
		level.Debug(logger).Log("kong_seeding_err", err)
	}

	var unsynchedUsers []model.User
	err = dbConn.WithContext(ctx).Model(&unsynchedUsers).Where("kong_id is null").Select()
	level.Error(logger).Log(
		"kong_seeding_err", err,
	)

	for _, user := range unsynchedUsers {

		kc, err := kongClient.Consumers.Create(ctx, &kong.Consumer{
			CustomID: String(fmt.Sprintf("%v", user.Id)),
			Username: String(user.Username),
			Tags:     StringSlice("seeded"),
		})
		if err != nil {
			level.Error(logger).Log(
				"username", user.Username,
				"user_id", user.Id,
				"kong_seeding_err", err,
			)
		}
		kongClient.ACLs.Create(ctx, kc.ID, &kong.ACLGroup{
			Group: &user.Role,
			Tags:  StringSlice("seeded"),
		})

		user.KongId = *kc.ID

		_, err = dbConn.WithContext(ctx).Model(user).Update()
		if err != nil {
			level.Error(logger).Log(
				"username", user.Username,
				"user_id", user.Id,
				"kong_consumer_id", *kc.ID,
				"seeding_err", err,
			)
		}
		level.Info(logger).Log(
			"username", user.Username,
			"user_id", user.Id,
			"kong_consumer_id", *kc.ID,
		)
	}
}

func SeedDB(cfg config.Config) {

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
	kongClient, err := kong.NewClient(String(cfg.Kong.BaseUrl), nil)
	if err != nil {
		level.Debug(logger).Log("kong_seeding_err", err)
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
		kc, err := kongClient.Consumers.Create(ctx, &kong.Consumer{
			CustomID: String(fmt.Sprintf("%v", user.Id)),
			Username: String(u.Username),
			Tags:     StringSlice("seeded"),
		})
		if err != nil {
			level.Error(logger).Log(
				"username", u.Username,
				"user_id", user.Id,
				"kong_seeding_err", err,
			)
		}
		user.KongId = *kc.ID

		_, err = dbConn.Model(user).Update()
		if err != nil {
			level.Error(logger).Log(
				"username", u.Username,
				"user_id", user.Id,
				"kong_consumer_id", *kc.ID,
				"seeding_err", err,
			)
		}
		level.Info(logger).Log(
			"username", u.Username,
			"user_id", user.Id,
			"kong_consumer_id", *kc.ID,
		)

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
