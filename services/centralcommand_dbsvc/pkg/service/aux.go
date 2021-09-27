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
