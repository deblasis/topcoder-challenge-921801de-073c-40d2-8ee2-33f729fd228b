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
