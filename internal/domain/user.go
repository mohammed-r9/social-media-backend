package domain

import "github.com/google/uuid"

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
