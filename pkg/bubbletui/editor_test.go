package bubbletui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestEditor returns an editor preloaded with content and the cursor parked
// at the top of the buffer.
func newTestEditor(t *testing.T, content string) Editor {
	t.Helper()

	e := NewEditor(keys.DefaultEditorKeyMap())
	e.SetWidth(80)
	e.SetHeight(10)
	e.editor.SetValue(content)
	e.editor.MoveToBegin()

	return e
}

// press sends a printable key the way a terminal reports one: Text carries the
// resulting character, which is what both key.Matches and the textarea read.
func press(e Editor, r rune) Editor {
	e, _ = e.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	return e
}

func pressAll(e Editor, s string) Editor {
	for _, r := range s {
		e = press(e, r)
	}

	return e
}

func pressEsc(e Editor) Editor {
	e, _ = e.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	return e
}

// moveTo parks the cursor on a logical row and column.
func moveTo(e *Editor, row, col int) {
	e.setCursor(row, col)
}

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

func TestFirstNonBlank(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{name: "no leading whitespace", line: "SELECT 1;", want: 0},
		{name: "leading spaces", line: "    SELECT 1;", want: 4},
		{name: "leading tab", line: "\tSELECT 1;", want: 1},
		{name: "blank line", line: "   ", want: 3},
		{name: "empty line", line: "", want: 0},
		{name: "counts runes not bytes", line: "  ñSELECT", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, firstNonBlank(tt.line))
		})
	}
}

func TestWordForward(t *testing.T) {
	tests := []struct {
		name               string
		content            string
		startRow, startCol int
		keys               string
		wantRow, wantCol   int
	}{
		{
			name:    "hops over words and treats punctuation as its own word",
			content: "SELECT foo, bar FROM baz;", startRow: 0, startCol: 0,
			keys: "w", wantRow: 0, wantCol: 7, // start of "foo"
		},
		{
			name:    "lands on a punctuation run",
			content: "SELECT foo, bar FROM baz;", startRow: 0, startCol: 0,
			keys: "ww", wantRow: 0, wantCol: 10, // the ","
		},
		{
			name:    "skips the comma and following space",
			content: "SELECT foo, bar FROM baz;", startRow: 0, startCol: 0,
			keys: "www", wantRow: 0, wantCol: 12, // start of "bar"
		},
		{
			name:    "reaches the last word",
			content: "SELECT foo, bar FROM baz;", startRow: 0, startCol: 0,
			keys: "wwwwww", wantRow: 0, wantCol: 24, // the ";"
		},
		{
			name:    "does not move past the last word in the buffer",
			content: "SELECT foo, bar FROM baz;", startRow: 0, startCol: 0,
			keys: "wwwwwww", wantRow: 0, wantCol: 24,
		},
		{
			name:    "crosses a line boundary",
			content: "foo\nbar", startRow: 0, startCol: 2,
			keys: "w", wantRow: 1, wantCol: 0,
		},
		{
			name:    "stops on an empty line rather than skipping it",
			content: "foo\n\nbar", startRow: 0, startCol: 0,
			keys: "w", wantRow: 1, wantCol: 0,
		},
		{
			name:    "moves past an empty line onto the next word",
			content: "foo\n\nbar", startRow: 0, startCol: 0,
			keys: "ww", wantRow: 2, wantCol: 0,
		},
		{
			name:    "treats parentheses as one-character words",
			content: "foo(bar)", startRow: 0, startCol: 0,
			keys: "www", wantRow: 0, wantCol: 7, // the ")"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEditor(t, tt.content)
			moveTo(&e, tt.startRow, tt.startCol)

			e = pressAll(e, tt.keys)

			assert.Equal(t, tt.content, e.editor.Value(), "word motions must not mutate the buffer")
			assert.Equal(t, tt.wantRow, e.editor.Line(), "cursor row")
			assert.Equal(t, tt.wantCol, e.editor.Column(), "cursor column")
		})
	}
}

