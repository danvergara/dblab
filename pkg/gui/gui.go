package gui

import (
	"errors"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/gocui"
)

// Gui wraps the gocui Gui object which handles rendering and events.
type Gui struct {
	g *gocui.Gui
	c *client.Client
}

// New builds a new gui handler.
func New(g *gocui.Gui, c *client.Client) *Gui {
	return &Gui{
		g: g,
		c: c,
	}
}

// Run setup the gui with keybindings and start the mainloop.
func (gui *Gui) Run() error {
	defer gui.g.Close()

	gui.g.Highlight = true
	gui.g.Cursor = true
	gui.g.Mouse = true
	gui.g.SelFrameColor = gocui.ColorGreen

	gui.g.SetManagerFunc(gui.layout)

	if err := gui.keybindings(); err != nil {
		return err
	}

	if err := gui.g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}
