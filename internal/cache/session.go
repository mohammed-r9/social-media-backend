package cache

import (
	"context"
	"log/slog"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo"
	"time"
)

type CachedSessionRepo struct {
	next   repo.SessionsRepository
	cache  Cache
	logger *slog.Logger
}

func NewCachedSessionRepo(next repo.SessionsRepository, cache Cache, logger *slog.Logger) *CachedSessionRepo {
	return &CachedSessionRepo{
		next:  next,
		cache: cache,
		logger: logger.With(
			"cache", "session-cache",
		),
	}
}

var _ repo.SessionsRepository = (*CachedSessionRepo)(nil)

func (c *CachedSessionRepo) CreateSession(ctx context.Context, params domain.CreateSessionParams) (domain.Session, error) {
	return c.next.CreateSession(ctx, params)
}

func (c *CachedSessionRepo) GetUserSession(ctx context.Context, sessionID string) (repo.UserSessionDTO, error) {
	var session domain.Session
	var user domain.User

	// session
	err := c.cache.Get(ctx, keySessionByID(sessionID), &session)
	if err != nil && err != apperrors.CacheMiss {
		c.logger.Warn("session cache read failed", "err", err, "session_id", sessionID)
	}

	if err == apperrors.CacheMiss {
		dto, err := c.next.GetUserSession(ctx, sessionID)
		if err != nil {
			c.logger.Error("db fallback failed", "err", err, "session_id", sessionID)
			return repo.UserSessionDTO{}, err
		}

		_ = c.cache.Set(ctx, keySessionByID(sessionID), dto.Session, 24*time.Hour)
		return dto, nil
	}

	// user
	err = c.cache.Get(ctx, keyUserByID(session.UserID.String()), &user)
	if err != nil && err != apperrors.CacheMiss {
		c.logger.Warn("user cache read failed", "err", err, "user_id", session.UserID)
	}

	if err == apperrors.CacheMiss {
		userSession, err := c.next.GetUserSession(ctx, sessionID)
		if err != nil {
			c.logger.Error("user db fallback failed", "err", err, "user_id", session.UserID)
			return repo.UserSessionDTO{}, err
		}
		return userSession, nil
	}

	return repo.UserSessionDTO{
		Session: session,
		User: struct {
			VerifiedAt *time.Time
		}{
			VerifiedAt: user.VerifiedAt,
		},
	}, nil
}
