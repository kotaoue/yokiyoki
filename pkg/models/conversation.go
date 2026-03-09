package models

import "time"

// ConversationType represents whether a conversation is on a PR or Issue
type ConversationType string

const (
	ConversationTypePR    ConversationType = "PR"
	ConversationTypeIssue ConversationType = "Issue"
)

// Conversation represents a PR or Issue with its associated comments
type Conversation struct {
	Repository string
	Type       ConversationType
	Number     int
	Title      string
	State      string
	Author     string
	URL        string
	CreatedAt  time.Time
	Comments   []Comment
}
