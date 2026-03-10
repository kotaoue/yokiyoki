package formatter

import (
	"encoding/json"
	"fmt"
	"time"

	"yokiyoki/pkg/models"
)

// CommitsJson handles JSON formatting of commit lists
type CommitsJson struct {
	commits []models.Commit
}

// NewCommitsJson creates a new CommitsJson formatter
func NewCommitsJson(commits []models.Commit) *CommitsJson {
	return &CommitsJson{commits: commits}
}

// Output outputs the commit list in JSON format
func (j *CommitsJson) Output(detailedStats bool) {
	if len(j.commits) == 0 {
		return
	}

	type commitRow struct {
		Repository string    `json:"repository"`
		SHA        string    `json:"sha"`
		Author     string    `json:"author"`
		Date       time.Time `json:"date"`
		Message    string    `json:"message"`
		Additions  *int      `json:"additions,omitempty"`
		Deletions  *int      `json:"deletions,omitempty"`
	}

	rows := make([]commitRow, 0, len(j.commits))
	for _, c := range j.commits {
		row := commitRow{
			Repository: c.Repository,
			SHA:        shortSHA(c.SHA),
			Author:     c.Author,
			Date:       c.Date,
			Message:    c.Message,
		}
		if detailedStats {
			additions := c.Additions
			deletions := c.Deletions
			row.Additions = &additions
			row.Deletions = &deletions
		}
		rows = append(rows, row)
	}

	out, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
	fmt.Println(string(out))
}
