package form

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muesli/termenv"
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

	// std data.
	hostInput     textinput.Model
	portInput     textinput.Model
	userInput     textinput.Model
	passwordInput textinput.Model
	databaseInput textinput.Model
	filePathInput textinput.Model
	limitInput    textinput.Model

	// ssl.
	modes []string
	ssl   string
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
		return updateSSL(msg, m)
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

// SSL returns the ssl name value.
func (m *Model) SSL() string {
	return m.ssl
}

// Limit returns the limit input value from the user.
func (m *Model) Limit() int {
	limit, err := strconv.Atoi(m.limitInput.Value())
	if err != nil {
		return 100
	}

	if limit <= 0 {
		return 100
	}

	return limit
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

	m := Model{
		// the supported drivers by the client.
		drivers: []string{"postgres", "mysql", "sqlite3"},
		// our default value.
		driver: "postgres",

		hostInput:     host,
		portInput:     port,
		userInput:     user,
		passwordInput: password,
		databaseInput: database,
		limitInput:    limit,
		filePathInput: filePath,
	}

	return m
}

// Run runs the menus programs to introduced the required data to connect with a database.
func Run() (command.Options, error) {
	m := initModel()
	if err := tea.NewProgram(&m).Start(); err != nil {
		return command.Options{}, err
	}

	opts := command.Options{
		Driver: m.driver,
		Host:   m.Host(),
		Port:   m.Port(),
		User:   m.User(),
		Pass:   m.Password(),
		DBName: m.Database(),
		SSL:    m.SSL(),
		Limit:  m.Limit(),
	}

	if m.driver == "sqlite3" {
		opts.URL = fmt.Sprintf("file:%s", m.FilePath())
	}

	return opts, nil
}

// IsEmpty checks if the given options objects is empty.
func IsEmpty(opts command.Options) bool {
	return cmp.Equal(opts, command.Options{}, cmpopts.IgnoreFields(command.Options{}, "SSL", "Limit"))
}
