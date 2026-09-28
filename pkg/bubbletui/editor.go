package bubbletui

import (
	"io"
	"os"
	"strings"
	"unicode"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/davecgh/go-spew/spew"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
)

func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "NORMAL"
	case InsertMode:
		return "INSERT"
	default:
		return ""
	}
}

// maxUndoHistory caps the undo stack so a long editing session can't grow the
// recorded history without bound.
const maxUndoHistory = 100

type executeQueryMsg struct {
	queriesToRun []string
}

type modeChangeMsg struct {
	mode Mode
}

// editorState is the buffer and cursor as they were before a mutation, so undo
// can put both back where the user left them.
type editorState struct {
	value string
	row   int
	col   int
}

type Editor struct {
	editor textarea.Model
	keyMap keys.EditorKeyMap

	mode       Mode
	register   string
	pendingCmd string

	undoStack []editorState
	redoStack []editorState

	dump io.Writer
}

func NewEditor(km keys.EditorKeyMap) Editor {
	isDark := compat.HasDarkBackground
	var dump *os.File

	if _, ok := os.LookupEnv("DBLAB_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("editor_messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	ta := textarea.New()
	ta.Placeholder = "Enter text..."
	s := textarea.DefaultStyles(isDark)
	s.Focused.Text = lipgloss.NewStyle().Foreground(mutedGreen)
	s.Blurred.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	ta.SetStyles(s)
	ta.Focus()

	return Editor{editor: ta, keyMap: km, dump: dump}
}

func (e *Editor) SetWidth(w int) {
	e.editor.SetWidth(w - 4)
}

func (e *Editor) SetHeight(h int) {
	e.editor.SetHeight(h - 2)
}

func (e *Editor) Blur() {
	e.editor.Blur()
}

func (e *Editor) Focus() tea.Cmd {
	return e.editor.Focus()
}

func (e Editor) Init() tea.Cmd {
	return nil
}

func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	if e.dump != nil {
		spew.Fdump(e.dump, msg)
	}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case querySelectedMsg:
		e.snapshot()
		e.editor.CursorEnd()
		if e.editor.Value() != "" {
			e.editor.InsertString("\n" + msg.QueryText + ";")
		} else {
			e.editor.InsertString(msg.QueryText + ";")
		}
		e.commit()
		return e, nil
	case tea.KeyPressMsg:
		if key.Matches(msg, e.keyMap.ExecuteQuery) {
			editorContent := e.editor.Value()

			queriesToRun := prepareQueriesForExecution(editorContent)
			if len(queriesToRun) == 0 {
				return e, nil
			}

			fireQueryCmd := func() tea.Msg {
				return executeQueryMsg{queriesToRun: queriesToRun}
			}

			return e, fireQueryCmd
		}

		if key.Matches(msg, e.keyMap.ExecuteSingleQuery) {
			value := e.editor.Value()

			if len(value) == 0 {
				return e, nil
			}

			query := queryAtCursor(value, e.editor.Line())
			if len(query) == 0 {
				return e, nil
			}
			queriesToRun := prepareQueriesForExecution(query)
			fireQueryCmd := func() tea.Msg {
				return executeQueryMsg{queriesToRun: queriesToRun}
			}
			return e, fireQueryCmd
		}

		switch e.mode {
		case NormalMode:
			char := msg.String()
			if e.pendingCmd != "" {
				switch e.pendingCmd {
				case "d":
					if char == "d" {
						e.snapshot()
						e.deleteCurrentLine()
						e.commit()
					}
					e.pendingCmd = ""
					return e, nil

				case "y":
					if char == "y" {
						e.yankCurrentLine()
					}
					e.pendingCmd = ""
					return e, nil
				}
			}

			switch char {
			case "d", "y":
				e.pendingCmd = char
				return e, nil
			case "p":
				e.snapshot()
				e.pasteAfter()
				e.commit()
				return e, nil
			case "x":
				e.snapshot()
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyDelete})
				e.commit()
				return e, cmd
			case "ctrl+d":
				e.snapshot()
				e.editor.Reset() // Clears text and cursor; undo brings it back.
				e.commit()
				return e, nil
			}

			switch {
			case key.Matches(msg, e.keyMap.LineStart):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyHome})
				return e, cmd
			case key.Matches(msg, e.keyMap.LineEnd):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
				return e, cmd

			case key.Matches(msg, e.keyMap.WordForward):
				// Note: if pendingCmd is "d" (delete-line pending), this falls
				// through the pendingCmd switch above without matching "d", so
				// the pending delete is silently cancelled here rather than
				// composing into "dw". Operator-pending word motions aren't
				// supported yet.
				e.wordForward()
				return e, nil

			case key.Matches(msg, e.keyMap.WordEnd):
				e.wordEnd()
				return e, nil

			case key.Matches(msg, e.keyMap.WordBackward):
				e.wordBackward()
				return e, nil

			case key.Matches(msg, e.keyMap.GoToBottom):
				// LineCount() returns the total number of lines.
				// Line() returns the current 0-indexed line position.
				lastLine := e.editor.LineCount() - 1
				for e.editor.Line() < lastLine {
					e.editor.CursorDown()
				}
				return e, nil
			case key.Matches(msg, e.keyMap.GoToTop):
				for e.editor.Line() > 0 {
					e.editor.CursorUp()
				}
				return e, nil
			case key.Matches(msg, e.keyMap.Insert):
				e.snapshot()
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.Append):
				e.snapshot()
				// SetCursorColumn clamps, so this stops at the end of the line.
				e.editor.SetCursorColumn(e.editor.Column() + 1)
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.AppendLineEnd):
				e.snapshot()
				e.editor.CursorEnd()
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.InsertLineStart):
				e.snapshot()
				e.editor.SetCursorColumn(firstNonBlank(e.currentLine()))
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.OpenLineBelow):
				e.snapshot()
				// Splitting the line at its end leaves an empty row below with
				// the cursor already on it.
				e.editor.CursorEnd()
				e.editor.InsertString("\n")
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.OpenLineAbove):
				e.snapshot()
				// Splitting at column 0 pushes the current line down; stepping
				// back up lands on the empty row that opened above it.
				e.editor.CursorStart()
				e.editor.InsertString("\n")
				e.editor.CursorUp()
				return e, e.enterInsertMode()

			case key.Matches(msg, e.keyMap.Undo):
				e.undo()
				return e, nil

			case key.Matches(msg, e.keyMap.Redo):
				e.redo()
				return e, nil

			case key.Matches(msg, e.keyMap.Left):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
				return e, cmd

			case key.Matches(msg, e.keyMap.Right):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyRight})
				return e, cmd

			case key.Matches(msg, e.keyMap.Down):
				e.editor.CursorDown()
				return e, nil

			case key.Matches(msg, e.keyMap.Up):
				e.editor.CursorUp()
				return e, nil
			}

			return e, nil
		case InsertMode:
			switch {
			case key.Matches(msg, e.keyMap.Normal):
				e.mode = NormalMode
				styles := e.editor.Styles()
				styles.Cursor.Blink = false
				e.editor.SetStyles(styles)
				e.editor, _ = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
				// Closes the edit opened on entering insert mode, so the whole
				// insert session undoes as a single step.
				e.commit()
				mode := e.mode
				fireModeChangeCmd := func() tea.Msg {
					return modeChangeMsg{mode: mode}
				}
				return e, fireModeChangeCmd
			}
		}
	}

	e.editor, cmd = e.editor.Update(msg)
	return e, cmd
}

