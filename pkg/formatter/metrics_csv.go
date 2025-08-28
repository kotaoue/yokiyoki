package formatter

import (
	"fmt"
	"strings"

	"yokiyoki/pkg/models"
)

// MetricsCsv handles CSV formatting of metrics
type MetricsCsv struct {
	metrics []models.Metrics
}

// NewMetricsCsv creates a new MetricsCsv formatter
func NewMetricsCsv(metrics []models.Metrics) *MetricsCsv {
	return &MetricsCsv{metrics: metrics}
}

// Output outputs metrics in CSV format
func (c *MetricsCsv) Output(byUser bool, detailedStats bool) {
	if len(c.metrics) == 0 {
		return
	}

	headers := c.header(byUser)
	fmt.Println(strings.Join(headers, ","))

	for _, m := range c.metrics {
		fmt.Println(c.toCSV(m, byUser))
	}
}

func (c *MetricsCsv) header(includeUser bool) []string {
	headers := []string{"Repository"}

	if includeUser {
		headers = append(headers, "User")
	}

	headers = append(headers,
		"Commits",
		"LinesAdded",
		"LinesDeleted",
		"PRsCreated",
		"PRsMerged",
		"PRMergeRate",
		"AvgPRMergeTime",
		"IssuesCreated",
		"IssuesClosed",
		"IssueResolveRate",
		"AvgIssueCloseTime",
		"OpenIssues")

	return headers
}

func (c *MetricsCsv) toCSV(m models.Metrics, includeUser bool) string {
	values := c.toSlice(m, includeUser)
	return strings.Join(values, ",")
}

func (c *MetricsCsv) toSlice(m models.Metrics, includeUser bool) []string {
	values := []string{m.Repository}

	if includeUser {
		values = append(values, m.User)
	}

	values = append(values,
		fmt.Sprintf("%d", m.Commits),
		fmt.Sprintf("%d", m.LinesAdded),
		fmt.Sprintf("%d", m.LinesDeleted),
		fmt.Sprintf("%d", m.PRsCreated),
		fmt.Sprintf("%d", m.PRsMerged),
		m.PRMergeRate,
		m.AvgPRMergeTime,
		fmt.Sprintf("%d", m.IssuesCreated),
		fmt.Sprintf("%d", m.IssuesClosed),
		m.IssueResolveRate,
		m.AvgIssueCloseTime,
		fmt.Sprintf("%d", m.OpenIssues))

	return values
}
