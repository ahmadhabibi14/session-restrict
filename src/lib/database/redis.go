package database

import (
	"os"
	"session-restrict/src/lib/logger"

	"github.com/go-redis/redis"
)

var ConnRd *redis.Client

func ConnectRedis() {
	ConnRd = redis.NewClient(&redis.Options{
		Addr:     os.Getenv(`REDIS_ADDR`),
		Password: os.Getenv(`REDIS_PASSWORD`),
		DB:       0,
	})
	err := ConnRd.Ping().Err()
	if err != nil {
		logger.Log.Fatal(err, `failed to connect to redis`)
	}
}
