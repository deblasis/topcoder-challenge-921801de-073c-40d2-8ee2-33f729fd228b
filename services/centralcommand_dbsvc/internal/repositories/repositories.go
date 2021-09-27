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
package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"github.com/google/uuid"
)

type ShipRepository interface {
	GetById(ctx context.Context, id string) (*model.Ship, error)
	Create(ctx context.Context, ship model.Ship) (*model.Ship, error)
	GetAll(ctx context.Context) ([]model.Ship, error)
}
type StationRepository interface {
	GetById(ctx context.Context, id string) (*model.Station, error)
	Create(ctx context.Context, station model.Station) (*model.Station, error)
	GetAll(ctx context.Context) ([]model.Station, error)
	GetAvailableForShip(ctx context.Context, shipId uuid.UUID) ([]model.Station, error)
}
type DockRepository interface {
	GetById(ctx context.Context, id string) (*model.Dock, error)
	Create(ctx context.Context, dock model.Dock) (*model.Dock, error)
	GetNextAvailableDockingStation(ctx context.Context, shipId uuid.UUID) (*model.NextAvailableDockingStation, error)
	LandShipToDock(ctx context.Context, shipId uuid.UUID, dockId uuid.UUID, duration int64) (*model.DockedShip, error)
}

type AuxRepository interface {
	Cleanup(ctx context.Context) error
}
