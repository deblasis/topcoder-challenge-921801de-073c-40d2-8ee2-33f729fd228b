package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ServiceName    = "deblasis-state-v1-CentralCommandDBService"
	AuxServiceName = "deblasis-state-v1-CentralCommandDBAuxService"
	Namespace      = "stc"
	Tags           = []string{}
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

func (s *centralCommandDBService) CreateShip(ctx context.Context, request *dtos.CreateShipRequest) (resp *dtos.CreateShipResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "CreateShip", "err", err)
		}
	}()
	verr := s.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.CreateShipResponse{
			Error: err,
		}, nil
	}

	existing, err := s.shipRepository.GetById(ctx, request.Id)
	if existing != nil {
		err = errs.NewError(http.StatusBadRequest, "ship already existing", errs.ErrCannotInsertAlreadyExistingEntity)
		return &dtos.CreateShipResponse{
			Error: err,
		}, nil
	}

	ret, err := s.shipRepository.Create(ctx, model.Ship{
		Id:     request.Id,
		Weight: request.Weight,
	})
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot create ship", err)
		return &dtos.CreateShipResponse{
			Error: err,
		}, nil
	}

	return &dtos.CreateShipResponse{
		Ship: converters.ShipToDto(ret),
	}, nil
}

func (s *centralCommandDBService) GetAllShips(ctx context.Context, request *dtos.GetAllShipsRequest) (resp *dtos.GetAllShipsResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "GetAllShips", "err", err)
		}
	}()
	verr := s.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.GetAllShipsResponse{
			Error: err,
		}, nil
	}

	ret, err := s.shipRepository.GetAll(ctx)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "unable to select ships", err)
		return &dtos.GetAllShipsResponse{
			Error: err,
		}, nil
	}

	return &dtos.GetAllShipsResponse{
		Ships: converters.ShipsToDto(ret),
	}, nil
}

func (u *centralCommandDBService) CreateStation(ctx context.Context, request *dtos.CreateStationRequest) (resp *dtos.CreateStationResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "CreateStation", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.CreateStationResponse{
			Error: err,
		}, nil
	}

	existing, err := u.stationRepository.GetById(ctx, request.Id)
	if existing != nil {
		err = errs.NewError(http.StatusBadRequest, "station already exists", errs.ErrCannotInsertAlreadyExistingEntity)
		return &dtos.CreateStationResponse{
			Error: err,
		}, nil
	}

	ret, err := u.stationRepository.Create(ctx, converters.StationToModel(dtos.Station(*request)))
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot insert station", err)
		return &dtos.CreateStationResponse{
			Error: err,
		}, nil
	}

	return &dtos.CreateStationResponse{
		Station: converters.StationToDto(ret),
	}, nil
}

func (u *centralCommandDBService) GetAllStations(ctx context.Context, request *dtos.GetAllStationsRequest) (resp *dtos.GetAllStationsResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetAllStations", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.GetAllStationsResponse{
			Error: err,
		}, nil
	}

	stations := make([]model.Station, 0)

	if request.ShipId != nil {
		stations, err = u.stationRepository.GetAvailableForShip(ctx, uuid.MustParse(*request.ShipId))
	} else {
		stations, err = u.stationRepository.GetAll(ctx)
	}

	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot select stations", err)
		return &dtos.GetAllStationsResponse{
			Error: err,
		}, nil
	}
	return &dtos.GetAllStationsResponse{
		Stations: converters.StationsToDto(stations),
	}, nil
}

func (u *centralCommandDBService) GetNextAvailableDockingStation(ctx context.Context, request *dtos.GetNextAvailableDockingStationRequest) (resp *dtos.GetNextAvailableDockingStationResponse, err error) {

	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetAllStations", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.GetNextAvailableDockingStationResponse{
			Error: err,
		}, nil
	}

	ret, err := u.dockRepository.GetNextAvailableDockingStation(ctx, uuid.MustParse(request.ShipId))
	if err != nil {
		// err = errs.NewError(http.StatusInternalServerError, "cannot get next available docking station", err)
		return &dtos.GetNextAvailableDockingStationResponse{
			Error: err,
		}, nil
	}
	return &dtos.GetNextAvailableDockingStationResponse{
		NextAvailableDockingStation: converters.NextAvailableDockingStationToDto(ret),
	}, nil
}

func (u *centralCommandDBService) LandShipToDock(ctx context.Context, request *dtos.LandShipToDockRequest) (resp *dtos.LandShipToDockResponse, err error) {

	defer func() {
		if !errs.IsNil(err) {
			level.Debug(u.logger).Log("method", "GetAllStations", "err", err)
		}
	}()
	verr := u.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &dtos.LandShipToDockResponse{
			Error: err,
		}, nil
	}
	_, err = u.dockRepository.LandShipToDock(ctx, uuid.MustParse(request.ShipId), uuid.MustParse(request.DockId), request.Duration)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot land ship to docking station", err)
		return &dtos.LandShipToDockResponse{
			Error: err,
		}, nil
	}
	return &dtos.LandShipToDockResponse{}, nil
}
