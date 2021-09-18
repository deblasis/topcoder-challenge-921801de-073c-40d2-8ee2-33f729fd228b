package dtos

import (
	"deblasis.net/space-traffic-control/common/errs"
)

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
	UsedCapacity float32 `json:"usedCapacity,omitempty"`
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
	Occupied int64 `json:"occupied,omitempty"`
	//“float - combined weight of all docked spaceships on this docking station”
	Weight float32 `json:"weight,omitempty"`
	//Reference to the Station entity
	Station *Station `json:"-" model:"-"`
}

type NextAvailableDockingStation struct {
	DockId                  string  `json:"dock_id,omitempty"`
	StationId               string  `json:"station_id,omitempty"`
	AvailableCapacity       float32 `json:"available_capacity,omitempty"`
	AvailableDocksAtStation int64   `json:"available_docks_at_station,omitempty"`
	SecondsUntilAvailable   int64   `json:"seconds_until_available,omitempty"`
}

type CreateShipRequest Ship
type CreateShipResponse struct {
	Ship  *Ship  `json:"ship"`
	Error string `json:"error,omitempty"`
}

type GetAllShipsRequest struct{}
type GetAllShipsResponse struct {
	Ships []Ship `json:"ships"`
	Error string `json:"error,omitempty"`
}

type CreateStationRequest Station
type CreateStationResponse struct {
	Station *Station `json:"station"`
	Error   string   `json:"error,omitempty"`
}

type GetAllStationsRequest struct{}
type GetAllStationsResponse struct {
	Stations []Station `json:"stations"`
	Error    string    `json:"error,omitempty"`
}

type GetNextAvailableDockingStationRequest struct {
	//"string - id of the ship"
	ShipId string `json:"ship_id"`
}

type GetNextAvailableDockingStationResponse struct {
	NextAvailableDockingStation *NextAvailableDockingStation `json:"next_available_docking_station"`
	Error                       string                       `json:"error,omitempty"`
}

func (r CreateShipResponse) Failed() error                     { return errs.Str2err(r.Error) }
func (r GetAllShipsResponse) Failed() error                    { return errs.Str2err(r.Error) }
func (r CreateStationResponse) Failed() error                  { return errs.Str2err(r.Error) }
func (r GetAllStationsResponse) Failed() error                 { return errs.Str2err(r.Error) }
func (r GetNextAvailableDockingStationResponse) Failed() error { return errs.Str2err(r.Error) }
