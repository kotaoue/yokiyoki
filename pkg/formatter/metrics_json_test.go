package formatter_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func TestNewMetricsJson(t *testing.T) {
	metrics := []models.Metrics{
		{
			Repository: "owner/repo",
			User:       "testuser",
			Commits:    10,
		},
	}

	j := formatter.NewMetricsJson(metrics)
	assert.NotNil(t, j)
}

func TestMetricsJson_Output(t *testing.T) {
	metrics := []models.Metrics{
		{
			Repository:        "owner/repo",
			User:              "testuser",
			Commits:           10,
			LinesAdded:        500,
			LinesDeleted:      200,
			PRsCreated:        5,
			PRsMerged:         4,
			PRMergeRate:       "80%",
			AvgPRMergeTime:    "2d 12h 30m",
			IssuesCreated:     3,
			IssuesClosed:      2,
			IssueResolveRate:  "67%",
			AvgIssueCloseTime: "1d 05h 15m",
			OpenIssues:        1,
		},
	}

	tests := []struct {
		name   string
		byUser bool
	}{
		{name: "without user", byUser: false},
		{name: "with user", byUser: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			j := formatter.NewMetricsJson(metrics)
			j.Output(tt.byUser, false)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			var result []map[string]interface{}
			err := json.Unmarshal([]byte(output), &result)
			assert.NoError(t, err)
			assert.Len(t, result, 1)
			assert.Equal(t, "owner/repo", result[0]["repository"])
			assert.Equal(t, float64(10), result[0]["commits"])

			if tt.byUser {
				assert.Equal(t, "testuser", result[0]["user"])
			} else {
				_, hasUser := result[0]["user"]
				assert.False(t, hasUser)
			}
		})
	}
}

func TestMetricsJson_Output_EmptyData(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	j := formatter.NewMetricsJson([]models.Metrics{})
	j.Output(false, false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Empty(t, output)
}
