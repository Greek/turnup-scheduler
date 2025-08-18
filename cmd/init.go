package cmd

import (
	"context"
	"flag"
	"log/slog"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/env"
	"turnup-scheduler/pkg/redis"
	"turnup-scheduler/pkg/scheduler"
	"turnup-scheduler/pkg/server"
)

var ctx = context.Background()
var (
	port = flag.Int("port", 50051, "The server port")
)

// Init starts the application by creating Redis connections
// and configurations.
func Init() {
	log := logging.BuildLogger("Init")
	flag.Parse()

	log.Info("Initializing...")
	env.LoadEnv()
	env.CheckEnv()

	redisClient := redis.InitRedis(ctx)
	sch := scheduler.CreateScheduler(ctx, redisClient)

	_, err := sch.CheckForInitialSnapshot()
	if err != nil {
		log.Error("Failed to execute CheckForInitialSnapshot", slog.Any("err", err))
	}

	redis.CreatePubsubListener(ctx, redisClient, sch)
	server.InitializeGrpcServer(*port, sch)
	select {}
}
