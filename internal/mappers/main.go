package mappers

import (
	"fmt"
	eventsTypes "turnup-scheduler/internal/datasource/events"
	"turnup-scheduler/internal/datasource/involved"
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
			Location:    event.Address.Location,
			Lat:         event.Address.Latitude,
			Long:        event.Address.Longitude,
			StartDate:   event.StartsOn,
			EndDate:     event.EndsOn,
			EventSource: eventsTypes.EventTypeInvolved,
		})
	}
	return standardEvents
}
