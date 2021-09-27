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
//go:build integration
// +build integration

package utils

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"

	centralcommand_dbsvc_v1 "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/transport"
	"google.golang.org/grpc"
)

func CleanupDB(ctx context.Context, logger log.Logger) error {
	var (
		conn *grpc.ClientConn
		err  error
	)
	target := "localhost:9383"
	if envTarget := os.Getenv("CENTRALCOMMAND_AUX_DBENDPOINT"); envTarget != "" {
		target = envTarget
	}

	conn, err = grpc.Dial(target, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return err
	}

	client := transport.NewAuxGrpcClient(conn, logger)
	resp, err := client.Cleanup(ctx, &centralcommand_dbsvc_v1.CleanupRequest{})
	if err != nil {
		return err
	}
	return resp.Error
}
