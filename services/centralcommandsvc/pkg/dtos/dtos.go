package dtos

import "deblasis.net/space-traffic-control/common/errors"

type Dock struct {
	//“string - id of the dock”
	Id string `json:"id,omitempty"`
	//Id of the station that hosts the dock
	StationId string `json:"station_id,omitempty"`
	//"Integer - total number of available ports"
	NumDockingPorts int64 `json:"numDockingPorts,omitempty"`
	//“Integer - number of docked spaceships on this docking station”
	Occupied int64 `json:"occupied,omitempty"`
	//“float - combined weight of all docked spaceships on this docking station”
	Weight float32 `json:"weight,omitempty"`
}

type Station struct {
	//"string - id of the shipping station"
	Id string `json:"id,omitempty"`

	//“float - total capacity”
	Capacity float32 `json:"capacity,omitempty"`

	//“float - total combined weight of all docked spaceships”
	//
	//COMPUTED on the database, so it should be treated as readonly
	UsedCapacity float32 `json:"usedCapacity,omitempty"`
	//Docks availavle at the station
	Docks []*Dock `json:"docks"`
}

type Ship struct {
	//"string - id of the ship"
	Id string `json:"id,omitempty"`
	//Can be 'docked' | 'in-flight'
	Status string `json:"status,omitempty" model:"-"`
	//validate:"required,oneof='in-flight' 'docked'"
	//'Float - weight of the spaceship'
	Weight float32 `json:"weight,omitempty"`
}

type RegisterShipRequest Ship

type RegisterShipResponse struct {
	Err string `json:"err,omitempty"`
}
type GetAllShipsRequest struct {
}
type GetAllShipsResponse struct {
	Ships []Ship `json:"ships"`
	Err   string `json:"err,omitempty"`
}
type RegisterStationRequest Station
type RegisterStationResponse struct {
	Err string `json:"err,omitempty"`
}
type GetAllStationsRequest struct {
}
type GetAllStationsResponse struct {
	Stations []Station `json:"stations"`
	Err      string    `json:"err,omitempty"`
}

// type SignupRequest struct {
// 	Username string `json:"username" validate:"required,notblank"`
// 	Password string `json:"password" validate:"required,notblank"`
// 	Role     string `json:"role" validate:"required,oneof=Ship Station Command"`
// }

// type SignupResponse struct {
// 	Token Token  `json:"token,omitempty"`
// 	Err   string `json:"err,omitempty"`
// }

// type Token struct {
// 	Token     string `json:"token,omitempty"`
// 	ExpiresAt int64  `json:"expires_at,omitempty"`
// }

// type LoginRequest struct {
// 	Username string `json:"username" validate:"required,notblank"`
// 	Password string `json:"password" validate:"required,notblank"`
// }

// type LoginResponse struct {
// 	Token Token  `json:"token,omitempty"`
// 	Err   string `json:"err,omitempty"`
// }

func (r RegisterShipResponse) Failed() error    { return errors.Str2err(r.Err) }
func (r GetAllShipsResponse) Failed() error     { return errors.Str2err(r.Err) }
func (r RegisterStationResponse) Failed() error { return errors.Str2err(r.Err) }
func (r GetAllStationsResponse) Failed() error  { return errors.Str2err(r.Err) }
