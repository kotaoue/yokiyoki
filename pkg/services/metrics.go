package services

import (
	"fmt"
	"sort"
	"time"

	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
)

// MetricsOptions represents configuration for metrics collection
type MetricsOptions struct {
	Period         *Chronometer
	ByUser         bool
	NormalizeUsers bool
	DetailedStats  bool
	SortBy         string
}

// Execute processes metrics collection with options
func Execute(repo models.Repository, options MetricsOptions) []models.Metrics {
	var metrics []models.Metrics
	if options.ByUser {
		metrics = executeByUser(repo, options)
	} else {
		result := executeForRepo(repo, options)
		metrics = []models.Metrics{result}
	}

	if options.SortBy != "" {
		sortMetrics(metrics, options.SortBy)
	}

	return metrics
}

func sortMetrics(metrics []models.Metrics, sortBy string) {
	sort.Slice(metrics, func(i, j int) bool {
		switch sortBy {
		case "repository,user":
			if metrics[i].Repository != metrics[j].Repository {
				return metrics[i].Repository < metrics[j].Repository
			}
			return metrics[i].User < metrics[j].User
		case "user,repository":
			if metrics[i].User != metrics[j].User {
				return metrics[i].User < metrics[j].User
			}
			return metrics[i].Repository < metrics[j].Repository
		default: // "repository"
			return metrics[i].Repository < metrics[j].Repository
		}
	})
}

func executeForRepo(repo models.Repository, options MetricsOptions) models.Metrics {
	repoFullName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)

	commits := repository.GetCommits(repo, options.Period.StartTime(), options.DetailedStats)
	prs := repository.GetPullRequests(repo, options.Period.StartTime())
	issues := repository.GetIssues(repo, options.Period.StartTime())

	return calculateMetricsFromData(repoFullName, "", commits, prs, issues, options.Period, options.DetailedStats)
}

func executeByUser(repo models.Repository, options MetricsOptions) []models.Metrics {
	repoFullName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)

	commits := repository.GetCommits(repo, options.Period.StartTime(), options.DetailedStats)
	prs := repository.GetPullRequests(repo, options.Period.StartTime())
	issues := repository.GetIssues(repo, options.Period.StartTime())

	userCommits := groupCommitsByUser(commits, options.NormalizeUsers)
	userPRs := groupPRsByUser(prs, options.NormalizeUsers)
	userIssues := groupIssuesByUser(issues, options.NormalizeUsers)

	users := extractUniqueUsers(userCommits, userPRs, userIssues)

	// ユーザーがいない場合は"-"で表示
	if len(users) == 0 {
		emptyMetrics := calculateMetricsFromData(repoFullName, "-", []models.Commit{}, []models.PullRequest{}, []models.Issue{}, options.Period, options.DetailedStats)
		return []models.Metrics{emptyMetrics}
	}

	return calculateUserMetrics(repoFullName, users, userCommits, userPRs, userIssues, options)
}

func userName(author string, normalizeUsers bool) string {
	username := models.NewUserName(author)
	return username.Name(normalizeUsers)
}

func groupCommitsByUser(commits []models.Commit, normalizeUsers bool) map[string][]models.Commit {
	userCommits := make(map[string][]models.Commit)
	for _, commit := range commits {
		author := userName(commit.Author, normalizeUsers)
		userCommits[author] = append(userCommits[author], commit)
	}
	return userCommits
}

func groupPRsByUser(prs []models.PullRequest, normalizeUsers bool) map[string][]models.PullRequest {
	userPRs := make(map[string][]models.PullRequest)
	for _, pr := range prs {
		author := userName(pr.Author, normalizeUsers)
		userPRs[author] = append(userPRs[author], pr)
	}
	return userPRs
}

func groupIssuesByUser(issues []models.Issue, normalizeUsers bool) map[string][]models.Issue {
	userIssues := make(map[string][]models.Issue)
	for _, issue := range issues {
		author := userName(issue.Author, normalizeUsers)
		userIssues[author] = append(userIssues[author], issue)
	}
	return userIssues
}

func extractUniqueUsers(userCommits map[string][]models.Commit, userPRs map[string][]models.PullRequest, userIssues map[string][]models.Issue) map[string]bool {
	users := make(map[string]bool)
	for user := range userCommits {
		users[user] = true
	}
	for user := range userPRs {
		users[user] = true
	}
	for user := range userIssues {
		users[user] = true
	}
	return users
}

