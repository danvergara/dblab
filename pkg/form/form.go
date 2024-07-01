package form

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muesli/termenv"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
)

const (
	defaultLimit = 100
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
	term         = termenv.ColorProfile()
)

// Model is a meta-model.
type Model struct {
	// menu management.
	cursor int
	steps  int

	// driver.
	drivers []string
	driver  string

	// ssl connection params.
	sslCertInput     textinput.Model
	sslKeyInput      textinput.Model
	sslPasswordInput textinput.Model
	sslRootcertInput textinput.Model

	// oracle specific.
	traceFileInput textinput.Model
	sslVerifyInput textinput.Model
	walletInput    textinput.Model

	// sql server.
	trustServerCertificateInput textinput.Model

	// std data.
	hostInput     textinput.Model
	portInput     textinput.Model
	userInput     textinput.Model
	passwordInput textinput.Model
	databaseInput textinput.Model
	filePathInput textinput.Model
	limitInput    textinput.Model

	// ssl.
	postgreSQLSSLModes []string
	mySQLSSLModes      []string
	oracleSSLModes     []string
	sqlServerSSLModes  []string
	sqliteSSLModes     []string
	sslMode            string
}

// Init initialize the meta-model.
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update update the view of the meta-model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if the pressed keys are esc or ctrl + c, finish the execution.
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.steps {
	case 0:
		return updateDriver(msg, m)
	case 1:
		return updateStd(msg, m)
	case 2:
		return updateSSLMode(msg, m)
	case 3:
		return updateSSLConn(msg, m)
	}

	return m, tea.Quit
}

// View displays the content on the terminal.
func (m *Model) View() string {
	var s string

	switch m.steps {
	case 0:
		s = driverView(m)
	case 1:
		s = standardView(m)
	case 2:
		s = sslView(m)
	case 3:
		s = sslConnView(m)
	}

	return fmt.Sprint(s)
}

// Host returns the host value.
func (m *Model) Host() string {
	return m.hostInput.Value()
}

// Port returns the Port value.
func (m *Model) Port() string {
	return m.portInput.Value()
}

// User returns the user value.
func (m *Model) User() string {
	return m.userInput.Value()
}

// Password returns the password value.
func (m *Model) Password() string {
	return m.passwordInput.Value()
}

// Database returns the database name value.
func (m *Model) Database() string {
	return m.databaseInput.Value()
}

// SSLMode returns the ssl mode name value.
func (m *Model) SSLMode() string {
	return m.sslMode
}

func (m *Model) SSLCert() string {
	return m.sslCertInput.Value()
}

func (m *Model) SSLKey() string {
	return m.sslKeyInput.Value()
}

func (m *Model) SSLPassword() string {
	return m.sslPasswordInput.Value()
}

func (m *Model) SSLRootcert() string {
	return m.sslRootcertInput.Value()
}

func (m *Model) SSLVerify() string {
	return m.sslVerifyInput.Value()
}

func (m *Model) TraceFile() string {
	return m.traceFileInput.Value()
}

func (m *Model) Wallet() string {
	return m.walletInput.Value()
}

func (m *Model) TrustServerCertificate() string {
	return m.trustServerCertificateInput.Value()
}

// Limit returns the limit input value from the user.
func (m *Model) Limit() (uint, error) {
	// if the user skipped the question, resort to default value
	if m.limitInput.Value() == "" {
		return defaultLimit, nil
	}
	limit, err := strconv.Atoi(m.limitInput.Value())
	if err != nil {
		return uint(0), err
	}

	if limit <= 0 {
		return uint(0), fmt.Errorf("invalid limit %d", limit)
	}

	return uint(limit), nil
}

