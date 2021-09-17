package authsvc_v1

import "deblasis.net/space-traffic-control/common/errs"

//TODO find a way to generate these and/or to move them elsewhere
func (r SignupResponse) Failed() error     { return errs.Str2err(r.Error) }
func (r LoginResponse) Failed() error      { return errs.Str2err(r.Error) }
func (r CheckTokenResponse) Failed() error { return errs.Str2err(r.Error) }
