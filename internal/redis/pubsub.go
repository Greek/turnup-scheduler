package redis

import (
	"context"
	"log"
	"strings"
	"time"
	"turnup-scheduler/pkg/scheduler"

	"github.com/redis/go-redis/v9"
)

func CreatePubsubListener(ctx context.Context, r redis.Client, sch *scheduler.Scheduler) {
	go func() {
		// Listen to all TTL expiry events
		pbs := r.Subscribe(ctx, "__keyevent@0__:expired")
		defer pbs.Close()

		ch := pbs.Channel()
		log.Print("[Pubsub] Listening for expired keys...")
		for msg := range ch {
			payloadSplit := strings.Split(msg.Payload, ":")
			currDate := strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")

			if payloadSplit[0] == "snapshot" || payloadSplit[0] == "timer" {
				log.Printf("[Pubsub] Snapshot %s:%s expired, creating new set of events.", payloadSplit[1], payloadSplit[2])
				sch.CreateSnapshot(currDate, payloadSplit[2], scheduler.CreateSnapshotOpts{Overwrite: true})
			}
		}
	}()
}
