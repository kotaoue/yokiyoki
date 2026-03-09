package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"yokiyoki/pkg/models"
)

// GetComments fetches all issue and PR comments for the given repository since the given time.
// The returned comments include an IssueNumber field that can be used to match them
// to specific PRs or Issues.
func GetComments(repo models.Repository, since time.Time) []models.Comment {
	sinceDate := since.Format("2006-01-02")
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/comments?since=%s", repo.Owner, repo.Name, sinceDate)
	rawComments, err := Executor(endpoint, repo, "comments")
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		return []models.Comment{}
	}

	var comments []models.Comment
	for _, raw := range rawComments {
		comment := parseComment(raw)
		if comment != nil {
			comments = append(comments, *comment)
		}
	}

	return comments
}

func parseComment(raw map[string]any) *models.Comment {
	body, ok := raw["body"].(string)
	if !ok {
		return nil
	}

	comment := models.Comment{
		Author:    parseUser(raw),
		Body:      body,
		CreatedAt: parseCreatedAt(raw),
	}

	if id, ok := raw["id"].(float64); ok {
		comment.ID = int(id)
	}

	if url, ok := raw["html_url"].(string); ok {
		comment.URL = url
	}

	if issueURL, ok := raw["issue_url"].(string); ok {
		comment.IssueNumber = extractNumberFromURL(issueURL)
	}

	return &comment
}

// extractNumberFromURL extracts the trailing numeric ID from a GitHub API URL.
// For example: "https://api.github.com/repos/owner/repo/issues/123" → 123
func extractNumberFromURL(apiURL string) int {
	parts := strings.Split(apiURL, "/")
	if len(parts) == 0 {
		return 0
	}

	numberStr := parts[len(parts)-1]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0
	}

	return number
}
