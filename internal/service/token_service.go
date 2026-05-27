package service

import (
	"context"
	"social-media-backend/internal/crypto/tokens"
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
) (tokens.ShortLivedToken, error) {
	if userID == uuid.Nil {
		return tokens.ShortLivedToken{}, domain.ErrInvalidUserID
	}

	token := tokens.GenerateEmailVerificationToken(userID)
	err := s.repo.StoreToken(ctx, tokens.StoreTokenParam{
		Token: token,
	})

	if err != nil {
		return tokens.ShortLivedToken{}, err
	}

	return token, nil
}

func (s *TokenService) GenerateAndStorePasswordResetToken(
	ctx context.Context,
	userID uuid.UUID,
) (tokens.ShortLivedToken, error) {

	if userID == uuid.Nil {
		return tokens.ShortLivedToken{}, domain.ErrInvalidUserID
	}

	token := tokens.GeneratePasswordResetToken(userID)
	err := s.repo.StoreToken(ctx, tokens.StoreTokenParam{
		Token: token,
	})

	if err != nil {
		return tokens.ShortLivedToken{}, err
	}

	return token, nil
}

type VerifyTokenParams struct {
	TokenPlainText string
	Scope          tokens.TokenScope
}

// VerifyToken verifies the token and returns a uuid of the owner if valid
func (s *TokenService) VerifyToken(ctx context.Context, params VerifyTokenParams) (uuid.UUID, error) {
	hash := tokens.HashToken(params.TokenPlainText)
	key := tokens.RedisKeyBuilder(hash, params.Scope)
	token, err := s.repo.GetToken(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}

	if token.UserID == uuid.Nil {
		return uuid.Nil, tokens.ErrInvalidToken
	}

	if token.Scope != params.Scope {
		return uuid.Nil, tokens.ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return uuid.Nil, tokens.ErrExpiredToken
	}

	return token.UserID, nil
}
