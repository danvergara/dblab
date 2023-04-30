package client

import (
	"fmt"
	"os"
	"testing"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	// postgres driver.
	_ "github.com/lib/pq"
	// sqlite driver.
	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
)

var (
	driver   string
	user     string
	password string
	host     string
	port     string
	name     string
	schema   string
)

func TestMain(m *testing.M) {
	driver = os.Getenv("DB_DRIVER")
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	name = os.Getenv("DB_NAME")
	schema = os.Getenv("DB_SCHEMA")

	os.Exit(m.Run())
}

func generateURL(driver string) string {
	switch driver {
	case drivers.POSTGRES:
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", driver, user, password, host, port, name)
	case drivers.MYSQL:
		return fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s", driver, user, password, host, port, name)
	case drivers.SQLITE:
		return name
	default:
		return ""
	}
}

func TestNewClientByURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	url := generateURL(driver)

	opts := command.Options{
		Driver: driver,
		URL:    url,
		Limit:  50,
	}

	c, err := New(opts)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, c)
}

func TestNewClientByUserData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		SSL:    "disable",
		Limit:  50,
		Schema: schema,
	}

	c, err := New(opts)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, c)
}

func TestNewClientPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  50,
	}

	c, err := New(opts)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, c)

	if err := c.DB().Ping(); err != nil {
		t.Error(err)
	}
}

func TestQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.Query("SELECT * FROM products;")

	assert.Len(t, r, 100)
	assert.Len(t, co, 3)
	assert.NoError(t, err)
}

func TestTableContent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.tableContent("products")

	assert.Len(t, r, int(opts.Limit))
	assert.Len(t, co, 3)
	assert.NoError(t, err)
}

func TestShowTables(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	tables, err := c.ShowTables()

	assert.Len(t, tables, 4)
	assert.NoError(t, err)
}

func TestTableStructure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.tableStructure("products")

	assert.Len(t, r, 3)
	assert.Len(t, co, 8)

	assert.NoError(t, err)
}

func TestConstraints(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.constraints("products")

	t.Logf("constraints columns %v", co)
	t.Logf("constraints content %v", r)

	assert.NoError(t, err)
	assert.Greater(t, len(r), 0)
	assert.Greater(t, len(co), 0)
}

func TestIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	r, co, err := c.indexes("products")

	assert.NoError(t, err)
	assert.Greater(t, len(r), 0)
	assert.Greater(t, len(co), 0)
}

func TestCountTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	count, err := c.tableCount("products")

	assert.Equal(t, count, 100)
	assert.NoError(t, err)
}

func TestMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	m, err := c.Metadata("products")

	assert.NoError(t, err)

	// Total count.
	assert.Equal(t, m.TotalPages, 1)

	// indexes.
	assert.Greater(t, len(m.Indexes.Rows), 0)
	assert.Greater(t, len(m.Indexes.Columns), 0)

	// constraints.
	assert.Greater(t, len(m.Constraints.Rows), 0)
	assert.Greater(t, len(m.Constraints.Columns), 0)

	// structure.
	assert.Len(t, m.Structure.Rows, 3)
	assert.Len(t, m.Structure.Columns, 8)

	// table content.
	assert.Len(t, m.TableContent.Rows, int(opts.Limit))
	assert.Len(t, m.TableContent.Columns, 3)
}

func TestTotalPages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	opts := command.Options{
		Driver: driver,
		User:   user,
		Pass:   password,
		Host:   host,
		Port:   port,
		DBName: name,
		Schema: schema,
		SSL:    "disable",
		Limit:  100,
	}

	c, _ := New(opts)

	_, err := c.Metadata("products")

	assert.NoError(t, err)

	// Total count.
	assert.Equal(t, c.TotalPages(), 1)
}
