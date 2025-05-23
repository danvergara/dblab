package client

import (
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb/v2"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	_ "modernc.org/sqlite"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/danvergara/dblab/pkg/pagination"
)

// databaseQuerier is an interface that indicates the methods
// a given type has to implement to interact with a database,
// to get specific data.
// This allows us to decouple the client from the database implementation and
// make adding new databases easier.
type databaseQuerier interface {
	ShowTables() (string, []interface{}, error)
	TableStructure(tableName string) (string, []interface{}, error)
	Constraints(tableName string) (string, []interface{}, error)
	Indexes(tableName string) (string, []interface{}, error)
	ShowDatabases() (string, []interface{}, error)
	ShowTablesPerDB(database string) (string, []interface{}, error)
}

// Client is used to store the pool of db connection.
type Client struct {
	db                *sqlx.DB
	databaseQuerier   databaseQuerier
	driver, schema    string
	paginationManager *pagination.Manager
	activeDatabase    string
	limit             uint
	showDataCatalog   bool
	dbs               map[string]*sqlx.DB
}

// New return an instance of the client.
func New(opts command.Options) (*Client, error) {
	conn, opts, err := connection.BuildConnectionFromOpts(opts)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(opts.Driver, conn)
	if err != nil {
		return nil, err
	}

	c := &Client{
		db:     db,
		driver: opts.Driver,
		limit:  opts.Limit,
		dbs:    make(map[string]*sqlx.DB),
	}

	if opts.Schema == "" {
		c.schema = "public"
	} else {
		c.schema = opts.Schema
	}

	// This is where an implementation of databaseQuerier is getting picked up.
	switch c.driver {
	case drivers.Postgres, drivers.PostgreSQL, drivers.PostgresSSH:
		c.databaseQuerier = newPostgres(c.schema)
	case drivers.MySQL:
		c.databaseQuerier = newMySQL()
	case drivers.SQLite:
		c.databaseQuerier = newSQLite()
	case drivers.Oracle:
		c.databaseQuerier = newOracle()
	case drivers.SQLServer:
		c.databaseQuerier = newMSSQL()
	case drivers.DuckDB:
		c.databaseQuerier = newDuckDB()
	default:
		return nil, fmt.Errorf("%s driver not supported", c.driver)
	}

	if opts.DBName == "" {
		switch c.driver {
		case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH, drivers.MySQL:
			c.showDataCatalog = true
			dbs, err := c.ShowDatabases()
			if err != nil {
				return nil, err
			}

			for _, d := range dbs {
				db, err := getDB(c.driver, conn, d)
				if err != nil {
					continue
				}

				c.dbs[d] = db
			}
		}
	}

	switch c.driver {
	case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH:
		if _, err = db.Exec(fmt.Sprintf("set search_path='%s'", c.schema)); err != nil {
			return nil, err
		}
	}

	pm, err := pagination.New(c.limit, 0, "")
	if err != nil {
		return nil, err
	}

	c.paginationManager = pm

	return c, nil
}

func (c *Client) SetActiveDatabase(database string) {
	c.activeDatabase = database
}

func (c *Client) ActiveDatabase() string {
	return c.activeDatabase
}

// DB Return the db attribute.
func (c *Client) DB() *sqlx.DB {
	return c.db
}

// Driver returns the driver of the database.
func (c *Client) Driver() string {
	return c.driver
}

func (c *Client) ShowDataCatalog() bool {
	return c.showDataCatalog
}

// Query returns performs the query and returns the result set and the column names.
func (c *Client) Query(q string, args ...interface{}) ([][]string, []string, error) {
	var (
		resultSet = [][]string{}
		db        *sqlx.DB
		ok        bool
	)

	db = c.db

	if c.activeDatabase != "" {
		switch c.driver {
		case drivers.Postgres, drivers.PostgreSQL, drivers.PostgresSSH, drivers.MySQL:
			db, ok = c.dbs[c.activeDatabase]
			if !ok {
				return nil, nil, fmt.Errorf(
					"connection with %s database not found",
					c.activeDatabase,
				)
			}
		}
	}

	// Runs the query extracting the content of the view calling the Buffer method.
	rows, err := db.Queryx(q, args...)
	if err != nil {
		return nil, nil, err
	}

	// Gets the names of the columns of the result set.
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		// cols is an []interface{} of all of the column results.
		cols, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}

		// Convert []interface{} into []string.
		s := make([]string, len(cols))
		for i, v := range cols {
			switch v.(type) {
			case string, []byte:
				s[i] = fmt.Sprintf("%s", v)
			case nil:
				s[i] = fmt.Sprint(v)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				s[i] = fmt.Sprintf("%d", v)
			case float32, float64:
				s[i] = fmt.Sprintf("%f", v)
			default:
				s[i] = fmt.Sprintf("%v", v)
			}
		}

		resultSet = append(resultSet, s)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return resultSet, columnNames, nil
}

