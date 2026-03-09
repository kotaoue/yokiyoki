package formatter

import (
	"fmt"
	"strings"
	"time"

	"yokiyoki/pkg/models"
)

// CommitsCsv handles CSV formatting of commit lists
type CommitsCsv struct {
	commits []models.Commit
}

// NewCommitsCsv creates a new CommitsCsv formatter
func NewCommitsCsv(commits []models.Commit) *CommitsCsv {
	return &CommitsCsv{commits: commits}
}

// Output outputs the commit list in CSV format
func (c *CommitsCsv) Output(detailedStats bool) {
	if len(c.commits) == 0 {
		return
	}

	headers := c.header(detailedStats)
	fmt.Println(strings.Join(headers, ","))

	for _, commit := range c.commits {
		fmt.Println(c.toCSV(commit, detailedStats))
	}
}

func (c *CommitsCsv) header(detailedStats bool) []string {
	headers := []string{
		"Repository",
		"SHA",
		"Author",
		"Date",
		"Message",
	}

	if detailedStats {
		headers = append(headers, "LinesAdded", "LinesDeleted")
	}

	return headers
}

func (c *CommitsCsv) toCSV(commit models.Commit, detailedStats bool) string {
	values := []string{
		commit.Repository,
		shortSHA(commit.SHA),
		commit.Author,
		commit.Date.Format(time.RFC3339),
		escapeCsvField(truncateMessage(commit.Message)),
	}

	if detailedStats {
		values = append(values,
			fmt.Sprintf("%d", commit.Additions),
			fmt.Sprintf("%d", commit.Deletions),
		)
	}

	return strings.Join(values, ",")
}

// escapeCsvField wraps a field in double quotes if it contains a comma, double
// quote, or newline, escaping any embedded double quotes.
func escapeCsvField(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
