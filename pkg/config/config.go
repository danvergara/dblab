package config

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kkyr/fig"
	"github.com/spf13/cobra"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
)

// Config struct is used to store the db connection data.
type Config struct {
	Database []Database
	User     string
	Pswd     string
	Host     string
	Port     string
	DBName   string
	Driver   string
	Limit    uint `fig:"limit" default:"100"`
}

type KeyBindingsConfig struct {
	KeyBindings KeyBindings
}

type Database struct {
	Name string `fig:"name"`

	Host     string
	Port     string
	DB       string
	User     string
	Password string
	Driver   string `validate:"required"`
	Schema   string

	// SSH Tunnel.
	SSHHost          string `fig:"ssh-host"`
	SSHPort          string `fig:"ssh-port"`
	SSHUser          string `fig:"ssh-user"`
	SSHPass          string `fig:"ssh-pass"`
	SSHKeyFile       string `fig:"ssh-key-file"`
	SSHKeyPassphrase string `fig:"ssh-key-pass"`

	// SSL connection params.
	SSL string `default:"disable"`

	SSLCert     string `fig:"sslcert"`
	SSLKey      string `fig:"sslkey"`
	SSLPassword string `fig:"sslpassword"`
	SSLRootcert string `fig:"sslrootcert"`

	// oracle specific.
	TraceFile string `fig:"trace"`
	SSLVerify string `fig:"ssl-verify"`
	Wallet    string `fig:"wallet"`

	// sql server.
	Encrypt                string `fig:"encrypt"`
	TrustServerCertificate string `fig:"trust-server-certificate"`
	ConnectionTimeout      string `fig:"connection-timeout"`
}

type KeyBindings struct {
	ExecuteQuery string `fig:"execute-query" default:"ctrl+e"`
	NextTab      string `fig:"next-tab"      default:"tab"`
	PrevTab      string `fig:"prev-tab"      default:"shift+tab"`
	PageTop      string `fig:"page-top"      default:"g"`
	PageBottom   string `fig:"page-bottom"   default:"G"`
	Navigation   NavigationBindgins
}

type NavigationBindgins struct {
	Up    string `fig:"up"    default:"ctrl+k"`
	Down  string `fig:"down"  default:"ctrl+j"`
	Left  string `fig:"left"  default:"ctrl+h"`
	Right string `fig:"right" default:"ctrl+l"`
}

// New returns a config instance the with db connection data inplace based on the flags of a cobra command.
func New(cmd *cobra.Command) *Config {
	conf := &Config{}

	cmd.PersistentFlags().StringVarP(&conf.User, "user", "", os.Getenv("DB_USER"), "DB user name")
	cmd.PersistentFlags().StringVarP(&conf.Pswd, "pswd", "", os.Getenv("DB_PASSWORD"), "DB pass")
	cmd.PersistentFlags().StringVarP(&conf.Port, "port", "", os.Getenv("DB_PORT"), "DB port")
	cmd.PersistentFlags().StringVarP(&conf.Host, "host", "", os.Getenv("DB_HOST"), "DB host")
	cmd.PersistentFlags().StringVarP(&conf.DBName, "name", "", os.Getenv("DB_NAME"), "DB name")
	cmd.PersistentFlags().
		StringVarP(&conf.Driver, "driver", "", os.Getenv("DB_DRIVER"), "DB driver")

	return conf
}

// Init reads in config file and returns a commands/Options instance.
func Init(configName string) (command.Options, error) {
	var opts command.Options
	var cfg Config
	var db Database

	configDir, err := os.UserConfigDir()
	if err != nil {
		return opts, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return opts, err
	}

	if err := fig.Load(&cfg, fig.File(".dblab.yaml"), fig.Dirs(".", home, configDir)); err != nil {
		return opts, err
	}

	if len(cfg.Database) == 0 {
		return opts, errors.New("empty database connection section on config file")
	}

	if configName != "" {
		for _, d := range cfg.Database {
			if configName == d.Name {
				db = d
			}
		}
	} else {
		db = cfg.Database[0]
	}

	opts = command.Options{
		Driver:                 db.Driver,
		Host:                   db.Host,
		Port:                   db.Port,
		User:                   db.User,
		Pass:                   db.Password,
		DBName:                 db.DB,
		Schema:                 db.Schema,
		Limit:                  cfg.Limit,
		SSL:                    db.SSL,
		SSLCert:                db.SSLCert,
		SSLKey:                 db.SSLKey,
		SSLPassword:            db.SSLPassword,
		SSLRootcert:            db.SSLRootcert,
		TraceFile:              db.TraceFile,
		SSLVerify:              db.SSLVerify,
		Wallet:                 db.Wallet,
		Encrypt:                db.Encrypt,
		TrustServerCertificate: db.TrustServerCertificate,
		ConnectionTimeout:      db.ConnectionTimeout,
		SSHHost:                db.SSHHost,
		SSHPort:                db.SSHPort,
		SSHUser:                db.SSHUser,
		SSHPass:                db.SSHPass,
		SSHKeyFile:             db.SSHKeyFile,
		SSHKeyPassphrase:       db.SSHKeyPassphrase,
	}

	return opts, nil
}

