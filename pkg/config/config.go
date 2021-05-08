package config

import (
	"flag"
	"fmt"
	"os"
)

// Config struct is used to store the db connection data.
type Config struct {
	dbUser     string
	dbPswd     string
	dbHost     string
	dbPort     string
	dbName     string
	dbDriver   string
	testDBHost string
	testDBName string
	apiPort    string
	migrate    string
}

// Get returns a config object with the db connection data already in place.
func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("DB_USER"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("DB_PASSWORD"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("DB_PORT"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("DB_HOST"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("DB_NAME"), "DB name")
	flag.StringVar(&conf.dbDriver, "dbdriver", os.Getenv("DB_DRIVER"), "DB driver")
	flag.StringVar(&conf.testDBHost, "testdbhost", os.Getenv("TEST_DB_HOST"), "test database host")
	flag.StringVar(&conf.testDBName, "testdbname", os.Getenv("TEST_DB_NAME"), "test database name")
	flag.StringVar(&conf.apiPort, "apiPort", os.Getenv("API_PORT"), "API Port")
	flag.StringVar(&conf.migrate, "migrate", "up", "specify if we should be migrating DB 'up' or 'down'")

	flag.Parse()

	return conf
}

// GetDBConnStr returns the connection string.
func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

// GetTestDBConnStr returns the test connection string.
func (c *Config) GetTestDBConnStr() string {
	return c.getDBConnStr(c.testDBHost, c.testDBName)
}

// getDBConnStr returns the connection string based on the provied host and db name.
func (c *Config) getDBConnStr(dbhost, dbname string) string {
	switch c.dbDriver {
	case "postgres":
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", c.dbDriver, c.dbUser, c.dbPswd, dbhost, c.dbPort, dbname)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true", c.dbUser, c.dbPswd, dbhost, c.dbPort, dbname)
	default:
		return ""
	}
}

// GetMigration return up or down string to instruct the program if it should migrate database up or down.
func (c *Config) GetMigration() string {
	return c.migrate
}

// Driver returns the db driver from config.
func (c *Config) Driver() string {
	return c.dbDriver
}
