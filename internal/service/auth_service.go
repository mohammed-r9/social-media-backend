package service

import (
	"context"
	"social-media-backend/internal/crypto/password"
	"social-media-backend/internal/crypto/tokens/stateful"
	"social-media-backend/internal/crypto/tokens/stateless"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo/postgres"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    postgres.UserRepository
	sessionRepo postgres.SessionsRepository
}

func NewAuthService(userRepo postgres.UserRepository, sessionRepo postgres.SessionsRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type LoginParams struct {
	Email      string
	Password   string
	DeviceName string
}

func (s *AuthService) Login(ctx context.Context, params LoginParams) (domain.AuthTokens, error) {
	if params.Email == "" || params.Password == "" {
		return domain.AuthTokens{}, domain.ErrInvalidCredentials
	}

	user, err := s.userRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return domain.AuthTokens{}, err
	}

	isValid, err := password.ComparePassword(password.ComparePasswordParams{
		Password:   params.Password,
		StoredHash: user.PasswordHash,
	})

	// TODO: ComparePassword should validate the hash format as well.
	if err != nil {
		return domain.AuthTokens{}, domain.ErrInvalidCredentials
	}

	if !isValid {
		return domain.AuthTokens{}, domain.ErrInvalidCredentials
	}

	tokens := stateful.GenerateSessionTokens()
	storedHashes := tokens.ToHash()
	sessionID := domain.GenerateSessionID()

	if err != nil {
		return domain.AuthTokens{}, err
	}

	_, err = s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:               sessionID,
		RefreshTokenHash: storedHashes.RefreshHash,
		CsrfTokenHash:    storedHashes.CsrfHash,
		DeviceName:       "temp",
	})

	if err != nil {
		return domain.AuthTokens{}, err
	}

	accessToken, err := stateless.GenerateAccessToken(user)
	if err != nil {
		return domain.AuthTokens{}, err
	}

	return domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: tokens.RefreshToken,
		CsrfToken:    tokens.CsrfToken,
		SessionID:    sessionID,
	}, nil
}

type RefreshParams struct {
	RefreshToken string
	CsrfToken    string
	SessionID    string
	UserID       uuid.UUID
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, params RefreshParams) (string, error) {
	session, err := s.sessionRepo.GetUserSession(ctx, domain.GetUserSessionParams{
		ID:     params.SessionID,
		UserID: params.UserID,
	})

	if err != nil {
		return "", err
	}

	if time.Now().After(session.Session.ExpiresAt) {
		return "", domain.ErrSessionExpired
	}

	if session.Session.RevokedAt != nil {
		return "", domain.ErrSessionRevoked
	}

	if err != nil {
		return "", err
	}

	isValid := stateful.CompareTokenToHash(stateful.CompareTokenToHashParams{
		PlainText:  params.RefreshToken,
		StoredHash: session.Session.RefreshTokenHash,
	})
	if !isValid {
		return "", domain.ErrInvalidToken
	}

	isValid = stateful.CompareTokenToHash(stateful.CompareTokenToHashParams{
		PlainText:  params.CsrfToken,
		StoredHash: session.Session.CsrfTokenHash,
	})

	if !isValid {
		return "", domain.ErrInvalidToken
	}

	return stateless.GenerateAccessToken(domain.User{
		ID:         params.UserID,
		VerifiedAt: session.VerifiedAt,
	})
}
