package utils

import (
	"time"
)

func CheckIfExpired(timestamp string) bool {
	eventTimed, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return true
	}
	diff := eventTimed.Sub(time.Now()).Hours()
	if diff < 12 {
		return false
	} else {
		return true
	}
}
