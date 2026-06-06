package service

import (
	"context"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo"

	"github.com/google/uuid"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UpdateUserPasswordParams struct {
	UserID      uuid.UUID
	NewPassword string
	OldPassword string
}

func (s *UserService) UpdateUserPassword(ctx context.Context, params UpdateUserPasswordParams) error {
	if params.UserID == uuid.Nil {
		return domain.ErrInvalidUserID
	}

	if params.NewPassword == params.OldPassword {
		return domain.ErrInvalidPassword
	}

	user, err := s.repo.GetUserByID(ctx, params.UserID)
	if err != nil {
		return err
	}

	isMatched, err := auth.ComparePassword(auth.ComparePasswordParams{
		Password:   params.OldPassword,
		StoredHash: user.PasswordHash,
	})
	if err != nil {
		return err
	}

	if !isMatched {
		return domain.ErrInvalidOldPassword
	}

	newPasswordHash, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		return err
	}

	id, err := s.repo.UpdateUserPassword(ctx, domain.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: newPasswordHash,
	})
	if id == uuid.Nil {
		return domain.ErrNoRowsAffected
	}

	return err
}
