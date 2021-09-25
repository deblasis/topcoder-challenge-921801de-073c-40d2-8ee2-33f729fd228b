package shippingstationsvc_v1

func (r RequestLandingResponse) Failed() error { return r.Error }
func (r LandingResponse) Failed() error        { return r.Error }
