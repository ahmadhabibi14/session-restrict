package database

import (
	"os"

	"github.com/go-redis/redis"
)

var ConnRd *redis.Client

func ConnectRedis() error {
	ConnRd = redis.NewClient(&redis.Options{
		Addr: os.Getenv(`REDIS_ADDR`),
		DB:   0,
	})
	return ConnRd.Ping().Err()
}
