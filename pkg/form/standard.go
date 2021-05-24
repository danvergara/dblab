package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle            = lipgloss.NewStyle()

	focusedSubmitButton = "[ " + focusedStyle.Render("Submit") + " ]"
	blurredSubmitButton = "[ " + blurredButtonStyle.Render("Submit") + " ]"
)

type standardModel struct {
	index         int
	hostInput     textinput.Model
	portInput     textinput.Model
	userInput     textinput.Model
	passwordInput textinput.Model
	databaseInput textinput.Model
	submitButton  string
}

// InitialStandardModel initializes a standardModel object.
func initialStandardModel() *standardModel {
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

	return &standardModel{
		index:         0,
		hostInput:     host,
		portInput:     port,
		userInput:     user,
		passwordInput: password,
		databaseInput: database,
		submitButton:  blurredSubmitButton,
	}
}

func (m *standardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *standardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			inputs := []textinput.Model{
				m.hostInput,
				m.portInput,
				m.userInput,
				m.passwordInput,
				m.databaseInput,
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.index == len(inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.index--
			} else {
				m.index++
			}

			if m.index > len(inputs) {
				m.index = 0
			} else if m.index < 0 {
				m.index = len(inputs)
			}

			for i := 0; i <= len(inputs)-1; i++ {
				if i == m.index {
					// Set focused state
					inputs[i].Focus()
					inputs[i].PromptStyle = focusedStyle
					inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				inputs[i].Blur()
				inputs[i].PromptStyle = noStyle
				inputs[i].TextStyle = noStyle
			}

			m.hostInput = inputs[0]
			m.portInput = inputs[1]
			m.userInput = inputs[2]
			m.passwordInput = inputs[3]
			m.databaseInput = inputs[4]

			if m.index == len(inputs) {
				m.submitButton = focusedSubmitButton
			} else {
				m.submitButton = blurredSubmitButton
			}

			return m, nil
		}
	}

	m, cmd = updateInputs(msg, m)
	return m, cmd
}

func updateInputs(msg tea.Msg, m *standardModel) (*standardModel, tea.Cmd) {
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

func (m *standardModel) View() string {
	s := "\nIntroduce the connection params:\n\n"

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

	s += "\n\n" + m.submitButton + "\n"
	return s
}
