package client

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
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

	query := "SELECT * FROM public.products;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM products;"
	}

	r, co, err := c.Query(query)
	suite.Len(r, 100)
	suite.Len(co, 3)
	suite.NoError(err)
}

func (suite *ClientTestSuite) TestReadOnly() {
	opts := command.Options{
		Driver:   suite.driver,
		User:     suite.user,
		Pass:     suite.password,
		Host:     suite.host,
		Port:     suite.port.Port(),
		DBName:   suite.dbName,
		Schema:   suite.dbSchema,
		SSL:      "disable",
		Limit:    100,
		ReadOnly: true,
	}
	c, _ := New(opts)
	_, _, err := c.Query(`INSERT INTO public.products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
	suite.Error(err)
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

	tableRef := TableRef{Name: "products", Schema: "public"}
	r, co, err := c.tableContent(tableRef)

	suite.Len(r, int(opts.Limit))
	suite.Len(co, 3)
	suite.NoError(err)
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

	tableRef := TableRef{Name: "products", Schema: "public"}
	r, co, err := c.constraints(tableRef)

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

	tableRef := TableRef{Name: "products", Schema: "public"}
	r, co, err := c.indexes(tableRef)
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

	tableRef := TableRef{Name: "products", Schema: "public"}
	m, err := c.Metadata(tableRef)
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

func (suite *ClientTestSuite) TestAsyncQuerySingleQuery() {
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
	suite.Require().NoError(err)

	query := "SELECT * FROM public.products;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM products;"
	}

	resultChan := c.AsyncQuery(context.Background(), []string{query}, 5)

	var results []QueryResult
	for r := range resultChan {
		results = append(results, r)
	}

	suite.Len(results, 1)
	suite.NoError(results[0].Error)
	suite.Equal(0, results[0].QueryIndex)
	suite.Len(results[0].Headers, 3)
	suite.Len(results[0].ResultSet, 100)
}

func (suite *ClientTestSuite) TestAsyncQueryMultipleQueries() {
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
	suite.Require().NoError(err)

	productsQuery := "SELECT * FROM public.products;"
	customersQuery := "SELECT * FROM public.customers;"
	if suite.driver == "mysql" {
		productsQuery = "SELECT * FROM products;"
		customersQuery = "SELECT * FROM customers;"
	}

	queries := []string{productsQuery, customersQuery}
	resultChan := c.AsyncQuery(context.Background(), queries, 5)

	resultsByIndex := make(map[int]QueryResult)
	for r := range resultChan {
		resultsByIndex[r.QueryIndex] = r
	}

	suite.Len(resultsByIndex, 2)

	prodResult := resultsByIndex[0]
	suite.NoError(prodResult.Error)
	suite.Len(prodResult.Headers, 3)
	suite.Len(prodResult.ResultSet, 100)

	custResult := resultsByIndex[1]
	suite.NoError(custResult.Error)
	suite.NotEmpty(custResult.Headers)
	suite.NotEmpty(custResult.ResultSet)
}

func (suite *ClientTestSuite) TestAsyncQueryWithInvalidQuery() {
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
	suite.Require().NoError(err)

	validQuery := "SELECT * FROM public.products;"
	if suite.driver == "mysql" {
		validQuery = "SELECT * FROM products;"
	}
	invalidQuery := "SELECT * FROM nonexistent_table_xyz;"

	queries := []string{validQuery, invalidQuery}
	resultChan := c.AsyncQuery(context.Background(), queries, 5)

	resultsByIndex := make(map[int]QueryResult)
	for r := range resultChan {
		resultsByIndex[r.QueryIndex] = r
	}

	suite.Len(resultsByIndex, 2)

	suite.NoError(resultsByIndex[0].Error)
	suite.Len(resultsByIndex[0].Headers, 3)
	suite.Len(resultsByIndex[0].ResultSet, 100)

	suite.Error(resultsByIndex[1].Error)
}

func (suite *ClientTestSuite) TestAsyncQueryContextCancellation() {
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
	suite.Require().NoError(err)

	query := "SELECT * FROM public.products;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM products;"
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	queries := []string{query, query, query}
	resultChan := c.AsyncQuery(ctx, queries, 5)

	var results []QueryResult
	for r := range resultChan {
		results = append(results, r)
	}

	suite.Len(results, len(queries))

	for _, r := range results {
		suite.Error(r.Error)
	}
}

func (suite *ClientTestSuite) TestAsyncQueryConcurrencyLimit() {
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
	suite.Require().NoError(err)

	query := "SELECT * FROM public.products;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM products;"
	}

	queries := []string{query, query, query, query, query}
	resultChan := c.AsyncQuery(context.Background(), queries, 1)

	resultsByIndex := make(map[int]QueryResult)
	for r := range resultChan {
		resultsByIndex[r.QueryIndex] = r
	}

	suite.Len(resultsByIndex, 5)
	for i := 0; i < 5; i++ {
		r, ok := resultsByIndex[i]
		suite.True(ok, "missing result for query index %d", i)
		suite.NoError(r.Error)
		suite.Len(r.Headers, 3)
		suite.Len(r.ResultSet, 100)
	}
}

func TestClietnTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
