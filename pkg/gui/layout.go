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
	options = []string{"Rows", "Structure", "Constraints"}
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

	if v, err := gui.g.SetView("tables", 0, int(0.16*float32(maxY)), int(0.19*float32(maxX)), int(0.95*float32(maxY))); err != nil {
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

	if v, err := gui.g.SetView("constraints", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.95*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Constraints"

		fmt.Fprintln(v, "Please select a table!")
	}

	if v, err := gui.g.SetView("structure", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.95*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Structure"

		fmt.Fprintln(v, "Please select a table!")
	}

	if v, err := gui.g.SetView("rows", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.95*float32(maxY))); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Rows"

		fmt.Fprintln(v, "Type the sql query above. Press Ctrl-c to quit.")
	}

	return nil
}

func moveDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "query" {
		_, err := g.SetCurrentView("rows")
		if err != nil {
			return err
		}
		_, err = g.SetViewOnTop("rows")
		if err != nil {
			return err
		}

		g.Highlight = true
		g.Cursor = true

		return err
	}

	_, err := g.SetCurrentView("view")

	return err
}

func nextView(from, to string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == from {
			_, err := g.SetCurrentView(to)

			g.Highlight = true
			g.Cursor = true

			return err
		}

		_, err := g.SetCurrentView(from)

		return err
	}
}

func setViewOnTop(from, to string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		if v == nil || v.Name() == from {
			tmpOptions := make([]string, len(options))
			copy(tmpOptions, options)

			for i, o := range tmpOptions {
				if strings.ToLower(o) == to {
					tmpOptions[i] = green.Sprint(o)
				}
			}

			nv, err := g.View("navigation")
			if err != nil {
				return err
			}
			nv.Clear()
			fmt.Fprint(nv, strings.Join(tmpOptions, "   "))

			return switchView(g, to)

		}

		if v == nil || v.Name() == to {
			tmpOptions := make([]string, len(options))
			copy(tmpOptions, options)

			for i, o := range tmpOptions {
				if strings.ToLower(o) == from {
					tmpOptions[i] = green.Sprint(o)
				}
			}

			nv, err := g.View("navigation")
			if err != nil {
				return err
			}
			nv.Clear()
			fmt.Fprint(nv, strings.Join(tmpOptions, "   "))
			return switchView(g, from)
		}

		return nil
	}
}

func switchView(g *gocui.Gui, v string) error {
	if _, err := g.SetViewOnTop(v); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(v); err != nil {
		return err
	}

	g.Highlight = true
	g.Cursor = true

	return nil
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

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx+1, cy); err != nil {
			if err := v.SetOrigin(ox+1, oy); err != nil {
				return err
			}
		}
	}

	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
			if err := v.SetOrigin(ox-1, oy); err != nil {
				return err
			}
		}
	}

	return nil
}

func navigation(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		word, _ := v.Word(cx, cy)

		tmpOptions := make([]string, len(options))
		copy(tmpOptions, options)

		if word != "" {
			for i, o := range tmpOptions {
				if o == word {
					tmpOptions[i] = green.Sprint(word)
				}
			}

			v.Clear()
			fmt.Fprint(v, strings.Join(tmpOptions, "   "))

			err := switchView(g, strings.ToLower(word))
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// quit is called to end the gui app.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
