package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/encoding"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var logger log.Logger

func NewHTTPHandler(e endpoints.EndpointSet, l log.Logger) http.Handler {
	logger = l

	r := mux.NewRouter().StrictSlash(false)
	r.Use(middlewares.JsonHeaderMiddleware)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encoding.EncodeError),
	}

	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		e.StatusEndpoint,
		healthcheck.DecodeHTTPServiceStatusRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/ships").Handler(httptransport.NewServer(
		e.CreateShipEndpoint,
		decodeHTTPCreateShipRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/ships").Handler(httptransport.NewServer(
		e.GetAllShipsEndpoint,
		decodeHTTPGetAllShipsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/stations").Handler(httptransport.NewServer(
		e.CreateStationEndpoint,
		decodeHTTPCreateStationRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/stations").Handler(httptransport.NewServer(
		e.GetAllStationsEndpoint,
		decodeHTTPGetAllStationsRequest,
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
	return req, nil
}

func decodeHTTPGetAllShipsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.GetAllShipsRequest
	return req, nil
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
	return req, nil
}

func decodeHTTPGetAllStationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.GetAllStationsRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encoding.EncodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
