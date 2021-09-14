package authsvc_v1

import "deblasis.net/space-traffic-control/common/errors"

//TODO find a way to generate these and/or to move them elsewhere
func (r SignupResponse) Failed() error     { return errors.Str2err(r.Error) }
func (r LoginResponse) Failed() error      { return errors.Str2err(r.Error) }
func (r CheckTokenResponse) Failed() error { return errors.Str2err(r.Error) }
