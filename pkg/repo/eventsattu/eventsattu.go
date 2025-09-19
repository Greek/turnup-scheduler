package eventsattu

import (
	"encoding/json"
	"fmt"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/lib/http"
)

type EventsAtTUEvent struct {
	Id             int              `json:"id"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	UrlName        string           `json:"urlname"`
	PhotoUrl       string           `json:"photo_url"`
	Geo            Geolocation      `json:"geo"`
	LocationName   string           `json:"location_name"`
	EventInstances []EventInstances `json:"event_instances"`
}

type Geolocation struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type EventInstance struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type EventInstances struct {
	EventInstance EventInstance `json:"event_instance"`
}

type EventObj struct {
	Event EventsAtTUEvent `json:"event"`
}

type EventsAtTUResponseWithEvents struct {
	Events []EventObj `json:"events"`
}

type GetAllEventsOpts struct {
	Take    int
	PerPage int
}

func GetAllEvents(opts GetAllEventsOpts) (EventsAtTUResponseWithEvents, error) {
	var data EventsAtTUResponseWithEvents
	rawData, err := http.GetHTTPData(constants.EVENTS_URL + "events/?days=" + fmt.Sprint(opts.Take) + "&pp=" + fmt.Sprint(opts.PerPage))
	if err != nil {
		return EventsAtTUResponseWithEvents{}, err
	}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return EventsAtTUResponseWithEvents{}, err
	}

	return data, nil
}
