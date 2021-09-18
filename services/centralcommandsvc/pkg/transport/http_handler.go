package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/common/transport_conf"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/endpoints"
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

	r.Methods("POST").Path("/ship").Handler(httptransport.NewServer(
		e.RegisterShipEndpoint,
		decodeHTTPRegisterShipRequest,
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
		e.RegisterStationEndpoint,
		decodeHTTPRegisterStationRequest,
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

	return r
}

func decodeHTTPRegisterShipRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.RegisterShipRequest
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
	return &pb.GetAllShipsRequest{}, nil
}

func decodeHTTPRegisterStationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.RegisterStationRequest
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
	return &pb.GetAllStationsRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		errs.EncodeErrorHTTP(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeHTTPGetNextAvailableDockingStationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.GetNextAvailableDockingStationRequest
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
