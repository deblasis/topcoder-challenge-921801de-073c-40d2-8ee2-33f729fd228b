package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/errs"
	"deblasis.net/space-traffic-control/common/middlewares"
	"deblasis.net/space-traffic-control/common/transport_conf"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/endpoints"
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

	// r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
	// 	e.StatusEndpoint,
	// 	healthcheck.DecodeHTTPServiceStatusRequest,
	// 	encodeResponse,
	// 	options...,
	// ))

	r.Methods("POST").Path("/users").Handler(httptransport.NewServer(
		e.CreateUserEndpoint,
		decodeHTTPCreateUserRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/users/{username}").Handler(httptransport.NewServer(
		e.GetUserByUsernameEndpoint,
		decodeHTTPGetUserByUsernameRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeHTTPGetUserByUsernameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.GetUserByUsernameRequest

	params := mux.Vars(r)
	username := params["username"]
	if username == "" {
		logger.Log("username is empty")
		return req, nil
	}

	req = dtos.GetUserByUsernameRequest{
		Username: username,
	}

	return &req, nil
}

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req dtos.CreateUserRequest
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
	if e, ok := response.(endpoint.Failer); ok && e != nil && e.Failed() != nil {
		errs.EncodeErrorHTTP(ctx, e.Failed(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}
