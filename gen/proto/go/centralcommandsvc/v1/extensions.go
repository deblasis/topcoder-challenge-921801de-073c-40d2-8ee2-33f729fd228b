package centralcommandsvc_v1

import (
	"deblasis.net/space-traffic-control/common/errs"
)

func (r RegisterShipResponse) Failed() error                   { return errs.Str2err(r.Error) }
func (r GetAllShipsResponse) Failed() error                    { return errs.Str2err(r.Error) }
func (r RegisterStationResponse) Failed() error                { return errs.Str2err(r.Error) }
func (r GetAllStationsResponse) Failed() error                 { return errs.Str2err(r.Error) }
func (r GetNextAvailableDockingStationResponse) Failed() error { return errs.Str2err(r.Error) }
func (r RegisterShipLandingResponse) Failed() error            { return errs.Str2err(r.Error) }
