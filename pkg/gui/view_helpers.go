package gui

import (
	"fmt"
	"strings"

	"github.com/danvergara/gocui"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// nextView the app to another view.
// the function makes sure the initial view
// is the same as the current view in the gui.
// Also checks if the next view is part of the hidden views.
// If so, the function will update the navigation view.
func nextView(from, to string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == from {
			if err := switchView(g, to); err != nil {
				return err
			}
		}

		if contains(options, cases.Title(language.English).String(to)) {
			if err := handleNavigationOptions(g, to); err == nil {
				return err
			}
		}

		return nil
	}
}

// switchView set the given view on top and
// makes it the current view
// in the state management of the gui.
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

// moveCursorHorizontally moves the cursor to the given direction one position.
func moveCursorHorizontally(direction string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var position int

		switch direction {
		case "right":
			position = 1
		case "left":
			position = -1
		}

		if v != nil {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()

			if err := v.SetCursor(cx+position, cy); err != nil {
				if position > 0 || ox > 0 {
					if err := v.SetOrigin(ox+position, oy); err != nil {
						return err
					}
				}
			}
		}

		return nil
	}
}

// moveCursorVertically moves the cursor vertically given a direction.
// the down position is handled in a special way to prevent the cursor keep going
// down when there's no characaters in the next lines.
func moveCursorVertically(direction string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()

			switch direction {
			case "up":
				if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
					if err := v.SetOrigin(ox, oy-1); err != nil {
						return err
					}
				}
			case "down":
				l, err := v.Line(cy + 1)
				if err != nil {
					return err
				}

				if l != "" {
					if err := v.SetCursor(cx, cy+1); err != nil {
						if err := v.SetOrigin(ox, oy+1); err != nil {
							return err
						}
					}
				}
			}
		}

		return nil
	}
}

// setViewOnTop sets a given view (defined as to) on top of the other (define as from).
func setViewOnTop(from, to string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil || v.Name() == from {
			if err := handleNavigationOptions(g, to); err != nil {
				return err
			}

			return switchView(g, to)
		}

		return nil
	}
}

// handleNavigationOptions manages the colors on the options on the Navigation view.
// colorized the given option with the green color, clears the view and
// then prints the colorized options again in the view.
func handleNavigationOptions(g *gocui.Gui, opt string) error {
	if opt == "" {
		return fmt.Errorf("empty option passed as parameter")
	}

	tmpOptions := make([]string, len(options))
	copy(tmpOptions, options)

	for i, o := range tmpOptions {
		if strings.EqualFold(o, opt) {
			tmpOptions[i] = green.Sprint(o)
		}
	}

	nv, err := g.View("navigation")
	if err != nil {
		return err
	}

	nv.Clear()
	fmt.Fprint(nv, strings.Join(tmpOptions, "   "))

	return nil
}

// navigation manages the navigation between the hidden menus.
// If the user clicks on an option on the navigation view,
// the corresponde view will be set on top and active.
func navigation(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		opt, _ := v.Word(cx, cy)

		// if err is equal to nil, proceed switching to the given view.
		if err := handleNavigationOptions(g, opt); err == nil {
			if err := switchView(g, strings.ToLower(opt)); err != nil {
				return err
			}
		}
	}

	return nil
}

// contains checks if a string is present in a slice.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// quit is called to end the gui app.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
