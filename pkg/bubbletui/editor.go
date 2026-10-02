package bubbletui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/davecgh/go-spew/spew"
	"github.com/ionut-t/goeditor"
	"github.com/ionut-t/goeditor/core"
)

// dblabSQLStyle names the chroma style registered below, used for the query
// editor's SQL syntax highlighting.
const dblabSQLStyle = "dblab-cyberpunk-sql"

// Register the dblab-cyberpunk-sql style the query editor highlights SQL with,
// built from the same neon palette used across the rest of the app (see
// cyberGreen, neonViolet, neonPurple and friends in bubbletui.go).
var _ = styles.Register(chroma.MustNewStyle(dblabSQLStyle, chroma.StyleEntries{
	chroma.Background:    "#E0E0E0",      // default text    → whiteText
	chroma.Keyword:       "#9D00FF bold", // SELECT/FROM/... → neonViolet
	chroma.NameBuiltin:   "#BF40BF bold", // VARCHAR/INT/... → neonPurple
	chroma.Name:          "#2ECC71",      // identifiers     → mutedGreen
	chroma.LiteralString: "#39FF14",      // 'strings'       → cyberGreen
	chroma.LiteralNumber: "#FF6600",      // numbers         → neonOrange
	chroma.Operator:      "#E0E0E0",      // + - * / = ...   → whiteText
	chroma.Punctuation:   "#E0E0E0",      // ; , ( ) ...     → whiteText
	chroma.Comment:       "#999999 italic",
	chroma.Error:         "#FF0000 bold",
}))

type executeQueryMsg struct {
	queriesToRun []string
}

type Editor struct {
	geditor goeditor.Model
	keyMap  keys.EditorKeyMap

	dump          io.Writer
	width, height int
}

func NewEditor(km keys.EditorKeyMap) (Editor, error) {
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
	geditor.SetLanguage("sql", dblabSQLStyle)
	geditor.WithSearchOptions(goeditor.SearchOptions{
		IgnoreCase: true,
		SmartCase:  true,
		Wrap:       true,
	})
	geditor.SetExtraWordChars('-')
	geditor.SetPlaceholder("Start typing...")

	if err := applyKeyMap(&geditor, km); err != nil {
		return Editor{}, err
	}

	return Editor{geditor: geditor, keyMap: km, dump: dump}, nil
}

// keyBinding pairs a configurable dblab key (bubbletea notation, e.g.
// "ctrl+r") with the goeditor command it should trigger (Vim notation, e.g.
// "<C-r>"), and the goeditor modes the mapping applies to.
type keyBinding struct {
	lhs   string
	rhs   string
	modes core.MapMode
}

// applyKeyMap wires the user-configurable editor keybindings into goeditor's
// own Vim emulation, so a key customised in the dblab config triggers the
// same goeditor command as its Vim default would. Bindings that already
// match goeditor's built-in default are left alone.
func applyKeyMap(geditor *goeditor.Model, km keys.EditorKeyMap) error {
	bindings := []keyBinding{
		// Motions: valid in normal, visual and operator-pending, as in Vim.
		{firstKey(km.Up), "k", core.MapAll},
		{firstKey(km.Down), "j", core.MapAll},
		{firstKey(km.Left), "h", core.MapAll},
		{firstKey(km.Right), "l", core.MapAll},
		{firstKey(km.LineStart), "0", core.MapAll},
		{firstKey(km.LineEnd), "$", core.MapAll},
		{firstKey(km.GoToTop), "gg", core.MapAll},
		{firstKey(km.GoToBottom), "G", core.MapAll},
		{firstKey(km.WordForward), "w", core.MapAll},
		{firstKey(km.WordEnd), "e", core.MapAll},
		{firstKey(km.WordBackward), "b", core.MapAll},
		// Mode switching: only make sense from normal mode.
		{firstKey(km.Insert), "i", core.MapNormal},
		{firstKey(km.Append), "a", core.MapNormal},
		{firstKey(km.AppendLineEnd), "A", core.MapNormal},
		{firstKey(km.InsertLineStart), "I", core.MapNormal},
		{firstKey(km.OpenLineBelow), "o", core.MapNormal},
		{firstKey(km.OpenLineAbove), "O", core.MapNormal},
		// Returning to normal mode happens from insert mode.
		{firstKey(km.Normal), "<Esc>", core.MapInsert},
		// History.
		{firstKey(km.Undo), "u", core.MapNormal},
		{firstKey(km.Redo), "<C-r>", core.MapNormal},
	}

	for _, b := range bindings {
		if b.lhs == "" {
			continue
		}

		lhs := vimKeyNotation(b.lhs)
		if lhs == b.rhs {
			// Already goeditor's own default; no mapping needed.
			continue
		}

		if err := geditor.Map(b.modes, lhs, b.rhs, true); err != nil {
			return fmt.Errorf("failed to map editor key %q to %q: %w", b.lhs, b.rhs, err)
		}
	}

	return nil
}

// firstKey returns the first key string configured for a binding, or "" if
// none is set.
func firstKey(b key.Binding) string {
	keys := b.Keys()
	if len(keys) == 0 {
		return ""
	}

	return keys[0]
}

// vimKeyNotation converts a bubbletea-style key string, as used throughout
// the dblab config (e.g. "esc", "ctrl+r", "shift+tab"), into the Vim
// notation goeditor's Map expects (e.g. "<Esc>", "<C-r>", "<S-Tab>").
// Plain printable keys (e.g. "k", "0", "$") are returned unchanged.
func vimKeyNotation(k string) string {
	namedKeys := map[string]string{
		"esc":       "Esc",
		"escape":    "Esc",
		"enter":     "CR",
		"return":    "CR",
		"tab":       "Tab",
		"space":     "Space",
		"backspace": "BS",
		"delete":    "Del",
		"insert":    "Insert",
		"up":        "Up",
		"down":      "Down",
		"left":      "Left",
		"right":     "Right",
		"home":      "Home",
		"end":       "End",
		"pgup":      "PageUp",
		"pgdown":    "PageDown",
	}

	parts := strings.Split(k, "+")
	name := strings.ToLower(parts[len(parts)-1])
	mods := parts[:len(parts)-1]

	named, isNamed := namedKeys[name]
	if !isNamed {
		if len(mods) == 0 {
			// A plain, printable key: Vim notation is the same string.
			return k
		}
		named = parts[len(parts)-1]
	}

	var b strings.Builder
	b.WriteString("<")
	for _, mod := range mods {
		switch strings.ToLower(mod) {
		case "ctrl":
			b.WriteString("C-")
		case "alt":
			b.WriteString("A-")
		case "shift":
			b.WriteString("S-")
		}
	}
	b.WriteString(named)
	b.WriteString(">")

	return b.String()
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
