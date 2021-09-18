package service

import (
	"context"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/consts"
	"deblasis.net/space-traffic-control/common/errs"
	ccpb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	cc "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/shippingstationsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName    = "deblasis-v1-ShippingStationService"
	Namespace      = "stc"
	Tags           = []string{}
	GrpcServerPort = 9482 //TODO config
)

type ShippingStationService interface {
	RequestLanding(ctx context.Context, request *pb.RequestLandingRequest) (*pb.RequestLandingResponse, error)
	Landing(ctx context.Context, request *pb.LandingRequest) (*pb.LandingResponse, error)
}

type shippingStationService struct {
	logger                     log.Logger
	validate                   *validator.Validate
	centralcommand_endpointset cc.EndpointSet
}

func NewShippingStationService(logger log.Logger, jwtConfig config.JWTConfig, centralcommand_endpointset cc.EndpointSet) ShippingStationService {
	return &shippingStationService{
		logger:                     logger,
		validate:                   common.GetValidator(),
		centralcommand_endpointset: centralcommand_endpointset,
	}
}

var (
	ErrShipAlreadyRegistered    = errors.New("this ship is already registered")
	ErrStationAlreadyRegistered = errors.New("this station is already registered")
)

func (s *shippingStationService) RequestLanding(ctx context.Context, request *pb.RequestLandingRequest) (*pb.RequestLandingResponse, error) {
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "RequestLanding",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "RequestLanding")

	role := common.ExtractUserRoleFromCtx(ctx)
	if role != consts.ROLE_SHIP {
		return nil, errs.ErrUnauthorized
	}

	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.RequestLandingResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		return &pb.RequestLandingResponse{Error: errs.ErrBadRequest.Error()}, nil
	}

	ret, err := s.centralcommand_endpointset.GetNextAvailableDockingStationEndpoint(ctx, &ccpb.GetNextAvailableDockingStationRequest{
		ShipId: userId,
	})
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.RequestLandingResponse{
			Error: ret.Failed().Error(),
		}, nil
	}

	


	//TODO refactor
	return &pb.RequestLandingResponse{
		Command:                    ,
		DockingStationIdOrDuration: ,
	}, nil
}

func (s *shippingStationService) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (*pb.RegisterStationResponse, error) {
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

func (s *shippingStationService) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (*pb.GetAllShipsResponse, error) {
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

func (s *shippingStationService) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (*pb.GetAllStationsResponse, error) {
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
