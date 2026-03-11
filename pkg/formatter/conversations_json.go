package formatter

import (
	"encoding/json"
	"fmt"
	"time"

	"yokiyoki/pkg/models"
)

// ConversationsJson handles JSON formatting of conversation comment lists.
type ConversationsJson struct {
	comments []models.Comment
}

// NewConversationsJson creates a new ConversationsJson formatter.
func NewConversationsJson(comments []models.Comment) *ConversationsJson {
	return &ConversationsJson{comments: comments}
}

// Output outputs the conversation list in JSON format.
func (j *ConversationsJson) Output() {
	if len(j.comments) == 0 {
		return
	}

	type commentRow struct {
		Repository string    `json:"repository"`
		Type       string    `json:"type"`
		Number     int       `json:"number"`
		Title      string    `json:"title"`
		Author     string    `json:"author"`
		CreatedAt  time.Time `json:"created_at"`
		Body       string    `json:"body"`
		URL        string    `json:"url"`
	}

	rows := make([]commentRow, 0, len(j.comments))
	for _, c := range j.comments {
		rows = append(rows, commentRow{
			Repository: c.Repository,
			Type:       c.Type,
			Number:     c.Number,
			Title:      c.Title,
			Author:     c.Author,
			CreatedAt:  c.CreatedAt,
			Body:       c.Body,
			URL:        c.URL,
		})
	}

	out, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
	fmt.Println(string(out))
}
