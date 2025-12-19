package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/logging"
	"turnup-scheduler/pkg/mappers"
	"turnup-scheduler/pkg/repo/events"
	"turnup-scheduler/pkg/repo/eventsattu"
	"turnup-scheduler/pkg/repo/involved"

	"github.com/redis/go-redis/v9"
)

const RETRY_INTERVAL = 600 * time.Millisecond

type CreateSnapshotOpts struct {
	Overwrite bool
}

func (s *Scheduler) GetSnapshot(date string, namespace string) ([]events.StandardEvent, error) {
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

func (s *Scheduler) CreateSnapshot(date string, namespace string, opts CreateSnapshotOpts) (string, error) {
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

	var involvedEventsResult involved.InvolvedResponseWithEvents
	var lastErr error
	for i := range 3 {
		involvedEventsResult, lastErr = involved.GetAllEvents(involved.GetAllEventsOpts{Take: 200})
		if lastErr == nil {
			break
		}
		log.Info("Retrying involved.GetAllEvents", slog.Int("attempt", i+1), slog.Any("err", lastErr))
		time.Sleep(RETRY_INTERVAL)
	}
	if lastErr != nil {
		log.Info("Failed to get events", slog.Any("err", lastErr))
		//TODO: Refactor this logic to try forever. Exiting is ok for now so long as our platform knows to restart on exit.
		os.Exit(1)
		// return "", lastErr
	}
	mappedInvolvedEvents := mappers.MapInvolvedEventsToStdEvent(involvedEventsResult.Value)

	log.Info("getting events@tu events")
	var eventsAtTUEventsResult eventsattu.EventsAtTUResponseWithEvents
	lastErr = nil
	for i := range 3 {
		eventsAtTUEventsResult, lastErr = eventsattu.GetAllEvents(eventsattu.GetAllEventsOpts{Take: 200, PerPage: 100})
		if lastErr == nil {
			break
		}
		log.Info("Retrying eventsattu.GetAllEvents", slog.Int("attempt", i+1), slog.Any("err", lastErr))
		time.Sleep(RETRY_INTERVAL)
	}
	if lastErr != nil {
		log.Info("Failed to get events", slog.Any("err", lastErr))
		//TODO: Refactor this logic to try forever. Exiting is ok for now so long as our platform knows to restart on exit.
		os.Exit(1)
		// return "", lastErr
	}
	mappedEventsAtTUEvents := mappers.MapEventAtTUEventsToStdEvent(eventsAtTUEventsResult)

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
