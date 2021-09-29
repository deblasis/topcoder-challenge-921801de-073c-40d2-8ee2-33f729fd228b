// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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
type DockResp struct {
	Id              string  `json:"id"`
	NumDockingPorts int     `json:"numDockingPorts"`
	Occupied        int     `json:"occupied"`
	Weight          float64 `json:"weight"`
}
type RegisterStationRequest struct {
	Capacity float64 `json:"capacity" faker:"oneof: 1, 3.14, 8, 256"`
	Docks    []*Dock `json:"docks"`
}
type RegisterStationResponse struct {
	Id           string      `json:"id"`
	Capacity     float64     `json:"capacity"`
	UsedCapacity float64     `json:"usedCapacity"`
	Docks        []*DockResp `json:"docks"`
}

type RegisterShipRequest struct {
	Weight float64 `json:"weight" faker:"oneof: 1, 3.14, 8, 256"`
}
type RegisterShipResponse struct {
	Id     string  `json:"id"`
	Weight float64 `json:"weight"`
	Status string  `json:"status"`
}

//SHIPPINGSTATION
type RequestLandingRequest struct {
	Time int `json:"time" faker:"oneof: 1, 5, 10, 20"`
}

type LandRequest struct {
	Time   int    `json:"time" faker:"oneof: 1, 5, 10, 20"`
	DockId string `json:"dockId"`
}
