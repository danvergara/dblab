package bubbletui

import (
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/stretchr/testify/assert"
)

func TestResulset_UpdateBeforeResize(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet, ResultSetTheme{})
	if tp, ok := rs.tablesMetadata[0].(*TablePanel); ok {
		cols := []table.Column{{Title: "id", Width: 15}, {Title: "name", Width: 15}}
		rows := []table.Row{
			{"1", "alice"},
			{"2", "bob"},
			{"3", "charlie"},
		}
		tp.table.SetColumns(cols)
		tp.table.SetRows(rows)
	}
	msg := tea.KeyPressMsg{Code: tea.KeyDown}
	assert.NotPanics(t, func() {
		rs.Update(msg)
	})
}

func TestCalculateTotalTableWidth(t *testing.T) {
	tests := []struct {
		name     string
		columns  []table.Column
		expected int
	}{
		{
			name:     "empty columns",
			columns:  nil,
			expected: 0,
		},
		{
			name: "single column",
			columns: []table.Column{
				{Title: "id", Width: 20},
			},
			expected: 22, // 20 + 2 padding
		},
		{
			name: "multiple columns",
			columns: []table.Column{
				{Title: "id", Width: 10},
				{Title: "name", Width: 20},
				{Title: "email", Width: 30},
			},
			expected: 66, // (10+2) + (20+2) + (30+2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTotalTableWidth(tt.columns)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPopulateTable(t *testing.T) {
	headers := []string{"id", "title", "a_very_long_header_title"}
	data := [][]string{
		{"1", "Hello", "Short"},
		{"2", "A rather long title that exceeds twenty chars", "Some content here that is also pretty long"},
	}

	cols, rows := populateTable(headers, data)

	assert.Len(t, cols, 3)
	assert.Len(t, rows, 2)

	// col 0: "id", max(2+5, 20) -> 20
	assert.Equal(t, "id", cols[0].Title)
	assert.Equal(t, 20, cols[0].Width)

	// col 1: "title", max(5+5, 20) -> 20 (bounded to keep views compact despite long data)
	assert.Equal(t, "title", cols[1].Title)
	assert.Equal(t, 20, cols[1].Width)

	// col 2: "a_very_long_header_title", max(24+5, 20) -> 29
	assert.Equal(t, "a_very_long_header_title", cols[2].Title)
	assert.Equal(t, len("a_very_long_header_title")+5, cols[2].Width)

	// Check mismatched column count in row does not panic
	mismatchedData := [][]string{
		{"1", "Hello", "Content", "ExtraCol1", "ExtraCol2"},
	}
	assert.NotPanics(t, func() {
		mCols, mRows := populateTable(headers, mismatchedData)
		assert.Len(t, mCols, 3)
		assert.Len(t, mRows, 1)
	})
}

func TestTablePanel_LastColumnNotTruncated(t *testing.T) {
	headers := []string{"col1", "col2", "col3", "col4", "col5", "last_col"}
	data := [][]string{
		{"val1", "val2", "val3", "val4", "val5", "unique_last_val_12345"},
	}

	cols, rows := populateTable(headers, data)
	panel := newTablePanel(10, 200, lipgloss.Color("#ffffff"), lipgloss.Color("#ffffff"), lipgloss.Color("#ffffff"), lipgloss.Color("#000000"))

	panel.table.SetWidth(calculateTotalTableWidth(cols))
	panel.table.SetColumns(cols)
	panel.table.SetRows(rows)

	view := panel.table.View()
	assert.Contains(t, view, "unique_last_val_12345")
}
