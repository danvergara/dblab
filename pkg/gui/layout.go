package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/danvergara/gocui"
	"github.com/fatih/color"
)

var (
	green   = color.New(color.FgGreen).Add(color.Bold)
	options = []string{"Rows", "Structure", "Constraints", "Indexes"}
)

// Layout is called for every screen re-render e.g. when the screen is resized.
func (gui *Gui) layout(g *gocui.Gui) error {
	maxX, maxY := gui.g.Size()

	if v, err := gui.g.SetView("banner", 0, 0, int(0.19*float32(maxX)), int(0.14*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.FrameColor = gocui.ColorMagenta
		myFigure := figure.NewFigure("dblab", "", true)
		figure.Write(v, myFigure)
	}

	if v, err := gui.g.SetView("tables", 0, int(0.16*float32(maxY)), int(0.19*float32(maxX)), int(0.94*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		if err := gui.showTables(); err != nil {
			return err
		}

		v.Title = "Tables"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := gui.g.SetView("navigation", int(0.2*float32(maxX)), 0, maxX-1, int(0.07*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Navigation"

		tmpOptions := make([]string, len(options))
		copy(tmpOptions, options)
		tmpOptions[0] = green.Sprint("Rows")

		fmt.Fprint(v, strings.Join(tmpOptions, "   "))
	}

	if v, err := gui.g.SetView("query", int(0.2*float32(maxX)), int(0.09*float32(maxY)), maxX-1, int(0.27*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "SQL Query"
		v.Editable = true
		v.Wrap = true
		v.Highlight = true

		if _, err := gui.g.SetCurrentView("query"); err != nil {
			return err
		}
	}

	if v, err := gui.g.SetView("indexes", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Indexes"

		fmt.Fprintln(v, "Please select a table!")
	}

	if v, err := gui.g.SetView("constraints", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Constraints"

		fmt.Fprintln(v, "Please select a table!")
	}

	if v, err := gui.g.SetView("structure", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Structure"

		fmt.Fprintln(v, "Please select a table!")
	}

	if v, err := gui.g.SetView("rows", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Rows"

		fmt.Fprintln(v, "Type the sql query above. Press Ctrl-c to quit.")
	}

	if v, err := gui.g.SetView("current-page", int(0.82*float32(maxX)), int(0.96*float32(maxY)), int(0.88*float32(maxX)), int(0.99*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		fmt.Fprintln(v, "100 rows")
	}

	if v, err := gui.g.SetView("prev-page", int(0.89*float32(maxX)), int(0.96*float32(maxY)), int(0.91*float32(maxX)), int(0.99*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		fmt.Fprint(v, " < ")
	}

	if v, err := gui.g.SetView("page", int(0.92*float32(maxX)), int(0.96*float32(maxY)), int(0.97*float32(maxX)), int(0.99*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		fmt.Fprint(v, " 1 of 100 ")
	}

	if v, err := gui.g.SetView("next-page", int(0.98*float32(maxX)), int(0.96*float32(maxY)), maxX-1, int(0.99*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		fmt.Fprint(v, " > ")
	}

	return nil
}
