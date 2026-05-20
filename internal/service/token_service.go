package service

import (
	"context"
	"social-media-backend/internal/crypto/tokens/stateful"
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

func (s *TokenService) GenerateAndStoreEmailVerificationToken(ctx context.Context,
	userID uuid.UUID) (stateful.ShortLivedToken, error) {
	if userID == uuid.Nil {
		return stateful.ShortLivedToken{}, domain.ErrInvalidUserID
	}

	token := stateful.GenerateEmailVerificationToken(userID)
	err := s.repo.StoreToken(ctx, domain.StoreTokenParam{
		Token: token,
	})

	if err != nil {
		return stateful.ShortLivedToken{}, err
	}

	return token, nil
}

type VerifyTokenParams struct {
	TokenPlainText string
	Scope          stateful.TokenScope
}

// VerifyToken verifies the token and returns a uuid of the owner if valid
func (s *TokenService) VerifyToken(ctx context.Context, params VerifyTokenParams) (uuid.UUID, error) {
	hash := stateful.HashToken(params.TokenPlainText)
	key := stateful.RedisKeyBuilder(hash, params.Scope)
	token, err := s.repo.GetToken(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}

	if token.UserID == uuid.Nil {
		return uuid.Nil, domain.ErrInvalidToken
	}

	if token.Scope != params.Scope {
		return uuid.Nil, domain.ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return uuid.Nil, domain.ErrExpiredToken
	}

	return token.UserID, nil
}
