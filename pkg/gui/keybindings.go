package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) keybindings() error {
	// navigation between panels.
	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlH, gocui.ModNone, nextView("query", "tables")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyCtrlL, gocui.ModNone, nextView("tables", "query")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlJ, gocui.ModNone, nextView("query", "rows")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyCtrlK, gocui.ModNone, nextView("rows", "query")); err != nil {
		return err
	}

	// SQL helpers
	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlSpace, gocui.ModNone, gui.inputQuery()); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyEnter, gocui.ModNone, gui.selectTable); err != nil {
		return err
	}

	// navigation directives for the tables panel.
	if err := gui.g.SetKeybinding("tables", gocui.KeyCtrlK, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyCtrlJ, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	// navigation directives for the  rows panel.
	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	// quit function event.
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
