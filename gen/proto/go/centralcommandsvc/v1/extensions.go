package centralcommandsvc_v1

func (r RegisterShipResponse) Failed() error                   { return r.Error }
func (r GetAllShipsResponse) Failed() error                    { return r.Error }
func (r RegisterStationResponse) Failed() error                { return r.Error }
func (r GetAllStationsResponse) Failed() error                 { return r.Error }
func (r GetNextAvailableDockingStationResponse) Failed() error { return r.Error }
func (r RegisterShipLandingResponse) Failed() error            { return r.Error }
