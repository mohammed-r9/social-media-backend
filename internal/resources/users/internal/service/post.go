package service

import (
	"context"
	"social-media-backend/internal/crypto/password"
	"social-media-backend/internal/domain"

	"github.com/google/uuid"
)

type RegisterParams struct {
	Name              string
	Email             string
	PassowrdPlainText string
}

func (s *Service) Register(ctx context.Context, params RegisterParams) (domain.User, error) {
	passowrdHash, err := password.HashPassword(params.PassowrdPlainText)
	if err != nil {
		return domain.User{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return domain.User{}, err
	}

	return s.repo.CreateUser(ctx, domain.CreateUserParams{
		ID:           userID,
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: passowrdHash,
	})
}
