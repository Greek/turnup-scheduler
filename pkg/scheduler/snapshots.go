package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/datasource/events"
	"turnup-scheduler/internal/datasource/eventsattu"
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

	allEvents := []events.StandardEvent{}

	log.Printf("[CreateSnapshot] Getting events...")
	log.Print("[CreateSnapshot] getting involved events")
	involvedEvents, err := involved.GetAllEvents(involved.GetAllEventsOpts{Take: 50})
	if err != nil {
		log.Printf("[CreateSnapshot] Failed to get events %+v", err)
		return "", err
	}
	mappedInvolvedEvents := mappers.MapInvolvedEventsToStdEvent(involvedEvents.Value)

	log.Print("[CreateSnapshot] getting events@tu events")
	eventsAtTUEvents, err := eventsattu.GetAllEvents(eventsattu.GetAllEventsOpts{Take: 50})
	if err != nil {
		log.Printf("[CreateSnapshot] Failed to get events %+v", err)
		return "", err
	}
	mappedEventsAtTUEvents := mappers.MapEventAtTUEventsToStdEvent(eventsAtTUEvents)

	allEvents = append(allEvents, mappedInvolvedEvents...)
	allEvents = append(allEvents, mappedEventsAtTUEvents...)

	marshaledEvents, err := json.Marshal(allEvents)
	if err != nil {
		log.Fatalf("[CreateSnapshot] Failed to marshal events: %v", err)
		return "", err
	}

	s.Redis.Del(s.Ctx, key)
	val, err = s.Redis.SetEx(s.Ctx, key, marshaledEvents, constants.FIFTEEN_MINUTES).Result()
	if err != nil {
		log.Fatalf("[CreateSnapshot] Failed to create snapshot %v", err)
		return "", err
	}
	log.Printf("[CreateSnapshot] Created new snapshot %s", key)

	return val, err
}
