package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*uuid.UUID, error)
}
