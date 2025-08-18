package redis

import (
	"context"
	realLog "log"
	"log/slog"
	"os"

	"turnup-scheduler/internal/env"
	"turnup-scheduler/internal/logging"

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
		realLog.Fatalf("Failed to connect to Redis %+v", err)
	}
	log.Info("Connected.")

	return *rdb
}
