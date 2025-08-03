package cmd

import (
	"context"
	"turnup-scheduler/internal/env"
	"turnup-scheduler/internal/redis"
	"turnup-scheduler/pkg/scheduler"
)

var ctx = context.Background()

// Init starts the application by creating Redis connections
// and configurations.
func Init() {
	env.LoadEnv()
	env.CheckEnv()

	redisClient := redis.InitRedis(ctx)
	sch := scheduler.CreateScheduler(ctx, redisClient)

	redis.CreatePubsubListener(ctx, redisClient, sch)
	sch.CheckForInitialSnapshot()
	select {}
}
