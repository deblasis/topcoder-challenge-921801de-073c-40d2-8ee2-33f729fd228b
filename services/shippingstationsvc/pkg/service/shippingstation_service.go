// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
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
	"github.com/go-kit/log"
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
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "RequestLanding", "err", err)
		}
	}()

	level.Info(s.logger).Log("handlingrequest", "RequestLanding",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "RequestLanding")

	role := common.ExtractUserRoleFromCtx(ctx)
	if role != consts.ROLE_SHIP {
		err = errs.NewError(http.StatusUnauthorized, "you are not a ship! You can't land here", errs.ErrUnauthorized)
		return nil, err
	}
	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.RequestLandingResponse{Error: errs.ToProtoV1(err)}, nil
	}

	request.Id = userId

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.RequestLandingResponse{
			Error: errs.ToProtoV1(err),
		}, nil
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
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "CreateShip", "err", err)
		}
	}()

	level.Info(s.logger).Log("handlingrequest", "Landing",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "Landing")

	request.ShipId = common.ExtractUserIdFromCtx(ctx)
	if request.ShipId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.LandingResponse{Error: errs.ToProtoV1(err)}, nil
	}

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.LandingResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	req := &ccpb.RegisterShipLandingRequest{
		ShipId:   request.ShipId,
		DockId:   request.DockId,
		Duration: request.Time,
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
