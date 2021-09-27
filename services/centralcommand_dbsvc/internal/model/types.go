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
package model

import "time"

type Ship struct {
	tableName struct{} `pg:"ships,select:ships_view"`

	//"string - id of the ship"
	Id string `json:"id,omitempty" db:"id" pg:"id,pk"`
	//Can be 'docked' | 'in-flight'
	Status string `json:"status,omitempty" db:"status"`
	//validate:"required,oneof='in-flight' 'docked'"
	//'Float - weight of the spaceship'
	Weight float32 `json:"weight,omitempty" db:"weight"`
}
type Station struct {
	tableName struct{} `pg:"stations,select:stations_view"`
	//"string - id of the shipping station"
	Id string `json:"id,omitempty" db:"id" pg:"id,pk"`

	//“float - total capacity”
	Capacity float32 `json:"capacity,omitempty" db:"capacity"`

	//“float - total combined weight of all docked spaceships”
	//
	//COMPUTED on the database, so it should be treated as readonly
	UsedCapacity *float32 `json:"usedCapacity,omitempty" db:"used_capacity"`
	//Docks availavle at the station
	Docks []*Dock `json:"docks" pg:"rel:has-many"`
}

type Dock struct {
	tableName struct{} `pg:"docks,select:docks_view"`

	//“string - id of the dock”
	Id string `json:"id,omitempty" db:"id" pg:"id,pk"`
	//Id of the station that hosts the dock
	StationId string `json:"station_id,omitempty" db:"station_id"`
	//"Integer - total number of available ports"
	NumDockingPorts int64 `json:"numDockingPorts,omitempty" db:"num_docking_ports"`
	//“Integer - number of docked spaceships on this docking station”
	Occupied *int64 `json:"occupied,omitempty" db:"occupied"`
	//“float - combined weight of all docked spaceships on this docking station”
	Weight float32 `json:"weight,omitempty" db:"weight"`
	//Reference to the Station entity
	Station *Station `json:"-" pg:"rel:has-one" model:"-"`
}

type DockedShip struct {
	tableName struct{} `pg:"docked_ships"`

	DockId string `json:"dock_id,omitempty" db:"dock_id" pg:"dock_id,pk"`
	ShipId string `json:"ship_id,omitempty" db:"ship_id" pg:"ship_id,pk"`

	DockedSince  time.Time `json:"docked_since,omitempty" db:"docked_since"`
	DockDuration int64     `json:"dock_duration,omitempty" db:"dock_duration"`

	Dock *Dock `json:"-" pg:"rel:has-one"`
	Ship *Ship `json:"-" pg:"rel:has-one"`
}

type NextAvailableDockingStation struct {
	DockId                    string   `json:"dock_id" db:"dock_id"`
	StationId                 string   `json:"station_id" db:"station_id"`
	ShipWeight                float32  `json:"ship_weight" db:"ship_weight"`
	AvailableCapacity         *float32 `json:"available_capacity" db:"available_capacity"`
	AvailableDocksAtStation   *int64   `json:"available_docks_at_station" db:"available_docks_at_station"`
	SecondsUntilNextAvailable *int64   `json:"seconds_until_next_available" db:"seconds_until_next_available"`
}
type AvailableStationsForShip struct {
	StationId       string   `json:"station_id" db:"station_id"`
	Capacity        float32  `json:"capacity" db:"capacity"`
	UsedCapacity    *float32 `json:"used_capacity" db:"used_capacity"`
	DockId          string   `json:"dock_id" db:"dock_id"`
	NumDockingPorts int64    `json:"numDockingPorts" db:"num_docking_ports"`
	Occupied        *int64   `json:"occupied" db:"occupied"`
	Weight          float32  `json:"weight" db:"weight"`
}

//TODO refactor
const (
	ShipsHaveLeftFunctionName                         = "ships_have_left"
	GetAvailableStationsForShipFunctionName           = "stations_available_for_ship"
	GetNextAvailableDockingStationForShipFunctionName = "get_next_available_docking_station_for_ship"
)
