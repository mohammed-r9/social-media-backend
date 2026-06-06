package cache

import (
	"context"
	"log/slog"
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
	var session repo.UserSessionDTO

	err := c.cache.Get(ctx, keySessionByID(sessionID), &session)
	if err == nil {
		return session, nil
	}

	if err != ErrCacheMiss {
		c.logger.Warn("cache read failed", "err", err, "session_id", sessionID)
	}

	session, err = c.next.GetUserSession(ctx, sessionID)
	if err != nil {
		c.logger.Error("db fallback failed", "err", err, "session_id", sessionID)
		return repo.UserSessionDTO{}, err
	}

	if !shouldCacheSession(session) {
		return session, nil
	}

	if err := c.cache.Set(ctx, keySessionByID(sessionID), session, time.Hour*24); err != nil {
		c.logger.Warn("cache set failed", "err", err, "session_id", sessionID)
	}

	return session, nil
}
