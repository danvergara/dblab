package form

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type sslModel struct {
	modes  []string
	cursor int
	result string
}

// initialSslModel returns a tea.Model that displays a list of modes based on the given driver.
func initialSslModel(driver string) *sslModel {
	var modes []string

	switch driver {
	case "postgres":
		modes = []string{"disable", "require", "verify-full"}
	case "mysql":
		modes = []string{"true", "false", "skip-verify", "preferred"}
	default:
		modes = []string{"disable", "require", "verify-full"}
	}

	return &sslModel{
		modes: modes,
	}
}

func (m *sslModel) Init() tea.Cmd {
	return nil
}

func (m *sslModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
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
		case "enter", " ":
			result := m.modes[m.cursor]
			m.result = result
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the drivers menu.
func (m *sslModel) View() string {
	// The header.
	s := "\nSelect the ssl mode:"
	var choices string
	// Iterate over the driver.
	for i, mode := range m.modes {
		choices += fmt.Sprintf("%s\n", checkbox(mode, m.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
}
