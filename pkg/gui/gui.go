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
func New(g *gocui.Gui, c *client.Client) (*Gui, error) {

	gui := Gui{
		g: g,
		c: c,
	}

	if err := gui.prepare(); err != nil {
		return &gui, err
	}

	return &gui, nil
}

func (gui *Gui) prepare() error {
	gui.g.Highlight = true
	gui.g.Cursor = true
	gui.g.Mouse = true
	gui.g.SelFrameColor = gocui.ColorGreen

	gui.setLayout()

	if err := gui.keybindings(); err != nil {
		return err
	}

	return nil
}

// Run setup the gui with keybindings and start the mainloop.
func (gui *Gui) Run() error {
	if err := gui.g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}

	return nil
}

// Gui returns a pointer of a gocui.Gui instance.
func (gui *Gui) Gui() *gocui.Gui {
	return gui.g
}
