package services_test

import (
	"testing"
	"time"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
	"yokiyoki/pkg/services"

	"github.com/stretchr/testify/assert"
)

func TestExecuteConversation(t *testing.T) {
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

	prCreatedAt := chronometer.StartTime().Add(12 * time.Hour).Format(time.RFC3339)
	mergedAt := chronometer.StartTime().Add(48 * time.Hour).Format(time.RFC3339)
	issueCreatedAt := chronometer.StartTime().Add(6 * time.Hour).Format(time.RFC3339)
	commentCreatedAt := chronometer.StartTime().Add(24 * time.Hour).Format(time.RFC3339)

	repository.Executor = func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
		switch resourceType {
		case "pull requests":
			return []map[string]any{
				{
					"number":     float64(1),
					"title":      "Test PR",
					"state":      "merged",
					"html_url":   "https://github.com/test/repo/pull/1",
					"created_at": prCreatedAt,
					"merged_at":  mergedAt,
					"user": map[string]any{
						"login": "alice",
					},
				},
			}, nil
		case "issues":
			return []map[string]any{
				{
					"number":     float64(2),
					"title":      "Test Issue",
					"state":      "open",
					"html_url":   "https://github.com/test/repo/issues/2",
					"created_at": issueCreatedAt,
					"user": map[string]any{
						"login": "bob",
					},
				},
			}, nil
		case "comments":
			return []map[string]any{
				{
					"id":         float64(100),
					"body":       "Looks good to me!",
					"html_url":   "https://github.com/test/repo/pull/1#issuecomment-100",
					"created_at": commentCreatedAt,
					"issue_url":  "https://api.github.com/repos/test/repo/issues/1",
					"user": map[string]any{
						"login": "bob",
					},
				},
				{
					"id":         float64(101),
					"body":       "Please check the tests",
					"html_url":   "https://github.com/test/repo/issues/2#issuecomment-101",
					"created_at": commentCreatedAt,
					"issue_url":  "https://api.github.com/repos/test/repo/issues/2",
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
		Owner: "test",
		Name:  "repo",
	}

	options := services.ConversationOptions{
		Period: chronometer,
	}

	conversations := services.ExecuteConversation(repo, options)

	assert.Len(t, conversations, 2)

	// Find PR conversation
	var prConv, issueConv *models.Conversation
	for i := range conversations {
		switch conversations[i].Type {
		case models.ConversationTypePR:
			prConv = &conversations[i]
		case models.ConversationTypeIssue:
			issueConv = &conversations[i]
		}
	}

	assert.NotNil(t, prConv)
	assert.Equal(t, "test/repo", prConv.Repository)
	assert.Equal(t, 1, prConv.Number)
	assert.Equal(t, "Test PR", prConv.Title)
	assert.Equal(t, "merged", prConv.State)
	assert.Equal(t, "alice", prConv.Author)
	assert.Len(t, prConv.Comments, 1)
	assert.Equal(t, "bob", prConv.Comments[0].Author)
	assert.Equal(t, "Looks good to me!", prConv.Comments[0].Body)

	assert.NotNil(t, issueConv)
	assert.Equal(t, "test/repo", issueConv.Repository)
	assert.Equal(t, 2, issueConv.Number)
	assert.Equal(t, "Test Issue", issueConv.Title)
	assert.Equal(t, "open", issueConv.State)
	assert.Equal(t, "bob", issueConv.Author)
	assert.Len(t, issueConv.Comments, 1)
	assert.Equal(t, "alice", issueConv.Comments[0].Author)
	assert.Equal(t, "Please check the tests", issueConv.Comments[0].Body)
}

func TestExecuteConversation_OutsidePeriod(t *testing.T) {
	chronometer, err := services.NewChronometer(services.ChronometerOption{
		Days: func() *int { d := 7; return &d }(),
	})
	assert.NoError(t, err)

	originalExecutor := repository.Executor
	defer func() {
		repository.Executor = originalExecutor
		repository.SetTestMode(false)
	}()

	repository.SetTestMode(true)

	// PR created before the period
	oldPRCreatedAt := chronometer.StartTime().Add(-48 * time.Hour).Format(time.RFC3339)

	repository.Executor = func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
		switch resourceType {
		case "pull requests":
			return []map[string]any{
				{
					"number":     float64(1),
					"title":      "Old PR",
					"state":      "open",
					"html_url":   "https://github.com/test/repo/pull/1",
					"created_at": oldPRCreatedAt,
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
		Owner: "test",
		Name:  "repo",
	}

	options := services.ConversationOptions{
		Period: chronometer,
	}

	conversations := services.ExecuteConversation(repo, options)

	// PR created before the period should not appear
	assert.Empty(t, conversations)
}
