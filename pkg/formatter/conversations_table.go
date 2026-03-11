package formatter

import (
	"fmt"
	"strings"

	"yokiyoki/pkg/models"
)

// ConversationsTable handles markdown table formatting of conversation comment lists.
type ConversationsTable struct {
	comments []models.Comment
}

// ConversationsTableColumn represents a table column configuration.
type ConversationsTableColumn struct {
	Header string
	Width  int
}

// NewConversationsTable creates a new ConversationsTable formatter.
func NewConversationsTable(comments []models.Comment) *ConversationsTable {
	return &ConversationsTable{comments: comments}
}

// Output outputs the conversation list in markdown table format.
func (t *ConversationsTable) Output() {
	tableData := t.buildTableData()
	columns := t.createColumns()
	t.calculateColumnWidths(columns, tableData)
	t.outputTable(tableData, columns)
}

func (t *ConversationsTable) buildTableData() [][]string {
	tableData := make([][]string, len(t.comments))
	for i, c := range t.comments {
		tableData[i] = t.toRow(c)
	}
	return tableData
}

func (t *ConversationsTable) toRow(c models.Comment) []string {
	return []string{
		c.Repository,
		c.Type,
		fmt.Sprintf("%d", c.Number),
		c.Title,
		c.Author,
		c.CreatedAt.Format("2006-01-02 15:04"),
		truncateMessage(c.Body),
	}
}

func (t *ConversationsTable) createColumns() []ConversationsTableColumn {
	return []ConversationsTableColumn{
		{Header: "Repository"},
		{Header: "Type"},
		{Header: "#"},
		{Header: "Title"},
		{Header: "Author"},
		{Header: "Date"},
		{Header: "Body"},
	}
}

func (t *ConversationsTable) calculateColumnWidths(columns []ConversationsTableColumn, tableData [][]string) {
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

func (t *ConversationsTable) outputTable(tableData [][]string, columns []ConversationsTableColumn) {
	t.printHeader(columns)
	t.printSeparator(columns)
	t.printRows(tableData, columns)
	fmt.Println()
}

func (t *ConversationsTable) printHeader(columns []ConversationsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf(" %-*s |", col.Width, col.Header)
	}
	fmt.Println()
}

func (t *ConversationsTable) printSeparator(columns []ConversationsTableColumn) {
	fmt.Print("|")
	for _, col := range columns {
		fmt.Printf("%s|", strings.Repeat("-", col.Width+2))
	}
	fmt.Println()
}

func (t *ConversationsTable) printRows(tableData [][]string, columns []ConversationsTableColumn) {
	for _, row := range tableData {
		t.printRow(row, columns)
	}
}

func (t *ConversationsTable) printRow(row []string, columns []ConversationsTableColumn) {
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
