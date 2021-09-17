package service

import (
	"context"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "deblasis-state-v1-CentralCommandDBService"
	Namespace   = "stc"
	Tags        = []string{}
)

type CentralCommandDBService interface {
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
		validate:          common.GetValidator(),
	}
}

func (u *centralCommandDBService) CreateShip(ctx context.Context, request *dtos.CreateShipRequest) (*dtos.CreateShipResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.CreateShipResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	existing, err := u.shipRepository.GetById(ctx, request.Id)
	if existing != nil {
		return &dtos.CreateShipResponse{
			Error: errs.ErrCannotInsertAlreadyExistingEntity.Error(),
		}, nil
	}

	ret, err := u.shipRepository.Create(ctx, model.Ship{
		Id:     request.Id,
		Weight: request.Weight,
	})
	if err != nil {
		return &dtos.CreateShipResponse{
			Error: errs.Err2str(err),
		}, nil
	}

	return &dtos.CreateShipResponse{
		Ship: converters.ShipToDto(ret),
	}, nil
}

func (u *centralCommandDBService) GetAllShips(ctx context.Context, request *dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.GetAllShipsResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := u.shipRepository.GetAll(ctx)
	if err != nil {
		return &dtos.GetAllShipsResponse{
			Error: errs.Err2str(err),
		}, nil
	}

	return &dtos.GetAllShipsResponse{
		Ships: converters.ShipsToDto(ret),
	}, nil
}

func (u *centralCommandDBService) CreateStation(ctx context.Context, request *dtos.CreateStationRequest) (*dtos.CreateStationResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.CreateStationResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	existing, err := u.stationRepository.GetById(ctx, request.Id)
	if existing != nil {
		return &dtos.CreateStationResponse{
			Error: errs.ErrCannotInsertAlreadyExistingEntity.Error(),
		}, nil
	}

	ret, err := u.stationRepository.Create(ctx, converters.StationToModel(dtos.Station(*request)))
	if err != nil {
		return &dtos.CreateStationResponse{
			Error: errs.Err2str(err),
		}, nil
	}

	return &dtos.CreateStationResponse{
		Station: converters.StationToDto(ret),
	}, nil
}

func (u *centralCommandDBService) GetAllStations(ctx context.Context, request *dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.GetAllStationsResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := u.stationRepository.GetAll(ctx)
	if err != nil {
		return &dtos.GetAllStationsResponse{
			Error: errs.Err2str(err),
		}, nil
	}
	return &dtos.GetAllStationsResponse{
		Stations: converters.StationsToDto(ret),
	}, nil
}
