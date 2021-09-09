package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/config"
	dbe "deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "centralcommandsvc.v1.CentralCommandService"
	Namespace   = "deblasis"
	Tags        = []string{}
)

type CentralCommandService interface {
	ServiceStatus(ctx context.Context) (int64, error)

	RegisterShip(ctx context.Context, request dtos.RegisterShipRequest) (*dtos.RegisterShipResponse, error)
	GetAllShips(ctx context.Context, request dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error)

	RegisterStation(ctx context.Context, request dtos.RegisterStationRequest) (*dtos.RegisterStationResponse, error)
	GetAllStations(ctx context.Context, request dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error)
}

type centralCommandService struct {
	logger             log.Logger
	validate           *validator.Validate
	db_svc_endpointset dbe.EndpointSet
}

func NewCentralCommandService(logger log.Logger, jwtConfig config.JWTConfig, db_svc_endpointset dbe.EndpointSet) CentralCommandService {
	return &centralCommandService{
		logger:             logger,
		validate:           validator.New(),
		db_svc_endpointset: db_svc_endpointset,
	}
}

func (u *centralCommandService) ServiceStatus(ctx context.Context) (int64, error) {
	level.Info(u.logger).Log("handling request", "ServiceStatus")
	defer level.Info(u.logger).Log("handled request", "ServiceStatus")
	return http.StatusOK, nil
}

func (s *centralCommandService) RegisterShip(ctx context.Context, request dtos.RegisterShipRequest) (*dtos.RegisterShipResponse, error) {
	level.Info(s.logger).Log("handling request", "RegisterShip")
	defer level.Info(s.logger).Log("handled request", "RegisterShip")
	//var resp dtos.RegisterShipResponse

	_, err := s.db_svc_endpointset.CreateShipEndpoint(ctx, converters.RegisterShipRequestToCreateShipRequestDBDto(request))
	if err != nil {
		return nil, err
	}

	return &dtos.RegisterShipResponse{}, nil
}

func (s *centralCommandService) RegisterStation(ctx context.Context, request dtos.RegisterStationRequest) (*dtos.RegisterStationResponse, error) {
	// level.Info(s.logger).Log("handling request", "RegisterStation")
	// defer level.Info(s.logger).Log("handled request", "RegisterStation")
	// var resp dtos.RegisterStationResponse

	// _, err := s.db_svc_endpointset.CreateStationEndpoint(ctx, request.ToDBDto())
	// if err != nil {
	// 	return resp, err
	// }

	// return resp, nil
	return &dtos.RegisterStationResponse{}, nil
}

func (u *centralCommandService) GetAllShips(ctx context.Context, request dtos.GetAllShipsRequest) (*dtos.GetAllShipsResponse, error) {
	level.Info(u.logger).Log("handling request", "GetAllShips")
	defer level.Info(u.logger).Log("handled request", "GetAllShips")
	ret, err := u.db_svc_endpointset.GetAllShips(ctx, converters.GetAllShipsRequestToDBDto(request))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve ships ")
	}
	return converters.DBDtoGetAllShipsResponseToDto(*ret), nil
}

func (u *centralCommandService) GetAllStations(ctx context.Context, request dtos.GetAllStationsRequest) (*dtos.GetAllStationsResponse, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllStations")
	defer level.Info(u.logger).Log("handled request", "GetAllStations")
	ret, err := u.db_svc_endpointset.GetAllStations(ctx, converters.GetAllStationsRequestToDBDto(request))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve stations ")
	}

	return converters.DBDtoGetAllStationsResponseToDto(*ret), nil
}

// func (u *userManager) Signup(ctx context.Context, request dtos.SignupRequest) (dtos.SignupResponse, error) {
// 	level.Info(u.logger).Log("handling request", "Signup")
// 	defer level.Info(u.logger).Log("handled request", "Signup")
// 	return http.StatusOK, nil
// }

// func (u *userManager) GetUserByUsername(ctx context.Context, username string) (dtos.User, error) {
// 	user, err := u.repository.GetUserByUsername(ctx, username)
// 	if err != nil {
// 		return dtos.User{}, errors.Wrap(err, "Failed to get user ")
// 	}
// 	return user, nil
// }
