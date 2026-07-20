package bubbletui

import (
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
)

type executeQueryMsg struct {
	queriesToRun []string
}

type Editor struct {
	editor     textarea.Model
	bindings   *command.TUIKeyMap
	mode       Mode
	register   string
	pendingCmd string
	dump       io.Writer
}

func NewEditor(kb *command.TUIKeyMap) Editor {
	var isDark = compat.HasDarkBackground
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

	return Editor{editor: ta, bindings: kb, dump: dump}
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
		e.editor.SetValue(msg.QueryText)
	case tea.KeyPressMsg:
		if key.Matches(msg, e.bindings.Editor.ExecuteQuery) {
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

		if key.Matches(msg, e.bindings.Editor.ExecuteSingleQuery) {
			value := e.editor.Value()

			if len(value) == 0 {
				return e, nil
			}

			query := queryAtCursor(value, e.editor.Line(), e.editor.Column())
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
						e.deleteCurrentLine()
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
				e.pasteAfter()
				return e, nil
			case "x":
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyDelete})
				return e, cmd
			case "0":
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyHome})
				return e, cmd
			case "$":
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
				return e, cmd
			case "ctrl+d":
				e.editor.Reset() // Clears text, cursor, and history
				return e, nil
			case "G":
				// LineCount() returns the total number of lines.
				// Line() returns the current 0-indexed line position.
				lastLine := e.editor.LineCount() - 1
				for e.editor.Line() < lastLine {
					e.editor.CursorDown()
				}
				return e, nil
			case "g":
				for e.editor.Line() > 0 {
					e.editor.CursorUp()
				}
				return e, nil
			}

			switch {
			case key.Matches(msg, e.bindings.Editor.Insert):
				e.mode = InsertMode
				styles := e.editor.Styles()
				styles.Cursor.Blink = true
				e.editor.SetStyles(styles)
				return e, nil

			case key.Matches(msg, e.bindings.Editor.Left):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
				return e, cmd

			case key.Matches(msg, e.bindings.Editor.Right):
				e.editor, cmd = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyRight})
				return e, cmd

			case key.Matches(msg, e.bindings.Editor.Down):
				e.editor.CursorDown()
				return e, nil

			case key.Matches(msg, e.bindings.Editor.Up):
				e.editor.CursorUp()
				return e, nil
			}

			return e, nil
		case InsertMode:
			switch {
			case key.Matches(msg, e.bindings.Editor.Normal):
				e.mode = NormalMode
				styles := e.editor.Styles()
				styles.Cursor.Blink = false
				e.editor.SetStyles(styles)
				e.editor, _ = e.editor.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
				return e, nil
			}
		}
	}

	e.editor, cmd = e.editor.Update(msg)
	return e, cmd
}

func (e Editor) View() tea.View {
	return tea.NewView(e.editor.View())
}

func queryAtCursor(content string, row, col int) string {
	lines := strings.Split(content, "\n")
	if row < 0 || row >= len(lines) {
		return strings.TrimSpace(content)
	}

	if col < 0 {
		col = 0
	} else if col > len(lines[row]) {
		col = len(lines[row])
	}

	cursorOffset := 0
	for i := 0; i < row; i++ {
		cursorOffset += len(lines[i]) + 1
	}
	cursorOffset += col

	start := 0
	end := len(content)
	for i := 0; i < cursorOffset && i < len(content); i++ {
		if content[i] == ';' {
			start = i + 1
		}
	}

	for i := cursorOffset; i < len(content); i++ {
		if content[i] == ';' {
			end = i
			break
		}
	}

	query := strings.TrimSpace(content[start:end])
	if query == "" {
		return strings.TrimSpace(content)
	}

	return query
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
