package transport

import (
	"context"
	"encoding/json"
	"net/http"

	errs "deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/api/model"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/service/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var logger log.Logger

func NewHTTPHandler(e endpoints.EndpointSet, l log.Logger) http.Handler {
	logger = l

	r := mux.NewRouter().StrictSlash(false)
	r.Use(jsonHeaderMiddleware)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		e.StatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
		options...,
	))

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
	var req model.GetUserByUsernameRequest

	params := mux.Vars(r)
	username := params["username"]
	if username == "" {
		logger.Log("username is empty")
		return req, nil
	}

	req = model.GetUserByUsernameRequest{
		Username: username,
	}

	return req, nil
}

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.CreateUserRequest
	if r.ContentLength == 0 {
		logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req model.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.deblasis.spacetrafficcontrol-v1+/json; charset=utf-8")
	switch err {
	case errs.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case errs.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func jsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/vnd.deblasis.spacetrafficcontrol-v1+/json; charset=utf-8")
		next.ServeHTTP(rw, r)
	})
}
