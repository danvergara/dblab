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

			if s == "enter" && m.cursor == len(inputs)-1 {
				m.steps = 2
				m.cursor = 0

				if m.driver == drivers.SQLite {
					return m, tea.Quit
				}

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

func updateSSLMode(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
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
			switch m.driver {
			case drivers.Postgres:
				if m.cursor < len(m.postgreSQLSSLModes)-1 {
					m.cursor++
				}
			case drivers.Oracle:
				if m.cursor < len(m.oracleSSLModes)-1 {
					m.cursor++
				}
			case drivers.SQLServer:
				if m.cursor < len(m.sqlServerSSLModes)-1 {
					m.cursor++
				}
			case drivers.MySQL:
				if m.cursor < len(m.mySQLSSLModes)-1 {
					m.cursor++
				}
			}
		case "enter":
			switch m.driver {
			case drivers.Postgres:
				m.sslMode = m.postgreSQLSSLModes[m.cursor]
			case drivers.MySQL:
				m.sslMode = m.mySQLSSLModes[m.cursor]
			case drivers.Oracle:
				m.sslMode = m.oracleSSLModes[m.cursor]
			case drivers.SQLServer:
				m.sslMode = m.sqlServerSSLModes[m.cursor]
			}

			m.steps = 3
			m.cursor = 0

			switch m.driver {
			case drivers.Postgres, drivers.Oracle, drivers.SQLServer, drivers.MySQL:
				if m.sslMode == "disable" || m.sslMode == "false" {
					return m, tea.Quit
				}
			}

			return m, nil
		}
	}

	return m, nil
}

func updateSSLConn(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var (
		cmd    tea.Cmd
		inputs []textinput.Model
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			inputs = sslConnInputs(m)

			s := msg.String()

			if s == "enter" && m.cursor == len(inputs)-1 {
				m.steps = 4
				m.cursor = 0

				return m, tea.Quit
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

			assignSSLConnInputValues(m, inputs)

			return m, nil
		}
	}

	m, cmd = updateSSLConnInputs(msg, m)
	return m, cmd
}

func sslConnInputs(m *Model) []textinput.Model {
	var inputs []textinput.Model

	switch m.driver {
	case drivers.Postgres:
		inputs = []textinput.Model{
			m.sslCertInput,
			m.sslKeyInput,
			m.sslPasswordInput,
			m.sslRootcertInput,
		}
	case drivers.Oracle:
		inputs = []textinput.Model{
			m.traceFileInput,
			m.sslVerifyInput,
			m.walletInput,
		}
	case drivers.SQLServer:
		inputs = []textinput.Model{
			m.trustServerCertificateInput,
		}
	}

	return inputs
}

func assignSSLConnInputValues(m *Model, inputs []textinput.Model) {
	switch m.driver {
	case drivers.Postgres:
		if len(inputs) == 4 {
			m.sslCertInput = inputs[0]
			m.sslKeyInput = inputs[1]
			m.sslPasswordInput = inputs[2]
			m.sslRootcertInput = inputs[3]
		}
	case drivers.Oracle:
		if len(inputs) == 3 {
			m.traceFileInput = inputs[0]
			m.sslVerifyInput = inputs[1]
			m.walletInput = inputs[2]
		}
	case drivers.SQLServer:
		if len(inputs) == 1 {
			m.trustServerCertificateInput = inputs[0]
		}
	}
}

func updateSSLConnInputs(msg tea.Msg, m *Model) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.driver {
	case drivers.Postgres:
		m.sslCertInput, cmd = m.sslCertInput.Update(msg)
		cmds = append(cmds, cmd)

		m.sslKeyInput, cmd = m.sslKeyInput.Update(msg)
		cmds = append(cmds, cmd)

		m.sslPasswordInput, cmd = m.sslPasswordInput.Update(msg)
		cmds = append(cmds, cmd)

		m.sslRootcertInput, cmd = m.sslRootcertInput.Update(msg)
		cmds = append(cmds, cmd)
	case drivers.Oracle:
		m.traceFileInput, cmd = m.traceFileInput.Update(msg)
		cmds = append(cmds, cmd)

		m.sslVerifyInput, cmd = m.sslVerifyInput.Update(msg)
		cmds = append(cmds, cmd)

		m.walletInput, cmd = m.walletInput.Update(msg)
		cmds = append(cmds, cmd)
	case drivers.SQLServer:
		m.trustServerCertificateInput, cmd = m.trustServerCertificateInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