func TestWordEnd(t *testing.T) {
	tests := []struct {
		name               string
		content            string
		startRow, startCol int
		keys               string
		wantRow, wantCol   int
	}{
		{
			name:    "jumps to the end of the current word",
			content: "SELECT foo", startRow: 0, startCol: 0,
			keys: "e", wantRow: 0, wantCol: 5, // end of "SELECT"
		},
		{
			name:    "jumps to the end of the next word",
			content: "SELECT foo", startRow: 0, startCol: 0,
			keys: "ee", wantRow: 0, wantCol: 9, // end of "foo"
		},
		{
			name:    "a one-character punctuation word ends on itself",
			content: "foo(bar)", startRow: 0, startCol: 2,
			keys: "e", wantRow: 0, wantCol: 3, // the "("
		},
		{
			name:    "an empty line's end is its own start",
			content: "foo\n\nbar", startRow: 0, startCol: 2,
			keys: "e", wantRow: 1, wantCol: 0,
		},
		{
			name:    "continues past an empty line to the next word's end",
			content: "foo\n\nbar", startRow: 0, startCol: 2,
			keys: "ee", wantRow: 2, wantCol: 2, // end of "bar"
		},
		{
			name:    "does not move past the end of the buffer",
			content: "foo", startRow: 0, startCol: 2,
			keys: "e", wantRow: 0, wantCol: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEditor(t, tt.content)
			moveTo(&e, tt.startRow, tt.startCol)

			e = pressAll(e, tt.keys)

			assert.Equal(t, tt.content, e.editor.Value(), "word motions must not mutate the buffer")
			assert.Equal(t, tt.wantRow, e.editor.Line(), "cursor row")
			assert.Equal(t, tt.wantCol, e.editor.Column(), "cursor column")
		})
	}
}

func TestWordBackward(t *testing.T) {
	tests := []struct {
		name               string
		content            string
		startRow, startCol int
		keys               string
		wantRow, wantCol   int
	}{
		{
			name:    "hops back to the start of the previous word",
			content: "SELECT foo, bar", startRow: 0, startCol: 14,
			keys: "b", wantRow: 0, wantCol: 12, // start of "bar"
		},
		{
			name:    "lands on a punctuation run going backward",
			content: "SELECT foo, bar", startRow: 0, startCol: 14,
			keys: "bb", wantRow: 0, wantCol: 10, // the ","
		},
		{
			name:    "reaches the first word",
			content: "SELECT foo, bar", startRow: 0, startCol: 14,
			keys: "bbb", wantRow: 0, wantCol: 7, // start of "foo"
		},
		{
			name:    "crosses a line boundary backward",
			content: "foo\nbar", startRow: 1, startCol: 0,
			keys: "b", wantRow: 0, wantCol: 0,
		},
		{
			name:    "stops on an empty line going backward",
			content: "foo\n\nbar", startRow: 2, startCol: 0,
			keys: "b", wantRow: 1, wantCol: 0,
		},
		{
			name:    "moves past an empty line onto the previous word",
			content: "foo\n\nbar", startRow: 2, startCol: 0,
			keys: "bb", wantRow: 0, wantCol: 0,
		},
		{
			name:    "treats parentheses as one-character words going backward",
			content: "foo(bar)", startRow: 0, startCol: 7,
			keys: "bbb", wantRow: 0, wantCol: 0, // start of "foo"
		},
		{
			name:    "does not move past the start of the buffer",
			content: "foo", startRow: 0, startCol: 0,
			keys: "b", wantRow: 0, wantCol: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEditor(t, tt.content)
			moveTo(&e, tt.startRow, tt.startCol)

			e = pressAll(e, tt.keys)

			assert.Equal(t, tt.content, e.editor.Value(), "word motions must not mutate the buffer")
			assert.Equal(t, tt.wantRow, e.editor.Line(), "cursor row")
			assert.Equal(t, tt.wantCol, e.editor.Column(), "cursor column")
		})
	}
}

