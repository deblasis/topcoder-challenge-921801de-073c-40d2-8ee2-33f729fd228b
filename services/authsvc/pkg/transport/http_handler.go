package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/encoding"
	"deblasis.net/space-traffic-control/common/healthcheck"
	"deblasis.net/space-traffic-control/common/middlewares"
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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
	return req, nil
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
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encoding.EncodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