func (e Editor) View() tea.View {
	return tea.NewView(e.editor.View())
}

// enterInsertMode switches to insert mode and announces the change so the
// status bar repaints. Callers position the cursor first.
func (e *Editor) enterInsertMode() tea.Cmd {
	e.mode = InsertMode

	styles := e.editor.Styles()
	styles.Cursor.Blink = true
	e.editor.SetStyles(styles)

	mode := e.mode

	return func() tea.Msg {
		return modeChangeMsg{mode: mode}
	}
}

// state captures the buffer and cursor as they are right now.
func (e Editor) state() editorState {
	return editorState{
		value: e.editor.Value(),
		row:   e.editor.Line(),
		col:   e.editor.Column(),
	}
}

// snapshot opens an edit, recording the buffer so the mutation that follows can
// be undone. Recording a new state drops the redo stack, as vim does: editing
// after an undo discards the branch that was undone.
func (e *Editor) snapshot() {
	e.undoStack = append(e.undoStack, e.state())

	if len(e.undoStack) > maxUndoHistory {
		e.undoStack = e.undoStack[len(e.undoStack)-maxUndoHistory:]
	}

	e.redoStack = nil
}

// commit closes the edit opened by snapshot, dropping the snapshot when the
// command left the buffer untouched. That keeps no-ops such as `i<esc>`, `p`
// with an empty register or `x` on an empty line from costing an undo step.
func (e *Editor) commit() {
	n := len(e.undoStack)

	if n > 0 && e.undoStack[n-1].value == e.editor.Value() {
		e.undoStack = e.undoStack[:n-1]
	}
}

func (e *Editor) undo() {
	if len(e.undoStack) == 0 {
		return
	}

	prev := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]
	e.redoStack = append(e.redoStack, e.state())

	e.restore(prev)
}

