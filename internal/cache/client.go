package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	c *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &RedisCache{
		c: c,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.c.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest any) error {
	data, err := r.c.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(data, dest)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.c.Del(ctx, key).Err()
}
