package bpbuild

import (
	"testing"
	"time"
)

func TestDurationAvg(t *testing.T) {
	tests := []struct {
		name      string
		durations []time.Duration
		want      time.Duration
	}{
		{
			name:      "single",
			durations: []time.Duration{time.Hour},
			want:      time.Hour,
		},
		{
			name:      "multiple",
			durations: []time.Duration{time.Hour, time.Hour},
			want:      time.Hour,
		},
		{
			name:      "mixed",
			durations: []time.Duration{time.Hour, 2 * time.Hour, 2 * time.Hour, 3 * time.Hour},
			want:      2 * time.Hour,
		},
		{
			name:      "empty",
			durations: []time.Duration{},
			want:      time.Duration(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := durationAvg(tt.durations); got != tt.want {
				t.Errorf("durationAvg() = %v, want %v", got, tt.want)
			}
		})
	}
}
