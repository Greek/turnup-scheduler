package redis

import (
	"context"
	"log"
	"os"

	"turnup-scheduler/internal/env"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context) redis.Client {
	env.CheckEnv()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("[Redis] Failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("[Redis] Failed to connect to Redis: %v", err)
	}
	log.Print("[Redis] Connected.")

	return *rdb
}
