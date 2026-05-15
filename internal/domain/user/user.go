package user

import "github.com/google/uuid"

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Name         string
	PassowrdHash string // temp
}
