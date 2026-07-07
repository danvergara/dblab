package bubbletui

import (
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestResulset_UpdateBeforeResize(t *testing.T) {
	kb := command.DefaultKeyMap()
	rs := NewResultSet(kb)
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
