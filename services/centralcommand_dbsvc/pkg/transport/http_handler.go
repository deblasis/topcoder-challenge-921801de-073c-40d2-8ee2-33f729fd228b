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
