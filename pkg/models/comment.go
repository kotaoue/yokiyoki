package models

import "time"

// Comment represents a comment (or initial body) in a PR or issue conversation.
type Comment struct {
	Repository string    `json:"repository"`
	Type       string    `json:"type"` // "pr" or "issue"
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Body       string    `json:"body"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
}
