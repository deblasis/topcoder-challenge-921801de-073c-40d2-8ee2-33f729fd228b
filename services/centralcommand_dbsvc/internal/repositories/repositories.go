package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
)

type ShipRepository interface {
	GetById(ctx context.Context, id string) (*model.Ship, error)
	//GetByStatus(ctx context.Context, status string) ([]model.Ship, error)
	Create(ctx context.Context, ship model.Ship) (*model.Ship, error)
	GetAll(ctx context.Context) ([]model.Ship, error)
}
type StationRepository interface {
	GetById(ctx context.Context, id string) (*model.Station, error)
	//GetByAvailableCapacityGreaterThanEqual(ctx context.Context, requiredCapacity float32) ([]model.Station, error)
	Create(ctx context.Context, station model.Station) (*model.Station, error)
	GetAll(ctx context.Context) ([]model.Station, error)
}
type DockRepository interface {
	GetById(ctx context.Context, id string) (*model.Dock, error)
	Create(ctx context.Context, dock model.Dock) (*model.Dock, error)
}
