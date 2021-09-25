package authsvc_v1

//TODO find a way to generate these and/or to move them elsewhere
func (r SignupResponse) Failed() error     { return r.Error }
func (r LoginResponse) Failed() error      { return r.Error }
func (r CheckTokenResponse) Failed() error { return r.Error }
