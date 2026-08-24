package cmd

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"turnup-scheduler/internal/env"
	"turnup-scheduler/internal/lib/redis"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/scheduler"
	"turnup-scheduler/pkg/server"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// Init starts the application by creating Redis connections
// and configurations.
func Init() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	sch.CreatePubsubListener(ctx, redisClient)
	if err := server.Run(ctx, *port, sch); err != nil {
		log.Error("Server stopped with error", slog.Any("err", err))
		os.Exit(1)
	}
}
