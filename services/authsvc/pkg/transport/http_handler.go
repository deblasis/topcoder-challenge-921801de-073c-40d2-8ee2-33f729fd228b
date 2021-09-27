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
package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var logger log.Logger

func NewHTTPHandler(e endpoints.EndpointSet, l log.Logger) http.Handler {
	logger = l

	r := mux.NewRouter().StrictSlash(false)
	r.Use(middlewares.JsonHeaderMiddleware)

	options := transport_conf.GetCommonHTTPServerOptions(l)
	// r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
	// 	e.StatusEndpoint,
	// 	healthcheck.DecodeHTTPServiceStatusRequest,
	// 	encodeResponse,
	// 	options...,
	// ))

	r.Methods("POST").Path("/signup").Handler(httptransport.NewServer(
		e.SignupEndpoint,
		decodeHTTPSignupRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		e.LoginEndpoint,
		decodeHTTPLoginRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPSignupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.SignupRequest
	if r.ContentLength == 0 {
		logger.Log("Post request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeHTTPLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.LoginRequest
	if r.ContentLength == 0 {
		logger.Log("Post request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		errs.EncodeErrorHTTP(ctx, e, w)
		return nil
	}
	if e, ok := response.(endpoint.Failer); ok && e != nil && !errs.IsNil(e.Failed()) {
		errs.EncodeErrorHTTP(ctx, e.Failed(), w)
		return nil
	}

	//TODO need to handle status when business error here and in all services...
	return json.NewEncoder(w).Encode(response)
}
