package repository

import (
	"context"
	"social-media-backend/internal/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params user.CreateUserParams) error
}
