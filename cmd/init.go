package cmd

import (
	"context"
	"log/slog"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/env"
	"turnup-scheduler/pkg/redis"
	"turnup-scheduler/pkg/scheduler"
)

var ctx = context.Background()

// Init starts the application by creating Redis connections
// and configurations.
func Init() {
	log := logging.BuildLogger("Init")
	log.Info("Initializing...")
	env.LoadEnv()
	env.CheckEnv()

	redisClient := redis.InitRedis(ctx)
	sch := scheduler.CreateScheduler(ctx, redisClient)

	_, err := sch.CheckForInitialSnapshot()
	if err != nil {
		log.Error("Failed to execute CheckForInitialSnapshot", slog.Any("err", err))
	}

	log.Info("Finished initializing.")
	redis.CreatePubsubListener(ctx, redisClient, sch)
	select {}
}