func (e *Editor) redo() {
	if len(e.redoStack) == 0 {
		return
	}

	next := e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]
	e.undoStack = append(e.undoStack, e.state())

	e.restore(next)
}

// restore puts the buffer back to s. SetValue leaves the cursor at the end of
// the text it inserts, so the cursor has to be placed afterwards.
func (e *Editor) restore(s editorState) {
	e.editor.SetValue(s.value)
	e.setCursor(s.row, s.col)
}

// setCursor moves the cursor to a logical row and column. CursorDown steps by
// visual line rather than logical line, so wrapped lines take more than one
// call; row is clamped because a row past the last line would never be reached
// and the loop would spin forever.
func (e *Editor) setCursor(row, col int) {
	row = max(0, min(row, e.editor.LineCount()-1))

	e.editor.MoveToBegin()
	for e.editor.Line() < row {
		e.editor.CursorDown()
	}

	e.editor.SetCursorColumn(col)
}

// currentLine returns the text of the line the cursor is on.
func (e Editor) currentLine() string {
	lines := strings.Split(e.editor.Value(), "\n")
	row := e.editor.Line()

	if row < 0 || row >= len(lines) {
		return ""
	}

	return lines[row]
}

// firstNonBlank returns the column of the first non-whitespace rune in line, or
// the end of the line when it is entirely blank. Columns are rune offsets,
// matching how the textarea indexes them.
func firstNonBlank(line string) int {
	runes := []rune(line)

	for i, r := range runes {
		if !unicode.IsSpace(r) {
			return i
		}
	}

	return len(runes)
}

// wordSlotKind classifies a single rune for the purposes of vim-style word
// motion: a word is a maximal run of word characters, or a maximal run of
// punctuation, so e.g. "foo(bar)" is four words: foo, (, bar, ).
const (
	slotBlank = iota
	slotWord
	slotPunct
	slotEmptyLine
	slotEOL
)

// wordClass classifies a single rune as blank, a word character (letter,
// digit, or underscore), or punctuation (anything else non-blank).
func wordClass(r rune) int {
	switch {
	case unicode.IsSpace(r):
		return slotBlank
	case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_':
		return slotWord
	default:
		return slotPunct
	}
}

// wordSlotKind classifies the position (row, col) in lines. An empty line is
// its own one-slot "word" (vim stops on a blank line rather than skipping
// over it); a column past the last rune of a non-empty line is treated like
// blank space, since that's where the line break lives.
func wordSlotKind(lines []string, row, col int) int {
	runeLine := []rune(lines[row])

	if len(runeLine) == 0 {
		return slotEmptyLine
	}

	if col >= len(runeLine) {
		return slotEOL
	}

	return wordClass(runeLine[col])
}

// nextWordSlot moves one slot forward, crossing into the next line once col
// runs past the end of the current one. ok is false at the very last slot in
// the buffer, so callers never wrap around.
func nextWordSlot(lines []string, row, col int) (int, int, bool) {
	runeLine := []rune(lines[row])

	if col < len(runeLine) {
		return row, col + 1, true
	}

	if row+1 >= len(lines) {
		return row, col, false
	}

	return row + 1, 0, true
}

// prevWordSlot moves one slot backward, landing on the previous line's
// end-of-line slot once col reaches 0. ok is false at the very first slot in
// the buffer.
func prevWordSlot(lines []string, row, col int) (int, int, bool) {
	if col > 0 {
		return row, col - 1, true
	}

	if row == 0 {
		return row, col, false
	}

	prevLine := []rune(lines[row-1])

	return row - 1, len(prevLine), true
}

// wordForwardTarget implements vim's `w`: skip the rest of the current
// word/punct run (if standing on one), then skip blanks until landing on a
// word, punct run, or an empty line. It leaves the position unchanged once
// there is no further word in the buffer, rather than resting on the
// past-the-last-character slot that nextWordSlot uses internally to detect
// the end of a line.
func wordForwardTarget(lines []string, row, col int) (int, int) {
	startRow, startCol := row, col
	kind := wordSlotKind(lines, row, col)

	if kind == slotWord || kind == slotPunct {
		for {
			nr, nc, ok := nextWordSlot(lines, row, col)
			if !ok {
				return startRow, startCol
			}

			if wordSlotKind(lines, nr, nc) != kind {
				row, col = nr, nc
				break
			}

			row, col = nr, nc
		}
	} else if kind == slotEmptyLine {
		nr, nc, ok := nextWordSlot(lines, row, col)
		if !ok {
			return startRow, startCol
		}

		row, col = nr, nc
	}

	for {
		k := wordSlotKind(lines, row, col)
		if k == slotWord || k == slotPunct || k == slotEmptyLine {
			return row, col
		}

		nr, nc, ok := nextWordSlot(lines, row, col)
		if !ok {
			return startRow, startCol
		}

		row, col = nr, nc
	}
}

