package service

import (
	"context"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ServiceName    = "deblasis-v1-CentralCommandService"
	Namespace      = "stc"
	Tags           = []string{}
	GrpcServerPort = 9482 //TODO config
)

type CentralCommandService interface {
	RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error)
	GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error)

	RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error)
	GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error)
	GetNextAvailableDockingStation(context.Context, *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error)
	RegisterShipLanding(context.Context, *pb.RegisterShipLandingRequest) (*pb.RegisterShipLandingResponse, error)
}

type centralCommandService struct {
	logger             log.Logger
	validate           *validator.Validate
	db_svc_endpointset dbe.EndpointSet
}

func NewCentralCommandService(logger log.Logger, jwtConfig config.JWTConfig, db_svc_endpointset dbe.EndpointSet) CentralCommandService {
	return &centralCommandService{
		logger:             logger,
		validate:           common.GetValidator(),
		db_svc_endpointset: db_svc_endpointset,
	}
}

func (s *centralCommandService) RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (*pb.RegisterShipResponse, error) {
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "RegisterShip",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "RegisterShip")

	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.RegisterShipResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		return &pb.RegisterShipResponse{Error: errs.ErrBadRequest.Error()}, nil
	}

	req := &dtos.CreateShipRequest{
		Id:     userId,
		Weight: request.Weight,
	}
	ret, err := s.db_svc_endpointset.CreateShip(ctx, req)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.RegisterShipResponse{
			Error: ret.Failed().Error(),
		}, nil
	}
	return converters.DBDtoCreateShipResponseToProto(*ret), nil
}

func (s *centralCommandService) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "RegisterStation",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "RegisterStation")

	//TODO refactor
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.RegisterStationResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		return &pb.RegisterStationResponse{Error: errs.ErrBadRequest.Error()}, nil
	}

	req := &dtos.CreateStationRequest{
		Id:       userId,
		Capacity: request.Capacity,
		Docks:    make([]*dtos.Dock, 0),
	}
	for _, s := range request.Docks {
		req.Docks = append(req.Docks, &dtos.Dock{
			NumDockingPorts: s.NumDockingPorts,
		})
	}

	ret, err := s.db_svc_endpointset.CreateStation(ctx, req)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.RegisterStationResponse{Error: ret.Failed().Error()}, nil
	}

	return converters.DBDtoCreateStationResponseToProto(*ret), nil
}

func (s *centralCommandService) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.GetAllShipsResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := s.db_svc_endpointset.GetAllShips(ctx,
		&dtos.GetAllShipsRequest{},
	)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.GetAllShipsResponse{
			Error: ret.Failed().Error(),
		}, nil
	}

	return converters.DBDtoGetAllShipsResponseToProto(*ret), nil
}

func (s *centralCommandService) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "GetAllStations")
	defer level.Info(s.logger).Log("handled request", "GetAllStations")

	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.GetAllStationsResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := s.db_svc_endpointset.GetAllStations(ctx,
		&dtos.GetAllStationsRequest{},
	)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.GetAllStationsResponse{
			Error: ret.Failed().Error(),
		}, nil
	}

	return converters.DBDtoGetAllStationsResponseToProto(*ret), nil
}

func (u *centralCommandService) GetNextAvailableDockingStation(ctx context.Context, request *pb.GetNextAvailableDockingStationRequest) (*pb.GetNextAvailableDockingStationResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.GetNextAvailableDockingStationResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	ret, err := u.db_svc_endpointset.GetNextAvailableDockingStation(ctx, &dtos.GetNextAvailableDockingStationRequest{ShipId: uuid.MustParse(request.ShipId).String()})
	if err != nil {
		return &pb.GetNextAvailableDockingStationResponse{
			Error: errs.Err2str(err),
		}, nil
	}
	//TODO converter
	return &pb.GetNextAvailableDockingStationResponse{
		NextAvailableDockingStation: &pb.NextAvailableDockingStation{
			DockId:                    ret.NextAvailableDockingStation.DockId,
			StationId:                 ret.NextAvailableDockingStation.StationId,
			ShipWeight:                ret.NextAvailableDockingStation.ShipWeight,
			AvailableCapacity:         ret.NextAvailableDockingStation.AvailableCapacity,
			AvailableDocksAtStation:   ret.NextAvailableDockingStation.AvailableDocksAtStation,
			SecondsUntilNextAvailable: ret.NextAvailableDockingStation.SecondsUntilNextAvailable,
		},
	}, nil
}

func (u *centralCommandService) RegisterShipLanding(ctx context.Context, request *pb.RegisterShipLandingRequest) (*pb.RegisterShipLandingResponse, error) {
	err := u.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.RegisterShipLandingResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	_, err = u.db_svc_endpointset.LandShipToDock(ctx, &dtos.LandShipToDockRequest{
		ShipId:   uuid.MustParse(request.ShipId).String(),
		DockId:   uuid.MustParse(request.DockId).String(),
		Duration: request.Duration,
	})
	if err != nil {
		return &pb.RegisterShipLandingResponse{
			Error: errs.Err2str(err),
		}, nil
	}
	return &pb.RegisterShipLandingResponse{}, nil
}
