package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danvergara/dblab/pkg/drivers"
)

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
	var (
		cmd    tea.Cmd
		inputs []textinput.Model
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			inputs = stdInputs(m)

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

			assignStdInputValues(m, inputs)

			return m, nil
		}
	}

	m, cmd = updateInputs(msg, m)
	return m, cmd
}

func stdInputs(m *Model) []textinput.Model {
	var inputs []textinput.Model

	if m.driver == drivers.SQLite {
		inputs = []textinput.Model{
			m.filePathInput,
		}
	} else {
		inputs = []textinput.Model{
			m.hostInput,
			m.portInput,
			m.userInput,
			m.passwordInput,
			m.databaseInput,
		}
	}

	inputs = append(inputs, m.limitInput)

	return inputs
}

func assignStdInputValues(m *Model, inputs []textinput.Model) {
	if m.driver == drivers.SQLite && len(inputs) == 2 {
		m.filePathInput = inputs[0]
		m.limitInput = inputs[1]
	} else if len(inputs) == 6 {
		{
			m.hostInput = inputs[0]
			m.portInput = inputs[1]
			m.userInput = inputs[2]
			m.passwordInput = inputs[3]
			m.databaseInput = inputs[4]
			m.limitInput = inputs[5]
		}
	}
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
			if len(m.modes) > 0 {
				m.ssl = m.modes[m.cursor]
			}
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

	if m.driver == drivers.SQLite {
		m.filePathInput, cmd = m.filePathInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
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
	}

	m.limitInput, cmd = m.limitInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
