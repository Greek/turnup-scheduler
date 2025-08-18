package lib

import (
	"strings"
	"time"
)

func BuildDate() string {
	return strings.ReplaceAll(time.Now().UTC().Format(time.DateOnly), "-", "")
}
