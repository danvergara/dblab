package app

import (
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/gui"
)

// App Struct.
type App struct {
	c *client.Client
	g *gui.Gui
}

// New bootstrap a new application.
func New(c *client.Client, g *gui.Gui) *App {
	return &App{
		c: c,
		g: g,
	}
}

// Run runs the application.
func (a *App) Run() error {
	if err := a.g.Run(); err != nil {
		return err
	}

	return nil
}
