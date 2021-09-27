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
package dtos

import (
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/google/uuid"
)

type User model.User
type GetUserByUsernameRequest struct {
	Username string `json:"username" validate:"required,notblank"`
}
type GetUserByIdRequest struct {
	Id string `json:"id" validate:"required,notblank"`
}
type GetUserResponse struct {
	User  *model.User `json:"user,omitempty"`
	Error error       `json:"error,omitempty"`
}

type CreateUserRequest User
type CreateUserResponse struct {
	Id    *uuid.UUID `json:"id"`
	Error error      `json:"error,omitempty"`
}

func (r GetUserResponse) Failed() error    { return r.Error }
func (r CreateUserResponse) Failed() error { return r.Error }
