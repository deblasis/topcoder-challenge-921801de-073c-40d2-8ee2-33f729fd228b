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
//go:build integration
// +build integration

package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/log"

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

func BarelyEqual(a, b float64) bool {
	ret := fmt.Sprintf("%.2f", float64(a)) == fmt.Sprintf("%.2f", float64(b))
	return ret
}

func RemoveItemAtIndex(slice *[]interface{}, i int) {
	s := *slice
	s[i] = s[len(s)-1]
	s[len(s)-1] = nil
	*slice = s[:len(s)-1]
}
