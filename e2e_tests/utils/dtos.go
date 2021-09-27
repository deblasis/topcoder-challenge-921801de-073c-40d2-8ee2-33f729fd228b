//go:build integration
// +build integration

package utils

//AUTH
type SignupRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
	Role     string `json:"role,omitempty" faker:"-"`
}
type LoginRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
}

//

//CENTRALCOMMAND
type Dock struct {
	NumDockingPorts int `json:"numDockingPorts" faker:"oneof: 1, 2, 5, 10"`
}
type RegisterStationRequest struct {
	Capacity float32 `json:"capacity" faker:"oneof: 1, 3.14, 8, 256"`
	Docks    []*Dock `json:"docks"`
}

type RegisterShipRequest struct {
	Weight float32 `json:"weight" faker:"oneof: 1, 3.14, 8, 256"`
}

//SHIPPINGSTATION
type RequestLandingRequest struct {
	Time int `json:"time" faker:"oneof: 1, 5, 10, 20"`
}

type LandRequest struct {
	Time   int    `json:"time" faker:"oneof: 1, 5, 10, 20"`
	DockId string `json:"dockId"`
}