func TestInsertModeEntries(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		row, col  int
		key       rune
		wantValue string
		wantRow   int
		wantCol   int
	}{
		{
			name:      "i inserts at the cursor",
			content:   "SELECT 1;",
			row:       0,
			col:       3,
			key:       'i',
			wantValue: "SELECT 1;",
			wantRow:   0,
			wantCol:   3,
		},
		{
			name:      "a inserts after the cursor",
			content:   "SELECT 1;",
			row:       0,
			col:       3,
			key:       'a',
			wantValue: "SELECT 1;",
			wantRow:   0,
			wantCol:   4,
		},
		{
			name:      "a at the end of the line stays on the line",
			content:   "SELECT 1;",
			row:       0,
			col:       9,
			key:       'a',
			wantValue: "SELECT 1;",
			wantRow:   0,
			wantCol:   9,
		},
		{
			name:      "A inserts at the end of the line",
			content:   "SELECT 1;\nSELECT 2;",
			row:       0,
			col:       2,
			key:       'A',
			wantValue: "SELECT 1;\nSELECT 2;",
			wantRow:   0,
			wantCol:   9,
		},
		{
			name:      "I inserts at the first non-blank character",
			content:   "    SELECT 1;",
			row:       0,
			col:       9,
			key:       'I',
			wantValue: "    SELECT 1;",
			wantRow:   0,
			wantCol:   4,
		},
		{
			name:      "o opens a line below",
			content:   "SELECT 1;\nSELECT 2;",
			row:       0,
			col:       3,
			key:       'o',
			wantValue: "SELECT 1;\n\nSELECT 2;",
			wantRow:   1,
			wantCol:   0,
		},
		{
			name:      "o on the last line appends a line",
			content:   "SELECT 1;",
			row:       0,
			col:       0,
			key:       'o',
			wantValue: "SELECT 1;\n",
			wantRow:   1,
			wantCol:   0,
		},
		{
			name:      "O opens a line above",
			content:   "SELECT 1;\nSELECT 2;",
			row:       1,
			col:       4,
			key:       'O',
			wantValue: "SELECT 1;\n\nSELECT 2;",
			wantRow:   1,
			wantCol:   0,
		},
		{
			name:      "O on the first line prepends a line",
			content:   "SELECT 1;",
			row:       0,
			col:       4,
			key:       'O',
			wantValue: "\nSELECT 1;",
			wantRow:   0,
			wantCol:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEditor(t, tt.content)
			moveTo(&e, tt.row, tt.col)

			e = press(e, tt.key)

			assert.Equal(t, InsertMode, e.mode, "should have entered insert mode")
			assert.Equal(t, tt.wantValue, e.editor.Value())
			assert.Equal(t, tt.wantRow, e.editor.Line(), "cursor row")
			assert.Equal(t, tt.wantCol, e.editor.Column(), "cursor column")
		})
	}
}

// Entering insert mode has to announce the change, otherwise the status bar
// keeps rendering NORMAL while the user types.
func TestInsertModeEntriesReportModeChange(t *testing.T) {
	for _, k := range []rune{'i', 'a', 'A', 'I', 'o', 'O'} {
		t.Run(string(k), func(t *testing.T) {
			e := newTestEditor(t, "SELECT 1;")

			e, cmd := e.Update(tea.KeyPressMsg{Code: k, Text: string(k)})
			require.NotNil(t, cmd, "expected a mode change command")

			msg, ok := cmd().(modeChangeMsg)
			require.True(t, ok, "expected a modeChangeMsg")
			assert.Equal(t, InsertMode, msg.mode)
			assert.Equal(t, InsertMode, e.mode)
		})
	}
}

func TestInsertModeEntriesTypeText(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;")
	moveTo(&e, 0, 0)

	e = press(e, 'A')
	e = pressAll(e, " -- note")

	assert.Equal(t, "SELECT 1; -- note", e.editor.Value())

	e = pressEsc(e)
	assert.Equal(t, NormalMode, e.mode)

	e = press(e, 'o')
	e = pressAll(e, "SELECT 2;")

	assert.Equal(t, "SELECT 1; -- note\nSELECT 2;", e.editor.Value())
}

func TestUndoRestoresDeletedLine(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;\nSELECT 2;\nSELECT 3;")
	moveTo(&e, 1, 4)

	e = pressAll(e, "dd")
	require.Equal(t, "SELECT 1;\nSELECT 3;", e.editor.Value())

	e = press(e, 'u')

	assert.Equal(t, "SELECT 1;\nSELECT 2;\nSELECT 3;", e.editor.Value())
	assert.Equal(t, 1, e.editor.Line(), "undo should restore the cursor row")
	assert.Equal(t, 4, e.editor.Column(), "undo should restore the cursor column")
}

func TestRedoReappliesUndoneEdit(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;\nSELECT 2;")
	moveTo(&e, 0, 0)

	e = pressAll(e, "dd")
	e = press(e, 'u')
	require.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value())

	e = press(e, 'U')

	assert.Equal(t, "SELECT 2;", e.editor.Value())
}

