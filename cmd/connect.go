package cmd

import (
	"fmt"

	"github.com/danvergara/dblab/pkg/app"
	"github.com/danvergara/dblab/pkg/bubbletui"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/config"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command.
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "dblab connect is a command used to re-use successful connections",
	Long: `dblab connect is a command that let the user re-user previous successful connections,
so the user does not have to type the creds every time. 
This command uses huh to show a form to list the database profiles listed stored in
$XDG_CONFIG_HOME/dblab/dblab.json.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called")
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var kb = command.DefaultKeyMap()

		opts, err := bubbletui.Run()
		if err != nil {
			return err
		}

		if keybindings {
			kb, err = config.SetupKeyMap()
			if err != nil {
				return err
			}
		}

		if err := connection.ValidateOpts(opts); err != nil {
			return err
		}

		app, err := app.New(opts, kb)
		if err != nil {
			return err
		}

		if err := app.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
