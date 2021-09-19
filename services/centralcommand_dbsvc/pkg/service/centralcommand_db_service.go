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
	"github.com/google/uuid"
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
	GetNextAvailableDockingStation(context.Context, *dtos.GetNextAvailableDockingStationRequest) (*dtos.GetNextAvailableDockingStationResponse, error)
	LandShipToDock(context.Context, *dtos.LandShipToDockRequest) (*dtos.LandShipToDockResponse, error)
}

type centralCommandDBService struct {
	shipRepository    repositories.ShipRepository
	stationRepository repositories.StationRepository
	dockRepository    repositories.DockRepository
	logger            log.Logger
	validate          *validator.Validate
}

func NewCentralCommandDBService(
	shipRepository repositories.ShipRepository,
	stationRepository repositories.StationRepository,
	dockRepository repositories.DockRepository,
	logger log.Logger,
) CentralCommandDBService {
	return &centralCommandDBService{
		shipRepository:    shipRepository,
		stationRepository: stationRepository,
		dockRepository:    dockRepository,
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

func (u *centralCommandDBService) GetNextAvailableDockingStation(ctx context.Context, request *dtos.GetNextAvailableDockingStationRequest) (*dtos.GetNextAvailableDockingStationResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.GetNextAvailableDockingStationResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := u.dockRepository.GetNextAvailableDockingStation(ctx, uuid.MustParse(request.ShipId))
	if err != nil {
		return &dtos.GetNextAvailableDockingStationResponse{
			Error: errs.Err2str(err),
		}, nil
	}
	return &dtos.GetNextAvailableDockingStationResponse{
		NextAvailableDockingStation: converters.NextAvailableDockingStationToDto(ret),
	}, nil
}

func (u *centralCommandDBService) LandShipToDock(ctx context.Context, request *dtos.LandShipToDockRequest) (*dtos.LandShipToDockResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &dtos.LandShipToDockResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}
	_, err = u.dockRepository.LandShipToDock(ctx, uuid.MustParse(request.ShipId), uuid.MustParse(request.DockId), request.Duration)
	if err != nil {
		return &dtos.LandShipToDockResponse{
			Error: errs.Err2str(err),
		}, nil
	}
	return &dtos.LandShipToDockResponse{}, nil
}
