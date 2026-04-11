package app

import (
	"github.com/danvergara/dblab/pkg/bubbletui"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/sshdb"
)

// App Struct.
type App struct {
	c  *client.Client
	sc *sshdb.SSHConfig
	m  *bubbletui.ModelV2
}

// New bootstrap a new application.
func New(opts command.Options) (*App, error) {
	var sc *sshdb.SSHConfig

	if opts.SSHHost != "" {
		sc = sshdb.New(
			sshdb.WithDBDriver(opts.Driver),
			sshdb.WithSShHost(opts.SSHHost),
			sshdb.WithSShPort(opts.SSHPort),
			sshdb.WithSSHUser(opts.SSHUser),
			sshdb.WithPass(opts.SSHPass),
			sshdb.WithSSHKeyFile(opts.SSHKeyFile),
			sshdb.WithSSHKeyPass(opts.SSHKeyPassphrase),
			sshdb.WithDBDURL(opts.URL),
		)

		if err := sc.SSHTunnel(); err != nil {
			return nil, err
		}
	}

	c, err := client.New(opts)
	if err != nil {
		return nil, err
	}

	m, err := bubbletui.NewModelV2(c, &opts.TUIKeyBindings)
	if err != nil {
		return nil, err
	}

	app := App{
		c:  c,
		m:  m,
		sc: sc,
	}

	return &app, nil
}

// Run runs the application.
func (a *App) Run() error {
	defer func() {
		if a.sc != nil {
			_ = a.sc.Close()
		}

		_ = a.c.DB().Close()
	}()

	if err := a.m.Run(); err != nil {
		return err
	}

	return nil
}
