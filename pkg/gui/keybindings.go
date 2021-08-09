package gui

import (
	"github.com/danvergara/gocui"
)

// keyBinding struct used to defines the multiple actions defined to interact with dblab.
type keyBinding struct {
	view    string
	key     interface{}
	mod     gocui.Modifier
	handler func(*gocui.Gui, *gocui.View) error
}

var keybindings []keyBinding = []keyBinding{
	{
		view:    "query",
		key:     gocui.KeyCtrlH,
		mod:     gocui.ModNone,
		handler: nextView("query", "tables"),
	},
	{
		view:    "tables",
		key:     gocui.KeyCtrlL,
		mod:     gocui.ModNone,
		handler: nextView("tables", "query"),
	},
	{
		view:    "query",
		key:     gocui.KeyCtrlJ,
		mod:     gocui.ModNone,
		handler: nextView("query", "rows"),
	},
	{
		view:    "rows",
		key:     gocui.KeyCtrlK,
		mod:     gocui.ModNone,
		handler: nextView("rows", "query"),
	},
	{
		view:    "rows",
		key:     gocui.KeyCtrlH,
		mod:     gocui.ModNone,
		handler: nextView("rows", "tables"),
	},
	{
		view:    "structure",
		key:     gocui.KeyCtrlH,
		mod:     gocui.ModNone,
		handler: nextView("structure", "tables"),
	},
	{
		view:    "structure",
		key:     gocui.KeyCtrlK,
		mod:     gocui.ModNone,
		handler: nextView("structure", "query"),
	},
	{
		view:    "constraints",
		key:     gocui.KeyCtrlH,
		mod:     gocui.ModNone,
		handler: nextView("constraints", "tables"),
	},
	{
		view:    "constraints",
		key:     gocui.KeyCtrlK,
		mod:     gocui.ModNone,
		handler: nextView("constraints", "query"),
	},
	{
		view:    "tables",
		key:     'k',
		mod:     gocui.ModNone,
		handler: moveCursorVertically("up"),
	},
	{
		view:    "tables",
		key:     gocui.KeyArrowUp,
		mod:     gocui.ModNone,
		handler: moveCursorVertically("up"),
	},
	{
		view:    "tables",
		key:     'j',
		mod:     gocui.ModNone,
		handler: moveCursorVertically("down"),
	},
	{
		view:    "tables",
		key:     gocui.KeyArrowDown,
		mod:     gocui.ModNone,
		handler: moveCursorVertically("down"),
	},
	{
		view:    "rows",
		key:     gocui.KeyArrowUp,
		mod:     gocui.ModNone,
		handler: moveCursorVertically("up"),
	},
	{
		view:    "rows",
		key:     'k',
		mod:     gocui.ModNone,
		handler: moveCursorVertically("up"),
	},
	{
		view:    "rows",
		key:     gocui.KeyArrowDown,
		mod:     gocui.ModNone,
		handler: moveCursorVertically("down"),
	},
	{
		view:    "rows",
		key:     'j',
		mod:     gocui.ModNone,
		handler: moveCursorVertically("down"),
	},
	{
		view:    "rows",
		key:     gocui.KeyArrowRight,
		mod:     gocui.ModNone,
		handler: moveCursorHorizontally("right"),
	},
	{
		view:    "rows",
		key:     'l',
		mod:     gocui.ModNone,
		handler: moveCursorHorizontally("right"),
	},
	{
		view:    "rows",
		key:     gocui.KeyArrowLeft,
		mod:     gocui.ModNone,
		handler: moveCursorHorizontally("left"),
	},
	{
		view:    "rows",
		key:     'h',
		mod:     gocui.ModNone,
		handler: moveCursorHorizontally("left"),
	},
	{
		view:    "constraints",
		key:     gocui.KeyCtrlF,
		mod:     gocui.ModNone,
		handler: setViewOnTop("constraints", "rows"),
	},
	{
		view:    "rows",
		key:     gocui.KeyCtrlF,
		mod:     gocui.ModNone,
		handler: setViewOnTop("rows", "constraints"),
	},
	{
		view:    "structure",
		key:     gocui.KeyCtrlS,
		mod:     gocui.ModNone,
		handler: setViewOnTop("structure", "rows"),
	},
	{
		view:    "rows",
		key:     gocui.KeyCtrlS,
		mod:     gocui.ModNone,
		handler: setViewOnTop("rows", "structure"),
	},
	{
		view:    "navigation",
		key:     gocui.MouseLeft,
		mod:     gocui.ModNone,
		handler: navigation,
	},
	{
		view:    "",
		key:     gocui.KeyCtrlC,
		mod:     gocui.ModNone,
		handler: quit,
	},
}

func (gui *Gui) keybindings() error {
	for _, k := range keybindings {
		if err := gui.g.SetKeybinding(k.view, k.key, k.mod, k.handler); err != nil {
			return err
		}
	}

	// SQL helpers
	if err := gui.g.SetKeybinding("query", gocui.KeyCtrlSpace, gocui.ModNone, gui.inputQuery()); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("tables", gocui.KeyEnter, gocui.ModNone, gui.metadata); err != nil {
		return err
	}

	return nil
}
