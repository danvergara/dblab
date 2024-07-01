package form

import (
	"fmt"

	"github.com/danvergara/dblab/pkg/drivers"
)

func driverView(m *Model) string {
	// The header.
	s := "Select the database driver:"
	var choices string
	// Iterate over the drivers.
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

	if m.driver == drivers.SQLite {
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
	var sslModes []string

	switch m.driver {
	case drivers.Postgres:
		sslModes = m.postgreSQLSSLModes
	case drivers.MySQL:
		sslModes = m.mySQLSSLModes
	case drivers.Oracle:
		sslModes = m.oracleSSLModes
	case drivers.SQLite:
		sslModes = m.sqliteSSLModes
	case drivers.SQLServer:
		sslModes = m.sqlServerSSLModes
	default:
		sslModes = m.postgreSQLSSLModes
	}

	// The header.
	s := "\nSelect the ssl mode (just press enter if you selected sqlite3):"
	var choices string
	// Iterate over the driver.
	for i, mode := range sslModes {
		choices += fmt.Sprintf("%s\n", checkbox(mode, m.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
}

func sslConnView(m *Model) string {
	s := "Introduce the SSL connection params:\n\n"

	inputs := sslConnViewInputs(m)

	for i := 0; i < len(inputs); i++ {
		s += inputs[i]
		if i < len(inputs)-1 {
			s += "\n"
		}
	}

	s += "\n"

	return s
}

func sslConnViewInputs(m *Model) []string {
	var inputs []string
	switch m.driver {
	case drivers.Postgres:
		if m.sslMode == "require" || m.sslMode == "verify-full" || m.sslMode == "verify-ca" {
			inputs = []string{
				m.sslCertInput.View(),
				m.sslKeyInput.View(),
				m.sslPasswordInput.View(),
				m.sslRootcertInput.View(),
			}
		}
	case drivers.Oracle:
		inputs = []string{
			m.traceFileInput.View(),
			m.sslVerifyInput.View(),
			m.walletInput.View(),
		}
	case drivers.SQLServer:
		if m.sslMode == "strict" || m.sslMode == "true" {
			inputs = []string{
				m.trustServerCertificateInput.View(),
			}
		}

	}

	return inputs
}
