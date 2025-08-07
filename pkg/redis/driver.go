package redis

import (
	"context"
	"log/slog"
	"os"

	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/env"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context) redis.Client {
	log := logging.BuildLogger("InitRedis")
	env.CheckEnv()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Error("Failed to parse Redis URL", slog.Any("err", err))
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Error("Failed to connect to Redis", slog.Any("err", err))
	}
	log.Info("Connected.")

	return *rdb
}
