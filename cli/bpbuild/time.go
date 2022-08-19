package bpbuild

import "time"

// durationAvg calculates avg for a given slice of durations.
func durationAvg(durations []time.Duration) time.Duration {
	if len(durations) < 1 {
		return time.Duration(0)
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	avg := total.Seconds() / float64(len(durations))
	return time.Duration(avg * float64(time.Second))
}

// getTimeFromStr parses string formatted MM-DD-YYY as time.
func getTimeFromStr(t string) (time.Time, error) {
	return time.Parse("01-02-2006", t)
}
