package form

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muesli/termenv"
)

var (
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle            = lipgloss.NewStyle()
	term               = termenv.ColorProfile()
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
	// Make sure these keys always quit
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

	return fmt.Sprintf("%s", s)
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

// Database returns te database name value.
func (m *Model) Database() string {
	return m.databaseInput.Value()
}

// SSL returns te ssl name value.
func (m *Model) SSL() string {
	return m.ssl
}

func updateDriver(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		// the "up" and "k" keys mve the cursor up.
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// the "down" and "j" keys move the cursor down.
		case "down", "j":
			if m.cursor < len(m.drivers)-1 {
				m.cursor++
			}
		case "enter":
			driver := m.drivers[m.cursor]
			m.driver = driver
			m.cursor = 0
			m.steps = 1
			return m, nil
		}
	}

	return m, nil
}

func updateStd(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			inputs := []textinput.Model{
				m.hostInput,
				m.portInput,
				m.userInput,
				m.passwordInput,
				m.databaseInput,
			}

			s := msg.String()

			// If so, exit.
			if s == "enter" && m.cursor == len(inputs)-1 {
				m.steps = 2
				m.cursor = 0
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor > len(inputs) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.cursor {
					// Set focused state.
					inputs[i].Focus()
					inputs[i].PromptStyle = focusedStyle
					inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state.
				inputs[i].Blur()
				inputs[i].PromptStyle = noStyle
				inputs[i].TextStyle = noStyle
			}

			m.hostInput = inputs[0]
			m.portInput = inputs[1]
			m.userInput = inputs[2]
			m.passwordInput = inputs[3]
			m.databaseInput = inputs[4]

			return m, nil
		}
	}

	m, cmd = updateInputs(msg, m)
	return m, cmd
}

func updateSSL(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		// These keys should exit the program.
		// the "up" and "k" keys mve the cursor up.
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// the "down" and "j" keys move the cursor down.
		case "down", "j":
			if m.cursor < len(m.modes)-1 {
				m.cursor++
			}
		case "enter":
			m.ssl = m.modes[m.cursor]
			m.steps = 3
			m.cursor = 0
			return m, tea.Quit
		}
	}

	return m, nil
}

func updateInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.hostInput, cmd = m.hostInput.Update(msg)
	cmds = append(cmds, cmd)

	m.portInput, cmd = m.portInput.Update(msg)
	cmds = append(cmds, cmd)

	m.userInput, cmd = m.userInput.Update(msg)
	cmds = append(cmds, cmd)

	m.passwordInput, cmd = m.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	m.databaseInput, cmd = m.databaseInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func driverView(m *Model) string {
	// The header.
	s := "Select the database driver:"
	var choices string
	// Iterate over the driver.
	for i, driver := range m.drivers {
		choices += fmt.Sprintf("%s\n", checkbox(driver, m.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
}

func standardView(m *Model) string {
	s := "Introduce the connection params:\n\n"

	inputs := []string{
		m.hostInput.View(),
		m.portInput.View(),
		m.userInput.View(),
		m.passwordInput.View(),
		m.databaseInput.View(),
	}

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"
	return s
}

func sslView(m *Model) string {
	switch m.driver {
	case "postgres":
		m.modes = []string{"disable", "require", "verify-full"}
	case "mysql":
		m.modes = []string{"true", "false", "skip-verify", "preferred"}
	default:
		m.modes = []string{"disable", "require", "verify-full"}
	}

	// The header.
	s := "\nSelect the ssl mode:"
	var choices string
	// Iterate over the driver.
	for i, mode := range m.modes {
		choices += fmt.Sprintf("%s\n", checkbox(mode, m.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
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
	host.Focus()
	host.PromptStyle = focusedStyle
	host.TextStyle = focusedStyle
	host.CharLimit = 200

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

	m := Model{
		// the supported drivers by the client.
		drivers: []string{"postgres", "mysql"},
		// our default value.
		driver: "postgres",

		hostInput:     host,
		portInput:     port,
		userInput:     user,
		passwordInput: password,
		databaseInput: database,
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
	}

	return opts, nil
}

// IsEmpty checks if the given options objects is empty.
func IsEmpty(opts command.Options) bool {
	return cmp.Equal(opts, command.Options{}, cmpopts.IgnoreFields(command.Options{}, "SSL"))
}
