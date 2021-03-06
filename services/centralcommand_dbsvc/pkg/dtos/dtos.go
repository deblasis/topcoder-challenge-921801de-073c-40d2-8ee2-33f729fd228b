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
package dtos

type Ship struct {
	//"string - id of the ship"
	Id string `json:"id"`
	//Can be 'docked' | 'in-flight'
	Status string `json:"status,omitempty" model:"-"`
	//validate:"required,oneof='in-flight' 'docked'"
	//'Float - weight of the spaceship'
	Weight float32 `json:"weight,omitempty" validate:"required"`
}

type Station struct {
	//"string - id of the shipping station"
	Id string `json:"id,omitempty"`

	//“float - total capacity”
	Capacity float32 `json:"capacity,omitempty" validate:"required"`

	//“float - total combined weight of all docked spaceships”
	//
	//COMPUTED on the database, so it should be treated as readonly
	UsedCapacity *float32 `json:"usedCapacity,omitempty"`
	//Docks available at the station
	Docks []*Dock `json:"docks" validate:"required"`
}

type Dock struct {
	//“string - id of the dock”
	Id string `json:"id,omitempty"`
	//Id of the station that hosts the dock
	StationId string `json:"station_id,omitempty"`
	//"Integer - total number of available ports"
	NumDockingPorts int64 `json:"numDockingPorts,omitempty" validate:"required"`
	//“Integer - number of docked spaceships on this docking station”
	Occupied *int64 `json:"occupied,omitempty"`
	//“float - combined weight of all docked spaceships on this docking station”
	Weight float32 `json:"weight,omitempty"`
	//Reference to the Station entity
	Station *Station `json:"-" model:"-"`
}

type NextAvailableDockingStation struct {
	DockId                    string   `json:"dock_id,omitempty"`
	StationId                 string   `json:"station_id,omitempty"`
	ShipWeight                float32  `json:"ship_weight,omitempty"`
	AvailableCapacity         *float32 `json:"available_capacity,omitempty"`
	AvailableDocksAtStation   *int64   `json:"available_docks_at_station,omitempty"`
	SecondsUntilNextAvailable *int64   `json:"seconds_until_next_available,omitempty"`
}

type CreateShipRequest Ship
type CreateShipResponse struct {
	Ship  *Ship `json:"ship"`
	Error error `json:"error,omitempty" model:"-"`
}

type GetAllShipsRequest struct{}
type GetAllShipsResponse struct {
	Ships []Ship `json:"ships"`
	Error error  `json:"error,omitempty" model:"-"`
}

type CreateStationRequest Station
type CreateStationResponse struct {
	Station *Station `json:"station"`
	Error   error    `json:"error,omitempty" model:"-"`
}

type GetAllStationsRequest struct {
	ShipId *string `json:"ship_id"`
}
type GetAllStationsResponse struct {
	Stations []Station `json:"stations"`
	Error    error     `json:"error,omitempty" model:"-"`
}

type GetNextAvailableDockingStationRequest struct {
	//"string - id of the ship"
	ShipId string `json:"ship_id"`
}

type GetNextAvailableDockingStationResponse struct {
	NextAvailableDockingStation *NextAvailableDockingStation `json:"next_available_docking_station"`
	Error                       error                        `json:"error,omitempty" model:"-"`
}

type LandShipToDockRequest struct {
	ShipId   string `json:"ship_id,omitempty" validate:"uuid4"`
	DockId   string `json:"dock_id,omitempty" validate:"uuid4"`
	Duration int64  `json:"duration,omitempty" validate:"required,notblank"`
}
type LandShipToDockResponse struct {
	Error error `json:"error,omitempty" model:"-"`
}

func (r CreateShipResponse) Failed() error                     { return r.Error }
func (r GetAllShipsResponse) Failed() error                    { return r.Error }
func (r CreateStationResponse) Failed() error                  { return r.Error }
func (r GetAllStationsResponse) Failed() error                 { return r.Error }
func (r GetNextAvailableDockingStationResponse) Failed() error { return r.Error }
func (r LandShipToDockResponse) Failed() error                 { return r.Error }
