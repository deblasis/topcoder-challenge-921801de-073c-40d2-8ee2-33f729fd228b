//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
