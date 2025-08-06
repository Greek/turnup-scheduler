package scheduler

import (
	"context"
	"log"
	"strings"
	"time"

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
	currDate := strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")

	_, err := s.GetSnapshot(currDate, "towson")
	if err == redis.Nil {
		log.Printf("[CheckForInitialSnapshot] Snapshot not found for snapshot:%s:%s. Creating new snapshot", currDate, "towson")
		_, err = s.CreateSnapshot(currDate, "towson", CreateSnapshotOpts{
			Overwrite: true,
		})
		if err != nil {
			return false, err.Error()
		}
	} else {
		log.Printf("[CheckForInitialSnapshot] Snapshot found for %s:%s Skipping.", currDate, "towson")
	}

	return true, ""
}
