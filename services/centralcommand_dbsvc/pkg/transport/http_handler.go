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
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var logger log.Logger

func NewHTTPHandler(e endpoints.EndpointSet, l log.Logger) http.Handler {
	logger = l

	r := mux.NewRouter().StrictSlash(false)
	r.Use(middlewares.JsonHeaderMiddleware)

	options := transport_conf.GetCommonHTTPServerOptions(l)

	r.Methods("POST").Path("/ship").Handler(httptransport.NewServer(
		e.CreateShipEndpoint,
		decodeHTTPCreateShipRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/ship/all").Handler(httptransport.NewServer(
		e.GetAllShipsEndpoint,
		decodeHTTPGetAllShipsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/station").Handler(httptransport.NewServer(
		e.CreateStationEndpoint,
		decodeHTTPCreateStationRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/station/all").Handler(httptransport.NewServer(
		e.GetAllStationsEndpoint,
		decodeHTTPGetAllStationsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/docks/nextavailable").Handler(httptransport.NewServer(
		e.GetNextAvailableDockingStationEndpoint,
		decodeHTTPGetNextAvailableDockingStationRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/land").Handler(httptransport.NewServer(
		e.LandShipToDockEndpoint,
		decodeHTTPLandShipToDockRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPCreateShipRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.CreateShipRequest
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

func decodeHTTPGetAllShipsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &dtos.GetAllShipsRequest{}, nil
}

func decodeHTTPCreateStationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.CreateStationRequest
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

func decodeHTTPGetAllStationsRequest(_ context.Context, r *http.Request) (interface{}, error) {

	qs := r.URL.Query().Get("ship_id")
	if qs != "" {
		return &dtos.GetAllStationsRequest{ShipId: &qs}, nil
	}

	return &dtos.GetAllStationsRequest{}, nil
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

	return json.NewEncoder(w).Encode(response)
}

func decodeHTTPGetNextAvailableDockingStationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.GetNextAvailableDockingStationRequest
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
func decodeHTTPLandShipToDockRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.LandShipToDockRequest
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
