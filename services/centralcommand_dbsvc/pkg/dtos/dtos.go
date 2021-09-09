package dtos

import (
	"deblasis.net/space-traffic-control/common/errors"
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
	//Docks availavle at the station
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

type CreateShipRequest Ship
type CreateShipResponse struct {
	Ship *Ship  `json:"ship"`
	Err  string `json:"err,omitempty"`
}

type GetAllShipsRequest struct{}
type GetAllShipsResponse struct {
	Ships []Ship `json:"ships"`
	Err   string `json:"err,omitempty"`
}

type CreateStationRequest Station
type CreateStationResponse struct {
	Station *Station `json:"station"`
	Err     string   `json:"err,omitempty"`
}

type GetAllStationsRequest struct{}
type GetAllStationsResponse struct {
	Stations []Station `json:"stations"`
	Err      string    `json:"err,omitempty"`
}

// ErrorMessage is for performing the error massage and returning by API
type ErrorMessage struct {
	Error []string `json:"error"`
}

// NewErrorMessage returns ErrorMessage by error string
func NewErrorMessage(err string) ErrorMessage {
	return ErrorMessage{Error: []string{err}}
}

func (r CreateShipResponse) Failed() error     { return errors.Str2err(r.Err) }
func (r GetAllShipsResponse) Failed() error    { return errors.Str2err(r.Err) }
func (r CreateStationResponse) Failed() error  { return errors.Str2err(r.Err) }
func (r GetAllStationsResponse) Failed() error { return errors.Str2err(r.Err) }
