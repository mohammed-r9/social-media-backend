package service

import (
	"context"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/repo"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	repo repo.TokenRepo
}

func NewTokenService(repo repo.TokenRepo) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) GenerateAndStoreEmailVerificationToken(
	ctx context.Context,
	userID uuid.UUID,
) (auth.ShortLivedToken, error) {
	if userID == uuid.Nil {
		return auth.ShortLivedToken{}, apperrors.InvalidUserID
	}

	token := auth.GenerateEmailVerificationToken(userID)
	err := s.repo.StoreToken(ctx, auth.StoreTokenParam{
		Token: token,
	})

	if err != nil {
		return auth.ShortLivedToken{}, err
	}

	return token, nil
}

func (s *TokenService) GenerateAndStorePasswordResetToken(
	ctx context.Context,
	userID uuid.UUID,
) (auth.ShortLivedToken, error) {

	if userID == uuid.Nil {
		return auth.ShortLivedToken{}, apperrors.InvalidUserID
	}

	token := auth.GeneratePasswordResetToken(userID)
	err := s.repo.StoreToken(ctx, auth.StoreTokenParam{
		Token: token,
	})

	if err != nil {
		return auth.ShortLivedToken{}, err
	}

	return token, nil
}

type VerifyTokenParams struct {
	TokenPlainText string
	Scope          auth.TokenScope
}

// VerifyToken verifies the token and returns a uuid of the owner if valid
func (s *TokenService) VerifyToken(ctx context.Context, params VerifyTokenParams) (uuid.UUID, error) {
	hash := auth.HashToken(params.TokenPlainText)
	key := auth.RedisTokenKeyBuilder(hash, params.Scope)
	token, err := s.repo.GetToken(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}

	if token.UserID == uuid.Nil {
		return uuid.Nil, apperrors.InvalidToken
	}

	if token.Scope != params.Scope {
		return uuid.Nil, apperrors.InvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return uuid.Nil, apperrors.ExpiredToken
	}

	return token.UserID, nil
}

type DeleteTokenParams struct {
	TokenPlainText string
	Scope          auth.TokenScope
}

func (s *TokenService) DeleteToken(ctx context.Context, params DeleteTokenParams) error {
	hash := auth.HashToken(params.TokenPlainText)
	key := auth.RedisTokenKeyBuilder(hash, params.Scope)
	return s.repo.DeleteToken(ctx, key)
}
