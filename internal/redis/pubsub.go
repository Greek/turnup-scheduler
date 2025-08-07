package redis

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/scheduler"

	"github.com/redis/go-redis/v9"
)

func CreatePubsubListener(ctx context.Context, r redis.Client, sch *scheduler.Scheduler) {
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
				key := payloadSplit[0] + payloadSplit[1] + ":" + payloadSplit[2]
				log.Info("Snapshot "+key+" expired, creating new set of events.", slog.String("key", key))
				sch.CreateSnapshot(currDate, payloadSplit[2], scheduler.CreateSnapshotOpts{Overwrite: true})
			}
		}
	}()
}
