package formatter

import (
	"fmt"
	"strings"

	"yokiyoki/pkg/models"
)

// MetricsTable handles markdown table formatting of metrics
type MetricsTable struct {
	metrics []models.Metrics
}

// MetricsTableColumn represents a table column configuration
type MetricsTableColumn struct {
	Header string
	Width  int
	Align  string // "left", "right"
}

// NewMetricsTable creates a new MetricsTable formatter
func NewMetricsTable(metrics []models.Metrics) *MetricsTable {
	return &MetricsTable{metrics: metrics}
}

// Output outputs metrics in markdown table format
func (t *MetricsTable) Output(byUser bool, detailedStats bool) {
	tableData := t.buildTableData(byUser, detailedStats)
	columns := t.createColumns(byUser, detailedStats)
	t.calculateColumnWidths(columns, tableData)
	t.outputTable(tableData, columns)
}

func (t *MetricsTable) buildTableData(byUser bool, detailedStats bool) [][]string {
	tableData := make([][]string, len(t.metrics))
	for i, m := range t.metrics {
		row := t.toMarkdownRow(m, byUser, detailedStats)
		tableData[i] = row
	}
	return tableData
}

func (t *MetricsTable) toMarkdownRow(m models.Metrics, byUser bool, detailedStats bool) []string {
	linesStr := t.formatLines(m)
	prsStr := t.formatPRs(m)
	issuesStr := t.formatIssues(m)
	avgPRMergeTime := t.formatTime(m.AvgPRMergeTime)
	avgIssueCloseTime := t.formatTime(m.AvgIssueCloseTime)

	row := []string{m.Repository}

	if byUser {
		row = append(row, m.User)
	}

	row = append(row,
		fmt.Sprintf("%d", m.Commits),
		prsStr,
		avgPRMergeTime,
		issuesStr,
		avgIssueCloseTime,
		fmt.Sprintf("%d", m.OpenIssues),
	)

	if detailedStats {
		row = append(row, linesStr)
	}

	return row
}

func (t *MetricsTable) createColumns(byUser bool, detailedStats bool) []MetricsTableColumn {
	columns := []MetricsTableColumn{
		{Header: "Repository", Align: "left"},
	}

	if byUser {
		columns = append(columns, MetricsTableColumn{Header: "User", Align: "left"})
	}

	columns = append(columns,
		MetricsTableColumn{Header: "Commits", Align: "right"},
		MetricsTableColumn{Header: "PR Merge Rate", Align: "left"},
		MetricsTableColumn{Header: "PR Merge Time", Align: "left"},
		MetricsTableColumn{Header: "Issue Resolve Rate", Align: "left"},
		MetricsTableColumn{Header: "Issue Resolve Time", Align: "left"},
		MetricsTableColumn{Header: "Active Issues", Align: "right"},
	)

	if detailedStats {
		columns = append(columns, MetricsTableColumn{Header: "Lines +/-", Align: "left"})
	}

	return columns
}

func (t *MetricsTable) calculateColumnWidths(columns []MetricsTableColumn, tableData [][]string) {
	for i, col := range columns {
		columns[i].Width = len(col.Header)
		for _, row := range tableData {
			if i >= len(row) || len(row[i]) <= columns[i].Width {
				continue
			}
			columns[i].Width = len(row[i])
		}
	}
}

func (t *MetricsTable) outputTable(tableData [][]string, columns []MetricsTableColumn) {
	t.printHeader(columns)
	t.printSeparator(columns)
	t.printRows(tableData, columns)
	fmt.Println()
}

func (t *MetricsTable) printHeader(columns []MetricsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf(" %-*s |", col.Width, col.Header)
	}
	fmt.Println()
}

func (t *MetricsTable) printSeparator(columns []MetricsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf("%s|", strings.Repeat("-", col.Width+2))
	}
	fmt.Println()
}

func (t *MetricsTable) printRows(tableData [][]string, columns []MetricsTableColumn) {
	for _, row := range tableData {
		t.printRow(row, columns)
	}
}

func (t *MetricsTable) printRow(row []string, columns []MetricsTableColumn) {
	fmt.Print("|")
	for i, col := range columns {
		cellValue := t.getCellValue(row, i)
		t.printCell(cellValue, col)
	}
	fmt.Println()
}

func (t *MetricsTable) printCell(value string, col MetricsTableColumn) {
	if col.Align == "right" {
		fmt.Printf(" %*s |", col.Width, value)
	} else {
		fmt.Printf(" %-*s |", col.Width, value)
	}
}

func (t *MetricsTable) getCellValue(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}

func (t *MetricsTable) formatLines(m models.Metrics) string {
	return fmt.Sprintf("+%d/-%d", m.LinesAdded, m.LinesDeleted)
}

func (t *MetricsTable) formatPRs(m models.Metrics) string {
	if m.PRsCreated == 0 {
		return "-/-"
	}

	if m.PRMergeRate == "None" {
		return fmt.Sprintf("%d/%d", m.PRsMerged, m.PRsCreated)
	}

	return fmt.Sprintf("%d/%d (%s)", m.PRsMerged, m.PRsCreated, m.PRMergeRate)
}

func (t *MetricsTable) formatIssues(m models.Metrics) string {
	if m.IssuesCreated == 0 {
		return "-/-"
	}

	if m.IssueResolveRate == "None" {
		return fmt.Sprintf("%d/%d", m.IssuesClosed, m.IssuesCreated)
	}

	return fmt.Sprintf("%d/%d (%s)", m.IssuesClosed, m.IssuesCreated, m.IssueResolveRate)
}

func (t *MetricsTable) formatTime(timeStr string) string {
	if timeStr == "None" {
		return "-"
	}
	return timeStr
}
