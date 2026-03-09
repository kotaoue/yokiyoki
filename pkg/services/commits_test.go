package services_test

import (
	"testing"
	"time"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
	"yokiyoki/pkg/services"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommits(t *testing.T) {
	chronometer, err := services.NewChronometer(services.ChronometerOption{
		Days: func() *int { d := 30; return &d }(),
	})
	assert.NoError(t, err)

	originalExecutor := repository.Executor
	defer func() {
		repository.Executor = originalExecutor
		repository.SetTestMode(false)
	}()

	repository.SetTestMode(true)

	commitDate := chronometer.StartTime().Add(24 * time.Hour)
	repository.Executor = func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
		return []map[string]any{
			{
				"sha": "abc1234567890",
				"commit": map[string]any{
					"message": "Fix bug in login flow",
					"author": map[string]any{
						"name": "alice",
						"date": commitDate.Format(time.RFC3339),
					},
				},
				"html_url": "https://github.com/test/url",
			},
		}, nil
	}

	repo := models.Repository{Owner: "test-owner", Name: "test-repo"}
	opts := services.CommitsOptions{
		Period:        chronometer,
		DetailedStats: false,
	}

	commits := services.ExecuteCommits(repo, opts)

	assert.Len(t, commits, 1)
	assert.Equal(t, "test-owner/test-repo", commits[0].Repository)
	assert.Equal(t, "abc1234567890", commits[0].SHA)
	assert.Equal(t, "alice", commits[0].Author)
	assert.Equal(t, "Fix bug in login flow", commits[0].Message)
}

func TestSortCommitsByDate(t *testing.T) {
	now := time.Now()
	commits := []models.Commit{
		{SHA: "old", Date: now.Add(-48 * time.Hour)},
		{SHA: "newest", Date: now},
		{SHA: "middle", Date: now.Add(-24 * time.Hour)},
	}

	services.SortCommitsByDate(commits)

	assert.Equal(t, "newest", commits[0].SHA)
	assert.Equal(t, "middle", commits[1].SHA)
	assert.Equal(t, "old", commits[2].SHA)
}
