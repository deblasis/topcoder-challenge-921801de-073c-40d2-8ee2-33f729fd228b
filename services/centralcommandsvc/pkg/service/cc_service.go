package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/config"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
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

	RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error)
	GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error)

	RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error)
	GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error)
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

func (s *centralCommandService) RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error) {
	level.Info(s.logger).Log("handling request", "RegisterShip")
	defer level.Info(s.logger).Log("handled request", "RegisterShip")

	req := converters.RegisterShipRequestToCreateShipRequestDBDto(request)
	ret, err := s.db_svc_endpointset.CreateShip(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to register ship ")
	}
	return converters.DBDtoCreateShipResponseToProto(*ret), nil
}

func (s *centralCommandService) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
	level.Info(s.logger).Log("handling request", "RegisterStation")
	defer level.Info(s.logger).Log("handled request", "RegisterStation")

	ret, err := s.db_svc_endpointset.CreateStation(ctx, converters.RegisterStationRequestToCreateStationRequestDBDto(request))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to register station ")
	}
	return converters.DBDtoCreateStationResponseToProto(*ret), nil
}

func (u *centralCommandService) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	level.Info(u.logger).Log("handling request", "GetAllShips")
	defer level.Info(u.logger).Log("handled request", "GetAllShips")
	ret, err := u.db_svc_endpointset.GetAllShips(ctx,
		&dtos.GetAllShipsRequest{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve ships ")
	}
	return converters.DBDtoGetAllShipsResponseToProto(*ret), nil
}

func (u *centralCommandService) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	//TODO use middleware
	level.Info(u.logger).Log("handling request", "GetAllStations")
	defer level.Info(u.logger).Log("handled request", "GetAllStations")
	ret, err := u.db_svc_endpointset.GetAllStations(ctx,
		&dtos.GetAllStationsRequest{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve stations ")
	}

	return converters.DBDtoGetAllStationsResponseToProto(*ret), nil
}

// func (u *userManager) Signup(ctx context.Context, request pb.SignupRequest) (pb.SignupResponse, error) {
// 	level.Info(u.logger).Log("handling request", "Signup")
// 	defer level.Info(u.logger).Log("handled request", "Signup")
// 	return http.StatusOK, nil
// }

// func (u *userManager) GetUserByUsername(ctx context.Context, username string) (pb.User, error) {
// 	user, err := u.repository.GetUserByUsername(ctx, username)
// 	if err != nil {
// 		return pb.User{}, errors.Wrap(err, "Failed to get user ")
// 	}
// 	return user, nil
// }
