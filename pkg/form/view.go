package form

import (
	"fmt"
	"strings"

	"github.com/danvergara/dblab/pkg/drivers"
)

func driverView(m *Model) string {
	// The header.
	s := "Select the database driver:"
	var choices strings.Builder
	// Iterate over the drivers.
	for i, driver := range m.drivers {
		choices.WriteString(fmt.Sprintf("%s\n", checkbox(driver, m.cursor == i)))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices.String())
}

func standardView(m *Model) string {
	var s strings.Builder
	s.WriteString("Introduce the connection params:\n\n")

	inputs := viewInputs(m)

	for i := range inputs {
		s.WriteString(inputs[i])
		if i < len(inputs)-1 {
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	return s.String()
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
	var choices strings.Builder
	// Iterate over the driver.
	for i, mode := range sslModes {
		choices.WriteString(fmt.Sprintf("%s\n", checkbox(mode, m.cursor == i)))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices.String())
}

func sslConnView(m *Model) string {
	var s strings.Builder
	s.WriteString("Introduce the SSL connection params:\n\n")

	inputs := sslConnViewInputs(m)

	for i := range inputs {
		s.WriteString(inputs[i])
		if i < len(inputs)-1 {
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")

	return s.String()
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
