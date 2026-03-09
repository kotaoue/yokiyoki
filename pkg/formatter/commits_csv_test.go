package formatter_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func TestNewCommitsCsv(t *testing.T) {
	csv := formatter.NewCommitsCsv(sampleCommits())
	assert.NotNil(t, csv)
}

func TestCommitsCsv_Output(t *testing.T) {
	commits := []models.Commit{
		{
			Repository: "owner/repo",
			SHA:        "abc1234567890",
			Author:     "testuser",
			Date:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Message:    "Fix bug",
			Additions:  10,
			Deletions:  5,
		},
	}

	tests := []struct {
		name          string
		detailedStats bool
		wantHeader    string
		wantData      string
	}{
		{
			name:          "without detailed stats",
			detailedStats: false,
			wantHeader:    "Repository,SHA,Author,Date,Message",
			wantData:      "owner/repo,abc1234,testuser,2024-01-15T10:30:00Z,Fix bug",
		},
		{
			name:          "with detailed stats",
			detailedStats: true,
			wantHeader:    "Repository,SHA,Author,Date,Message,LinesAdded,LinesDeleted",
			wantData:      "owner/repo,abc1234,testuser,2024-01-15T10:30:00Z,Fix bug,10,5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			csv := formatter.NewCommitsCsv(commits)
			csv.Output(tt.detailedStats)

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

func TestCommitsCsv_Output_EmptyData(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	csv := formatter.NewCommitsCsv([]models.Commit{})
	csv.Output(false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	assert.Empty(t, output)
}

func TestCommitsCsv_Output_MessageWithComma(t *testing.T) {
	commits := []models.Commit{
		{
			Repository: "owner/repo",
			SHA:        "abc1234567890",
			Author:     "testuser",
			Date:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Message:    "Fix bug, add feature",
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	csv := formatter.NewCommitsCsv(commits)
	csv.Output(false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := strings.TrimSpace(buf.String())

	lines := strings.Split(output, "\n")
	// Message containing comma should be quoted
	assert.Contains(t, lines[1], `"Fix bug, add feature"`)
}
