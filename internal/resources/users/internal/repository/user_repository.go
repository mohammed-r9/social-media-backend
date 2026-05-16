package repo

import (
	"context"
	"social-media-backend/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
}
