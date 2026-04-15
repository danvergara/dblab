package bubbletui

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
)

type executeQueryMsg struct {
	Query string
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
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(mutedGreen)
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
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
	case tea.KeyMsg:
		if key.Matches(msg, e.bindings.Editor.ExecuteQuery) {
			query := e.editor.Value()
			if strings.TrimSpace(query) == "" {
				return e, nil
			}

			fireQueryCmd := func() tea.Msg {
				return executeQueryMsg{Query: query}
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
				e.editor, cmd = e.editor.Update(tea.KeyMsg{Type: tea.KeyDelete})
				return e, cmd
			case "0":
				e.editor, cmd = e.editor.Update(tea.KeyMsg{Type: tea.KeyHome})
				return e, cmd
			case "$":
				e.editor, cmd = e.editor.Update(tea.KeyMsg{Type: tea.KeyEnd})
				return e, cmd
			}

			switch {
			case key.Matches(msg, e.bindings.Editor.Insert):
				e.mode = InsertMode
				e.editor.Cursor.SetMode(cursor.CursorBlink)
				return e, nil

			case key.Matches(msg, e.bindings.Editor.Left):
				e.editor, cmd = e.editor.Update(tea.KeyMsg{Type: tea.KeyLeft})
				return e, cmd

			case key.Matches(msg, e.bindings.Editor.Right):
				e.editor, cmd = e.editor.Update(tea.KeyMsg{Type: tea.KeyRight})
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
				e.editor.Cursor.SetMode(cursor.CursorStatic)
				// Optional: move back one space on escape
				e.editor, _ = e.editor.Update(tea.KeyMsg{Type: tea.KeyLeft})
				return e, nil
			}
		}
	}

	e.editor, cmd = e.editor.Update(msg)
	return e, cmd
}

func (e Editor) View() string {
	return e.editor.View()
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
