package service

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"social-media-backend/internal/adapters/storage"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo"

	"github.com/google/uuid"
)

type UserService struct {
	repo        repo.UserRepository
	fileStorage storage.Storage
}

func NewUserService(repo repo.UserRepository, fileStorage storage.Storage) *UserService {
	return &UserService{
		repo:        repo,
		fileStorage: fileStorage,
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
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, apperrors.InvalidUserID
	}

	if user.Profile.AvatarKey != "" {
		url, err := s.fileStorage.GetURL(ctx, user.Profile.AvatarKey)
		if err != nil {
			return domain.User{}, err
		}
		user.UpdateAvatarUrl(url)
	}

	return user, nil
}

func (s *UserService) UpdateUserAvatar(
	ctx context.Context,
	userID uuid.UUID,
	imageFile *multipart.FileHeader,
	contentType string) (string, error) {
	if userID == uuid.Nil {
		return "", apperrors.InvalidUserID
	}

	file, err := imageFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var img image.Image

	switch contentType {
	case storage.ContentTypeJPEG.String():
		img, err = jpeg.Decode(file)
	case storage.ContentTypePNG.String():
		img, err = png.Decode(file)
	default:
		return "", apperrors.InvalidMime
	}

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: 85,
	})
	if err != nil {
		return "", err
	}
	objKey := storage.GenereateObjectKey()

	err = s.fileStorage.Upload(ctx, objKey, &buf, storage.ContentTypeJPEG)
	if err != nil {
		return "", err
	}
	err = s.repo.UpdateUserAvatar(ctx, userID, objKey)
	if err != nil {
		s.fileStorage.Delete(ctx, objKey)
		return "", err
	}

	// not a big deal if it fails
	imgUrl, _ := s.fileStorage.GetURL(ctx, objKey)

	return imgUrl, nil
}
