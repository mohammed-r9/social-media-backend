package redis

import (
	"context"
	"encoding/json"
	"social-media-backend/internal/crypto/tokens/stateful"
	"social-media-backend/internal/domain"
)

type TokenRepo interface {
	StoreToken(ctx context.Context, params domain.StoreTokenParam) error
	GetToken(ctx context.Context, key string) (stateful.ShortLivedToken, error)
}

func (r *RedisRepo) StoreToken(ctx context.Context, params domain.StoreTokenParam) error {
	token, err := json.Marshal(params.Token)
	if err != nil {
		return err
	}

	key := stateful.RedisKeyBuilder(params.Token.Hash, params.Token.Scope)
	cmd := r.c.Set(ctx, key, token, params.Token.Ttl)

	return cmd.Err()
}

func (r *RedisRepo) GetToken(ctx context.Context, key string) (stateful.ShortLivedToken, error) {
	var token stateful.ShortLivedToken

	cmd := r.c.Get(ctx, key)
	if cmd.Err() != nil {
		return token, cmd.Err()
	}

	data, err := cmd.Result()
	if err != nil {
		return token, err
	}

	if data == "" {
		return token, domain.ErrTokenNotFound
	}

	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return token, err
	}

	return token, nil
}
