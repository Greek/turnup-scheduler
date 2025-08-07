package events

type EventType string

const (
	EventTypeInvolved   EventType = "involved"
	EventTypeEventsAtTU EventType = "events"
)

type StandardEvent struct {
	Id             any       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	OriginalUrl    string    `json:"original_url"`
	OrganizationId int       `json:"organization_id"`
	Location       string    `json:"location"`
	CoverImage     string    `json:"cover_image"`
	Lat            string    `json:"lat"`
	Long           string    `json:"long"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	EventSource    EventType `json:"event_source"`
}
