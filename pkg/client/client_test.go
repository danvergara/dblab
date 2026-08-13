package client

import (
	"context"
	"encoding/json"
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

// jsonViewPayload is the document used across the JSON view tests. It covers every
// token type the highlighter styles: keys, strings, numbers, booleans and null.
const jsonViewPayload = `{"name":"dblab","tags":["sql","tui"],"count":42,"ok":true,"extra":null}`

// jsonExpr returns a driver specific expression yielding a single JSON typed column,
// so the JSON view can be exercised without a JSON column in the schema.
func (suite *ClientTestSuite) jsonExpr(payload string) string {
	switch suite.driver {
	case drivers.Postgres:
		return fmt.Sprintf("'%s'::jsonb", payload)
	case drivers.MySQL:
		return fmt.Sprintf("CAST('%s' AS JSON)", payload)
	default:
		// SQLite is dynamically typed, so a plain string literal is enough.
		return fmt.Sprintf("'%s'", payload)
	}
}

// actorTable returns the actor table reference for the active driver.
func (suite *ClientTestSuite) actorTable() string {
	if suite.driver == drivers.Postgres {
		return "public.actor"
	}
	return "actor"
}

// connOptions returns the connection options used by the JSON view tests.
func (suite *ClientTestSuite) connOptions() command.Options {
	return command.Options{
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
}

// runAsyncQueries drains the result channel into a map keyed by query index.
func (suite *ClientTestSuite) runAsyncQueries(c *Client, queries ...string) map[int]QueryResult {
	resultChan := c.AsyncQuery(context.Background(), queries, 5)

	resultsByIndex := make(map[int]QueryResult)
	for r := range resultChan {
		resultsByIndex[r.QueryIndex] = r
	}

	return resultsByIndex
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

func (suite *ClientTestSuite) TestAsyncQueryJSONView() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	query := fmt.Sprintf("SELECT %s AS doc | json", suite.jsonExpr(jsonViewPayload))

	results := suite.runAsyncQueries(c, query)
	suite.Require().Len(results, 1)

	r := results[0]
	suite.Require().NoError(r.Error)
	suite.Equal(JSONQuery, r.QueryType)
	suite.NotEmpty(r.JSONData)

	// The engines normalize JSON (jsonb reorders keys, MySQL compacts whitespace),
	// so compare the decoded document rather than the raw bytes.
	var got, want map[string]any
	suite.Require().NoError(json.Unmarshal(r.JSONData, &got))
	suite.Require().NoError(json.Unmarshal([]byte(jsonViewPayload), &want))
	suite.Equal(want, got)

	// The suffix is stripped before execution but kept on the result, since
	// QueryResult.Query is what feeds the query history.
	suite.Equal(query, r.Query)

	// The table oriented fields stay empty on a JSON view.
	suite.Empty(r.ResultSet)
	suite.Empty(r.Headers)
}

func (suite *ClientTestSuite) TestAsyncQueryJSONViewSuffixVariants() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	expr := suite.jsonExpr(jsonViewPayload)

	for _, suffix := range []string{"| json", "| JSON", "| Json", "| json   "} {
		suite.Run(suffix, func() {
			results := suite.runAsyncQueries(c, fmt.Sprintf("SELECT %s AS doc %s", expr, suffix))
			suite.Require().Len(results, 1)

			r := results[0]
			suite.Require().NoError(r.Error)
			suite.Equal(JSONQuery, r.QueryType)
			suite.NotEmpty(r.JSONData)
		})
	}
}

func (suite *ClientTestSuite) TestAsyncQueryJSONViewMultipleColumns() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	query := fmt.Sprintf("SELECT %s AS doc, 1 AS n | json", suite.jsonExpr(jsonViewPayload))

	results := suite.runAsyncQueries(c, query)
	suite.Require().Len(results, 1)

	r := results[0]
	suite.Require().Error(r.Error)
	suite.Equal(JSONQuery, r.QueryType)
	suite.Contains(r.Error.Error(), "requires exactly 1 column")
	suite.Empty(r.JSONData)
}

func (suite *ClientTestSuite) TestAsyncQueryJSONViewNonJSONColumn() {
	if suite.driver == drivers.SQLite {
		suite.T().Skip("SQLite is dynamically typed, every column type is accepted")
	}

	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	results := suite.runAsyncQueries(c, "SELECT 1 AS n | json")
	suite.Require().Len(results, 1)

	r := results[0]
	suite.Require().Error(r.Error)
	suite.Equal(JSONQuery, r.QueryType)
	suite.Contains(r.Error.Error(), "requires a JSON column")
	suite.Empty(r.JSONData)
}

func (suite *ClientTestSuite) TestAsyncQueryJSONViewNoRows() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	query := fmt.Sprintf(
		"SELECT %s AS doc FROM %s WHERE 1 = 0 | json",
		suite.jsonExpr(jsonViewPayload),
		suite.actorTable(),
	)

	results := suite.runAsyncQueries(c, query)
	suite.Require().Len(results, 1)

	r := results[0]
	suite.Require().Error(r.Error)
	suite.Equal(JSONQuery, r.QueryType)
	suite.Contains(r.Error.Error(), "no data returned")
	suite.Empty(r.JSONData)
}

func (suite *ClientTestSuite) TestAsyncQueryWithoutJSONSuffix() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	query := fmt.Sprintf("SELECT %s AS doc", suite.jsonExpr(jsonViewPayload))

	results := suite.runAsyncQueries(c, query)
	suite.Require().Len(results, 1)

	r := results[0]
	suite.Require().NoError(r.Error)
	suite.Equal(NormalQuery, r.QueryType)
	suite.Nil(r.JSONData)
	suite.Len(r.Headers, 1)
	suite.Len(r.ResultSet, 1)
}

func (suite *ClientTestSuite) TestAsyncQueryJSONViewMixedBatch() {
	c, err := New(suite.connOptions())
	suite.Require().NoError(err)

	jsonQuery := fmt.Sprintf("SELECT %s AS doc | json", suite.jsonExpr(jsonViewPayload))
	actorQuery := fmt.Sprintf("SELECT * FROM %s;", suite.actorTable())

	results := suite.runAsyncQueries(c, jsonQuery, actorQuery)
	suite.Require().Len(results, 2)

	jsonResult := results[0]
	suite.NoError(jsonResult.Error)
	suite.Equal(JSONQuery, jsonResult.QueryType)
	suite.NotEmpty(jsonResult.JSONData)

	actorResult := results[1]
	suite.NoError(actorResult.Error)
	suite.Equal(NormalQuery, actorResult.QueryType)
	suite.Nil(actorResult.JSONData)
	suite.Len(actorResult.Headers, 4)
	suite.Len(actorResult.ResultSet, 200)
}

func TestClietnTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
