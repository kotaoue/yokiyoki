package formatter_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/models"
)

func sampleComments() []models.Comment {
	return []models.Comment{
		{
			Repository: "owner/repo",
			Type:       "pr",
			Number:     1,
			Title:      "Fix login bug",
			Author:     "alice",
			Body:       "This PR fixes the login issue.",
			URL:        "https://github.com/owner/repo/pull/1",
			CreatedAt:  time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			Repository: "owner/repo",
			Type:       "issue",
			Number:     2,
			Title:      "Login broken",
			Author:     "bob",
			Body:       "The login page is broken.",
			URL:        "https://github.com/owner/repo/issues/2",
			CreatedAt:  time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC),
		},
	}
}

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewConversationsTable(t *testing.T) {
	table := formatter.NewConversationsTable(sampleComments())
	assert.NotNil(t, table)
}

func TestConversationsTable_Output(t *testing.T) {
	output := captureOutput(func() {
		table := formatter.NewConversationsTable(sampleComments())
		table.Output()
	})

	assert.Contains(t, output, "Repository")
	assert.Contains(t, output, "Type")
	assert.Contains(t, output, "Author")
	assert.Contains(t, output, "Body")
	assert.Contains(t, output, "owner/repo")
	assert.Contains(t, output, "alice")
	assert.Contains(t, output, "Fix login bug")
	assert.Contains(t, output, "This PR fixes the login issue.")
}

func TestConversationsTable_Output_Empty(t *testing.T) {
	output := captureOutput(func() {
		table := formatter.NewConversationsTable([]models.Comment{})
		table.Output()
	})

	assert.Contains(t, output, "Repository")
	assert.Contains(t, output, "Author")
}

func TestNewConversationsCsv(t *testing.T) {
	csv := formatter.NewConversationsCsv(sampleComments())
	assert.NotNil(t, csv)
}

func TestConversationsCsv_Output(t *testing.T) {
	output := captureOutput(func() {
		csv := formatter.NewConversationsCsv(sampleComments())
		csv.Output()
	})

	assert.Contains(t, output, "Repository,Type,Number,Title,Author,Date,Body,URL")
	assert.Contains(t, output, "owner/repo")
	assert.Contains(t, output, "alice")
	assert.Contains(t, output, "pr")
}

func TestConversationsCsv_Output_Empty(t *testing.T) {
	output := captureOutput(func() {
		csv := formatter.NewConversationsCsv([]models.Comment{})
		csv.Output()
	})

	assert.Empty(t, output)
}

func TestNewConversationsJson(t *testing.T) {
	j := formatter.NewConversationsJson(sampleComments())
	assert.NotNil(t, j)
}

func TestConversationsJson_Output(t *testing.T) {
	output := captureOutput(func() {
		j := formatter.NewConversationsJson(sampleComments())
		j.Output()
	})

	assert.Contains(t, output, `"repository"`)
	assert.Contains(t, output, `"type"`)
	assert.Contains(t, output, `"author"`)
	assert.Contains(t, output, `"body"`)
	assert.Contains(t, output, "owner/repo")
	assert.Contains(t, output, "alice")
}

func TestConversationsJson_Output_Empty(t *testing.T) {
	output := captureOutput(func() {
		j := formatter.NewConversationsJson([]models.Comment{})
		j.Output()
	})

	assert.Empty(t, output)
}
