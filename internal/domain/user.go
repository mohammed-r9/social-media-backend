package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidUserID      = errors.New("invalid user id")
)

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string // temp
}

type User struct {
	ID    uuid.UUID
	Email string
	Name  string
}

type UserRepository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
}
