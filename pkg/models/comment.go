package models

import "time"

// Comment represents a comment on a PR or Issue
type Comment struct {
	ID          int       `json:"id"`
	Author      string    `json:"author"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	URL         string    `json:"url"`
	IssueNumber int       `json:"issue_number"`
}
