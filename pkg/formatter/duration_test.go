package formatter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "zero duration",
			duration: 0,
			want:     "0d 00h 00m",
		},
		{
			name:     "1 hour",
			duration: time.Hour,
			want:     "0d 01h 00m",
		},
		{
			name:     "1 day",
			duration: 24 * time.Hour,
			want:     "1d 00h 00m",
		},
		{
			name:     "1 day 2 hours 30 minutes",
			duration: 24*time.Hour + 2*time.Hour + 30*time.Minute,
			want:     "1d 02h 30m",
		},
		{
			name:     "5 days 12 hours 45 minutes",
			duration: 5*24*time.Hour + 12*time.Hour + 45*time.Minute,
			want:     "5d 12h 45m",
		},
		{
			name:     "30 minutes",
			duration: 30 * time.Minute,
			want:     "0d 00h 30m",
		},
		{
			name:     "23 hours 59 minutes",
			duration: 23*time.Hour + 59*time.Minute,
			want:     "0d 23h 59m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.FormatDuration(tt.duration)
			assert.Equal(t, tt.want, got)
		})
	}
}
