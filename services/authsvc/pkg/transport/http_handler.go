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
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/authsvc/pkg/endpoints"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

var logger log.Logger

func NewHTTPHandler(e endpoints.EndpointSet, l log.Logger) http.Handler {
	logger = l

	r := mux.NewRouter().StrictSlash(false)
	r.Use(middlewares.JsonHeaderMiddleware)

	options := transport_conf.GetCommonHTTPServerOptions(l)

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

	return json.NewEncoder(w).Encode(response)
}
