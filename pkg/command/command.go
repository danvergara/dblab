package command

import (
	"os"

	"charm.land/bubbles/v2/key"
)

// Options is a struct that stores the provided commands by the user.
type Options struct {
	Driver string `json:"driver"`
	URL    string `json:"url"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"-"`
	DBName string `json:"db_name"`
	Schema string `json:"schema"`
	Limit  uint   `json:"limit"`
	Socket string `json:"socket"`
	SSL    string `json:"ssl"`
	// SSH.
	SSHHost          string `json:"ssh_host"`
	SSHPort          string `json:"ssh_port"`
	SSHUser          string `json:"ssh_user"`
	SSHPass          string `json:"-"`
	SSHKeyFile       string `json:"ssh_key_file"`
	SSHKeyPassphrase string `json:"ssh_key_passphrase"`
	// SSL connection params.
	SSLCert     string `json:"ssl_cert"`
	SSLKey      string `json:"ssl_key"`
	SSLPassword string `json:"-"`
	SSLRootcert string `json:"ssl_rootcert"`
	// oracle specific.
	TraceFile string `json:"trace_file"`
	SSLVerify string `json:"ssl_verify"`
	Wallet    string `json:"wallet"`
	// sql server.
	Encrypt                string `json:"encrypt"`
	TrustServerCertificate string `json:"trust_server_certificate"`
	ConnectionTimeout      string `json:"connection_timeout"`
	// Read Only mode.
	ReadOnly bool `json:"read_only"`
}

type TUIKeyMap struct {
	NextTab         key.Binding
	PrevTab         key.Binding
	PageTop         key.Binding
	PageBottom      key.Binding
	EndOfLine       key.Binding
	BeginningOfLine key.Binding
	Help            key.Binding
	Quit            key.Binding
	Navigation      TUINavigationKeyMap
	Editor          EditorKeyMap
}

func (k TUIKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k TUIKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextTab, k.PrevTab, k.PageTop, k.PageBottom, k.EndOfLine, k.BeginningOfLine},
		{k.Navigation.Up, k.Navigation.Down, k.Navigation.Left, k.Navigation.Right, k.Help, k.Quit},
		{k.Editor.Up, k.Editor.Down, k.Editor.Left, k.Editor.Right, k.Editor.Insert, k.Editor.Normal, k.Editor.ExecuteQuery, k.Editor.ExecuteSingleQuery},
	}
}

type EditorKeyMap struct {
	// Normal Mode Navigation.
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Mode Switching.
	Insert key.Binding
	Normal key.Binding

	// Actions.
	ExecuteQuery       key.Binding
	ExecuteSingleQuery key.Binding
}

type TUINavigationKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

func DefaultKeyMap() *TUIKeyMap {
	return &TUIKeyMap{
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab (result set view)"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab (result set view)"),
		),
		PageTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to top (sidebar database graph)"),
		),
		PageBottom: key.NewBinding(
			key.WithKeys("G"), // Capital 'G' for shift+g
			key.WithHelp("G", "go to bottom (sidebar database graph)"),
		),
		BeginningOfLine: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "navigate all the way to the left of the table"),
		),
		EndOfLine: key.NewBinding(
			key.WithKeys("$"),
			key.WithHelp("$", "navigate all the way to the right of the table"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Navigation: TUINavigationKeyMap{
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
		},
		Editor: EditorKeyMap{
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
		},
	}
}

// SetDefault returns a Options struct and fills the empty
// values with environment variables if any.
func SetDefault(opts Options) Options {
	if opts.URL == "" {
		opts.URL = os.Getenv("DATABASE_URL")
	}

	if opts.Driver == "" {
		opts.Driver = os.Getenv("DB_DRIVER")
	}

	if opts.Host == "" {
		opts.Host = os.Getenv("DB_HOST")
	}

	if opts.User == "" {
		opts.User = os.Getenv("DB_USER")
	}

	if opts.Pass == "" {
		opts.Pass = os.Getenv("DB_PASSWORD")
	}

	if opts.DBName == "" {
		opts.DBName = os.Getenv("DB_NAME")
	}

	if opts.Port == "" {
		opts.Port = os.Getenv("DB_PORT")
	}

	if opts.Schema == "" {
		opts.Schema = os.Getenv("DB_SCHEMA")
	}

	return opts
}
