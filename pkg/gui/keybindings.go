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

	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlJ, gocui.ModNone, moveDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyCtrlK, gocui.ModNone, nextView("rows", "query")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyCtrlH, gocui.ModNone, nextView("rows", "tables")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("structure", gocui.KeyCtrlH, gocui.ModNone, nextView("structure", "tables")); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("structure", gocui.KeyCtrlK, gocui.ModNone, nextView("structure", "query")); err != nil {
		return err
	}

	// SQL helpers
	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlSpace, gocui.ModNone, gui.inputQuery()); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyEnter, gocui.ModNone, gui.renderStructure); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyEnter, gocui.ModNone, gui.selectTable); err != nil {
		return err
	}

	// navigation directives for the tables panel.
	if err := gui.g.SetKeybinding("tables", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	// navigation directives for the  rows panel.
	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", 'l', gocui.ModNone, cursorRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", 'h', gocui.ModNone, cursorLeft); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("structure", gocui.KeyCtrlS, gocui.ModNone, setViewOnTop); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyCtrlS, gocui.ModNone, setViewOnTop); err != nil {
		return err
	}

	// quit function event.
	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
