package lib_test

import (
	"testing"
	"time"

	"turnup-scheduler/internal/lib"
)

func TestBuildDate(t *testing.T) {
	got := lib.BuildDate()
	want := time.Now().UTC().Format("20060102")
	if got != want {
		t.Errorf("BuildDate() = %v, want %v", got, want)
	}
}
