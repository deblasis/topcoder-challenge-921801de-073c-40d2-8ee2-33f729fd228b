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
package endpoints

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	"github.com/go-kit/kit/endpoint"
)

// GetUserByUsername(ctx context.Context, username string) (model.User, error)
func (s EndpointSet) GetUserByUsername(ctx context.Context, request *dtos.GetUserByUsernameRequest) (*dtos.GetUserResponse, error) {
	resp, err := s.GetUserByUsernameEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetUserResponse)
	return response, nil
}

// GetUserById(ctx context.Context, id string) (model.User, error)
func (s EndpointSet) GetUserById(ctx context.Context, request *dtos.GetUserByIdRequest) (*dtos.GetUserResponse, error) {
	resp, err := s.GetUserByIdEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.GetUserResponse)
	return response, nil
}

// CreateUser(ctx context.Context, user *model.User) (int64, error)
func (s EndpointSet) CreateUser(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	resp, err := s.CreateUserEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*dtos.CreateUserResponse)
	return response, nil
}

var (
	_ endpoint.Failer = dtos.GetUserResponse{}
	_ endpoint.Failer = dtos.CreateUserResponse{}
)