// An insert session is a single undo step, as it is in vim: everything typed
// between entering insert mode and pressing esc comes back in one go.
func TestUndoGroupsWholeInsertSession(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;")
	moveTo(&e, 0, 9)

	e = press(e, 'a')
	e = pressAll(e, " FROM users")
	e = pressEsc(e)
	require.Equal(t, "SELECT 1; FROM users", e.editor.Value())

	e = press(e, 'u')

	assert.Equal(t, "SELECT 1;", e.editor.Value())
}

func TestUndoOfOpenLineBelow(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;")

	e = press(e, 'o')
	e = pressAll(e, "SELECT 2;")
	e = pressEsc(e)
	require.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value())

	e = press(e, 'u')

	assert.Equal(t, "SELECT 1;", e.editor.Value())
}

func TestUndoRestoresClearedBuffer(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;\nSELECT 2;")

	e, _ = e.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	require.Equal(t, "", e.editor.Value())

	e = press(e, 'u')

	assert.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value())
}

func TestUndoRestoresPaste(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;\nSELECT 2;")
	moveTo(&e, 0, 0)

	e = pressAll(e, "yy")
	e = press(e, 'p')
	require.Equal(t, "SELECT 1;\nSELECT 1;\nSELECT 2;", e.editor.Value())

	e = press(e, 'u')

	assert.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value())
}

// No-op commands must not swallow an undo step, or `u` after them appears to do
// nothing at all.
func TestNoOpCommandsDoNotConsumeUndoSteps(t *testing.T) {
	tests := []struct {
		name          string
		edit          func(Editor) Editor
		wantAfterEdit string
		noop          func(Editor) Editor
	}{
		{
			name:          "insert mode entered and left without typing",
			edit:          func(e Editor) Editor { return pressAll(e, "dd") },
			wantAfterEdit: "SELECT 2;",
			noop:          func(e Editor) Editor { return pressEsc(press(e, 'i')) },
		},
		{
			name: "paste with an empty register",
			// Edits through insert mode, so the register stays empty: dd and yy
			// would both fill it and make the paste a real edit.
			edit: func(e Editor) Editor {
				e = press(e, 'A')
				e = pressAll(e, " -- note")

				return pressEsc(e)
			},
			wantAfterEdit: "SELECT 1; -- note\nSELECT 2;",
			noop:          func(e Editor) Editor { return press(e, 'p') },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newTestEditor(t, "SELECT 1;\nSELECT 2;")
			moveTo(&e, 0, 0)

			e = tt.edit(e)
			require.Equal(t, tt.wantAfterEdit, e.editor.Value())

			e = tt.noop(e)
			require.Equal(t, tt.wantAfterEdit, e.editor.Value(),
				"the no-op should leave the buffer alone")

			e = press(e, 'u')

			assert.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value(),
				"undo should still reach the edit")
		})
	}
}

func TestUndoAndRedoWithEmptyHistoryAreNoOps(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;")

	e = press(e, 'u')
	assert.Equal(t, "SELECT 1;", e.editor.Value())

	e = press(e, 'U')
	assert.Equal(t, "SELECT 1;", e.editor.Value())
}

// Editing after an undo discards the branch that was undone, so redo has
// nothing left to reapply.
func TestEditAfterUndoClearsRedo(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;\nSELECT 2;\nSELECT 3;")
	moveTo(&e, 0, 0)

	e = pressAll(e, "dd")
	e = press(e, 'u')
	require.Equal(t, "SELECT 1;\nSELECT 2;\nSELECT 3;", e.editor.Value())

	// A different edit replaces the undone branch.
	moveTo(&e, 2, 0)
	e = pressAll(e, "dd")
	require.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value())

	e = press(e, 'U')

	assert.Equal(t, "SELECT 1;\nSELECT 2;", e.editor.Value(), "redo should be empty")
}

func TestUndoHistoryIsBounded(t *testing.T) {
	e := newTestEditor(t, "SELECT 1;")

	for i := 0; i < maxUndoHistory*2; i++ {
		e = press(e, 'o')
		e = pressEsc(e)
	}

	assert.Len(t, e.undoStack, maxUndoHistory)
}
