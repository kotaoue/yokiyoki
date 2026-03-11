package services_test

import (
	"testing"
	"time"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
	"yokiyoki/pkg/services"

	"github.com/stretchr/testify/assert"
)

func TestExecuteConversations(t *testing.T) {
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

	prDate := chronometer.StartTime().Add(24 * time.Hour)

	// The mock executor is called for:
	// 1. GetPullRequests (endpoint contains "/pulls")
	// 2. GetIssues (endpoint contains "/issues?")
	// 3. GetComments for PR #1 (endpoint contains "/issues/1/comments")
	// 4. GetComments for Issue #2 (endpoint contains "/issues/2/comments")
	repository.Executor = func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
		switch {
		case contains(endpoint, "pulls?state=all"):
			return []map[string]any{
				{
					"number":     float64(1),
					"title":      "Test PR",
					"state":      "open",
					"body":       "PR description",
					"html_url":   "https://github.com/test/repo/pull/1",
					"created_at": prDate.Format(time.RFC3339),
					"user":       map[string]any{"login": "alice"},
				},
			}, nil
		case contains(endpoint, "issues?state=all"):
			return []map[string]any{
				{
					"number":     float64(2),
					"title":      "Test Issue",
					"state":      "open",
					"body":       "Issue description",
					"html_url":   "https://github.com/test/repo/issues/2",
					"created_at": prDate.Format(time.RFC3339),
					"user":       map[string]any{"login": "bob"},
					"labels":     []any{},
				},
			}, nil
		default:
			// Comments endpoint
			return []map[string]any{
				{
					"body":       "A reply comment",
					"html_url":   "https://github.com/test/repo/issues/1#issuecomment-1",
					"created_at": prDate.Add(time.Hour).Format(time.RFC3339),
					"user":       map[string]any{"login": "charlie"},
				},
			}, nil
		}
	}

	repo := models.Repository{Owner: "test-owner", Name: "test-repo"}
	opts := services.ConversationsOptions{Period: chronometer}

	comments := services.ExecuteConversations(repo, opts)

	// Should have: PR body, PR comment, Issue body, Issue comment = 4 entries
	assert.Len(t, comments, 4)

	// All should be tagged with the repository name
	for _, c := range comments {
		assert.Equal(t, "test-owner/test-repo", c.Repository)
	}
}

func TestSortCommentsByDate(t *testing.T) {
	now := time.Now()
	comments := []models.Comment{
		{Author: "c", CreatedAt: now},
		{Author: "a", CreatedAt: now.Add(-2 * time.Hour)},
		{Author: "b", CreatedAt: now.Add(-1 * time.Hour)},
	}

	services.SortCommentsByDate(comments)

	assert.Equal(t, "a", comments[0].Author)
	assert.Equal(t, "b", comments[1].Author)
	assert.Equal(t, "c", comments[2].Author)
}

// contains reports whether s contains the given substring.
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > len(sub) && containsString(s, sub))
}

func containsString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
