package command

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
)

// Options is a struct that stores the provided commands by the user.
type Options struct {
	Driver string
	URL    string
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
	// PostgreSQL only.
	Schema string
	Limit  uint
	Socket string
	SSL    string
	// SSH.
	SSHHost          string
	SSHPort          string
	SSHUser          string
	SSHPass          string
	SSHKeyFile       string
	SSHKeyPassphrase string
	// SSL connection params.
	SSLCert     string
	SSLKey      string
	SSLPassword string
	SSLRootcert string
	// oracle specific.
	TraceFile string
	SSLVerify string
	Wallet    string
	// sql server.
	Encrypt                string
	TrustServerCertificate string
	ConnectionTimeout      string
	// TUI keybidings.
	TUIKeyBindings TUIKeyMap
}

// UpdateKeybindings method updates the TUIKeyBindings field, since the keybidings configuration parted ways with the connection configuration.
func (o *Options) UpdateKeybindings(k TUIKeyMap) {
	o.TUIKeyBindings = k
}

type TUIKeyMap struct {
	ExecuteQuery    key.Binding
	NextTab         key.Binding
	PrevTab         key.Binding
	PageTop         key.Binding
	PageBottom      key.Binding
	EndOfLine       key.Binding
	BeginningOfLine key.Binding
	Navigation      TUINavigationKeyMap
	Editor          EditorKeyMap
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
	ExecuteQuery key.Binding
}

type TUINavigationKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

func DefaultKeyMap() TUIKeyMap {
	return TUIKeyMap{
		ExecuteQuery: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "execute query"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),
		PageTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to top"),
		),
		PageBottom: key.NewBinding(
			key.WithKeys("G"), // Capital 'G' for shift+g
			key.WithHelp("G", "go to bottom"),
		),
		BeginningOfLine: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "navigate all the way to the left of the table"),
		),
		EndOfLine: key.NewBinding(
			key.WithKeys("$"),
			key.WithHelp("$", "navigate all the way to the right of the table"),
		),
		Navigation: TUINavigationKeyMap{
			Up: key.NewBinding(
				key.WithKeys("ctrl+k"),
				key.WithKeys("ctrl+k", "Toggle to the panel above"),
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
				key.WithHelp("k", "up"),
			),
			Down: key.NewBinding(
				key.WithKeys("j"),
				key.WithHelp("j", "down"),
			),
			Left: key.NewBinding(
				key.WithKeys("h"),
				key.WithHelp("h", "left"),
			),
			Right: key.NewBinding(
				key.WithKeys("l"),
				key.WithHelp("l", "right"),
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
				key.WithHelp("ctrl+e", "execute query"),
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
