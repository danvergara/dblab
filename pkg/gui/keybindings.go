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

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
