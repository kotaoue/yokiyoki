package formatter_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func sampleCommits() []models.Commit {
	return []models.Commit{
		{
			Repository: "owner/repo",
			SHA:        "abc1234567890",
			Author:     "testuser",
			Date:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			Message:    "Fix bug in login flow",
			Additions:  10,
			Deletions:  5,
		},
		{
			Repository: "owner/repo",
			SHA:        "def9876543210",
			Author:     "otheruser",
			Date:       time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC),
			Message:    "Add new feature",
			Additions:  50,
			Deletions:  3,
		},
	}
}

func TestNewCommitsTable(t *testing.T) {
	table := formatter.NewCommitsTable(sampleCommits())
	assert.NotNil(t, table)
}

func TestCommitsTable_Output(t *testing.T) {
	commits := sampleCommits()

	tests := []struct {
		name          string
		detailedStats bool
		wantContains  []string
	}{
		{
			name:          "without detailed stats",
			detailedStats: false,
			wantContains: []string{
				"Repository",
				"SHA",
				"Author",
				"Date",
				"Message",
				"owner/repo",
				"abc1234",
				"testuser",
				"2024-01-15",
				"Fix bug in login flow",
			},
		},
		{
			name:          "with detailed stats",
			detailedStats: true,
			wantContains: []string{
				"Repository",
				"Lines +/-",
				"owner/repo",
				"abc1234",
				"+10/-5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			table := formatter.NewCommitsTable(commits)
			table.Output(tt.detailedStats)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			for _, want := range tt.wantContains {
				assert.Contains(t, output, want)
			}
		})
	}
}

func TestCommitsTable_Output_Empty(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	table := formatter.NewCommitsTable([]models.Commit{})
	table.Output(false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Header row should still appear
	assert.Contains(t, output, "Repository")
	assert.Contains(t, output, "SHA")
}
