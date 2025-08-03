package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/datasource/events"
	"turnup-scheduler/internal/datasource/involved"
	"turnup-scheduler/internal/mappers"

	"github.com/redis/go-redis/v9"
)

type CreateSnapshotOpts struct {
	Overwrite bool
}

func (s Scheduler) GetSnapshot(date string, namespace string) ([]events.StandardEvent, error) {
	key := constants.KEY_LATEST_SNAPSHOT + date + ":" + namespace

	val, err := s.Redis.Get(s.Ctx, key).Result()
	if err != nil {
		return nil, redis.Nil
	}

	var events []events.StandardEvent
	if err := json.Unmarshal([]byte(val), &events); err != nil {
		fmt.Printf("[GetSnapshot] error unmarshaling snapshot: %v\n", err)
		return nil, err
	}

	return events, nil
}

func (s Scheduler) CreateSnapshot(date string, namespace string, opts CreateSnapshotOpts) (string, error) {
	log.Printf("[CreateSnapshot] Creating...")
	key := constants.KEY_LATEST_SNAPSHOT + date + ":" + namespace

	val, err := s.Redis.Get(s.Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[CreateSnapshot]\tKey not found, continuing.")
	}
	if val != "" && !opts.Overwrite {
		return "", errors.New("[CreateSnapshot] key already exists")
	}

	log.Printf("[CreateSnapshot] Getting events...")
	events, err := involved.GetAllEvents(involved.GetAllEventsOpts{Take: 50})
	if err != nil {
		log.Printf("[CreateSnapshot] Failed to get events %+v", err)
		return "", err
	}

	mappedEvents := mappers.MapInvolvedEventsToStdEvent(events.Value)
	mappedEventsBytes, err := json.Marshal(mappedEvents)
	if err != nil {
		log.Printf("[CreateSnapshot] Failed to marshal events %+v", err)
		return "", err
	}

	s.Redis.Del(s.Ctx, key)
	val, err = s.Redis.SetEx(s.Ctx, key, mappedEventsBytes, constants.FIFTEEN_MINUTES).Result()
	if err != nil {
		log.Fatalf("[CreateSnapshot] Failed to create snapshot %v", err)
		return "", err
	}
	log.Printf("[CreateSnapshot] Created new snapshot %s", key)

	return val, err
}
