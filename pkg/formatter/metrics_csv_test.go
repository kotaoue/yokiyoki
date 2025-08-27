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

func TestNewMetricsCsv(t *testing.T) {
	metrics := []models.Metrics{
		{
			Repository: "owner/repo",
			User:       "testuser",
			Commits:    10,
		},
	}

	csv := formatter.NewMetricsCsv(metrics)
	assert.NotNil(t, csv)
}

func TestMetricsCsv_Output(t *testing.T) {
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
		name       string
		byUser     bool
		wantHeader string
		wantData   string
	}{
		{
			name:       "without user",
			byUser:     false,
			wantHeader: "Repository,Commits,LinesAdded,LinesDeleted,PRsCreated,PRsMerged,PRMergeRate,AvgPRMergeTime,IssuesCreated,IssuesClosed,IssueResolveRate,AvgIssueCloseTime,OpenIssues",
			wantData:   "owner/repo,10,500,200,5,4,80%,2d 12h 30m,3,2,67%,1d 05h 15m,1",
		},
		{
			name:       "with user",
			byUser:     true,
			wantHeader: "Repository,User,Commits,LinesAdded,LinesDeleted,PRsCreated,PRsMerged,PRMergeRate,AvgPRMergeTime,IssuesCreated,IssuesClosed,IssueResolveRate,AvgIssueCloseTime,OpenIssues",
			wantData:   "owner/repo,testuser,10,500,200,5,4,80%,2d 12h 30m,3,2,67%,1d 05h 15m,1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			csv := formatter.NewMetricsCsv(metrics)
			csv.Output(tt.byUser, false)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := strings.TrimSpace(buf.String())

			lines := strings.Split(output, "\n")
			assert.Equal(t, tt.wantHeader, lines[0])
			assert.Equal(t, tt.wantData, lines[1])
		})
	}
}

func TestMetricsCsv_Output_EmptyData(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	csv := formatter.NewMetricsCsv([]models.Metrics{})
	csv.Output(false, false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	assert.Empty(t, output)
}
