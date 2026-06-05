package cache

import (
	"context"
	"log/slog"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/repo"
	"time"

	"github.com/google/uuid"
)

type CachedUserRepo struct {
	next   repo.UserRepository
	cache  Cache
	logger *slog.Logger
}

func NewChachedUserRepo(next repo.UserRepository, cache Cache, logger *slog.Logger) *CachedUserRepo {
	return &CachedUserRepo{
		next:  next,
		cache: cache,
		logger: logger.With(
			"cache", "user-cache",
		),
	}
}

var _ repo.UserRepository = (*CachedUserRepo)(nil)

func (c *CachedUserRepo) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	return c.next.CreateUser(ctx, params)
}

func (c *CachedUserRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return c.next.GetUserByEmail(ctx, email)
}

func (c *CachedUserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var user domain.User

	err := c.cache.Get(ctx, keyUserByID(userID.String()), &user)
	if err == nil {
		return user, nil
	}

	user, err = c.next.GetUserByID(ctx, userID)
	if err != nil {
		c.logger.Error("db fallback failed", "err", err, "user_id", userID)
		return domain.User{}, err
	}

	if err := c.cache.Set(ctx, keyUserByID(userID.String()), user, time.Hour); err != nil {
		c.logger.Error("cache set failed", "err", err, "user_id", userID)
	}

	return user, nil
}

func (c *CachedUserRepo) UpdateUserPassword(ctx context.Context, params domain.UpdatePasswordParams) (uuid.UUID, error) {
	if err := c.cache.Delete(ctx, keyUserByID(params.ID.String())); err != nil {
		c.logger.Error("cache delete failed", "err", err, "user_id", params.ID)
	}

	userID, err := c.next.UpdateUserPassword(ctx, params)
	if err != nil {
		c.logger.Error("update password failed", "err", err, "user_id", params.ID)
		return uuid.Nil, err
	}

	return userID, nil
}

func (c *CachedUserRepo) VerifyUserEmail(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	user, err := c.next.VerifyUserEmail(ctx, userID)
	if err != nil {
		c.logger.Error("verify email failed", "err", err, "user_id", userID)
		return uuid.Nil, err
	}

	if err := c.cache.Delete(ctx, keyUserByID(userID.String())); err != nil {
		c.logger.Error("cache delete failed", "err", err, "user_id", userID)
	}

	return user, nil
}

func (c *CachedUserRepo) UpdateSelfProfile(ctx context.Context, params domain.UpdateProfileParams) (domain.Profile, error) {
	if err := c.cache.Delete(ctx, keyUserByID(params.UserID.String())); err != nil {
		c.logger.Error("cache delete failed", "err", err, "user_id", params.UserID)
	}

	profile, err := c.next.UpdateSelfProfile(ctx, params)
	if err != nil {
		c.logger.Error("update profile failed", "err", err, "user_id", params.UserID)
		return domain.Profile{}, err
	}

	return profile, nil
}
