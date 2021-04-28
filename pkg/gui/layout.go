package gui

import (
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

// Layout is called for every screen re-render e.g. when the screen is resized.
func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := gui.g.Size()

	if v, err := gui.g.SetView("tables", 0, 0, int(0.2*float32(maxX)), maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Tables"
	}

	if v, err := gui.g.SetView("query", int(0.2*float32(maxX)), 0, maxX, maxY-40); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "SQL Query"

		v.Editable = true
		v.Wrap = true

		if _, err := gui.g.SetCurrentView("query"); err != nil {
			return err
		}
	}

	if v, err := gui.g.SetView("rows", int(0.2*float32(maxX)), maxY-40, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// init table
		table := tablewriter.NewWriter(v)
		table.SetCenterSeparator("|")
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.Render()

		v.Title = "Rows"
	}

	return nil
}

// Quit is called to end the gui app.
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
