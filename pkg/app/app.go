package app

import (
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/sshdb"
	"github.com/danvergara/dblab/pkg/tui"
)

// App Struct.
type App struct {
	t *tui.Tui
	c *client.Client
}

// New bootstrap a new application.
func New(opts command.Options) (*App, error) {
	var sshConfig *sshdb.SSHConfig

	if opts.SSHHost != "" {
		sshConfig = sshdb.New(
			sshdb.WithDBDriver(opts.Driver),
			sshdb.WithSShHost(opts.SSHHost),
			sshdb.WithSShPort(opts.SSHPort),
			sshdb.WithSSHUser(opts.SSHUser),
			sshdb.WithPass(opts.SSHPass),
			sshdb.WithSSHKeyFile(opts.SSHKeyFile),
			sshdb.WithSSHKeyPass(opts.SSHKeyPassphrase),
			sshdb.WithDBDURL(opts.URL),
		)

		if err := sshConfig.SSHTunnel(); err != nil {
			return nil, err
		}

		defer sshConfig.Close()
	}

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
