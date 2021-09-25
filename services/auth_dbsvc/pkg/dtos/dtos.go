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
