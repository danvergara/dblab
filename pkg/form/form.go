package form

import tea "github.com/charmbracelet/bubbletea"

// Run runs the menus programs to introduced the required data to connect with a database.
func Run() error {
	if err := tea.NewProgram(initializeDriverModel).Start(); err != nil {
		return err
	}

	if err := tea.NewProgram(initialStandardModel()).Start(); err != nil {
		return err
	}

	driver := initializeDriverModel.Driver()
	if err := tea.NewProgram(initialSslModel(driver)).Start(); err != nil {
		return err
	}

	return nil
}
