package config

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/kkyr/fig"
	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"

	// drivers.
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/mattn/go-sqlite3"
)

// Config struct is used to store the db connection data.
type Config struct {
	Database struct {
		Host     string
		Port     string
		DB       string `validate:"required"`
		User     string
		Password string
		Driver   string `validate:"required"`
		SSL      string `default:"disable"`
	}
	User   string
	Pswd   string
	Host   string
	Port   string
	DBName string
	Driver string
	Limit  int `fig:"limit" default:"100"`
}

// New returns a config instance the with db connection data inplace based on the flags of a cobra command.
func New(cmd *cobra.Command) *Config {
	conf := &Config{}

	cmd.PersistentFlags().StringVarP(&conf.User, "user", "", os.Getenv("DB_USER"), "DB user name")
	cmd.PersistentFlags().StringVarP(&conf.Pswd, "pswd", "", os.Getenv("DB_PASSWORD"), "DB pass")
	cmd.PersistentFlags().StringVarP(&conf.Port, "port", "", os.Getenv("DB_PORT"), "DB port")
	cmd.PersistentFlags().StringVarP(&conf.Host, "host", "", os.Getenv("DB_HOST"), "DB host")
	cmd.PersistentFlags().StringVarP(&conf.DBName, "name", "", os.Getenv("DB_NAME"), "DB name")
	cmd.PersistentFlags().StringVarP(&conf.Driver, "driver", "", os.Getenv("DB_DRIVER"), "DB driver")

	return conf
}

// Init reads in config file and returns a commands/Options instance.
func Init() (command.Options, error) {
	var opts command.Options
	var cfg Config

	home, err := os.UserHomeDir()
	if err != nil {
		return opts, err
	}

	if err := fig.Load(&cfg, fig.File(".dblab.yaml"), fig.Dirs(".", home)); err != nil {
		return opts, err
	}

	opts = command.Options{
		Driver: cfg.Database.Driver,
		Host:   cfg.Database.Host,
		Port:   cfg.Database.Port,
		User:   cfg.Database.User,
		Pass:   cfg.Database.Password,
		DBName: cfg.Database.DB,
		SSL:    cfg.Database.SSL,
		Limit:  cfg.Limit,
	}

	return opts, nil
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
	case "sqlite3":
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
	case "postgres", "mysql":
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

// getDBConnStr returns the connection string based on the provied host and db name.
func (c *Config) getDBConnStr(dbhost, dbname string) string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", c.Driver, c.User, c.Pswd, dbhost, c.Port, dbname)
	case "mysql":
		return fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", c.Driver, c.User, c.Pswd, dbhost, c.Port, dbname)
	case "sqlite3":
		return c.DBName
	default:
		return ""
	}
}

// getSQLXConnStr returns the connection string based on the provied host and db name.
func (c *Config) getSQLXConnStr(dbhost, dbname string) string {
	switch c.Driver {
	case "postgres":
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", c.Driver, c.User, c.Pswd, dbhost, c.Port, dbname)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pswd, dbhost, c.Port, dbname)
	case "sqlite3":
		return c.DBName
	default:
		return ""
	}
}
