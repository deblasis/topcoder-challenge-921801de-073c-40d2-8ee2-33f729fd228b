package service

import (
	"context"

	"deblasis.net/space-traffic-control/common"
	"deblasis.net/space-traffic-control/common/config"
	"deblasis.net/space-traffic-control/common/consts"
	"deblasis.net/space-traffic-control/common/errs"
	ccpb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	pb "deblasis.net/space-traffic-control/gen/proto/go/shippingstationsvc/v1"
	cc "deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ServiceName = "deblasis-v1-ShippingStationService"
	Namespace   = "stc"
	Tags        = []string{}
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

	ret, err := s.centralcommand_endpointset.GetNextAvailableDockingStation(ctx, &ccpb.GetNextAvailableDockingStationRequest{
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

	response := &pb.RequestLandingResponse{}
	next := ret.NextAvailableDockingStation

	if next.AvailableCapacity >= next.ShipWeight && next.AvailableDocksAtStation >= 1 {
		response.Command = pb.RequestLandingResponse_LAND
		response.DockingStationIdOrDuration = &pb.RequestLandingResponse_DockingStationId{DockingStationId: next.DockId}
	} else {
		response.Command = pb.RequestLandingResponse_WAIT
		response.DockingStationIdOrDuration = &pb.RequestLandingResponse_Duration{Duration: next.SecondsUntilNextAvailable}
	}
	return response, nil
}

func (s *shippingStationService) Landing(ctx context.Context, request *pb.LandingRequest) (*pb.LandingResponse, error) {
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "Landing",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "Landing")

	//TODO refactor
	err := s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return &pb.LandingResponse{
			Error: errors.Wrap(validationErrors, "Validation failed").Error(),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		return &pb.LandingResponse{Error: errs.ErrBadRequest.Error()}, nil
	}

	req := &ccpb.RegisterShipLandingRequest{
		ShipId:   userId,
		DockId:   request.DockId,
		Duration: request.Duration,
	}

	ret, err := s.centralcommand_endpointset.RegisterShipLanding(ctx, req)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if ret.Failed() != nil {
		return &pb.LandingResponse{Error: ret.Failed().Error()}, nil
	}

	return &pb.LandingResponse{}, nil
}
