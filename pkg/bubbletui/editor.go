package bubbletui

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

type executeQueryMsg struct {
	Query string
}

type Editor struct {
	editor   textarea.Model
	bindings *command.TUIKeyBindings
	dump     io.Writer
}

func NewEditor(kb *command.TUIKeyBindings) Editor {
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
		switch {
		case key.Matches(msg, e.bindings.ExecuteQuery):
			query := e.editor.Value()
			if strings.TrimSpace(query) == "" {
				return e, nil
			}

			fireQueryCmd := func() tea.Msg {
				return executeQueryMsg{Query: query}
			}

			return e, fireQueryCmd
		}

		e.editor, cmd = e.editor.Update(msg)
	}

	return e, cmd
}

func (e Editor) View() string {
	return e.editor.View()
}
