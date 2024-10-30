package app

import (
	"github.com/danvergara/gocui"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/tui"
)

// App Struct.
type App struct {
	t *tui.Tui
	c *client.Client
}

// New bootstrap a new application.
func New(g *gocui.Gui, opts command.Options) (*App, error) {
	c, err := client.New(opts)
	if err != nil {
		return nil, err
	}

	t, err := tui.New(c)
	if err != nil {
		return nil, err
	}

	app := App{
		t: t,
		c: c,
	}

	return &app, nil
}

// Run runs the application.
func (a *App) Run() error {
	defer func() {
		_ = a.c.DB().Close()
	}()

	if err := a.t.Run(); err != nil {
		return err
	}

	return nil
}
