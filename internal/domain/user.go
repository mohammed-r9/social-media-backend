package domain

import (
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
	ErrInvalidPassword     = errors.New("invalid user password")
	ErrInvalidOldPassword  = errors.New("invalid user old password")
	ErrInvalidUserName     = errors.New("invalid user name")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrNoRowsAffected      = errors.New("no rows affected")
)

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Username     string
	PasswordHash string
}

type UpdatePasswordParams struct {
	ID           uuid.UUID
	PasswordHash string
}

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	Name         string
	Phone        string
	PasswordHash string
	IsSuspended  bool
	VerifiedAt   *time.Time
}
