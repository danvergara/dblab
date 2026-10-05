package bubbletui

import (
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResulset_UpdateBeforeResize(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)
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

func TestResultSet_EmptyIndexesAndConstraints(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)
	rs.SetSize(80, 24)

	meta := &client.Metadata{
		TableContent: client.Table{
			Columns: []string{"id", "title"},
			Rows:    [][]string{{"1", "Item 1"}},
		},
		Structure: client.Table{
			Columns: []string{"Field", "Type"},
			Rows:    [][]string{{"id", "int"}},
		},
		Indexes: client.Table{
			Columns: []string{"index_name", "column_name"},
			Rows:    [][]string{},
		},
		Constraints: client.Table{
			Columns: []string{"constraint_name", "constraint_type"},
			Rows:    [][]string{},
		},
	}

	rs, _ = rs.Update(metadataSuccessMsg{metadata: meta, isTable: true})

	// Indexes tab (index 2) should be a TextPanel showing NoIndexesMsg.
	idxPanel, ok := rs.tablesMetadata[2].(*TextPanel)
	require.True(t, ok, "indexes panel should be a *TextPanel when table has no indexes")
	assert.Equal(t, NoIndexesMsg, idxPanel.View().Content)

	// Constraints tab (index 3) should be a TextPanel showing NoConstraintsMsg.
	constPanel, ok := rs.tablesMetadata[3].(*TextPanel)
	require.True(t, ok, "constraints panel should be a *TextPanel when table has no constraints")
	assert.Equal(t, NoConstraintsMsg, constPanel.View().Content)

	// Navigate to Indexes tab (tab index 2).
	tabKey := tea.KeyPressMsg{Code: tea.KeyTab}
	rs, _ = rs.Update(tabKey) // activeTab = 1 (Columns)
	assert.Equal(t, 1, rs.activeTab)
	rs, _ = rs.Update(tabKey) // activeTab = 2 (Indexes)
	assert.Equal(t, 2, rs.activeTab)
	assert.Contains(t, rs.viewport.View(), NoIndexesMsg)

	// Navigate to Constraints tab (tab index 3).
	rs, _ = rs.Update(tabKey) // activeTab = 3 (Constraints)
	assert.Equal(t, 3, rs.activeTab)
	assert.Contains(t, rs.viewport.View(), NoConstraintsMsg)

	// Resizing does not panic or alter placeholder panels.
	assert.NotPanics(t, func() {
		rs.SetSize(100, 30)
	})

	// Scrolling actions do not panic on TextPanel.
	assert.NotPanics(t, func() {
		rs.Update(tea.KeyPressMsg{Code: 'g'})
		rs.Update(tea.KeyPressMsg{Code: 'G'})
		rs.Update(tea.KeyPressMsg{Code: '$'})
	})
}

func TestResultSet_PopulatedIndexesAndConstraints(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)
	rs.SetSize(80, 24)

	meta := &client.Metadata{
		TableContent: client.Table{
			Columns: []string{"id"},
			Rows:    [][]string{{"1"}},
		},
		Structure: client.Table{
			Columns: []string{"Field"},
			Rows:    [][]string{{"id"}},
		},
		Indexes: client.Table{
			Columns: []string{"index_name", "column_name"},
			Rows:    [][]string{{"idx_id", "id"}},
		},
		Constraints: client.Table{
			Columns: []string{"constraint_name", "constraint_type"},
			Rows:    [][]string{{"pk_id", "PRIMARY KEY"}},
		},
	}

	rs, _ = rs.Update(metadataSuccessMsg{metadata: meta, isTable: true})

	_, okIdx := rs.tablesMetadata[2].(*TablePanel)
	assert.True(t, okIdx, "indexes panel should be a *TablePanel when indexes exist")

	_, okConst := rs.tablesMetadata[3].(*TablePanel)
	assert.True(t, okConst, "constraints panel should be a *TablePanel when constraints exist")
}

