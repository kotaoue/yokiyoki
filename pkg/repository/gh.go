package repository

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"yokiyoki/pkg/models"
)

// Executor is the function used for executing commands (can be overridden for testing)
var Executor executeFunc = execute

// executeFunc defines the function signature for executing commands
type executeFunc func(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error)

// GetPullRequests fetches pull requests for the given repository using GitHub CLI
// If since is zero time, fetches all pull requests. Otherwise filters by creation date.
func GetPullRequests(repo models.Repository, since time.Time) []models.PullRequest {
	// Use search-based filtering for date ranges (except in test environment)
	if !since.IsZero() && !isTestEnvironment() {
		return getPullRequestsWithSearch(repo, since)
	}

	endpoint := fmt.Sprintf("/repos/%s/%s/pulls?state=all", repo.Owner, repo.Name)
	rawPRs, err := Executor(endpoint, repo, "pull requests")
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		return []models.PullRequest{}
	}

	var prs []models.PullRequest
	for _, raw := range rawPRs {
		pr := models.PullRequest{
			Number:    int(raw["number"].(float64)),
			Title:     raw["title"].(string),
			State:     raw["state"].(string),
			URL:       raw["html_url"].(string),
			Author:    parseUser(raw),
			CreatedAt: parseCreatedAt(raw),
		}

		if mergedTime := parseTimeField(raw, "merged_at"); mergedTime != nil {
			pr.MergedAt = mergedTime
		}

		if closedTime := parseTimeField(raw, "closed_at"); closedTime != nil {
			pr.ClosedAt = closedTime
		}

		if additions, ok := raw["additions"].(float64); ok {
			pr.Additions = int(additions)
		}
		if deletions, ok := raw["deletions"].(float64); ok {
			pr.Deletions = int(deletions)
		}

		prs = append(prs, pr)
	}

	return prs
}

// GetCommits fetches all commits for the given repository since the period start date using GitHub CLI
func GetCommits(repo models.Repository, from time.Time, detailedStats bool) []models.Commit {
	since := from.Format("2006-01-02")
	endpoint := fmt.Sprintf("/repos/%s/%s/commits?since=%s", repo.Owner, repo.Name, since)
	rawCommits, err := Executor(endpoint, repo, "commits")
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		return []models.Commit{}
	}

	var commits []models.Commit
	for _, raw := range rawCommits {
		authorName, authorDate := parseCommitAuthor(raw)
		additions, deletions := getDetailedStats(raw, repo, detailedStats)

		commit := models.Commit{
			SHA:       raw["sha"].(string),
			Message:   raw["commit"].(map[string]any)["message"].(string),
			URL:       raw["html_url"].(string),
			Author:    authorName,
			Date:      authorDate,
			Additions: additions,
			Deletions: deletions,
		}

		commits = append(commits, commit)
	}

	return commits
}

// GetIssues fetches issues for the given repository using GitHub CLI
// If since is zero time, fetches all issues. Otherwise filters by creation date.
func GetIssues(repo models.Repository, since time.Time) []models.Issue {
	// Use search-based filtering for date ranges (except in test environment)
	if !since.IsZero() && !isTestEnvironment() {
		return getIssuesWithSearch(repo, since)
	}

	endpoint := fmt.Sprintf("/repos/%s/%s/issues?state=all", repo.Owner, repo.Name)
	rawIssues, err := Executor(endpoint, repo, "issues")
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		return []models.Issue{}
	}

	var issues []models.Issue
	for _, raw := range rawIssues {
		if _, exists := raw["pull_request"]; exists {
			continue
		}

		issue := models.Issue{
			Number:    int(raw["number"].(float64)),
			Title:     raw["title"].(string),
			State:     raw["state"].(string),
			Author:    parseUser(raw),
			CreatedAt: parseCreatedAt(raw),
			Labels:    parseLabels(raw),
		}

		if closedTime := parseTimeField(raw, "closed_at"); closedTime != nil {
			issue.ClosedAt = closedTime
		}

		issues = append(issues, issue)
	}

	return issues
}

func execute(endpoint string, repo models.Repository, resourceType string) ([]map[string]any, error) {
	cmd := exec.Command("gh", "api", endpoint, "--paginate")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s for %s/%s: %w", resourceType, repo.Owner, repo.Name, err)
	}

	var rawData []map[string]any
	err = json.Unmarshal(output, &rawData)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s for %s/%s: %w", resourceType, repo.Owner, repo.Name, err)
	}

	fmt.Printf("Found %d %s for %s/%s\n", len(rawData), resourceType, repo.Owner, repo.Name)
	return rawData, nil
}