// wordEndTarget implements vim's `e`: always step forward at least one slot,
// skip blanks, then extend to the end of the word/punct run reached (a blank
// line's "end" is itself). It leaves the position unchanged once there is no
// further word in the buffer, for the same reason wordForwardTarget does.
func wordEndTarget(lines []string, row, col int) (int, int) {
	startRow, startCol := row, col

	nr, nc, ok := nextWordSlot(lines, row, col)
	if !ok {
		return startRow, startCol
	}

	row, col = nr, nc

	for {
		k := wordSlotKind(lines, row, col)
		if k == slotWord || k == slotPunct || k == slotEmptyLine {
			break
		}

		nr, nc, ok := nextWordSlot(lines, row, col)
		if !ok {
			return startRow, startCol
		}

		row, col = nr, nc
	}

	kind := wordSlotKind(lines, row, col)
	if kind == slotEmptyLine {
		return row, col
	}

	for {
		nr, nc, ok := nextWordSlot(lines, row, col)
		if !ok || wordSlotKind(lines, nr, nc) != kind {
			break
		}

		row, col = nr, nc
	}

	return row, col
}

// wordBackwardTarget implements vim's `b`, mirroring wordForwardTarget using
// prevWordSlot.
func wordBackwardTarget(lines []string, row, col int) (int, int) {
	nr, nc, ok := prevWordSlot(lines, row, col)
	if !ok {
		return row, col
	}

	row, col = nr, nc

	for {
		k := wordSlotKind(lines, row, col)
		if k == slotWord || k == slotPunct || k == slotEmptyLine {
			break
		}

		nr, nc, ok := prevWordSlot(lines, row, col)
		if !ok {
			break
		}

		row, col = nr, nc
	}

	kind := wordSlotKind(lines, row, col)
	if kind == slotWord || kind == slotPunct {
		for {
			nr, nc, ok := prevWordSlot(lines, row, col)
			if !ok || wordSlotKind(lines, nr, nc) != kind {
				break
			}

			row, col = nr, nc
		}
	}

	return row, col
}

func (e *Editor) wordForward() {
	lines := strings.Split(e.editor.Value(), "\n")
	row, col := wordForwardTarget(lines, e.editor.Line(), e.editor.Column())
	e.setCursor(row, col)
}

func (e *Editor) wordEnd() {
	lines := strings.Split(e.editor.Value(), "\n")
	row, col := wordEndTarget(lines, e.editor.Line(), e.editor.Column())
	e.setCursor(row, col)
}

func (e *Editor) wordBackward() {
	lines := strings.Split(e.editor.Value(), "\n")
	row, col := wordBackwardTarget(lines, e.editor.Line(), e.editor.Column())
	e.setCursor(row, col)
}

// queryAtCursor returns the text of the line the cursor is currently on.
func queryAtCursor(content string, currentIndex int) string {
	lines := strings.Split(content, "\n")
	if currentIndex >= 0 && currentIndex < len(lines) {
		return lines[currentIndex]
	}

	return ""
}

func (e *Editor) yankCurrentLine() {
	lines := strings.Split(e.editor.Value(), "\n")
	row := e.editor.Line()

	if row >= 0 && row < len(lines) {
		e.register = lines[row]
	}
}

func (e *Editor) deleteCurrentLine() {
	lines := strings.Split(e.editor.Value(), "\n")
	row := e.editor.Line()

	if row >= 0 && row < len(lines) {
		e.register = lines[row]

		lines = append(lines[:row], lines[row+1:]...)

		e.editor.SetValue(strings.Join(lines, "\n"))

		targetRow := row
		if targetRow >= len(lines) {
			targetRow = len(lines) - 1
		}

		targetRow = max(0, targetRow)

		for e.editor.Line() > 0 {
			e.editor.CursorUp()
		}

		for e.editor.Line() < targetRow {
			e.editor.CursorDown()
		}

		e.editor.CursorStart()
	}
}

func (e *Editor) pasteAfter() {
	if e.register == "" {
		return
	}

	lines := strings.Split(e.editor.Value(), "\n")
	row := e.editor.Line()

	if row >= 0 && row < len(lines) {
		lines = append(lines[:row+1], append([]string{e.register}, lines[row+1:]...)...)
		e.editor.SetValue(strings.Join(lines, "\n"))
	}
}
