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

func (s *Scheduler) CreatePubsubListener(ctx context.Context, r redis.Client) {
	log := logging.BuildLogger("CreatePubsubListener")

	go func() {
		// Listen to all TTL expiry events
		pbs := r.Subscribe(ctx, "__keyevent@0__:expired")
		defer pbs.Close()

		ch := pbs.Channel()
		log.Info("Listening for expired keys...")
		for msg := range ch {
			payloadSplit := strings.Split(msg.Payload, ":")
			currDate := strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")

			if payloadSplit[0] == "snapshot" || payloadSplit[0] == "timer" {
				log.Info("Snapshot "+msg.Payload+" expired, creating new set of events.", slog.String("key", msg.Payload))
				s.CreateSnapshot(currDate, payloadSplit[2], CreateSnapshotOpts{Overwrite: true})
			}
		}
	}()
}

// CheckForInitialSnapshot
func (s *Scheduler) CheckForInitialSnapshot() (bool, error) {
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
			return false, err
		}
	}

	return true, nil
}