func parseTimeField(raw map[string]any, fieldName string) *time.Time {
	timeStr, ok := raw[fieldName].(string)
	if !ok || timeStr == "" {
		return nil
	}

	date, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}

	return &date
}

func parseCreatedAt(raw map[string]any) time.Time {
	timeStr, ok := raw["created_at"].(string)
	if !ok {
		return time.Time{}
	}

	date, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}

	return date
}

func parseUser(raw map[string]any) string {
	user, ok := raw["user"].(map[string]any)
	if !ok {
		return ""
	}

	login, ok := user["login"].(string)
	if !ok {
		return ""
	}

	return login
}

func parseLabels(raw map[string]any) []string {
	labels, ok := raw["labels"].([]any)
	if !ok {
		return nil
	}

	var result []string
	for _, label := range labels {
		labelMap, ok := label.(map[string]any)
		if !ok {
			continue
		}

		name, ok := labelMap["name"].(string)
		if !ok {
			continue
		}

		result = append(result, name)
	}

	return result
}

func parseCommitAuthor(raw map[string]any) (string, time.Time) {
	commit, ok := raw["commit"].(map[string]any)
	if !ok {
		return "", time.Time{}
	}

	author, ok := commit["author"].(map[string]any)
	if !ok {
		return "", time.Time{}
	}

	var name string
	var date time.Time

	if n, ok := author["name"].(string); ok {
		name = n
	}

	if dateStr, ok := author["date"].(string); ok {
		if d, err := time.Parse(time.RFC3339, dateStr); err == nil {
			date = d
		}
	}

	return name, date
}

func getDetailedStats(raw map[string]any, repo models.Repository, detailedStats bool) (int, int) {
	if !detailedStats {
		return 0, 0
	}

	sha, ok := raw["sha"].(string)
	if !ok {
		return 0, 0
	}

	return getCommitStats(repo, sha)
}

func getCommitStats(repo models.Repository, sha string) (int, int) {
	cmd := exec.Command("gh", "api", fmt.Sprintf("/repos/%s/%s/commits/%s", repo.Owner, repo.Name, sha))
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	var statsData map[string]any
	if err := json.Unmarshal(output, &statsData); err != nil {
		return 0, 0
	}

	stats, ok := statsData["stats"].(map[string]any)
	if !ok {
		return 0, 0
	}

	var additions, deletions int
	if add, ok := stats["additions"].(float64); ok {
		additions = int(add)
	}
	if del, ok := stats["deletions"].(float64); ok {
		deletions = int(del)
	}

	return additions, deletions
}

func getPullRequestsWithSearch(repo models.Repository, since time.Time) []models.PullRequest {
	searchQuery := buildDateRangeQuery(since)
	rawPRs := fetchPRsWithGHCommand(repo, searchQuery)
	return parsePRsFromJSON(rawPRs, repo.Owner, repo.Name)
}

func parseUserFromJSON(raw map[string]any) string {
	if author, ok := raw["author"].(map[string]any); ok {
		if login, ok := author["login"].(string); ok {
			return login
		}
	}
	return ""
}

func parseCreatedAtFromJSON(raw map[string]any) time.Time {
	if timeStr, ok := raw["createdAt"].(string); ok {
		if date, err := time.Parse(time.RFC3339, timeStr); err == nil {
			return date
		}
	}
	return time.Time{}
}

func parseTimeFieldFromJSON(raw map[string]any, fieldName string) *time.Time {
	if timeStr, ok := raw[fieldName].(string); ok && timeStr != "" {
		if date, err := time.Parse(time.RFC3339, timeStr); err == nil {
			return &date
		}
	}
	return nil
}

func getIssuesWithSearch(repo models.Repository, since time.Time) []models.Issue {
	searchQuery := buildDateRangeQuery(since)
	rawIssues := fetchIssuesWithGHCommand(repo, searchQuery)
	return parseIssuesFromJSON(rawIssues, repo.Owner, repo.Name)
}

var isInTestMode = false

func SetTestMode(testMode bool) {
	isInTestMode = testMode
}

func isTestEnvironment() bool {
	// Check if we're running in test mode
	return isInTestMode
}
