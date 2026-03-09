package formatter

import (
	"fmt"
	"strings"

	"yokiyoki/pkg/models"
)

// CommitsTable handles markdown table formatting of commit lists
type CommitsTable struct {
	commits []models.Commit
}

// CommitsTableColumn represents a table column configuration
type CommitsTableColumn struct {
	Header string
	Width  int
	Align  string // "left", "right"
}

// NewCommitsTable creates a new CommitsTable formatter
func NewCommitsTable(commits []models.Commit) *CommitsTable {
	return &CommitsTable{commits: commits}
}

// Output outputs the commit list in markdown table format
func (t *CommitsTable) Output(detailedStats bool) {
	tableData := t.buildTableData(detailedStats)
	columns := t.createColumns(detailedStats)
	t.calculateColumnWidths(columns, tableData)
	t.outputTable(tableData, columns)
}

func (t *CommitsTable) buildTableData(detailedStats bool) [][]string {
	tableData := make([][]string, len(t.commits))
	for i, c := range t.commits {
		tableData[i] = t.toRow(c, detailedStats)
	}
	return tableData
}

func (t *CommitsTable) toRow(c models.Commit, detailedStats bool) []string {
	row := []string{
		c.Repository,
		shortSHA(c.SHA),
		c.Author,
		c.Date.Format("2006-01-02 15:04"),
		truncateMessage(c.Message),
	}

	if detailedStats {
		row = append(row, fmt.Sprintf("+%d/-%d", c.Additions, c.Deletions))
	}

	return row
}

func (t *CommitsTable) createColumns(detailedStats bool) []CommitsTableColumn {
	columns := []CommitsTableColumn{
		{Header: "Repository", Align: "left"},
		{Header: "SHA", Align: "left"},
		{Header: "Author", Align: "left"},
		{Header: "Date", Align: "left"},
		{Header: "Message", Align: "left"},
	}

	if detailedStats {
		columns = append(columns, CommitsTableColumn{Header: "Lines +/-", Align: "left"})
	}

	return columns
}

func (t *CommitsTable) calculateColumnWidths(columns []CommitsTableColumn, tableData [][]string) {
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

func (t *CommitsTable) outputTable(tableData [][]string, columns []CommitsTableColumn) {
	t.printHeader(columns)
	t.printSeparator(columns)
	t.printRows(tableData, columns)
	fmt.Println()
}

func (t *CommitsTable) printHeader(columns []CommitsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf(" %-*s |", col.Width, col.Header)
	}
	fmt.Println()
}

func (t *CommitsTable) printSeparator(columns []CommitsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf("%s|", strings.Repeat("-", col.Width+2))
	}
	fmt.Println()
}

func (t *CommitsTable) printRows(tableData [][]string, columns []CommitsTableColumn) {
	for _, row := range tableData {
		t.printRow(row, columns)
	}
}

func (t *CommitsTable) printRow(row []string, columns []CommitsTableColumn) {
	fmt.Print("|")
	for i, col := range columns {
		value := ""
		if i < len(row) {
			value = row[i]
		}
		fmt.Printf(" %-*s |", col.Width, value)
	}
	fmt.Println()
}

// shortSHA returns the first 7 characters of a commit SHA.
func shortSHA(sha string) string {
	if len(sha) >= 7 {
		return sha[:7]
	}
	return sha
}

// truncateMessage returns the first line of a commit message, truncated to 72 characters.
// If the first line exceeds 72 characters, it is cut at 69 and "..." is appended.
func truncateMessage(msg string) string {
	firstLine := strings.SplitN(msg, "\n", 2)[0]
	if len(firstLine) > 72 {
		return firstLine[:69] + "..."
	}
	return firstLine
}
