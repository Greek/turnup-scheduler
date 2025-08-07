package cmd

import (
	"context"
	"turnup-scheduler/pkg/env"
	"turnup-scheduler/pkg/redis"
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

	sch.CheckForInitialSnapshot()
	redis.CreatePubsubListener(ctx, redisClient, sch)
	select {}
}
