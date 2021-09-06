package dtos

import "deblasis.net/space-traffic-control/common/errors"

type SignupRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
	Role     string `json:"role" validate:"required,oneof=Ship Station Command"`
}

type SignupResponse struct {
	Token Token  `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}

type Token struct {
	Token     string `json:"token,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,notblank"`
	Password string `json:"password" validate:"required,notblank"`
}

type LoginResponse struct {
	Token Token  `json:"token,omitempty"`
	Err   string `json:"err,omitempty"`
}

func (r SignupResponse) Failed() error { return errors.Str2err(r.Err) }
func (r LoginResponse) Failed() error  { return errors.Str2err(r.Err) }
