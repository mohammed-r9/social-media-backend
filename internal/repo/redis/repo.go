package rdrepo

import (
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	c *redis.Client
}

func NewRedisRepository(c *redis.Client) *RedisRepo {
	return &RedisRepo{
		c: c,
	}
}
