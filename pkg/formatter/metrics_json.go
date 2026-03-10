package formatter

import (
	"encoding/json"
	"fmt"

	"yokiyoki/pkg/models"
)

// MetricsJson handles JSON formatting of metrics
type MetricsJson struct {
	metrics []models.Metrics
}

// NewMetricsJson creates a new MetricsJson formatter
func NewMetricsJson(metrics []models.Metrics) *MetricsJson {
	return &MetricsJson{metrics: metrics}
}

// Output outputs metrics in JSON format
func (j *MetricsJson) Output(byUser bool, detailedStats bool) {
	if len(j.metrics) == 0 {
		return
	}

	type metricsRow struct {
		Repository        string `json:"repository"`
		User              string `json:"user,omitempty"`
		Commits           int    `json:"commits"`
		LinesAdded        int    `json:"lines_added"`
		LinesDeleted      int    `json:"lines_deleted"`
		PRsCreated        int    `json:"prs_created"`
		PRsMerged         int    `json:"prs_merged"`
		PRMergeRate       string `json:"pr_merge_rate"`
		AvgPRMergeTime    string `json:"avg_pr_merge_time"`
		IssuesCreated     int    `json:"issues_created"`
		IssuesClosed      int    `json:"issues_closed"`
		IssueResolveRate  string `json:"issue_resolve_rate"`
		AvgIssueCloseTime string `json:"avg_issue_close_time"`
		OpenIssues        int    `json:"open_issues"`
	}

	rows := make([]metricsRow, 0, len(j.metrics))
	for _, m := range j.metrics {
		row := metricsRow{
			Repository:        m.Repository,
			Commits:           m.Commits,
			LinesAdded:        m.LinesAdded,
			LinesDeleted:      m.LinesDeleted,
			PRsCreated:        m.PRsCreated,
			PRsMerged:         m.PRsMerged,
			PRMergeRate:       m.PRMergeRate,
			AvgPRMergeTime:    m.AvgPRMergeTime,
			IssuesCreated:     m.IssuesCreated,
			IssuesClosed:      m.IssuesClosed,
			IssueResolveRate:  m.IssueResolveRate,
			AvgIssueCloseTime: m.AvgIssueCloseTime,
			OpenIssues:        m.OpenIssues,
		}
		if byUser {
			row.User = m.User
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
