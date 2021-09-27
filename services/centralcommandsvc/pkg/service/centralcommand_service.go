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
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	dbe "deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (s *centralCommandService) RegisterShip(ctx context.Context, request *pb.RegisterShipRequest) (resp *pb.RegisterShipResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "RegisterShip",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "RegisterShip", "err", err)

	verr := s.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.RegisterShipResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.RegisterShipResponse{Error: errs.ToProtoV1(err)}, nil
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
	if !errs.IsNil(ret.Failed()) {
		level.Debug(s.logger).Log("failed", ret.Failed())
		return &pb.RegisterShipResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}
	return converters.DBDtoCreateShipResponseToProto(*ret), nil
}

func (s *centralCommandService) RegisterStation(ctx context.Context, request *pb.RegisterStationRequest) (resp *pb.RegisterStationResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "RegisterStation",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "RegisterStation", "err", err)

	//TODO refactor
	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.RegisterStationResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	userId := common.ExtractUserIdFromCtx(ctx)
	if userId == "" {
		err = errs.NewError(http.StatusBadRequest, "id is empty", errs.ErrValidationFailed)
		return &pb.RegisterStationResponse{Error: errs.ToProtoV1(err)}, nil
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
	level.Error(s.logger).Log("ret", ret.Error)

	if !errs.IsNil(ret.Failed()) {
		level.Debug(s.logger).Log("failed", ret.Failed())
		return &pb.RegisterStationResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}

	return converters.DBDtoCreateStationResponseToProto(*ret), nil
}

func (s *centralCommandService) GetAllShips(ctx context.Context, request *pb.GetAllShipsRequest) (resp *pb.GetAllShipsResponse, err error) {
	level.Info(s.logger).Log("handlingrequest", "GetAllShips",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "GetAllShips", "err", err)

	verr := s.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.GetAllShipsResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	ret, err := s.db_svc_endpointset.GetAllShips(ctx,
		&dtos.GetAllShipsRequest{},
	)
	if err != nil {
		level.Debug(s.logger).Log("err", err)
		return nil, err
	}
	if !errs.IsNil(ret.Failed()) {
		return &pb.GetAllShipsResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}

	return converters.DBDtoGetAllShipsResponseToProto(*ret), nil
}

func (s *centralCommandService) GetAllStations(ctx context.Context, request *pb.GetAllStationsRequest) (resp *pb.GetAllStationsResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "GetAllStations",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "GetAllStations", "err", err)

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.GetAllStationsResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	req := &dtos.GetAllStationsRequest{}
	userId := common.ExtractUserIdFromCtx(ctx)
	role := common.ExtractUserRoleFromCtx(ctx)

	if role == consts.ROLE_SHIP {
		req.ShipId = &userId
	}

	ret, err := s.db_svc_endpointset.GetAllStations(ctx, req)
	if err != nil {
		return nil, err
	}
	if !errs.IsNil(ret.Failed()) {
		return &pb.GetAllStationsResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}

	return converters.DBDtoGetAllStationsResponseToProto(*ret), nil
}

func (s *centralCommandService) GetNextAvailableDockingStation(ctx context.Context, request *pb.GetNextAvailableDockingStationRequest) (resp *pb.GetNextAvailableDockingStationResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "GetNextAvailableDockingStation",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "GetNextAvailableDockingStation", "err", err)

	err = s.validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.GetNextAvailableDockingStationResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	ret, err := s.db_svc_endpointset.GetNextAvailableDockingStation(ctx, &dtos.GetNextAvailableDockingStationRequest{ShipId: uuid.MustParse(request.ShipId).String()})
	if err != nil {
		return nil, err
	}
	if !errs.IsNil(ret.Failed()) {
		return &pb.GetNextAvailableDockingStationResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}
	return converters.DBDtoGetNextAvailableDockingStationResponseToProto(ret), nil
}

func (s *centralCommandService) RegisterShipLanding(ctx context.Context, request *pb.RegisterShipLandingRequest) (resp *pb.RegisterShipLandingResponse, err error) {
	//TODO use middleware
	level.Info(s.logger).Log("handlingrequest", "RegisterShipLanding",
		"userId", ctx.Value(common.ContextKeyUserId),
		"role", ctx.Value(common.ContextKeyUserRole),
	)
	defer level.Info(s.logger).Log("handledrequest", "RegisterShipLanding", "err", err)

	verr := s.validate.Struct(request)
	if verr != nil {
		validationErrors := verr.(validator.ValidationErrors)
		err = errs.NewError(http.StatusBadRequest, "validation failed", validationErrors)
		return &pb.RegisterShipLandingResponse{
			Error: errs.ToProtoV1(err),
		}, nil
	}

	ret, err := s.db_svc_endpointset.LandShipToDock(ctx, &dtos.LandShipToDockRequest{
		ShipId:   uuid.MustParse(request.ShipId).String(),
		DockId:   uuid.MustParse(request.DockId).String(),
		Duration: request.Duration,
	})
	if err != nil {
		return nil, err
	}
	if !errs.IsNil(ret.Failed()) {
		return &pb.RegisterShipLandingResponse{
			Error: errs.ToProtoV1(ret.Error),
		}, nil
	}

	return &pb.RegisterShipLandingResponse{}, nil
}
