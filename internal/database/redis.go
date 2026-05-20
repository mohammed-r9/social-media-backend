package database

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "test",
		DB:       0,
	})

	return rdb
}
