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
