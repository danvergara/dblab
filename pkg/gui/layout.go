package gui

import (
	"github.com/jroimartin/gocui"
)

// Layout is called for every screen re-render e.g. when the screen is resized.
func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := gui.g.Size()

	if v, err := gui.g.SetView("tables", 0, 0, int(0.2*float32(maxX)), maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := gui.showTables(); err != nil {
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

		v.Title = "Rows"
	}

	return nil
}

func setQueryView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "tables" {
		_, err := g.SetCurrentView("query")
		return err
	}

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	_, err := g.SetCurrentView("tables")
	return err
}

func setTablesView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "query" {
		_, err := g.SetCurrentView("tables")

		g.Highlight = true
		g.Cursor = true
		g.SelFgColor = gocui.ColorGreen

		return err
	}

	_, err := g.SetCurrentView("query")
	return err
}

func setRowsView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "query" {
		_, err := g.SetCurrentView("rows")

		g.Highlight = true
		g.Cursor = true
		g.SelFgColor = gocui.ColorGreen

		return err
	}

	_, err := g.SetCurrentView("query")
	return err
}

func setQueryViewFromRows(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "rows" {
		_, err := g.SetCurrentView("query")
		return err
	}

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	_, err := g.SetCurrentView("rows")
	return err
}

// quit is called to end the gui app.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
