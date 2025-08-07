package scheduler

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"turnup-scheduler/internal/logging"

	"github.com/redis/go-redis/v9"
)

type Scheduler struct {
	Ctx   context.Context
	Redis redis.Client
}

func CreateScheduler(ctx context.Context, redisClient redis.Client) *Scheduler {
	return &Scheduler{
		Ctx:   ctx,
		Redis: redisClient,
	}
}

// CheckForInitialSnapshot
func (s Scheduler) CheckForInitialSnapshot() (bool, string) {
	logger := logging.BuildLogger("CheckForInitialSnapshot")
	currDate := strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")
	namespace := "towson"

	_, err := s.GetSnapshot(currDate, namespace)
	if err == redis.Nil {
		logger.Info("Snapshot not found. Creating new snapshot", slog.String("currDate", currDate), slog.String("namespace", namespace))
		_, err = s.CreateSnapshot(currDate, namespace, CreateSnapshotOpts{
			Overwrite: true,
		})
		if err != nil {
			return false, err.Error()
		}
	} else {
		logger.Info("Snapshot found", slog.String("currDate", currDate), slog.String("namespace", namespace))
	}

	return true, ""
}