// Table represents a SQL table.
type Table struct {
	name    string
	Rows    [][]string
	Columns []string
}

func (t *Table) Name() string {
	return t.name
}

// Metadata sums up the most relevant data from a table.
type Metadata struct {
	TableContent Table
	Structure    Table
	Constraints  Table
	Indexes      Table
	TotalPages   int
}

// Metadata returns the most relevant data from a given table.
func (c *Client) Metadata(tableName string) (*Metadata, error) {
	tcRows, tcColumns, err := c.tableContent(tableName)
	if err != nil {
		return nil, err
	}

	sRows, sColumns, err := c.tableStructure(tableName)
	if err != nil {
		return nil, err
	}

	cRows, cColumns, err := c.constraints(tableName)
	if err != nil {
		return nil, err
	}

	iRows, iColumns, err := c.indexes(tableName)
	if err != nil {
		return nil, err
	}

	m := Metadata{
		TableContent: Table{
			Rows:    tcRows,
			Columns: tcColumns,
		},
		Structure: Table{
			Rows:    sRows,
			Columns: sColumns,
		},
		Constraints: Table{
			Rows:    cRows,
			Columns: cColumns,
		},
		Indexes: Table{
			Rows:    iRows,
			Columns: iColumns,
		},
	}

	return &m, nil
}

func (c *Client) ShowTablesPerDB(database string) ([]string, error) {
	var (
		query string
		err   error
		args  []interface{}
		db    *sqlx.DB
		ok    bool
	)

	tables := make([]string, 0)

	query, args, err = c.databaseQuerier.ShowTablesPerDB(database)
	if err != nil {
		return nil, err
	}

	switch c.driver {
	case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH, drivers.MySQL:
		db, ok = c.dbs[database]
		if !ok {
			return nil, fmt.Errorf("connection with %s database not found", database)
		}
	default:
		db = c.db
	}

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// ShowTables list all the tables in the database on the tables panel.
func (c *Client) ShowTables() ([]string, error) {
	var (
		query string
		err   error
		args  []interface{}
	)

	tables := make([]string, 0)

	query, args, err = c.databaseQuerier.ShowTables()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// ShowDatabases returns a list of the databases the user has access to.
func (c *Client) ShowDatabases() ([]string, error) {
	var (
		query string
		err   error
		args  []interface{}
	)

	databases := make([]string, 0)

	query, args, err = c.databaseQuerier.ShowDatabases()
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}

		databases = append(databases, database)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

// TableContent returns all the rows of a table.
func (c *Client) tableContent(tableName string) ([][]string, []string, error) {
	var query string

	switch c.driver {
	case drivers.Postgres, drivers.PostgreSQL, drivers.PostgresSSH:
		query = fmt.Sprintf(
			"SELECT * FROM %q LIMIT %d OFFSET %d;",
			tableName,
			c.paginationManager.Limit(),
			c.paginationManager.Offset(),
		)
	case drivers.Oracle:
		query = fmt.Sprintf(
			"SELECT * FROM %s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			tableName,
			c.paginationManager.Offset(),
			c.paginationManager.Limit(),
		)
	case drivers.SQLServer:
		query = fmt.Sprintf(
			"SELECT * FROM %s ORDER BY (SELECT NULL) OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			tableName,
			c.paginationManager.Offset(),
			c.paginationManager.Limit(),
		)
	default:
		query = fmt.Sprintf(
			"SELECT * FROM %s LIMIT %d OFFSET %d;",
			tableName,
			c.paginationManager.Limit(),
			c.paginationManager.Offset(),
		)
	}

	return c.Query(query)
}

// tableStructure returns the structure of the table columns.
func (c *Client) tableStructure(tableName string) ([][]string, []string, error) {
	var (
		query string
		err   error
		args  []interface{}
	)

	query, args, err = c.databaseQuerier.TableStructure(tableName)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(query, args...)
}

// constraints returns the resultet of from information_schema.table_constraints.
func (c *Client) constraints(tableName string) ([][]string, []string, error) {
	sql, args, err := c.databaseQuerier.Constraints(tableName)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(sql, args...)
}

// indexes returns a resulset with the information of the indexes given a table name.
func (c *Client) indexes(tableName string) ([][]string, []string, error) {
	query, args, err := c.databaseQuerier.Indexes(tableName)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(query, args...)
}

func getDB(driver, connString, database string) (*sqlx.DB, error) {
	var newConnString string

	switch driver {
	case drivers.MySQL:
		newConnString = strings.Replace(connString, "/", fmt.Sprintf("/%s", database), 1)
	default:
		u, err := url.Parse(connString)
		if err != nil {
			return nil, err
		}

		u.Path = "/" + database
		newConnString = u.String()
	}

	db, err := sqlx.Open(driver, newConnString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
