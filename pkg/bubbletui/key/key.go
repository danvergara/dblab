package key

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Help           key.Binding
	Quit           key.Binding
	Navigation     NavigationKeyMap
	Editor         EditorKeyMap
	SiebarViewport SiebarViewportKeyMap
	ResultSet      ResultSetKeyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Navigation:     DefaultNavitationKeyMap(),
		Editor:         DefaultEditorKeyMap(),
		SiebarViewport: DefaultSiebarViewportKeyMap(),
		ResultSet:      DefaultResultSetKeyMap(),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SiebarViewport.GoToBottom, k.SiebarViewport.GoToTop},
		{k.ResultSet.NextTab, k.ResultSet.PrevTab, k.ResultSet.LineStart, k.ResultSet.LineEnd, k.ResultSet.GoToTop, k.ResultSet.GoToBottom},
		{k.Editor.Up, k.Editor.Down, k.Editor.Left, k.Editor.Right, k.Editor.Insert, k.Editor.Normal, k.Editor.ExecuteQuery, k.Editor.ExecuteSingleQuery},
		{k.Navigation.Up, k.Navigation.Down, k.Navigation.Left, k.Navigation.Right, k.Help, k.Quit},
	}
}

type EditorKeyMap struct {
	// Normal Mode Navigation.
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Internal navigation.
	LineStart  key.Binding
	LineEnd    key.Binding
	GoToTop    key.Binding
	GoToBottom key.Binding

	// Mode Switching.
	Insert key.Binding
	Normal key.Binding

	// Actions.
	ExecuteQuery       key.Binding
	ExecuteSingleQuery key.Binding
}

func DefaultEditorKeyMap() EditorKeyMap {
	return EditorKeyMap{
		Up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "move right"),
		),
		// --- Internal navigation ---
		GoToTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to top (sidebar database graph)"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("G"), // Capital 'G' for shift+g
			key.WithHelp("G", "go to bottom (sidebar database graph)"),
		),
		LineStart: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "navigate all the way to the left of the table"),
		),
		LineEnd: key.NewBinding(
			key.WithKeys("$"),
			key.WithHelp("$", "navigate all the way to the right of the table"),
		),
		// --- Mode Switching ---
		Insert: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert mode"),
		),
		Normal: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "normal mode"),
		),
		// --- Actions ---
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute queries in the editor"),
		),
		ExecuteSingleQuery: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "execute single query"),
		),
	}
}

type NavigationKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

func DefaultNavitationKeyMap() NavigationKeyMap {
	return NavigationKeyMap{
		Up: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "Toggle to the panel above"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "Toggle to the panel below"),
		),
		Left: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "Toggle to the panel on the left"),
		),
		Right: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "Toggle to the panel on the right"),
		),
	}
}

type SiebarViewportKeyMap struct {
	GoToTop    key.Binding
	GoToBottom key.Binding
}

func DefaultSiebarViewportKeyMap() SiebarViewportKeyMap {
	return SiebarViewportKeyMap{
		GoToTop: key.NewBinding(
			key.WithKeys("shift+k"),
			key.WithHelp("shift+k", "go to top (sidebar database graph)"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("shift+j"),
			key.WithHelp("shift+j", "go to bottom (sidebar database graph)"),
		),
	}
}

type ResultSetKeyMap struct {
	NextTab    key.Binding
	PrevTab    key.Binding
	LineStart  key.Binding
	LineEnd    key.Binding
	GoToTop    key.Binding
	GoToBottom key.Binding
}

func DefaultResultSetKeyMap() ResultSetKeyMap {
	return ResultSetKeyMap{
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab (result set view)"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab (result set view)"),
		),
		GoToTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to top (sidebar database graph)"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("G"), // Capital 'G' for shift+g
			key.WithHelp("G", "go to bottom (sidebar database graph)"),
		),
		LineStart: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "navigate all the way to the left of the table"),
		),
		LineEnd: key.NewBinding(
			key.WithKeys("$"),
			key.WithHelp("$", "navigate all the way to the right of the table"),
		),
	}
}
