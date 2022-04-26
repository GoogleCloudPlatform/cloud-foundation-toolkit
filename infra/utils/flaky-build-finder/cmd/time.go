package cmd

import (
	"time"
)

func getTimeFromStr(t string) (time.Time, error) {
	return time.Parse("01-02-2006", t)
}

func formatTimeCB(t time.Time) string {
	return t.Format(time.RFC3339)
}
