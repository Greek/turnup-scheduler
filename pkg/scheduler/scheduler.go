package scheduler

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Scheduler struct {
	Ctx      context.Context
	Redis    redis.Client
	CurrDate string
}

func CreateScheduler(ctx context.Context, redisClient redis.Client) *Scheduler {
	currTime := strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")
	return &Scheduler{
		Ctx:      ctx,
		Redis:    redisClient,
		CurrDate: currTime,
	}
}

// CheckForInitialSnapshot
func (s Scheduler) CheckForInitialSnapshot() (bool, string) {
	_, err := s.GetSnapshot(s.CurrDate, "towson")
	if err == redis.Nil {
		log.Printf("[CheckForInitialSnapshot] Snapshot not found for snapshot:%s:%s. Creating new snapshot", s.CurrDate, "towson")
		_, err = s.CreateSnapshot(s.CurrDate, "towson", CreateSnapshotOpts{
			Overwrite: true,
		})
		if err != nil {
			return false, err.Error()
		}
	}

	return true, ""
}
