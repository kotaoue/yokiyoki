package formatter

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"yokiyoki/pkg/models"
)

const maxBodyLength = 100

// ConversationFormatter handles text-based output of conversation data
type ConversationFormatter struct {
	conversations []models.Conversation
}

// NewConversationFormatter creates a new ConversationFormatter
func NewConversationFormatter(conversations []models.Conversation) *ConversationFormatter {
	return &ConversationFormatter{conversations: conversations}
}

// Output prints the conversations to stdout grouped by repository
func (f *ConversationFormatter) Output() {
	byRepo := f.groupByRepository()

	for _, repoName := range f.orderedRepositories() {
		convs := byRepo[repoName]
		header := fmt.Sprintf("%s (%d conversations)", repoName, len(convs))
		fmt.Println(header)
		fmt.Println(strings.Repeat("-", len(header)))

		for _, conv := range convs {
			f.printConversation(conv)
		}

		fmt.Println()
	}
}

func (f *ConversationFormatter) printConversation(conv models.Conversation) {
	fmt.Printf("\n%s #%d [%s] by %s - %s (%s)\n",
		conv.Type,
		conv.Number,
		conv.State,
		conv.Author,
		conv.Title,
		conv.CreatedAt.Format("2006-01-02"),
	)

	if len(conv.Comments) == 0 {
		fmt.Println("  (no comments)")
		return
	}

	for _, comment := range conv.Comments {
		body := truncateBody(comment.Body)
		fmt.Printf("  %s (%s): %s\n",
			comment.Author,
			comment.CreatedAt.Format("2006-01-02 15:04"),
			body,
		)
	}
}

func (f *ConversationFormatter) groupByRepository() map[string][]models.Conversation {
	result := make(map[string][]models.Conversation)
	for _, conv := range f.conversations {
		result[conv.Repository] = append(result[conv.Repository], conv)
	}
	return result
}

// orderedRepositories returns repository names in the order they first appear in conversations.
func (f *ConversationFormatter) orderedRepositories() []string {
	repos := make([]string, 0)
	seen := make(map[string]bool)
	for _, conv := range f.conversations {
		if !seen[conv.Repository] {
			repos = append(repos, conv.Repository)
			seen[conv.Repository] = true
		}
	}
	return repos
}

// truncateBody shortens a comment body to maxBodyLength runes, appending "..." if truncated.
func truncateBody(body string) string {
	// Normalize newlines for display
	body = strings.ReplaceAll(body, "\r\n", " ")
	body = strings.ReplaceAll(body, "\n", " ")

	if utf8.RuneCountInString(body) <= maxBodyLength {
		return body
	}

	runes := []rune(body)
	return string(runes[:maxBodyLength]) + "..."
}
