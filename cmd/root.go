package cmd

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"log"
	"os"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/gui"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Define all the global flags.
var (
	cfgFile string
	driver  string
	url     string
	host    string
	port    string
	user    string
	pass    string
	db      string
	ssl     string
)

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dblab",
		Short: "Interactive databse client",
		Long:  `dblab is a terminal UI based interactive database client for Postgres, MySQL and SQLite.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			opts := command.Options{
				Driver: driver,
				URL:    url,
				Host:   host,
				Port:   port,
				User:   user,
				Pass:   pass,
				DBName: db,
				SSL:    ssl,
			}

			if opts.Host == "" && opts.Port == "" && opts.User == "" && opts.Pass == "" && opts.DBName == "" && opts.Driver == "" && opts.URL == "" {
				return fmt.Errorf("empty values required to open a session in database")
			}

			g, err := gocui.NewGui(gocui.OutputNormal)
			if err != nil {
				log.Panicln(err)
			}
			defer g.Close()

			g.SetManagerFunc(gui.Layout)

			if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
				log.Panicln(err)
			}

			if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
				log.Panicln(err)
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dblab.yaml)")

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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".dblab" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dblab")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
