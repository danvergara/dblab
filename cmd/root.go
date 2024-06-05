package cmd

import (
	"github.com/danvergara/gocui"
	"github.com/spf13/cobra"

	"github.com/danvergara/dblab/pkg/app"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/config"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/danvergara/dblab/pkg/form"
)

// Define all the global flags.
var (
	cfg         bool
	cfgName     string
	driver      string
	url         string
	host        string
	port        string
	user        string
	pass        string
	schema      string
	db          string
	ssl         string
	limit       uint
	socket      string
	sslcert     string
	sslkey      string
	sslpassword string
	sslrootcert string
	// oracle specific.
	traceFile string
	sslVerify string
	wallet    string
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
				opts, err = config.Init(cfgName)
				if err != nil {
					return err
				}
			} else {
				opts = command.Options{
					Driver:      driver,
					URL:         url,
					Host:        host,
					Port:        port,
					User:        user,
					Pass:        pass,
					DBName:      db,
					Schema:      schema,
					SSL:         ssl,
					Limit:       limit,
					Socket:      socket,
					SSLCert:     sslcert,
					SSLKey:      sslkey,
					SSLPassword: sslpassword,
					SSLRootcert: sslrootcert,
					SSLVerify:   sslVerify,
					TraceFile:   traceFile,
					Wallet:      wallet,
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

			gcui, err := gocui.NewGui(gocui.OutputNormal)
			if err != nil {
				return err
			}

			app, err := app.New(gcui, opts)
			if err != nil {
				return err
			}

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

	// config file flag.
	rootCmd.PersistentFlags().
		BoolVarP(&cfg, "config", "", false, "Get the connection data from a config file (default locations are: current directory, $HOME/.dblab.yaml or $XDG_CONFIG_HOME/.dblab.yaml)")
	// cfg-name is used to indicate the name of the config section to be used to establish a
	// connection with desired database.
	// default: if empty, the first item of the databases options is gonna be selected.
	rootCmd.Flags().StringVarP(&cfgName, "cfg-name", "", "", "Database config name section")

	// global flags used to open a database connection.
	rootCmd.Flags().StringVarP(&driver, "driver", "", "", "Database driver")
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Database connection string")
	rootCmd.Flags().StringVarP(&host, "host", "", "", "Server host name or IP")
	rootCmd.Flags().StringVarP(&port, "port", "", "", "Server port")
	rootCmd.Flags().StringVarP(&user, "user", "", "", "Database user")
	rootCmd.Flags().StringVarP(&pass, "pass", "", "", "Password for user")
	rootCmd.Flags().StringVarP(&db, "db", "", "", "Database name")
	rootCmd.Flags().StringVarP(&schema, "schema", "", "", "Database schema (postgres only)")
	rootCmd.Flags().StringVarP(&ssl, "ssl", "", "", "SSL mode")
	rootCmd.Flags().
		UintVarP(&limit, "limit", "", 100, "Size of the result set from the table content query (should be greater than zero, otherwise the app will error out)")
	rootCmd.Flags().StringVarP(&socket, "socket", "", "", "Path to a Unix socket file")
	rootCmd.Flags().StringVarP(
		&sslcert,
		"sslcert",
		"",
		"",
		"This parameter specifies the file name of the client SSL certificate, replacing the default ~/.postgresql/postgresql.crt",
	)
	rootCmd.Flags().StringVarP(
		&sslkey,
		"sslkey",
		"",
		"",
		"This parameter specifies the location for the secret key used for the client certificate. It can either specify a file name that will be used instead of the default ~/.postgresql/postgresql.key, or it can specify a key obtained from an external “engine”",
	)
	rootCmd.Flags().
		StringVarP(&sslpassword, "sslpassword", "", "", "This parameter specifies the password for the secret key specified in sslkey")
	rootCmd.Flags().StringVarP(
		&sslrootcert,
		"sslrootcert",
		"",
		"",
		"This parameter specifies the name of a file containing SSL certificate authority (CA) certificate(s) The default is ~/.postgresql/root.crt",
	)
	rootCmd.Flags().
		StringVarP(&sslVerify, "ssl-verify", "", "", "[enable|disable] or [true|false] enable ssl verify for the server")
	rootCmd.Flags().StringVarP(&traceFile, "trace-file", "", "", "File name for trace log")
	rootCmd.Flags().StringVarP(&wallet, "wallet", "", "", "Path for auto-login oracle wallet")
}
