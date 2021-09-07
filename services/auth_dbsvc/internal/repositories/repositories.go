package repositories

import (
	"context"

	"deblasis.net/space-traffic-control/services/auth_dbsvc/internal/model"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
}
