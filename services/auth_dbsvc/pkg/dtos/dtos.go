package dtos

import (
	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
)

type GetUserByUsernameRequest struct {
	Username string `json:"username" validate:"required,notblank"`
}
type GetUserByUsernameResponse struct {
	User model.User `json:"user"`
	Err  string     `json:"err,omitempty"`
}

type CreateUserRequest model.User
type CreateUserResponse struct {
	Id  int64  `json:"id"`
	Err string `json:"err,omitempty"`
}

// ErrorMessage is for performing the error massage and returning by API
type ErrorMessage struct {
	Error []string `json:"error"`
}

// NewErrorMessage returns ErrorMessage by error string
func NewErrorMessage(err string) ErrorMessage {
	return ErrorMessage{Error: []string{err}}
}

func (r GetUserByUsernameResponse) Failed() error { return errors.Str2err(r.Err) }
func (r CreateUserResponse) Failed() error        { return errors.Str2err(r.Err) }
