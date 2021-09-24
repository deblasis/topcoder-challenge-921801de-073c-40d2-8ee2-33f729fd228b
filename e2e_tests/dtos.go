//go:build integration
// +build integration

package e2e_tests

//AUTH
type signupRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
	Role     string `json:"role,omitempty" faker:"-"`
}
type loginRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
}

//

//CENTRALCOMMAND
type dock struct {
	NumDockingPorts int `json:"numDockingPorts" faker:"oneof: 1, 2, 5, 10"`
}
type registerStationRequest struct {
	Capacity float32 `json:"capacity" faker:"oneof: 1, 3.14, 8, 256"`
	Docks    []*dock `json:"docks"`
}

type registerShipRequest struct {
	Weight float32 `json:"weight" faker:"oneof: 1, 3.14, 8, 256"`
}
