package client

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	mysqltest "github.com/testcontainers/testcontainers-go/modules/mysql"
	postgrestest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "modernc.org/sqlite"

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

	switch suite.driver {
	case drivers.Postgres:
		pgContainer, err := postgrestest.Run(suite.ctx,
			"danvergara/sakila:postgres-18.4",
			postgrestest.WithDatabase(suite.dbName),
			postgrestest.WithUsername(suite.user),
			postgrestest.WithPassword(suite.password),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(1).WithStartupTimeout(30*time.Second)),
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

	case drivers.MySQL:
		mysqlContainer, err := mysqltest.Run(suite.ctx,
			"danvergara/sakila:mysql-8.4",
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
	}

}

func (suite *ClientTestSuite) TearDownSuite() {
	if suite.container != nil {
		if err := suite.container.Terminate(suite.ctx); err != nil {
			suite.T().Fatalf("failed to terminate container: %s", err)
		}
	}
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

	query := "SELECT * FROM public.actor;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM actor;"
	}

	r, co, err := c.Query(query)
	suite.Len(r, 200)
	suite.Len(co, 4)
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
	_, _, err := c.Query(`INSERT INTO public.actor(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
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

	tableRef := TableRef{Name: "actor", Schema: "public"}
	r, co, err := c.tableContent(tableRef)

	suite.Len(r, int(opts.Limit))
	suite.Len(co, 4)
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

	tableRef := TableRef{Name: "actor", Schema: "public"}
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

	tableRef := TableRef{Name: "actor", Schema: "public"}
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

	tableRef := TableRef{Name: "actor", Schema: "public"}
	m, err := c.Metadata(tableRef)
	suite.NoError(err)
	suite.NotNil(m)

	// indexes.
	suite.Greater(len(m.Indexes.Rows), 0)
	suite.Greater(len(m.Indexes.Columns), 0)

	// constraints.
	suite.Greater(len(m.Constraints.Rows), 0)
	suite.Greater(len(m.Constraints.Columns), 0)

	switch suite.driver {
	case drivers.Postgres:
		suite.Len(m.Structure.Columns, 5)
		suite.Len(m.Structure.Rows, 4)
	case drivers.MySQL:
		suite.Len(m.Structure.Columns, 6)
		suite.Len(m.Structure.Rows, 4)
	default:
		suite.Len(m.Structure.Columns, 6)
	}

	// table content.
	suite.Len(m.TableContent.Rows, int(opts.Limit))
	suite.Len(m.TableContent.Columns, 4)
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

	query := "SELECT * FROM public.actor;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM actor;"
	}

	resultChan := c.AsyncQuery(context.Background(), []string{query}, 5)

	var results []QueryResult
	for r := range resultChan {
		results = append(results, r)
	}

	suite.Len(results, 1)
	suite.NoError(results[0].Error)
	suite.Equal(0, results[0].QueryIndex)
	suite.Len(results[0].Headers, 4)
	suite.Len(results[0].ResultSet, 200)
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

	productsQuery := "SELECT * FROM public.actor;"
	customersQuery := "SELECT * FROM public.film;"
	if suite.driver == "mysql" {
		productsQuery = "SELECT * FROM actor;"
		customersQuery = "SELECT * FROM film;"
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
	suite.Len(prodResult.Headers, 4)
	suite.Len(prodResult.ResultSet, 200)

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

	validQuery := "SELECT * FROM public.actor;"
	if suite.driver == "mysql" {
		validQuery = "SELECT * FROM actor;"
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
	suite.Len(resultsByIndex[0].Headers, 4)
	suite.Len(resultsByIndex[0].ResultSet, 200)

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

	query := "SELECT * FROM public.actor;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM actor;"
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

	query := "SELECT * FROM public.actor;"
	if suite.driver == "mysql" {
		query = "SELECT * FROM actor;"
	}

	queries := []string{query, query, query, query, query}
	resultChan := c.AsyncQuery(context.Background(), queries, 1)

	resultsByIndex := make(map[int]QueryResult)
	for r := range resultChan {
		resultsByIndex[r.QueryIndex] = r
	}

	suite.Len(resultsByIndex, 5)
	for i := range 5 {
		r, ok := resultsByIndex[i]
		suite.True(ok, "missing result for query index %d", i)
		suite.NoError(r.Error)
		suite.Len(r.Headers, 4)
		suite.Len(r.ResultSet, 200)
	}
}

func TestClietnTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
