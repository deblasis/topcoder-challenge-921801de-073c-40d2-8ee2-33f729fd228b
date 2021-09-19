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
}
type DockRepository interface {
	GetById(ctx context.Context, id string) (*model.Dock, error)
	Create(ctx context.Context, dock model.Dock) (*model.Dock, error)
	GetNextAvailableDockingStation(ctx context.Context, shipId uuid.UUID) (*model.NextAvailableDockingStation, error)
	LandShipToDock(ctx context.Context, shipId uuid.UUID, dockId uuid.UUID, duration int64) (*model.DockedShip, error)
}
