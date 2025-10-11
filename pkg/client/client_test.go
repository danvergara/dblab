package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	mysqltest "github.com/testcontainers/testcontainers-go/modules/mysql"
	postgrestest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "modernc.org/sqlite"

	"github.com/danvergara/dblab/db/seeds"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
)

type ClientTestSuite struct {
	suite.Suite
	container testcontainers.Container
	ctx       context.Context
	driver    string
	user      string
	password  string
	dbName    string
	dbSchema  string
	host      string
	port      nat.Port
	db        *sqlx.DB
}

func (suite *ClientTestSuite) SetupSuite() {
	if testing.Short() {
		suite.T().Skip("skipping integration tests in short mode.")
	}

	suite.driver = os.Getenv("DB_DRIVER")
	suite.user = os.Getenv("DB_USER")
	suite.password = os.Getenv("DB_PASSWORD")
	suite.dbName = os.Getenv("DB_NAME")
	suite.dbSchema = os.Getenv("DB_SCHEMA")

	suite.ctx = context.Background()

	// var err error
	var dsn string

	switch suite.driver {
	case drivers.Postgres:
		pgContainer, err := postgrestest.Run(suite.ctx,
			"postgres:17-alpine",
			postgrestest.WithDatabase(suite.dbName),
			postgrestest.WithUsername(suite.user),
			postgrestest.WithPassword(suite.password),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).WithStartupTimeout(30*time.Second)),
		)
		require.NoError(suite.T(), err)
		suite.container = pgContainer

		suite.host, err = pgContainer.Host(suite.ctx)
		require.NoError(suite.T(), err)
		suite.port, err = pgContainer.MappedPort(suite.ctx, "5432")
		require.NoError(suite.T(), err)

		sqlxDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			suite.user, suite.password, suite.host, suite.port.Port(), suite.dbName)

		suite.db, err = sqlx.Connect("postgres", sqlxDSN)
		require.NoError(suite.T(), err)

		dsn = sqlxDSN
	case drivers.MySQL:
		mysqlContainer, err := mysqltest.Run(suite.ctx,
			"mysql:9",
			mysqltest.WithDatabase(suite.dbName),
			mysqltest.WithUsername(suite.user),
			mysqltest.WithPassword(suite.password),
			testcontainers.WithWaitStrategy(
				wait.ForLog("port: 3306  MySQL Community Server - GPL").
					WithOccurrence(1).WithStartupTimeout(20*time.Second)),
		)
		require.NoError(suite.T(), err)
		suite.container = mysqlContainer

		suite.host, err = mysqlContainer.Host(suite.ctx)
		require.NoError(suite.T(), err)
		suite.port, err = mysqlContainer.MappedPort(suite.ctx, "3306")
		require.NoError(suite.T(), err)

		sqlxDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			suite.user, suite.password, suite.host, suite.port.Port(), suite.dbName)

		suite.db, err = sqlx.Connect("mysql", sqlxDSN)
		require.NoError(suite.T(), err)

		dsn = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s",
			suite.user, suite.password, suite.host, suite.port.Port(), suite.dbName)
	}

	// Run migrations
	_ = suite.runMigrations(dsn)
	// require.NoError(suite.T(), err)

	// Run seeds
	suite.runSeeds()
}

func (suite *ClientTestSuite) TearDownSuite() {
	if suite.container != nil {
		if err := suite.container.Terminate(suite.ctx); err != nil {
			suite.T().Fatalf("failed to terminate container: %s", err)
		}
	}
}

func (suite *ClientTestSuite) runSeeds() {
	seeds.Execute(suite.db, suite.driver)
}

