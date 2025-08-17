package involved

import (
	"encoding/json"
	"fmt"
	"turnup-scheduler/internal/constants"
	"turnup-scheduler/internal/lib/http"
)

type InvolvedEvent struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId int    `json:"organizationId"`
	ImagePath      string `json:"imagePath"`
	ImageUrl       string `json:"imageUrl"`
	// Address        InvolvedAddress `json:"address"`
	Location  string `json:"location"`
	StartsOn  string `json:"startsOn"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	EndsOn    string `json:"endsOn"`
}

type InvolvedAddress struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Latitude  string `json:"latitude"`

}

type InvolvedResponseWithEvents struct {
	Value []InvolvedEvent `json:"value"`
}

type GetAllEventsOpts struct{ Take int }

func GetAllEvents(opts GetAllEventsOpts) (InvolvedResponseWithEvents, error) {
	var data InvolvedResponseWithEvents
	rawData, err := http.GetHTTPData(constants.INVOLVED_URL + "/discovery/event/search?take=" + fmt.Sprint(opts.Take))
	if err != nil {
		return InvolvedResponseWithEvents{}, err
	}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return InvolvedResponseWithEvents{}, err
	}

	return data, nil
}
