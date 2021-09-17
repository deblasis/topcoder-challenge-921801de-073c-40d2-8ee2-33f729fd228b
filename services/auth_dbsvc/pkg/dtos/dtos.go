package dtos

import (
	"deblasis.net/space-traffic-control/common/errs"
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
	Error string      `json:"error,omitempty"`
}

type CreateUserRequest User
type CreateUserResponse struct {
	Id    *uuid.UUID `json:"id"`
	Error string     `json:"error,omitempty"`
}

// ErrorMessage is for performing the error massage and returning by API
type ErrorMessage struct {
	Error []string `json:"error"`
}

// NewErrorMessage returns ErrorMessage by error string
func NewErrorMessage(err string) ErrorMessage {
	return ErrorMessage{Error: []string{err}}
}

func (r GetUserResponse) Failed() error    { return errs.Str2err(r.Error) }
func (r CreateUserResponse) Failed() error { return errs.Str2err(r.Error) }