func SetupKeybindings() (command.TUIKeyBindings, error) {
	var kbc KeyBindingsConfig
	var tkb command.TUIKeyBindings

	configDir, err := os.UserConfigDir()
	if err != nil {
		return tkb, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return tkb, err
	}

	if err := fig.Load(&kbc, fig.File(".dblab.yaml"), fig.Dirs(".", home, configDir)); err != nil {
		return tkb, err
	}

	tkb = command.TUIKeyBindings{
		ExecuteQuery: key.NewBinding(key.WithKeys(kbc.KeyBindings.ExecuteQuery), key.WithHelp(kbc.KeyBindings.ExecuteQuery, "execute query")),
		NextTab:      key.NewBinding(key.WithKeys(kbc.KeyBindings.NextTab), key.WithHelp(kbc.KeyBindings.NextTab, "next tab")),
		PrevTab:      key.NewBinding(key.WithKeys(kbc.KeyBindings.PrevTab), key.WithHelp(kbc.KeyBindings.PrevTab, "previous tab")),
		PageTop:      key.NewBinding(key.WithKeys(kbc.KeyBindings.PageTop), key.WithHelp(kbc.KeyBindings.PageTop, "go to top")),
		PageBottom:   key.NewBinding(key.WithKeys(kbc.KeyBindings.PageBottom), key.WithHelp(kbc.KeyBindings.PageBottom, "go to bottom")),
		Navigation: command.TUINavigationBindgins{
			Up:    key.NewBinding(key.WithKeys(kbc.KeyBindings.Navigation.Up), key.WithHelp(kbc.KeyBindings.Navigation.Up, "Toggle to the panel above")),
			Down:  key.NewBinding(key.WithKeys(kbc.KeyBindings.Navigation.Down), key.WithHelp(kbc.KeyBindings.Navigation.Down, "Toggle to the panel below")),
			Left:  key.NewBinding(key.WithKeys(kbc.KeyBindings.Navigation.Left), key.WithHelp(kbc.KeyBindings.Navigation.Left, "Toggle to the panel on the left")),
			Right: key.NewBinding(key.WithKeys(kbc.KeyBindings.Navigation.Right), key.WithHelp(kbc.KeyBindings.Navigation.Right, "Toggle to the panel on the right")),
		},
	}

	return tkb, nil
}

// Open returns a db connection using the data from the config object.
func (c *Config) Open() (*sql.DB, error) {
	db, err := sql.Open(c.Driver, c.GetDBConnStr())
	if err != nil {
		fmt.Printf("Error Opening DB: %v \n", err)
		return nil, err
	}

	return db, err
}

// MigrateInstance returns a migrate instance based on the given driver.
func (c *Config) MigrateInstance() (*migrate.Migrate, error) {
	db, err := c.Open()
	if err != nil {
		return nil, err
	}

	switch c.Driver {
	case drivers.SQLite:
		dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			fmt.Printf("instance error: %v \n", err)
			return nil, err
		}

		fileSource, err := (&file.File{}).Open("file://db/migrations")
		if err != nil {
			fmt.Printf("opening file error: %v \n", err)
			return nil, err
		}

		m, err := migrate.NewWithInstance("file", fileSource, c.DBName, dbDriver)
		if err != nil {
			fmt.Printf("migrate error: %v \n", err)
			return nil, err
		}

		return m, nil
	case drivers.Postgres, drivers.MySQL, drivers.SQLServer:
		m, err := migrate.New("file://db/migrations", c.GetDBConnStr())
		if err != nil {
			fmt.Printf("migrate error: %v \n", err)
			return nil, err
		}
		return m, nil
	default:
		return nil, err
	}
}

// Get returns a config object with the db connection data already in place.
func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.User, "dbuser", os.Getenv("DB_USER"), "DB user name")
	flag.StringVar(&conf.Pswd, "dbpswd", os.Getenv("DB_PASSWORD"), "DB pass")
	flag.StringVar(&conf.Port, "dbport", os.Getenv("DB_PORT"), "DB port")
	flag.StringVar(&conf.Host, "dbhost", os.Getenv("DB_HOST"), "DB host")
	flag.StringVar(&conf.DBName, "dbname", os.Getenv("DB_NAME"), "DB name")
	flag.StringVar(&conf.Driver, "dbdriver", os.Getenv("DB_DRIVER"), "DB driver")

	return conf
}

// GetDBConnStr returns the connection string.
func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.Host, c.DBName)
}

// GetSQLXDBConnStr returns the connection string.
func (c *Config) GetSQLXDBConnStr() string {
	return c.getSQLXConnStr(c.Host, c.DBName)
}

// getDBConnStr returns the connection string based on the provided host and db name.
func (c *Config) getDBConnStr(dbhost, dbname string) string {
	switch c.Driver {
	case drivers.Postgres:
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			c.Driver,
			c.User,
			c.Pswd,
			dbhost,
			c.Port,
			dbname,
		)
	case drivers.MySQL:
		return fmt.Sprintf(
			"%s://%s:%s@tcp(%s:%s)/%s",
			c.Driver,
			c.User,
			c.Pswd,
			dbhost,
			c.Port,
			dbname,
		)
	case drivers.SQLite:
		return c.DBName
	case drivers.SQLServer:
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s?database=%s",
			c.Driver,
			c.User,
			c.Pswd,
			dbhost,
			c.Port,
			dbname,
		)
	default:
		return ""
	}
}

// getSQLXConnStr returns the connection string based on the provided host and db name.
func (c *Config) getSQLXConnStr(dbhost, dbname string) string {
	switch c.Driver {
	case drivers.Postgres:
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			c.Driver,
			c.User,
			c.Pswd,
			dbhost,
			c.Port,
			dbname,
		)
	case drivers.MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pswd, dbhost, c.Port, dbname)
	case drivers.SQLite:
		return c.DBName
	case drivers.SQLServer:
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s?database=%s",
			c.Driver,
			c.User,
			c.Pswd,
			dbhost,
			c.Port,
			dbname,
		)
	default:
		return ""
	}
}
