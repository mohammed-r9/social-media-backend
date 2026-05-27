package service

import (
	"context"
	"log/slog"
	"social-media-backend/internal/adapters/mailer"
	"social-media-backend/internal/crypto/password"
	"social-media-backend/internal/crypto/tokens/stateful"
	"social-media-backend/internal/crypto/tokens/stateless"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo/postgres"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo     postgres.UserRepository
	sessionRepo  postgres.SessionsRepository
	tokenService TokenService
	mailer       mailer.EmailService
	logger       *slog.Logger
}

func NewAuthService(userRepo postgres.UserRepository,
	sessionRepo postgres.SessionsRepository,
	tokenService TokenService, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
		logger: logger.With(
			"service", "user-service",
		),
	}
}

type LoginParams struct {
	Email      string
	Password   string
	DeviceName string
}

func (s *AuthService) Register(ctx context.Context, params RegisterParams) (domain.User, error) {
	log := s.logger.With("op", "create_user")

	if params.PassowrdPlainText == "" {
		log.Warn("validation failed", "field", "password")
		return domain.User{}, domain.ErrInvalidPassword
	}

	if params.Email == "" {
		log.Warn("validation failed", "field", "email")
		return domain.User{}, domain.ErrInvalidUserEmail
	}

	if params.Name == "" {
		log.Warn("validation failed", "field", "name")
		return domain.User{}, domain.ErrInvalidUserName
	}
	if params.Username == "" {
		log.Warn("validation failed", "field", "username")
		return domain.User{}, domain.ErrInvalidUsername
	}

	passowrdHash, err := password.HashPassword(params.PassowrdPlainText)
	if err != nil {
		log.Error("password hashing failed", "err", err)
		return domain.User{}, err
	}

	userID, err := uuid.NewV7()
	if err != nil {
		log.Error("generating uuid failed", "err", err)
		return domain.User{}, err
	}
	user, err := s.userRepo.CreateUser(ctx, domain.CreateUserParams{
		ID:           userID,
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: passowrdHash,
		Username:     params.Username,
	})
	if err != nil {
		log.Error("user creation failed", "err", err)
		return domain.User{}, err
	}
	token, err := s.tokenService.GenerateAndStoreEmailVerificationToken(ctx, user.ID)
	if err != nil {
		log.Error("email verification token generation failed", "err", err)
		return domain.User{}, err
	}

	err = s.mailer.SendVerificationEmail(user.Email, token.Raw)
	if err != nil {
		log.Error("sending email verification failed", "err", err)
		return domain.User{}, err
	}

	return user, nil
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

	_, err = s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:               sessionID,
		UserID:           user.ID,
		RefreshTokenHash: storedHashes.RefreshHash,
		CsrfTokenHash:    storedHashes.CsrfHash,
		DeviceName:       "temp",
		ExpiresAt:        time.Now().Add(stateful.REFRESH_TTL),
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
	CsrfToken    *string
	SessionID    string
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, params RefreshParams) (string, error) {
	dto, err := s.sessionRepo.GetUserSession(ctx, params.SessionID)

	if err != nil {
		return "", err
	}

	if time.Now().After(dto.Session.ExpiresAt) {
		return "", domain.ErrSessionExpired
	}

	if dto.Session.RevokedAt != nil {
		return "", domain.ErrSessionRevoked
	}

	isValid := stateful.CompareTokenToHash(stateful.CompareTokenToHashParams{
		PlainText:  params.RefreshToken,
		StoredHash: dto.Session.RefreshTokenHash,
	})
	if !isValid {
		return "", domain.ErrTokenMismatch
	}

	if params.CsrfToken != nil {
		isValid = stateful.CompareTokenToHash(stateful.CompareTokenToHashParams{
			PlainText:  *params.CsrfToken,
			StoredHash: dto.Session.CsrfTokenHash,
		})

		if !isValid {
			return "", domain.ErrTokenMismatch
		}
	}

	return stateless.GenerateAccessToken(domain.User{
		ID:         dto.Session.UserID,
		VerifiedAt: dto.User.VerifiedAt,
	})
}

func (s *AuthService) VerifyUserEmail(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrInvalidToken
	}

	userID, err := s.tokenService.VerifyToken(ctx, VerifyTokenParams{
		TokenPlainText: token,
		Scope:          stateful.ScopeEmailVerification,
	})
	if err != nil {
		return err
	}
	err = s.userRepo.VerifyUserEmail(ctx, userID)
	return err
}

type UpdateUserPasswordParams struct {
	UserID      uuid.UUID
	NewPassword string
	OldPassword string
}

func (s *AuthService) UpdateUserPassword(ctx context.Context, params UpdateUserPasswordParams) error {
	if params.UserID == uuid.Nil {
		return domain.ErrInvalidUserID
	}

	if params.NewPassword == params.OldPassword {
		return domain.ErrInvalidPassword
	}

	user, err := s.userRepo.GetUserByID(ctx, params.UserID)
	if err != nil {
		return err
	}

	isMatched, err := password.ComparePassword(password.ComparePasswordParams{
		Password:   params.OldPassword,
		StoredHash: user.PasswordHash,
	})
	if err != nil {
		return err
	}

	if !isMatched {
		return domain.ErrInvalidOldPassword
	}

	newPasswordHash, err := password.HashPassword(params.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateUserPassword(ctx, domain.UpdatePasswordParams{
		ID:           user.ID,
		PasswordHash: newPasswordHash,
	})
}

type ResetUserPasswordParams struct {
	Token       string
	NewPassword string
}

func (s *AuthService) ResetUserPassword(ctx context.Context, params ResetUserPasswordParams) error {
	if params.Token == "" {
		return domain.ErrInvalidToken
	}

	userID, err := s.tokenService.VerifyToken(ctx, VerifyTokenParams{
		TokenPlainText: params.Token,
		Scope:          stateful.ScopePasswordReset,
	})
	if err != nil {
		return err
	}

	hash, err := password.HashPassword(params.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateUserPassword(ctx, domain.UpdatePasswordParams{
		ID:           userID,
		PasswordHash: hash,
	})
}

func (s *AuthService) AskForPasswordReset(ctx context.Context, email string) error {
	if email == "" {
		return domain.ErrInvalidUserEmail
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := s.tokenService.GenerateAndStorePasswordResetToken(ctx, user.ID)
	if err != nil {
		return err
	}
	return s.mailer.SendPasswordResetEmail(email, token.Raw)
}
