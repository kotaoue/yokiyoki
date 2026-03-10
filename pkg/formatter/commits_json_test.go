package formatter_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func TestNewCommitsJson(t *testing.T) {
	j := formatter.NewCommitsJson(sampleCommits())
	assert.NotNil(t, j)
}

func TestCommitsJson_Output(t *testing.T) {
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
	}{
		{name: "without detailed stats", detailedStats: false},
		{name: "with detailed stats", detailedStats: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			j := formatter.NewCommitsJson(commits)
			j.Output(tt.detailedStats)

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
			assert.Equal(t, "abc1234", result[0]["sha"])
			assert.Equal(t, "testuser", result[0]["author"])
			assert.Equal(t, "Fix bug", result[0]["message"])

			if tt.detailedStats {
				assert.Equal(t, float64(10), result[0]["additions"])
				assert.Equal(t, float64(5), result[0]["deletions"])
			} else {
				_, hasAdditions := result[0]["additions"]
				assert.False(t, hasAdditions)
				_, hasDeletions := result[0]["deletions"]
				assert.False(t, hasDeletions)
			}
		})
	}
}

func TestCommitsJson_Output_EmptyData(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	j := formatter.NewCommitsJson([]models.Commit{})
	j.Output(false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Empty(t, output)
}
