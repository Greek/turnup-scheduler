package constants

import "time"

const (
	INVOLVED_URL        = "https://involved.towson.edu/api/"
	EVENTS_URL          = "https://events.towson.edu/api/2/"
	KEY_LATEST_SNAPSHOT = "snapshot:"
	FIFTEEN_MINUTES     = time.Duration(15 * time.Minute)
)