func (suite *ClientTestSuite) runMigrations(dsn string) error {
	// Get absolute path to migrations.
	migrationsPath, err := filepath.Abs("../../db/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Up()
}

func (suite *ClientTestSuite) generateURL() string {
	switch suite.driver {
	case drivers.Postgres:
		return fmt.Sprintf(
			"%s://%s:%s@%s:%s/%s?sslmode=disable",
			suite.driver,
			suite.user,
			suite.password,
			suite.host,
			suite.port.Port(),
			suite.dbName,
		)
	case drivers.MySQL:
		return fmt.Sprintf(
			"%s://%s:%s@tcp(%s:%s)/%s",
			suite.driver,
			suite.user,
			suite.password,
			suite.host,
			suite.port.Port(),
			suite.dbName,
		)
	case drivers.SQLite:
		return suite.dbName
	default:
		return ""
	}
}

func (suite *ClientTestSuite) TestNewClientByURL() {
	url := suite.generateURL()

	opts := command.Options{
		Driver: suite.driver,
		URL:    url,
		Limit:  50,
	}

	c, err := New(opts)
	suite.NoError(err)
	suite.NotNil(c)
}

func (suite *ClientTestSuite) TestNewClientByUserData() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		SSL:    "disable",
		Limit:  50,
		Schema: suite.dbSchema,
	}

	c, err := New(opts)
	suite.NoError(err)
	suite.NotNil(c)
}

func (suite *ClientTestSuite) TestNewClientPing() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  50,
	}

	c, err := New(opts)
	suite.NoError(err)
	suite.NotNil(c)
	err = c.DB().Ping()
	suite.NoError(err)
}

func (suite *ClientTestSuite) TestQuery() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.Query("SELECT * FROM products;")
	suite.Len(r, 100)
	suite.Len(co, 3)
	suite.NoError(err)
}

func (suite *ClientTestSuite) TestTableContent() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.tableContent("products")

	suite.Len(r, int(opts.Limit))
	suite.Len(co, 3)
	suite.NoError(err)
}

func (suite *ClientTestSuite) TestShowTables() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, err := New(opts)
	suite.NoError(err)

	tables, err := c.ShowTables()
	suite.NoError(err)
	suite.Len(tables, 4)
}

func (suite *ClientTestSuite) TestShowTablesPerDB() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	tables, err := c.ShowTablesPerDB(suite.dbName)
	suite.Len(tables, 4)
	suite.NoError(err)
}

func (suite *ClientTestSuite) TestTableStructure() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.tableStructure("products")
	suite.NoError(err)
	suite.Len(r, 3)

	switch suite.driver {
	case drivers.Postgres:
		suite.Len(co, 8)
	case drivers.MySQL:
		suite.Len(co, 6)
	default:
		suite.Len(co, 8)
	}
}

func (suite *ClientTestSuite) TestConstraints() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.constraints("products")

	suite.T().Logf("constraints columns %v", co)
	suite.T().Logf("constraints content %v", r)

	suite.NoError(err)
	suite.NotEmpty(r)
	suite.NotEmpty(co)
}

func (suite *ClientTestSuite) TestIndexes() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.indexes("products")
	suite.NoError(err)
	suite.NotEmpty(r)
	suite.NotEmpty(co)
}

func (suite *ClientTestSuite) TestMetadata() {
	opts := command.Options{
		Driver: suite.driver,
		User:   suite.user,
		Pass:   suite.password,
		Host:   suite.host,
		Port:   suite.port.Port(),
		DBName: suite.dbName,
		Schema: suite.dbSchema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	m, err := c.Metadata("products")
	suite.NoError(err)
	suite.NotNil(m)

	// indexes.
	suite.Greater(len(m.Indexes.Rows), 0)
	suite.Greater(len(m.Indexes.Columns), 0)

	// constraints.
	suite.Greater(len(m.Constraints.Rows), 0)
	suite.Greater(len(m.Constraints.Columns), 0)

	// structure.
	suite.Len(m.Structure.Rows, 3)

	switch suite.driver {
	case drivers.Postgres:
		suite.Len(m.Structure.Columns, 8)
	case drivers.MySQL:
		suite.Len(m.Structure.Columns, 6)
	default:
		suite.Len(m.Structure.Columns, 6)
	}

	// table content.
	suite.Len(m.TableContent.Rows, int(opts.Limit))
	suite.Len(m.TableContent.Columns, 3)
}

func TestCleitnTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
