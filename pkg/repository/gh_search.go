package repository

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"yokiyoki/pkg/models"
)

// buildDateRangeQuery creates GitHub search query for date range
func buildDateRangeQuery(since time.Time) string {
	sinceDate := since.Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("created:%s..%s", sinceDate, today)
}

// fetchPRsWithGHCommand executes gh pr list command with search
func fetchPRsWithGHCommand(repo models.Repository, searchQuery string) []map[string]any {
	cmd := exec.Command("gh", "pr", "list",
		"-R", fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
		"--state", "all",
		"--search", searchQuery,
		"--limit", "1000",
		"--json", "number,title,state,author,createdAt,mergedAt,closedAt,url,additions,deletions")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Warning: could not fetch PRs with search for %s/%s: %v\n", repo.Owner, repo.Name, err)
		return []map[string]any{}
	}

	var rawPRs []map[string]any
	err = json.Unmarshal(output, &rawPRs)
	if err != nil {
		fmt.Printf("Warning: could not parse PR search results for %s/%s: %v\n", repo.Owner, repo.Name, err)
		return []map[string]any{}
	}

	return rawPRs
}

// fetchIssuesWithGHCommand executes gh issue list command with search
func fetchIssuesWithGHCommand(repo models.Repository, searchQuery string) []map[string]any {
	cmd := exec.Command("gh", "issue", "list",
		"-R", fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
		"--state", "all",
		"--search", searchQuery,
		"--limit", "1000",
		"--json", "number,title,state,author,createdAt,closedAt,labels")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Warning: could not fetch issues with search for %s/%s: %v\n", repo.Owner, repo.Name, err)
		return []map[string]any{}
	}

	var rawIssues []map[string]any
	err = json.Unmarshal(output, &rawIssues)
	if err != nil {
		fmt.Printf("Warning: could not parse issue search results for %s/%s: %v\n", repo.Owner, repo.Name, err)
		return []map[string]any{}
	}

	return rawIssues
}

// parsePRsFromJSON converts JSON data to PullRequest models
func parsePRsFromJSON(rawPRs []map[string]any, owner, name string) []models.PullRequest {
	fmt.Printf("Found %d pull requests for %s/%s\n", len(rawPRs), owner, name)

	var prs []models.PullRequest
	total := len(rawPRs)
	fmt.Printf("Processing %d pull requests...\n", total)

	for i, raw := range rawPRs {
		showProgress("Pull requests", i, total)

		pr := models.PullRequest{
			Number:    int(raw["number"].(float64)),
			Title:     raw["title"].(string),
			State:     raw["state"].(string),
			URL:       raw["url"].(string),
			Author:    parseUserFromJSON(raw),
			CreatedAt: parseCreatedAtFromJSON(raw),
		}

		if mergedTime := parseTimeFieldFromJSON(raw, "mergedAt"); mergedTime != nil {
			pr.MergedAt = mergedTime
		}

		if closedTime := parseTimeFieldFromJSON(raw, "closedAt"); closedTime != nil {
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

// parseIssuesFromJSON converts JSON data to Issue models
func parseIssuesFromJSON(rawIssues []map[string]any, owner, name string) []models.Issue {
	fmt.Printf("Found %d issues for %s/%s\n", len(rawIssues), owner, name)

	var issues []models.Issue
	total := len(rawIssues)
	fmt.Printf("Processing %d issues...\n", total)

	for i, raw := range rawIssues {
		showProgress("Issues", i, total)

		issue := models.Issue{
			Number:    int(raw["number"].(float64)),
			Title:     raw["title"].(string),
			State:     raw["state"].(string),
			Author:    parseUserFromJSON(raw),
			CreatedAt: parseCreatedAtFromJSON(raw),
			Labels:    parseLabelsFromJSON(raw),
		}

		if closedTime := parseTimeFieldFromJSON(raw, "closedAt"); closedTime != nil {
			issue.ClosedAt = closedTime
		}

		issues = append(issues, issue)
	}

	return issues
}

func parseLabelsFromJSON(raw map[string]any) []string {
	if labels, ok := raw["labels"].([]any); ok {
		var result []string
		for _, label := range labels {
			if labelMap, ok := label.(map[string]any); ok {
				if name, ok := labelMap["name"].(string); ok {
					result = append(result, name)
				}
			}
		}
		return result
	}
	return nil
}

// showProgress displays processing progress
func showProgress(resourceType string, current, total int) {
	if total > 10 && (current%50 == 0 || current == total-1) {
		fmt.Printf("\r%s: %d/%d", resourceType, current+1, total)
		if current == total-1 {
			fmt.Printf(" - completed\n")
		}
	}
}
