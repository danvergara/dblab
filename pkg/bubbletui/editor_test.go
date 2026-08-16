package bubbletui

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"

	"github.com/stretchr/testify/assert"
)

func TestQueryAtCursor(t *testing.T) {
	content := "SELECT 1;\nSELECT 2;\nSELECT 3;"

	tests := []struct {
		name string
		row  int
		want string
	}{
		{name: "first query at start", row: 0, want: "SELECT 1;"},
		{name: "second query in middle", row: 1, want: "SELECT 2;"},
		{name: "second query on semicolon", row: 1, want: "SELECT 2;"},
		{name: "third query at semicolon", row: 2, want: "SELECT 3;"},
		{name: "third query after end of line", row: 3, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryAtCursor(content, tt.row)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInsertQueriesAndExecute(t *testing.T) {
	queries := []string{"SELECT * FROM users", "SELECT * FROM products", "SELECT * FROM logs"}
	editor, err := createEditorWithQueries(queries)
	assert.Nil(t, err)
	// Run all queries
	editor, cmd := editor.Update(tea.KeyPressMsg{Mod: tea.ModCtrl, Code: 'e'})
	q, ok := cmd().(executeQueryMsg)
	assert.True(t, ok, "Message should be a executeQueryMsg")
	for i, expectedQuery := range queries {
		assert.Equal(t, q.queriesToRun[i],expectedQuery)
	}
}

func TestNavigateInEditorAndRunSingleQuery(t *testing.T) {
	queries := []string{"SELECT * FROM users", "SELECT * FROM products", "SELECT * FROM logs"}
	editor, err := createEditorWithQueries(queries)
	assert.Nil(t, err)
	assert.Equal(t, editor.editor.Line(), 2, "Editor not in expected line")
	editor, _ = editor.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, editor.editor.Line(), 1, "Editor didn't go up")
	// Run second query
	editor, cmd := editor.Update(tea.KeyPressMsg{Mod: tea.ModCtrl, Code: 'r'})
	q, ok := cmd().(executeQueryMsg)
	assert.True(t, ok, "Message should be a executeQueryMsg")
	assert.Len(t, q.queriesToRun, 1, "Should execute a single query.")
	assert.Equal(t, q.queriesToRun[0], queries[1], "Should run the second query")
}

func createEditorWithQueries(queries []string) (Editor, error) {
	qs := strings.Join(queries, ";\n")
	km := keys.DefaultEditorKeyMap()
	km.ExecuteQuery = key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "execute queries in the editor"),
	)
	km.ExecuteSingleQuery = key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "execute queries in the editor"),
	)
	editor := NewEditor(km)
	editor, cmd := editor.Update(tea.KeyPressMsg{Text: "i", Code: 'i'})
	_, ok := cmd().(modeChangeMsg)
	if !ok {
		return Editor{}, errors.New("Cannot change editor mode")
	}
	editor, _ = editor.Update(tea.PasteMsg{Content: qs})
	if editor.editor.Value() != qs {
		return Editor{}, errors.New("Editor content doesn't match: " + editor.editor.Value())
	}
	return editor, nil
}
