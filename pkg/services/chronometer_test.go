package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/services"
)

func TestNewChronometer(t *testing.T) {
	days := 7
	before := time.Now()
	chronometer, err := services.NewChronometer(services.ChronometerOption{
		Days: &days,
	})
	assert.NoError(t, err)
	after := time.Now()

	// Check that start is approximately 7 days before now
	expectedStart := before.AddDate(0, 0, -days)
	assert.True(t, chronometer.StartTime().After(expectedStart.Add(-time.Second)) && chronometer.StartTime().Before(after.AddDate(0, 0, -days).Add(time.Second)))

	// Check that end is approximately now
	assert.True(t, chronometer.EndTime().After(before.Add(-time.Second)) && chronometer.EndTime().Before(after.Add(time.Second)))
}

func TestNewChronometerFromDates(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
		wantError bool
	}{
		{
			name:      "valid dates",
			startDate: "2024-01-01",
			endDate:   "2024-01-31",
			wantError: false,
		},
		{
			name:      "invalid start date",
			startDate: "invalid",
			endDate:   "2024-01-31",
			wantError: true,
		},
		{
			name:      "invalid end date",
			startDate: "2024-01-01",
			endDate:   "invalid",
			wantError: true,
		},
		{
			name:      "empty dates",
			startDate: "",
			endDate:   "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chronometer, err := services.NewChronometer(services.ChronometerOption{
				StartDate: &tt.startDate,
				EndDate:   &tt.endDate,
			})

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			expectedStart, _ := time.Parse("2006-01-02", tt.startDate)
			expectedEnd, _ := time.Parse("2006-01-02", tt.endDate)
			expectedEnd = expectedEnd.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

			assert.True(t, chronometer.StartTime().Equal(expectedStart))
			assert.True(t, chronometer.EndTime().Equal(expectedEnd))
		})
	}
}

func TestChronometer_Contains(t *testing.T) {
	startDate := "2024-01-01"
	endDate := "2024-01-31"
	chronometer, err := services.NewChronometer(services.ChronometerOption{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	assert.NoError(t, err)

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "date within period",
			date: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "date at start",
			date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "date before start",
			date: time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			want: false,
		},
		{
			name: "date after end",
			date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chronometer.Contains(tt.date)
			assert.Equal(t, tt.want, got)
		})
	}
}
