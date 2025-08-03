package redis

import (
	"context"
	"log"
	"log/slog"
	"os"

	"turnup-scheduler/internal/env"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context) redis.Client {
	env.CheckEnv()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	slog.Info("Connected to Redis")

	return *rdb
}
