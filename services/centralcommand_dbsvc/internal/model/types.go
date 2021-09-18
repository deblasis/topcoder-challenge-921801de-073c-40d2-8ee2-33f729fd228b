package model

type Ship struct {
	tableName struct{} `pg:"ships,select:ships_view"`

	//"string - id of the ship"
	Id string `json:"id,omitempty" db:"id"`
	//Can be 'docked' | 'in-flight'
	Status string `json:"status,omitempty" db:"status"`
	//validate:"required,oneof='in-flight' 'docked'"
	//'Float - weight of the spaceship'
	Weight float32 `json:"weight,omitempty" db:"weight"`
}
type Station struct {
	tableName struct{} `pg:"stations,select:stations_view"`
	//"string - id of the shipping station"
	Id string `json:"id,omitempty" db:"id"`

	//“float - total capacity”
	Capacity float32 `json:"capacity,omitempty" db:"capacity"`

	//“float - total combined weight of all docked spaceships”
	//
	//COMPUTED on the database, so it should be treated as readonly
	UsedCapacity float32 `json:"usedCapacity,omitempty" db:"used_capacity"`
	//Docks availavle at the station
	Docks []*Dock `json:"docks" pg:"rel:has-many"`
}

type Dock struct {
	tableName struct{} `pg:"docks,select:docks_view"`

	//“string - id of the dock”
	Id string `json:"id,omitempty" db:"id"`
	//Id of the station that hosts the dock
	StationId string `json:"station_id,omitempty" db:"station_id"`
	//"Integer - total number of available ports"
	NumDockingPorts int64 `json:"numDockingPorts,omitempty" db:"num_docking_ports"`
	//“Integer - number of docked spaceships on this docking station”
	Occupied int64 `json:"occupied,omitempty" db:"occupied"`
	//“float - combined weight of all docked spaceships on this docking station”
	Weight float32 `json:"weight,omitempty" db:"weight"`
	//Reference to the Station entity
	Station *Station `json:"-" pg:"rel:has-one" model:"-"`
}

type DockedShip struct {
	tableName struct{} `pg:"docked_ships"`

	DockId string `json:"dock_id,omitempty" db:"dock_id"`
	ShipId string `json:"ship_id,omitempty" db:"ship_id"`

	Dock *Dock `json:"-" pg:"rel:has-one"`
	Ship *Ship `json:"-" pg:"rel:has-one"`
}

type NextAvailableDockingStation struct {
	DockId                  string  `json:"dock_id,omitempty" db:"dock_id"`
	StationId               string  `json:"station_id,omitempty" db:"station_id"`
	AvailableCapacity       float32 `json:"available_capacity,omitempty" db:"available_capacity"`
	AvailableDocksAtStation int64   `json:"available_docks_at_station,omitempty" db:"available_docks_at_station"`
	SecondsUntilAvailable   int64   `json:"seconds_until_available" db:"seconds_until_available"`
}

//TODO refactor
const (
	ShipsHaveLeftFunctionName                         = "ships_have_left"
	GetNextAvailableDockingStationForShipFunctionName = "get_next_available_docking_station_for_ship"
)
