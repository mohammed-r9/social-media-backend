package service

import (
	"context"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/domain"
	rdrepo "social-media-backend/internal/repo/redis"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	repo rdrepo.TokenRepo
}

func NewTokenService(repo rdrepo.TokenRepo) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) GenerateAndStoreEmailVerificationToken(
	ctx context.Context,
	userID uuid.UUID,
) (auth.ShortLivedToken, error) {
	if userID == uuid.Nil {
		return auth.ShortLivedToken{}, domain.ErrInvalidUserID
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
		return auth.ShortLivedToken{}, domain.ErrInvalidUserID
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
	key := auth.RedisKeyBuilder(hash, params.Scope)
	token, err := s.repo.GetToken(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}

	if token.UserID == uuid.Nil {
		return uuid.Nil, auth.ErrInvalidToken
	}

	if token.Scope != params.Scope {
		return uuid.Nil, auth.ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return uuid.Nil, auth.ErrExpiredToken
	}

	return token.UserID, nil
}
