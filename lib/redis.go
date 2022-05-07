package lib

import (
	redis "github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Cmd *redis.Client
}

var redisDB *redis.Client

func GetRedis() *redis.Client {
	if redisDB != nil {
		return redisDB
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		PoolSize: 10,
	})

	return rdb
}
