package services

import (
	"fmt"
	"sort"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
)

// ConversationsOptions represents configuration for conversation collection.
type ConversationsOptions struct {
	Period *Chronometer
}

// ExecuteConversations fetches PR and issue conversations (initial bodies and
// comments) for the given repository, filtered to the configured period.
// Each PR/issue that was created within the period is included together with
// all its comments.
func ExecuteConversations(repo models.Repository, opts ConversationsOptions) []models.Comment {
	repoFullName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	var allComments []models.Comment

	// Collect PR conversations.
	prs := repository.GetPullRequests(repo, opts.Period.StartTime())
	for _, pr := range prs {
		if !prInPeriod(pr, opts.Period) {
			continue
		}

		// Include the PR description as the opening comment.
		if pr.Body != "" {
			allComments = append(allComments, models.Comment{
				Repository: repoFullName,
				Type:       "pr",
				Number:     pr.Number,
				Title:      pr.Title,
				Author:     pr.Author,
				Body:       pr.Body,
				URL:        pr.URL,
				CreatedAt:  pr.CreatedAt,
			})
		}

		// Append all follow-up comments on the PR.
		prComments := repository.GetComments(repo, pr.Number)
		for _, c := range prComments {
			c.Repository = repoFullName
			c.Type = "pr"
			c.Number = pr.Number
			c.Title = pr.Title
			allComments = append(allComments, c)
		}
	}

	// Collect Issue conversations.
	issues := repository.GetIssues(repo, opts.Period.StartTime())
	for _, issue := range issues {
		if !issueInPeriod(issue, opts.Period) {
			continue
		}

		// Include the issue description as the opening comment.
		if issue.Body != "" {
			allComments = append(allComments, models.Comment{
				Repository: repoFullName,
				Type:       "issue",
				Number:     issue.Number,
				Title:      issue.Title,
				Author:     issue.Author,
				Body:       issue.Body,
				URL:        issue.URL,
				CreatedAt:  issue.CreatedAt,
			})
		}

		// Append all follow-up comments on the issue.
		issueComments := repository.GetComments(repo, issue.Number)
		for _, c := range issueComments {
			c.Repository = repoFullName
			c.Type = "issue"
			c.Number = issue.Number
			c.Title = issue.Title
			allComments = append(allComments, c)
		}
	}

	SortCommentsByDate(allComments)
	return allComments
}

// SortCommentsByDate sorts comments by creation date in ascending order (oldest first).
func SortCommentsByDate(comments []models.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
}

func prInPeriod(pr models.PullRequest, period *Chronometer) bool {
	if period.Contains(pr.CreatedAt) {
		return true
	}
	if pr.ClosedAt != nil && period.Contains(*pr.ClosedAt) {
		return true
	}
	if pr.MergedAt != nil && period.Contains(*pr.MergedAt) {
		return true
	}
	return false
}

func issueInPeriod(issue models.Issue, period *Chronometer) bool {
	if period.Contains(issue.CreatedAt) {
		return true
	}
	if issue.ClosedAt != nil && period.Contains(*issue.ClosedAt) {
		return true
	}
	return false
}
