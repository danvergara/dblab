package gui

import (
	"github.com/jroimartin/gocui"
)

func (gui *Gui) keybindings() error {
	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlH, gocui.ModNone, setQueryView); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("query", gocui.KeyEnter, gocui.ModNone, gui.runQuery()); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyCtrlL, gocui.ModNone, setTablesView); err != nil {
		return err
	}

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

	if err := gui.g.SetKeybinding("tables", gocui.KeyEnter, gocui.ModNone, gui.selectTable); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlJ, gocui.ModNone, setRowsView); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyCtrlK, gocui.ModNone, setQueryViewFromRows); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("rows", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
