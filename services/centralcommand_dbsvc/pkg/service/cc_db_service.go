package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
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
	ServiceStatus(context.Context) (int64, error)

	CreateShip(context.Context, *dtos.CreateShipRequest) (*dtos.CreateShipResponse, error)
	GetAllShips(context.Context, *dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error)

	CreateStation(context.Context, *dtos.CreateStationRequest) (*dtos.CreateStationResponse, error)
	GetAllStations(context.Context, *dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error)
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

func (u *centralCommandDBService) CreateShip(ctx context.Context, request *dtos.CreateShipRequest) (*dtos.CreateShipResponse, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "CreateShip")
	defer level.Info(u.logger).Log("handled request", "CreateShip")
	ret, err := u.shipRepository.Create(ctx, model.Ship{
		Weight: request.Weight,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create ship ")
	}

	return &dtos.CreateShipResponse{
		Ship: converters.ShipToDto(ret),
	}, nil
}

func (u *centralCommandDBService) GetAllShips(ctx context.Context, request *dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllShips")
	defer level.Info(u.logger).Log("handled request", "GetAllShips")
	ret, err := u.shipRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve ships ")
	}

	return &dtos.GetAllShipsResponse{
		Ships: converters.ShipsToDto(ret),
	}, nil
}

func (u *centralCommandDBService) CreateStation(ctx context.Context, request *dtos.CreateStationRequest) (*dtos.CreateStationResponse, error) {
	level.Info(u.logger).Log("handling request", "CreateStation")
	defer level.Info(u.logger).Log("handled request", "CreateStation")
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, errors.Wrap(validationErrors, "Failed to create station")
	}

	ret, err := u.stationRepository.Create(ctx, converters.StationToModel(dtos.Station(*request)))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create station")
	}

	return &dtos.CreateStationResponse{
		Station: converters.StationToDto(ret),
	}, nil
}

func (u *centralCommandDBService) GetAllStations(ctx context.Context, request *dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllStations")
	defer level.Info(u.logger).Log("handled request", "GetAllStations")
	ret, err := u.stationRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve stations ")
	}

	return &dtos.GetAllStationsResponse{
		Stations: converters.StationsToDto(ret),
	}, nil
}
