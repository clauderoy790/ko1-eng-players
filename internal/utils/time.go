package utils

import (
	"fmt"
	"log"
	"time"
	_ "time/tzdata"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}
}

func Now() *time.Time {
	t := time.Now().In(location)
	return &t
}

// todo here implement
func TimeToString(t *time.Time) string {
	return t.In(location).Format(time.RFC1123)
}

func StringToTime(timeStr string) *time.Time {
	t, err := time.ParseInLocation(time.RFC1123, timeStr, location)
	if err != nil {
		log.Fatalf("Failed to parse time: %s with location: %v", timeStr, location)
		return nil
	}
	return &t
}

func DeltaDisplayTime(now, lastSeen *time.Time) string {
	duration := now.Sub(*lastSeen)

	if duration < time.Minute {
		return "a few seconds ago"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, pluralize(days))
	} else {
		weeks := int(duration.Hours() / (24 * 7))
		return fmt.Sprintf("%d week%s ago", weeks, pluralize(weeks))
	}
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
