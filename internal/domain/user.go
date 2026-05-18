package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyTaken   = errors.New("email already taken")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidUserEmail    = errors.New("invalid user email")
	ErrUnverifiedUserEmail = errors.New("unverified user email")
)

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PasswordHash string // temp
}

type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Phone        string
	PasswordHash string
	IsSuspended  bool
	VerifiedAt   *time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
