package services

import (
	"fmt"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
)

// ConversationOptions represents configuration for conversation collection
type ConversationOptions struct {
	Period *Chronometer
}

// ExecuteConversation fetches PR and Issue conversations for the given repository
func ExecuteConversation(repo models.Repository, options ConversationOptions) []models.Conversation {
	repoFullName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)

	prs := repository.GetPullRequests(repo, options.Period.StartTime())
	issues := repository.GetIssues(repo, options.Period.StartTime())
	comments := repository.GetComments(repo, options.Period.StartTime())

	commentsByNumber := groupCommentsByIssueNumber(comments)

	var conversations []models.Conversation

	for _, pr := range prs {
		if !options.Period.Contains(pr.CreatedAt) {
			continue
		}
		conv := models.Conversation{
			Repository: repoFullName,
			Type:       models.ConversationTypePR,
			Number:     pr.Number,
			Title:      pr.Title,
			State:      pr.State,
			Author:     pr.Author,
			URL:        pr.URL,
			CreatedAt:  pr.CreatedAt,
			Comments:   commentsByNumber[pr.Number],
		}
		conversations = append(conversations, conv)
	}

	for _, issue := range issues {
		if !options.Period.Contains(issue.CreatedAt) {
			continue
		}
		conv := models.Conversation{
			Repository: repoFullName,
			Type:       models.ConversationTypeIssue,
			Number:     issue.Number,
			Title:      issue.Title,
			State:      issue.State,
			Author:     issue.Author,
			URL:        issue.URL,
			CreatedAt:  issue.CreatedAt,
			Comments:   commentsByNumber[issue.Number],
		}
		conversations = append(conversations, conv)
	}

	return conversations
}

// groupCommentsByIssueNumber groups comments by the issue/PR number they belong to
func groupCommentsByIssueNumber(comments []models.Comment) map[int][]models.Comment {
	result := make(map[int][]models.Comment)
	for _, c := range comments {
		result[c.IssueNumber] = append(result[c.IssueNumber], c)
	}
	return result
}
