package service

import (
	"context"
	"net/http"

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

func (s *shippingStationService) RequestLanding(ctx context.Context, request *pb.RequestLandingRequest) (resp *pb.RequestLandingResponse, err error) {
	defer func() {
		if err != nil {
			level.Debug(s.logger).Log("method", "RequestLanding", "err", err)
		}
	}()
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "RequestLanding",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "RequestLanding")

	role := common.ExtractUserRoleFromCtx(ctx)
	if role != consts.ROLE_SHIP {
		err = errs.NewError(http.StatusUnauthorized, "you are not a ship! You can't land here", errs.ErrUnauthorized)
		//TODO check if this should be a domain error
		return nil, err
	}

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.RequestLandingResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.RequestLandingResponse{Error: errs.ToProtoV1(err)}, nil
	}

	ret, err := s.centralcommand_endpointset.GetNextAvailableDockingStation(ctx, &ccpb.GetNextAvailableDockingStationRequest{
		ShipId: userId,
	})
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if !errs.IsNil(ret.Failed()) {
		return &pb.RequestLandingResponse{
			Error: ret.Error,
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

func (s *shippingStationService) Landing(ctx context.Context, request *pb.LandingRequest) (resp *pb.LandingResponse, err error) {
	defer func() {
		if err != nil {
			level.Debug(s.logger).Log("method", "CreateShip", "err", err)
		}
	}()
	//TODO use middleware
	level.Info(s.logger).Log("handling request", "Landing",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handled request", "Landing")

	//TODO refactor
	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.LandingResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.LandingResponse{Error: errs.ToProtoV1(err)}, nil
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
	if !errs.IsNil(ret.Failed()) {
		return &pb.LandingResponse{Error: ret.Error}, nil
	}

	return &pb.LandingResponse{}, nil
}
