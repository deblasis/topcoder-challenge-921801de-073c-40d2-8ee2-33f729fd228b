package shippingstationsvc_v1

import (
	"deblasis.net/space-traffic-control/common/errs"
)

func (r RequestLandingResponse) Failed() error { return errs.Str2err(r.Error) }
func (r LandingResponse) Failed() error        { return errs.Str2err(r.Error) }
