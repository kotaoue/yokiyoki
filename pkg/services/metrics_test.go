package services_test

import (
	"testing"
	"time"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
	"yokiyoki/pkg/services"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	chronometer, err := services.NewChronometer(services.ChronometerOption{
		Days: func() *int { d := 30; return &d }(),
	})
	assert.NoError(t, err)

	originalExecutor := repository.Executor
	defer func() { repository.Executor = originalExecutor }()

	repository.Executor = func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
		switch resourceType {
		case "commits":
			return []map[string]any{
				{
					"sha": "abc123",
					"commit": map[string]any{
						"message": "Test commit",
						"author": map[string]any{
							"name": "alice",
							"date": chronometer.StartTime().Add(24 * time.Hour).Format(time.RFC3339),
						},
					},
					"html_url": "https://github.com/test/url",
					"stats": map[string]any{
						"additions": float64(10),
						"deletions": float64(5),
					},
				},
			}, nil
		case "pull requests":
			mergedAt := chronometer.StartTime().Add(48 * time.Hour).Format(time.RFC3339)
			return []map[string]any{
				{
					"number":     float64(1),
					"title":      "Test PR",
					"state":      "merged",
					"html_url":   "https://github.com/test/pr",
					"created_at": chronometer.StartTime().Add(12 * time.Hour).Format(time.RFC3339),
					"merged_at":  mergedAt,
					"user": map[string]any{
						"login": "alice",
					},
				},
			}, nil
		case "issues":
			return []map[string]any{
				{
					"number":     float64(1),
					"title":      "Test Issue",
					"state":      "open",
					"created_at": chronometer.StartTime().Add(6 * time.Hour).Format(time.RFC3339),
					"user": map[string]any{
						"login": "alice",
					},
				},
			}, nil
		default:
			return []map[string]any{}, nil
		}
	}

	repo := models.Repository{
		Owner: "test-owner",
		Name:  "test-repo",
	}

	options := services.MetricsOptions{
		Period:         chronometer,
		ByUser:         false,
		NormalizeUsers: false,
		DetailedStats:  false,
		SortBy:         "",
	}

	metrics := services.Execute(repo, options)
	expected := []models.Metrics{{
		Repository:        "test-owner/test-repo",
		User:              "",
		Commits:           1,
		LinesAdded:        0,
		LinesDeleted:      0,
		PRsCreated:        1,
		PRsMerged:         1,
		PRMergeRate:       "100%",
		AvgPRMergeTime:    "1d 12h 00m",
		IssuesCreated:     1,
		IssuesClosed:      0,
		IssueResolveRate:  "0%",
		AvgIssueCloseTime: "None",
		OpenIssues:        1,
	}}
	assert.Equal(t, expected, metrics)
}
