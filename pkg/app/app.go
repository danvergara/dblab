package app

import (
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/gui"
	"github.com/danvergara/gocui"
)

// App Struct.
type App struct {
	g *gui.Gui
	c *client.Client
}

// New bootstrap a new application.
func New(g *gocui.Gui, opts command.Options) (*App, error) {
	c, err := client.New(opts)
	if err != nil {
		return nil, err
	}

	gcui, err := gui.New(g, c)
	if err != nil {
		return nil, err
	}

	app := App{
		g: gcui,
		c: c,
	}

	return &app, nil
}

// Run runs the application.
func (a *App) Run() error {

	defer func() {
		_ = a.c.DB().Close()
		a.g.Gui().Close()
	}()

	if err := a.g.Run(); err != nil {
		return err
	}

	return nil
}
