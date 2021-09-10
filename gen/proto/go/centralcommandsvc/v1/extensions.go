package centralcommandsvc_v1

import "deblasis.net/space-traffic-control/common/errors"

//TODO find a way to generate these and/or to move them elsewhere
func (r RegisterShipResponse) Failed() error    { return errors.Str2err(r.Error) }
func (r GetAllShipsResponse) Failed() error     { return errors.Str2err(r.Error) }
func (r RegisterStationResponse) Failed() error { return errors.Str2err(r.Error) }
func (r GetAllStationsResponse) Failed() error  { return errors.Str2err(r.Error) }
