package form

import "fmt"

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

	inputs := viewInputs(m)

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"
	return s
}

func viewInputs(m *Model) []string {
	var inputs []string

	if m.driver == "sqlite3" {
		inputs = []string{
			m.filePathInput.View(),
			m.limitInput.View(),
		}
	} else {
		inputs = []string{
			m.hostInput.View(),
			m.portInput.View(),
			m.userInput.View(),
			m.passwordInput.View(),
			m.databaseInput.View(),
			m.limitInput.View(),
		}
	}

	return inputs
}

func sslView(m *Model) string {
	switch m.driver {
	case "postgres":
		m.modes = []string{"disable", "require", "verify-full"}
	case "mysql":
		m.modes = []string{"true", "false", "skip-verify", "preferred"}
	case "sqlite3":
		m.modes = []string{}
	default:
		m.modes = []string{"disable", "require", "verify-full"}
	}

	// The header.
	s := "\nSelect the ssl mode (just press enter if you selected sqlite3):"
	var choices string
	// Iterate over the driver.
	for i, mode := range m.modes {
		choices += fmt.Sprintf("%s\n", checkbox(mode, m.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
}
