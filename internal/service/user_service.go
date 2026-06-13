package service

import (
	"context"
	"social-media-backend/internal/adapters/s3"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo"

	"github.com/google/uuid"
)

type UserService struct {
	repo repo.UserRepository
	file s3.FileStorage
}

func NewUserService(repo repo.UserRepository, fileStorage s3.FileStorage) *UserService {
	return &UserService{
		repo: repo,
		file: fileStorage,
	}
}

type UpdateUserPasswordParams struct {
	UserID      uuid.UUID
	NewPassword string
	OldPassword string
}

func (s *UserService) UpdateUserPassword(ctx context.Context, params UpdateUserPasswordParams) error {
	if params.UserID == uuid.Nil {
		return apperrors.InvalidUserID
	}

	if params.NewPassword == params.OldPassword {
		return apperrors.InvalidPassword
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
		return apperrors.InvalidOldPassword
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
		return apperrors.NoRowsAffected
	}

	return err
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	if userID == uuid.Nil {
		return domain.User{}, apperrors.InvalidUserID
	}
	return s.repo.GetUserByID(ctx, userID)
}
