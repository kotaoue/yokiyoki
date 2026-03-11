package formatter

import (
	"fmt"
	"strings"
	"time"

	"yokiyoki/pkg/models"
)

// ConversationsCsv handles CSV formatting of conversation comment lists.
type ConversationsCsv struct {
	comments []models.Comment
}

// NewConversationsCsv creates a new ConversationsCsv formatter.
func NewConversationsCsv(comments []models.Comment) *ConversationsCsv {
	return &ConversationsCsv{comments: comments}
}

// Output outputs the conversation list in CSV format.
func (c *ConversationsCsv) Output() {
	if len(c.comments) == 0 {
		return
	}

	headers := []string{"Repository", "Type", "Number", "Title", "Author", "Date", "Body", "URL"}
	fmt.Println(strings.Join(headers, ","))

	for _, comment := range c.comments {
		fmt.Println(c.toCSV(comment))
	}
}

func (c *ConversationsCsv) toCSV(comment models.Comment) string {
	values := []string{
		comment.Repository,
		comment.Type,
		fmt.Sprintf("%d", comment.Number),
		escapeCsvField(comment.Title),
		comment.Author,
		comment.CreatedAt.Format(time.RFC3339),
		escapeCsvField(comment.Body),
		comment.URL,
	}
	return strings.Join(values, ",")
}
