package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "centralcommand_dbsvc.v1.CentralCommandDBService"
	Namespace   = "deblasis"
	Tags        = []string{}
)

type CentralCommandDBService interface {
	ServiceStatus(ctx context.Context) (int64, error)

	CreateShip(ctx context.Context, ship model.Ship) (*model.Ship, error)
	GetAllShips(ctx context.Context) ([]*model.Ship, error)

	CreateStation(ctx context.Context, station model.Station) (*model.Station, error)
	GetAllStations(ctx context.Context) ([]*model.Station, error)
}

type centralCommandDBService struct {
	shipRepository    repositories.ShipRepository
	stationRepository repositories.StationRepository
	logger            log.Logger
	validate          *validator.Validate
}

func NewCentralCommandDBService(
	shipRepository repositories.ShipRepository,
	stationRepository repositories.StationRepository,
	logger log.Logger,
) CentralCommandDBService {
	return &centralCommandDBService{
		shipRepository:    shipRepository,
		stationRepository: stationRepository,
		logger:            logger,
		validate:          validator.New(),
	}
}

func (u *centralCommandDBService) ServiceStatus(ctx context.Context) (int64, error) {
	level.Info(u.logger).Log("handling request", "ServiceStatus")
	defer level.Info(u.logger).Log("handled request", "ServiceStatus")
	return http.StatusOK, nil
}

func (u *centralCommandDBService) CreateShip(ctx context.Context, ship model.Ship) (*model.Ship, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "CreateShip")
	defer level.Info(u.logger).Log("handled request", "CreateShip")
	ret, err := u.shipRepository.Create(ctx, ship)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to create ship ")
	}
	return ret, nil
}

func (u *centralCommandDBService) GetAllShips(ctx context.Context) ([]*model.Ship, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllShips")
	defer level.Info(u.logger).Log("handled request", "GetAllShips")
	ret, err := u.shipRepository.GetAll(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to retrieve ships ")
	}
	return ret, nil
}

func (u *centralCommandDBService) CreateStation(ctx context.Context, station model.Station) (*model.Station, error) {
	level.Info(u.logger).Log("handling request", "CreateStation")
	defer level.Info(u.logger).Log("handled request", "CreateStation")
	err := u.validate.Struct(station)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, errors.Wrap(validationErrors, "Failed to create station")
	}

	ret, err := u.stationRepository.Create(ctx, station)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to create station")
	}
	return ret, nil
}

func (u *centralCommandDBService) GetAllStations(ctx context.Context) ([]*model.Station, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllStations")
	defer level.Info(u.logger).Log("handled request", "GetAllStations")
	ret, err := u.stationRepository.GetAll(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "Failed to retrieve ships ")
	}
	return ret, nil
}
