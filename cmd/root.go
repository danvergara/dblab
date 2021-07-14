package cmd

import (
	"github.com/danvergara/dblab/pkg/app"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/config"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/danvergara/dblab/pkg/form"
	"github.com/danvergara/dblab/pkg/gui"
	"github.com/danvergara/gocui"
	"github.com/spf13/cobra"
)

// Define all the global flags.
var (
	cfg    bool
	driver string
	url    string
	host   string
	port   string
	user   string
	pass   string
	db     string
	ssl    string
)

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dblab",
		Short: "Interactive database client",
		Long:  `dblab is a terminal UI based interactive database client for Postgres and MySQL.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts command.Options
			var err error

			if cfg {
				opts, err = config.Init()
				if err != nil {
					return err
				}
			} else {
				opts = command.Options{
					Driver: driver,
					URL:    url,
					Host:   host,
					Port:   port,
					User:   user,
					Pass:   pass,
					DBName: db,
					SSL:    ssl,
				}

				if form.IsEmpty(opts) {
					opts, err = form.Run()
					if err != nil {
						return err
					}
				}
			}

			if err := connection.ValidateOpts(opts); err != nil {
				return err
			}

			c, err := client.New(opts)
			if err != nil {
				return err
			}

			gcui, err := gocui.NewGui(gocui.OutputNormal)
			if err != nil {
				return err
			}

			g := gui.New(gcui, c)

			app := app.New(g)

			if err := app.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = NewRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&cfg, "config", "", false, "get the connection data from a config file (default is $HOME/.dblab.yaml or the current directory)")

	// global flags used to open a database connection.
	rootCmd.Flags().StringVarP(&driver, "driver", "", "", "Database driver")
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Database connection string")
	rootCmd.Flags().StringVarP(&host, "host", "", "", "Server host name or IP")
	rootCmd.Flags().StringVarP(&port, "port", "", "", "Server port")
	rootCmd.Flags().StringVarP(&user, "user", "", "", "Database user")
	rootCmd.Flags().StringVarP(&pass, "pass", "", "", "Password for user")
	rootCmd.Flags().StringVarP(&db, "db", "", "", "Database name")
	rootCmd.Flags().StringVarP(&ssl, "ssl", "", "", "SSL mode")
}
