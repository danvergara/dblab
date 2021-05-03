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

func nextView(from, to string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == from {
			_, err := g.SetCurrentView(to)

			g.Highlight = true
			g.Cursor = true
			g.SelFgColor = gocui.ColorGreen

			return err
		}

		_, err := g.SetCurrentView(from)

		return err
	}
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		l, err := v.Line(cy + 1)
		if err != nil {
			return err
		}
		if l != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// quit is called to end the gui app.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
