package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/auth"
	"social-media-backend/internal/utils"

	"github.com/redis/go-redis/v9"
)

type TokenRepo interface {
	StoreToken(ctx context.Context, params auth.StoreTokenParam) error
	GetToken(ctx context.Context, key string) (StoredToken, error)
	DeleteToken(ctx context.Context, key string) error
}

var _ TokenRepo = (*RedisRepo)(nil)

func (r *RedisRepo) StoreToken(ctx context.Context, params auth.StoreTokenParam) error {
	tokenBytes, err := json.Marshal(params.Token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	key := utils.RedisTokenKeyBuilder(params.Token.Hash, params.Token.Scope)
	ttl := params.Token.Ttl

	if err := r.c.Set(ctx, key, tokenBytes, ttl).Err(); err != nil {
		return fmt.Errorf("redis set failed for key %s: %w", key, err)
	}

	return nil
}

// the key should be formatted as "prefix:token_hash"
func (r *RedisRepo) GetToken(ctx context.Context, key string) (StoredToken, error) {
	var token StoredToken

	data, err := r.c.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return token, apperrors.TokenNotFound
		}
		return token, fmt.Errorf("redis get failed: %w", err)
	}

	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return token, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return token, nil
}

func (r *RedisRepo) DeleteToken(ctx context.Context, key string) error {
	res, err := r.c.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis delete failed: %w", err)
	}

	if res == 0 {
		return apperrors.TokenNotFound
	}

	return nil
}