func TestResultSet_EmptyIndexes_WithConstraints(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)

	meta := &client.Metadata{
		TableContent: client.Table{
			Columns: []string{"id"},
			Rows:    [][]string{{"1"}},
		},
		Structure: client.Table{
			Columns: []string{"Field"},
			Rows:    [][]string{{"id"}},
		},
		Indexes: client.Table{
			Columns: []string{"index_name", "column_name"},
			Rows:    nil,
		},
		Constraints: client.Table{
			Columns: []string{"constraint_name", "constraint_type"},
			Rows:    [][]string{{"pk_id", "PRIMARY KEY"}},
		},
	}

	rs, _ = rs.Update(metadataSuccessMsg{metadata: meta, isTable: true})

	idxPanel, okIdx := rs.tablesMetadata[2].(*TextPanel)
	require.True(t, okIdx, "indexes panel should be *TextPanel")
	assert.Equal(t, NoIndexesMsg, idxPanel.View().Content)

	_, okConst := rs.tablesMetadata[3].(*TablePanel)
	assert.True(t, okConst, "constraints panel should be *TablePanel")
}

func TestResultSet_WithIndexes_EmptyConstraints(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)

	meta := &client.Metadata{
		TableContent: client.Table{
			Columns: []string{"id"},
			Rows:    [][]string{{"1"}},
		},
		Structure: client.Table{
			Columns: []string{"Field"},
			Rows:    [][]string{{"id"}},
		},
		Indexes: client.Table{
			Columns: []string{"index_name", "column_name"},
			Rows:    [][]string{{"idx_id", "id"}},
		},
		Constraints: client.Table{
			Columns: []string{"constraint_name", "constraint_type"},
			Rows:    nil,
		},
	}

	rs, _ = rs.Update(metadataSuccessMsg{metadata: meta, isTable: true})

	_, okIdx := rs.tablesMetadata[2].(*TablePanel)
	assert.True(t, okIdx, "indexes panel should be *TablePanel")

	constPanel, okConst := rs.tablesMetadata[3].(*TextPanel)
	require.True(t, okConst, "constraints panel should be *TextPanel")
	assert.Equal(t, NoConstraintsMsg, constPanel.View().Content)
}

func TestResultSet_SwitchBetweenTables(t *testing.T) {
	var km = keys.DefaultKeyMap()
	rs := NewResultSet(km.ResultSet)
	rs.SetSize(80, 24)

	emptyMeta := &client.Metadata{
		TableContent: client.Table{Columns: []string{"id"}, Rows: [][]string{{"1"}}},
		Structure:    client.Table{Columns: []string{"Field"}, Rows: [][]string{{"id"}}},
		Indexes:      client.Table{Columns: []string{"index_name"}, Rows: nil},
		Constraints:  client.Table{Columns: []string{"constraint_name"}, Rows: nil},
	}

	populatedMeta := &client.Metadata{
		TableContent: client.Table{Columns: []string{"id"}, Rows: [][]string{{"1"}}},
		Structure:    client.Table{Columns: []string{"Field"}, Rows: [][]string{{"id"}}},
		Indexes:      client.Table{Columns: []string{"index_name"}, Rows: [][]string{{"idx_id"}}},
		Constraints:  client.Table{Columns: []string{"constraint_name"}, Rows: [][]string{{"pk_id"}}},
	}

	// 1. Initial table without indexes or constraints
	rs, _ = rs.Update(metadataSuccessMsg{metadata: emptyMeta, isTable: true})
	_, ok1Idx := rs.tablesMetadata[2].(*TextPanel)
	_, ok1Const := rs.tablesMetadata[3].(*TextPanel)
	assert.True(t, ok1Idx)
	assert.True(t, ok1Const)

	// 2. Switch to table with indexes and constraints
	rs, _ = rs.Update(metadataSuccessMsg{metadata: populatedMeta, isTable: true})
	_, ok2Idx := rs.tablesMetadata[2].(*TablePanel)
	_, ok2Const := rs.tablesMetadata[3].(*TablePanel)
	assert.True(t, ok2Idx)
	assert.True(t, ok2Const)

	// 3. Switch back to table without indexes or constraints
	rs, _ = rs.Update(metadataSuccessMsg{metadata: emptyMeta, isTable: true})
	idxPanel, ok3Idx := rs.tablesMetadata[2].(*TextPanel)
	constPanel, ok3Const := rs.tablesMetadata[3].(*TextPanel)
	assert.True(t, ok3Idx)
	assert.True(t, ok3Const)
	assert.Equal(t, NoIndexesMsg, idxPanel.View().Content)
	assert.Equal(t, NoConstraintsMsg, constPanel.View().Content)
}