// FilePath returns the path to the database file (just in sqlite3) value.
func (m *Model) FilePath() string {
	return m.filePathInput.Value()
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func initModel() Model {
	host := textinput.NewModel()
	host.Placeholder = "Host"
	host.PromptStyle = focusedStyle
	host.TextStyle = focusedStyle
	host.CharLimit = 200
	host.Focus()

	port := textinput.NewModel()
	port.Placeholder = "Port"
	port.CharLimit = 200

	user := textinput.NewModel()
	user.Placeholder = "Username"
	user.CharLimit = 200

	password := textinput.NewModel()
	password.Placeholder = "Password"
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '*'
	password.CharLimit = 200

	database := textinput.NewModel()
	database.Placeholder = "Database"
	database.CharLimit = 200

	limit := textinput.NewModel()
	limit.Placeholder = "Limit"
	limit.CharLimit = 200

	filePath := textinput.NewModel()
	filePath.Placeholder = "File Path"
	filePath.CharLimit = 1000
	filePath.Focus()

	sslCert := textinput.NewModel()
	sslCert.Placeholder = "Client SSL certificate"
	sslCert.CharLimit = 1000
	sslCert.Focus()

	sslKey := textinput.NewModel()
	sslKey.Placeholder = "The location for the secret key used for the client certificate"
	sslKey.CharLimit = 1000

	sslPassword := textinput.NewModel()
	sslPassword.Placeholder = "The password for the secret key"
	sslPassword.CharLimit = 1000

	sslRootCert := textinput.NewModel()
	sslRootCert.Placeholder = "The name of a file containing SSL certificate authority (CA) certificate(s)"
	sslRootCert.CharLimit = 1000

	sslVerify := textinput.NewModel()
	sslVerify.Placeholder = "SSL Verify"
	sslVerify.CharLimit = 200

	traceFile := textinput.NewModel()
	traceFile.Placeholder = "Trace file"
	traceFile.CharLimit = 1000

	wallet := textinput.NewModel()
	wallet.Placeholder = "Path to wallet"
	wallet.CharLimit = 1000

	trustServerCertificate := textinput.NewModel()
	trustServerCertificate.Placeholder = "Server certificate is checked or not"
	trustServerCertificate.CharLimit = 1000
	trustServerCertificate.Focus()

	m := Model{
		// the supported drivers by the client.
		drivers: []string{"postgres", "mysql", "sqlite", "oracle", "sqlserver"},
		// our default value.
		driver: "postgres",

		sslMode:            "disable",
		postgreSQLSSLModes: []string{"disable", "require", "verify-full", "verify-ca"},
		mySQLSSLModes:      []string{"true", "false", "skip-verify", "preferred"},
		oracleSSLModes:     []string{"enable", "disable"},
		sqlServerSSLModes:  []string{"strict", "disable", "false", "true"},
		sqliteSSLModes:     []string{},

		hostInput:                   host,
		portInput:                   port,
		userInput:                   user,
		passwordInput:               password,
		databaseInput:               database,
		limitInput:                  limit,
		filePathInput:               filePath,
		sslCertInput:                sslCert,
		sslKeyInput:                 sslKey,
		sslPasswordInput:            sslPassword,
		sslRootcertInput:            sslRootCert,
		sslVerifyInput:              sslVerify,
		traceFileInput:              traceFile,
		walletInput:                 wallet,
		trustServerCertificateInput: trustServerCertificate,
	}

	return m
}

// Run runs the menus programs to introduced the required data to connect with a database.
func Run() (command.Options, error) {
	m := initModel()
	if err := tea.NewProgram(&m).Start(); err != nil {
		return command.Options{}, err
	}

	limit, err := m.Limit()
	if err != nil {
		return command.Options{}, err
	}

	opts := command.Options{
		Driver:                 m.driver,
		Host:                   m.Host(),
		Port:                   m.Port(),
		User:                   m.User(),
		Pass:                   m.Password(),
		DBName:                 m.Database(),
		SSL:                    m.SSLMode(),
		SSLCert:                m.SSLCert(),
		SSLKey:                 m.SSLKey(),
		SSLPassword:            m.SSLPassword(),
		SSLRootcert:            m.SSLRootcert(),
		SSLVerify:              m.SSLVerify(),
		TraceFile:              m.TraceFile(),
		Wallet:                 m.Wallet(),
		TrustServerCertificate: m.TrustServerCertificate(),
		Limit:                  limit,
	}

	if m.driver == drivers.SQLServer {
		opts.Encrypt = m.SSLMode()
	}

	if m.driver == "sqlite" {
		opts.URL = fmt.Sprintf("file:%s", m.FilePath())
	}

	return opts, nil
}

// IsEmpty checks if the given options objects is empty.
func IsEmpty(opts command.Options) bool {
	return cmp.Equal(
		opts,
		command.Options{},
		cmpopts.IgnoreFields(
			command.Options{},
			"SSL",
			"Limit",
			"Socket",
			"SSL",
			"SSLCert",
			"SSLKey",
			"SSLPassword",
			"SSLRootcert",
			"TraceFile",
			"SSLVerify",
			"Wallet",
			"TrustServerCertificate",
		),
	)
}
