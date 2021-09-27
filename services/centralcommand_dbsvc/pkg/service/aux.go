//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package service

import (
	"context"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/repositories"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

type CentralCommandDBAuxService interface {
	Cleanup(context.Context, *pb.CleanupRequest) (*pb.CleanupResponse, error)
}

type centralCommandDBAuxService struct {
	auxRepository repositories.AuxRepository
	logger        log.Logger
}

func NewCentralCommandDBAuxService(
	auxRepository repositories.AuxRepository,
	logger log.Logger,
) CentralCommandDBAuxService {
	return &centralCommandDBAuxService{
		auxRepository: auxRepository,
		logger:        logger,
	}
}

func (s *centralCommandDBAuxService) Cleanup(ctx context.Context, request *pb.CleanupRequest) (resp *pb.CleanupResponse, err error) {
	defer func() {
		if !errs.IsNil(err) {
			level.Debug(s.logger).Log("method", "Cleanup", "err", err)
		}
	}()

	err = s.auxRepository.Cleanup(ctx)
	if err != nil {
		err = errs.NewError(http.StatusInternalServerError, "cannot cleanup", err)
		return &pb.CleanupResponse{Error: errs.ToProtoV1(err)}, nil
	}

	return &pb.CleanupResponse{}, nil
}
