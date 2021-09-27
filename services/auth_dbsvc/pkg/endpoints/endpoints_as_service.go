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
