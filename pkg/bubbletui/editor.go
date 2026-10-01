package bubbletui

import (
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/davecgh/go-spew/spew"
	"github.com/ionut-t/goeditor"
)

type executeQueryMsg struct {
	queriesToRun []string
}

type Editor struct {
	geditor goeditor.Model
	keyMap  keys.EditorKeyMap

	dump          io.Writer
	width, height int
}

func NewEditor(km keys.EditorKeyMap) Editor {
	isDark := lipgloss.HasDarkBackground(os.Stdout, os.Stderr)
	var dump *os.File

	if _, ok := os.LookupEnv("DBLAB_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("editor_messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	geditor := goeditor.New(0, 0)
	geditor.Focus()
	geditor.SetCursorMode(goeditor.CursorBlink)
	geditor.SetLanguage("sql", languageTheme(isDark))
	geditor.WithSearchOptions(goeditor.SearchOptions{
		IgnoreCase: true,
		SmartCase:  true,
		Wrap:       true,
	})
	geditor.SetExtraWordChars('-')
	geditor.SetPlaceholder("Start typing...")

	return Editor{geditor: geditor, keyMap: km, dump: dump}
}

func (e *Editor) SetSize(w, h int) {
	if e.width == w && e.height == h {
		return
	}
	e.width, e.height = w, h
	e.geditor.SetSize(w-4, h-2)
}

func (e *Editor) Blur() {
	e.geditor.Blur()
}

func (e *Editor) Focus() {
	e.geditor.Focus()
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
		if e.geditor.GetCurrentContent() != "" {
			e.geditor.SetContent(e.geditor.GetCurrentContent() + "\n" + msg.QueryText + ";")
		} else {
			e.geditor.SetContent(msg.QueryText + ";")
		}
		_ = e.geditor.SetCursorPositionEnd()
		editorModel, cmd := e.geditor.Update(msg)
		e.geditor = editorModel
		return e, cmd
	case tea.KeyPressMsg:
		if key.Matches(msg, e.keyMap.ExecuteQuery) {
			editorContent := e.geditor.GetCurrentContent()

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
			value := e.geditor.GetCurrentContent()

			if len(value) == 0 {
				return e, nil
			}

			query := queryAtCursor(value, e.geditor.GetCursorPosition().Row)
			if len(query) == 0 {
				return e, nil
			}
			queriesToRun := prepareQueriesForExecution(query)
			fireQueryCmd := func() tea.Msg {
				return executeQueryMsg{queriesToRun: queriesToRun}
			}
			return e, fireQueryCmd
		}
	}

	var cmds []tea.Cmd
	editorModel, cmd := e.geditor.Update(msg)
	cmds = append(cmds, cmd)
	e.geditor = editorModel

	return e, tea.Batch(cmds...)
}

func (e Editor) View() tea.View {
	return tea.NewView(e.geditor.View())
}

// queryAtCursor returns the text of the line the cursor is currently on.
func queryAtCursor(content string, currentIndex int) string {
	lines := strings.Split(content, "\n")
	if currentIndex >= 0 && currentIndex < len(lines) {
		return lines[currentIndex]
	}

	return ""
}

func languageTheme(isDark bool) string {
	if isDark {
		return "catppuccin-mocha"
	}

	return "catppuccin-latte"
}
