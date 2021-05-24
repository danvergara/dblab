package form

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var term = termenv.ColorProfile()

type driverModel struct {
	drivers []string
	cursor  int
	driver  string
}

// initializeDriverModel is the initial model used to select one driver.
var initializeDriverModel = &driverModel{
	// the supported drivers by the client.
	drivers: []string{"postgres", "mysql"},
	// our default value.
	driver: "postgres",
}

// Init the drivers menu.
func (d *driverModel) Init() tea.Cmd {
	return nil
}

// Update manage the drivers menu.
func (d *driverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return d, tea.Quit
		// the "up" and "k" keys mve the cursor up.
		case "up", "k":
			if d.cursor > 0 {
				d.cursor--
			}
		// the "down" and "j" keys move the cursor down.
		case "down", "j":
			if d.cursor < len(d.drivers)-1 {
				d.cursor++
			}
		case "enter", " ":
			driver := d.drivers[d.cursor]
			d.driver = driver
			return d, tea.Quit
		}
	}

	return d, nil
}

// View renders the drivers menu.
func (d *driverModel) View() string {
	// The header.
	s := "Select the database driver:"
	var choices string
	// Iterate over the driver.
	for i, driver := range d.drivers {
		choices += fmt.Sprintf("%s\n", checkbox(driver, d.cursor == i))
	}

	return fmt.Sprintf("%s\n\n%s", s, choices)
}

func (d *driverModel) Driver() string {
	return d.driver
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
