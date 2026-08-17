package keys

import (
	"charm.land/bubbles/v2/key"
	"github.com/danvergara/dblab/pkg/config"
)

type KeyMap struct {
	Help       key.Binding
	Quit       key.Binding
	History    key.Binding
	Navigation NavigationKeyMap
	Editor     EditorKeyMap
	Sidebar    SidebarKeyMap
	ResultSet  ResultSetKeyMap
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
		History: key.NewBinding(
			key.WithKeys("alt+h"),
			key.WithHelp("alt+h", "query history"),
		),
		Navigation: DefaultNavitationKeyMap(),
		Editor:     DefaultEditorKeyMap(),
		Sidebar:    DefaultSidebarKeyMap(),
		ResultSet:  DefaultResultSetKeyMap(),
	}
}

func ReadKeyMapFromConfig() (KeyMap, error) {
	cfg, err := config.GetKeyMap()
	if err != nil {
		return KeyMap{}, err
	}

	km := KeyMap{
		Help:    key.NewBinding(key.WithKeys(cfg.KeyBindings.Help), key.WithHelp(cfg.KeyBindings.Help, "toggle help")),
		Quit:    key.NewBinding(key.WithKeys(cfg.KeyBindings.Quit), key.WithHelp(cfg.KeyBindings.Quit, "quit")),
		History: key.NewBinding(key.WithKeys(cfg.KeyBindings.History), key.WithHelp(cfg.KeyBindings.History, "query history")),
		Navigation: NavigationKeyMap{
			Up:    key.NewBinding(key.WithKeys(cfg.KeyBindings.Navigation.Up), key.WithHelp(cfg.KeyBindings.Navigation.Up, "Toggle to the panel above")),
			Down:  key.NewBinding(key.WithKeys(cfg.KeyBindings.Navigation.Down), key.WithHelp(cfg.KeyBindings.Navigation.Down, "Toggle to the panel below")),
			Left:  key.NewBinding(key.WithKeys(cfg.KeyBindings.Navigation.Left), key.WithHelp(cfg.KeyBindings.Navigation.Left, "Toggle to the panel on the left")),
			Right: key.NewBinding(key.WithKeys(cfg.KeyBindings.Navigation.Right), key.WithHelp(cfg.KeyBindings.Navigation.Right, "Toggle to the panel on the right")),
		},
		Editor: EditorKeyMap{
			Up:                 key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Up), key.WithHelp(cfg.KeyBindings.Editor.Up, "move up (editor)")),
			Down:               key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Down), key.WithHelp(cfg.KeyBindings.Editor.Down, "move down (editor)")),
			Left:               key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Left), key.WithHelp(cfg.KeyBindings.Editor.Left, "move left (editor)")),
			Right:              key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Right), key.WithHelp(cfg.KeyBindings.Editor.Right, "move right (editor)")),
			Insert:             key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Insert), key.WithHelp(cfg.KeyBindings.Editor.Insert, "insert mode (editor)")),
			Normal:             key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.Normal), key.WithHelp(cfg.KeyBindings.Editor.Normal, "normal mode (editor)")),
			LineStart:          key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.LineStart), key.WithHelp(cfg.KeyBindings.Editor.LineStart, "line start (editor)")),
			LineEnd:            key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.LineEnd), key.WithHelp(cfg.KeyBindings.Editor.LineEnd, "line end (editor)")),
			GoToTop:            key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.GoToTop), key.WithHelp(cfg.KeyBindings.Editor.GoToTop, "go top (editor)")),
			GoToBottom:         key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.GoToBottom), key.WithHelp(cfg.KeyBindings.Editor.GoToBottom, "go bottom (editor)")),
			ExecuteQuery:       key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.ExecuteQuery), key.WithHelp(cfg.KeyBindings.Editor.ExecuteQuery, "execute queries in the editor (editor)")),
			ExecuteSingleQuery: key.NewBinding(key.WithKeys(cfg.KeyBindings.Editor.ExecuteSingleQuery), key.WithHelp(cfg.KeyBindings.Editor.ExecuteSingleQuery, "execute single query (editor)")),
		},
		Sidebar: SidebarKeyMap{
			GoToTop:    key.NewBinding(key.WithKeys(cfg.KeyBindings.Sidebar.GoToTop), key.WithHelp(cfg.KeyBindings.Sidebar.GoToTop, "go top in the (sidebar)")),
			GoToBottom: key.NewBinding(key.WithKeys(cfg.KeyBindings.Sidebar.GoToBottom), key.WithHelp(cfg.KeyBindings.Sidebar.GoToBottom, "go bottom (sidebar)")),
		},
		ResultSet: ResultSetKeyMap{
			PrevTab:    key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.PrevTab), key.WithHelp(cfg.KeyBindings.ResultSet.PrevTab, "prev tab (result set) ")),
			NextTab:    key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.NextTab), key.WithHelp(cfg.KeyBindings.ResultSet.NextTab, "next tab (result set)")),
			LineStart:  key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.LineStart), key.WithHelp(cfg.KeyBindings.ResultSet.LineStart, "line start (result set)")),
			LineEnd:    key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.LineEnd), key.WithHelp(cfg.KeyBindings.ResultSet.LineEnd, "line end (result set)")),
			GoToTop:    key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.GoToTop), key.WithHelp(cfg.KeyBindings.ResultSet.GoToTop, "go top (result set)")),
			GoToBottom: key.NewBinding(key.WithKeys(cfg.KeyBindings.ResultSet.GoToBottom), key.WithHelp(cfg.KeyBindings.ResultSet.GoToBottom, "go bottom (result set)")),
		},
	}

	return km, nil
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Sidebar.GoToBottom, k.Sidebar.GoToTop},
		{k.ResultSet.NextTab, k.ResultSet.PrevTab, k.ResultSet.LineStart, k.ResultSet.LineEnd, k.ResultSet.GoToTop, k.ResultSet.GoToBottom, k.Editor.ExecuteQuery, k.Editor.ExecuteSingleQuery},
		{k.Editor.Up, k.Editor.Down, k.Editor.Left, k.Editor.Right, k.Editor.Insert, k.Editor.Normal, k.Editor.GoToBottom, k.Editor.GoToTop},
		{k.Editor.LineStart, k.Editor.LineEnd, k.Navigation.Up, k.Navigation.Down, k.Navigation.Left, k.Navigation.Right, k.Help, k.Quit},
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

type SidebarKeyMap struct {
	GoToTop    key.Binding
	GoToBottom key.Binding
}

func DefaultSidebarKeyMap() SidebarKeyMap {
	return SidebarKeyMap{
		GoToTop: key.NewBinding(
			key.WithKeys("alt+k"),
			key.WithHelp("alt+k", "go to top (sidebar database graph)"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("alt+j"),
			key.WithHelp("alt+j", "go to bottom (sidebar database graph)"),
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
