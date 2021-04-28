package gui

import "github.com/jroimartin/gocui"

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g *gocui.Gui
}

// New builds a new gui handler
func New(g *gocui.Gui) *Gui {
	return &Gui{
		g: g,
	}
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	defer gui.g.Close()

	gui.g.SetManagerFunc(gui.layout)

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}

	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