func calculateUserMetrics(
	repoFullName string,
	users map[string]bool,
	userCommits map[string][]models.Commit,
	userPRs map[string][]models.PullRequest,
	userIssues map[string][]models.Issue,
	options MetricsOptions,
) []models.Metrics {
	var result []models.Metrics
	for user := range users {
		if user != "" {
			userMetrics := calculateMetricsFromData(
				repoFullName,
				user,
				userCommits[user],
				userPRs[user],
				userIssues[user],
				options.Period,
				options.DetailedStats,
			)
			result = append(result, userMetrics)
		}
	}
	return result
}

func filterCommitsInPeriod(commits []models.Commit, period *Chronometer) []models.Commit {
	var filtered []models.Commit
	for _, commit := range commits {
		if period.Contains(commit.Date) {
			filtered = append(filtered, commit)
		}
	}
	return filtered
}

func filterPRsInPeriod(prs []models.PullRequest, period *Chronometer) []models.PullRequest {
	var filtered []models.PullRequest
	for _, pr := range prs {
		switch {
		case period.Contains(pr.CreatedAt):
			filtered = append(filtered, pr)
		case pr.ClosedAt != nil && period.Contains(*pr.ClosedAt):
			filtered = append(filtered, pr)
		case pr.MergedAt != nil && period.Contains(*pr.MergedAt):
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

func filterIssuesInPeriod(issues []models.Issue, period *Chronometer) []models.Issue {
	var filtered []models.Issue
	for _, issue := range issues {
		switch {
		case period.Contains(issue.CreatedAt):
			filtered = append(filtered, issue)
		case issue.ClosedAt != nil && period.Contains(*issue.ClosedAt):
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func analyzeIssues(issues []models.Issue, period *Chronometer) (int, int, int, []time.Duration) {
	issuesCreated := len(issues)
	issuesClosed := 0
	openIssues := 0
	var closeTimes []time.Duration

	for _, issue := range issues {
		if issue.ClosedAt != nil && period.Contains(*issue.ClosedAt) {
			issuesClosed++
			closeTimes = append(closeTimes, issue.ClosedAt.Sub(issue.CreatedAt))
		} else if issue.State == "open" {
			openIssues++
		}
	}

	return issuesCreated, issuesClosed, openIssues, closeTimes
}

func calculateMetricsFromData(
	repo, user string,
	commits []models.Commit,
	prs []models.PullRequest,
	issues []models.Issue,
	period *Chronometer,
	detailedStats bool,
) models.Metrics {
	filteredCommits := filterCommitsInPeriod(commits, period)
	filteredPRs := filterPRsInPeriod(prs, period)
	filteredIssues := filterIssuesInPeriod(issues, period)

	commitCount := len(filteredCommits)
	linesAdded := 0
	linesDeleted := 0

	if detailedStats {
		for _, commit := range filteredCommits {
			linesAdded += commit.Additions
			linesDeleted += commit.Deletions
		}
	}

	prsCreated := len(filteredPRs)
	prsMerged := 0
	var mergeTimes []time.Duration
	for _, pr := range filteredPRs {
		if pr.MergedAt != nil {
			prsMerged++
			if period.Contains(*pr.MergedAt) {
				mergeTimes = append(mergeTimes, pr.MergedAt.Sub(pr.CreatedAt))
			}
		}
	}

	issuesCreated, issuesClosed, openIssues, closeTimes := analyzeIssues(filteredIssues, period)

	avgPRMergeTime := calculateAverageTime(mergeTimes)
	avgIssueCloseTime := calculateAverageTime(closeTimes)
	prMergeRate := calculateRate(prsMerged, prsCreated)
	issueResolveRate := calculateRate(issuesClosed, issuesCreated)

	return models.Metrics{
		Repository:        repo,
		User:              user,
		Commits:           commitCount,
		LinesAdded:        linesAdded,
		LinesDeleted:      linesDeleted,
		PRsCreated:        prsCreated,
		PRsMerged:         prsMerged,
		PRMergeRate:       prMergeRate,
		AvgPRMergeTime:    avgPRMergeTime,
		IssuesCreated:     issuesCreated,
		IssuesClosed:      issuesClosed,
		IssueResolveRate:  issueResolveRate,
		AvgIssueCloseTime: avgIssueCloseTime,
		OpenIssues:        openIssues,
	}
}

func calculateAverageTime(times []time.Duration) string {
	if len(times) == 0 {
		return "None"
	}

	var total time.Duration
	for _, d := range times {
		total += d
	}
	avg := total / time.Duration(len(times))
	return formatter.FormatDuration(avg)
}

func calculateRate(completed, total int) string {
	if total == 0 {
		return "None"
	}

	rate := float64(completed) / float64(total) * 100
	return fmt.Sprintf("%.0f%%", rate)
}
