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
