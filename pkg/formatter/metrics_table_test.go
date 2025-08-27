package formatter_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func TestNewMetricsTable(t *testing.T) {
	metrics := []models.Metrics{
		{
			Repository: "owner/repo",
			User:       "testuser",
			Commits:    10,
		},
	}

	table := formatter.NewMetricsTable(metrics)
	assert.NotNil(t, table)
}

func TestMetricsTable_Output(t *testing.T) {
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
		name          string
		byUser        bool
		detailedStats bool
		wantContains  []string
	}{
		{
			name:          "by user without detailed stats",
			byUser:        true,
			detailedStats: false,
			wantContains: []string{
				"Repository",
				"User",
				"owner/repo",
				"testuser",
				"10",
				"4/5 (80%)",
			},
		},
		{
			name:          "by user with detailed stats",
			byUser:        true,
			detailedStats: true,
			wantContains: []string{
				"Repository",
				"User",
				"Lines +/-",
				"owner/repo",
				"testuser",
				"10",
				"4/5 (80%)",
				"+500/-200",
			},
		},
		{
			name:          "by repo without detailed stats",
			byUser:        false,
			detailedStats: false,
			wantContains: []string{
				"Repository",
				"owner/repo",
				"10",
				"4/5 (80%)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			table := formatter.NewMetricsTable(metrics)
			table.Output(tt.byUser, tt.detailedStats)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, want := range tt.wantContains {
				assert.Contains(t, output, want)
			}

			if tt.byUser {
				assert.Contains(t, output, "User")
			} else {
				// Make sure User column is not included when byUser is false
				lines := strings.Split(output, "\n")
				headerLine := lines[0]
				assert.NotContains(t, headerLine, "| User |")
			}
		})
	}
}
