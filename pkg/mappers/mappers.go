package mappers

import (
	"fmt"
	eventsTypes "turnup-scheduler/pkg/datasource/events"
	"turnup-scheduler/pkg/datasource/eventsattu"
	"turnup-scheduler/pkg/datasource/involved"
)

func MapInvolvedEventsToStdEvent(events []involved.InvolvedEvent) []eventsTypes.StandardEvent {
	standardEvents := []eventsTypes.StandardEvent{}

	for _, event := range events {
		imagePath := ""
		if event.ImagePath != "" {
			imagePath = "https://se-images.campuslabs.com/clink/images/" + event.ImagePath
		} else {
			imagePath = event.ImageUrl
		}

		standardEvents = append(standardEvents, eventsTypes.StandardEvent{
			Id:          event.Id,
			Name:        event.Name,
			Description: event.Description,
			OriginalUrl: "https://involved.towson.edu/event/" + fmt.Sprint(event.Id),
			CoverImage:  imagePath,
			Location:    event.Location,
			Lat:         event.Latitude,
			Long:        event.Longitude,
			StartDate:   event.StartsOn,
			EndDate:     event.EndsOn,
			EventSource: eventsTypes.EventTypeInvolved,
		})
	}
	return standardEvents
}

func MapEventAtTUEventsToStdEvent(events eventsattu.EventsAtTUResponseWithEvents) []eventsTypes.StandardEvent {
	standardEvents := []eventsTypes.StandardEvent{}

	for _, entry := range events.Events {
		standardEvents = append(standardEvents, eventsTypes.StandardEvent{
			Id:          entry.Event.Id,
			Name:        entry.Event.Title,
			Description: entry.Event.Description,
			OriginalUrl: "https://events.towson.edu/event/" + fmt.Sprint(entry.Event.UrlName),
			Location:    entry.Event.LocationName,
			CoverImage:  entry.Event.PhotoUrl,
			Lat:         entry.Event.Geo.Latitude,
			Long:        entry.Event.Geo.Longitude,
			StartDate:   entry.Event.EventInstances[0].EventInstance.Start,
			EndDate:     entry.Event.EventInstances[0].EventInstance.End,
			EventSource: eventsTypes.EventTypeEventsAtTU,
		})
	}
	return standardEvents
}
