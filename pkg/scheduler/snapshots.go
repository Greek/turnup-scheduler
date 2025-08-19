package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/mappers"
	"turnup-scheduler/pkg/repo/events"
	"turnup-scheduler/pkg/repo/eventsattu"
	"turnup-scheduler/pkg/repo/involved"

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
	log := logging.BuildLogger("CreateSnapshot")
	log.Info("Creating...")
	key := constants.KEY_LATEST_SNAPSHOT + date + ":" + namespace

	val, err := s.Redis.Get(s.Ctx, key).Result()
	if err == redis.Nil {
		log.Info("Key not found, continuing.")
	}
	if val != "" && !opts.Overwrite {
		return "", errors.New("key already exists")
	}

	allEvents := []events.StandardEvent{}

	log.Info("Getting events...")
	log.Info("getting involved events")
	involvedEvents, err := involved.GetAllEvents(involved.GetAllEventsOpts{Take: 50})
	if err != nil {
		log.Info("Failed to get events", slog.Any("err", err))
		return "", err
	}
	mappedInvolvedEvents := mappers.MapInvolvedEventsToStdEvent(involvedEvents.Value)

	log.Info("getting events@tu events")
	eventsAtTUEvents, err := eventsattu.GetAllEvents(eventsattu.GetAllEventsOpts{Take: 50})
	if err != nil {
		log.Info("Failed to get events", slog.Any("err", err))
		return "", err
	}
	mappedEventsAtTUEvents := mappers.MapEventAtTUEventsToStdEvent(eventsAtTUEvents)

	allEvents = append(allEvents, mappedInvolvedEvents...)
	allEvents = append(allEvents, mappedEventsAtTUEvents...)

	marshaledEvents, err := json.Marshal(allEvents)
	if err != nil {
		log.Error("Failed to marshal events", slog.Any("err", err))
		return "", err
	}

	s.Redis.Del(s.Ctx, key)
	err = s.Redis.SetEx(s.Ctx, key, marshaledEvents, constants.FIFTEEN_MINUTES).Err()
	if err != nil {
		log.Error("Failed to create snapshot", slog.Any("err", err))
		return "", err
	}
	log.Info("Created new snapshot", slog.String("key", key))

	return string(marshaledEvents), err
}
