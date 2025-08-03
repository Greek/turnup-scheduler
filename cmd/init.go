package cmd

import (
	"context"
	"time"
	"turnup-scheduler/internal/env"
	"turnup-scheduler/internal/redis"
	"turnup-scheduler/pkg/scheduler"
)

var ctx = context.Background()

// Init starts the application by creating Redis connections
// and configurations.
func Init() {
	env.CheckEnv()
	redisClient := redis.InitRedis(ctx)

	sch := scheduler.CreateScheduler(ctx, redisClient)

	redis.CreatePubsubListener(ctx, redisClient, sch)
	sch.CheckForInitialSnapshot()
	sch.Redis.SetEx(sch.Ctx, "hello", "world", 1*time.Second)
	select {}
}
