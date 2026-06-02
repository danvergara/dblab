package client

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	_ "modernc.org/sqlite"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/danvergara/dblab/pkg/pagination"
)

type TableRef struct {
	Schema string
	Name   string
}

type DBNode struct {
	ID         string
	Name       string
	Type       string
	ParentID   string
	ParentName string
	Children   []*DBNode
}

// databaseQuerier is an interface that indicates the methods
// a given type has to implement to interact with a database,
// to get specific data.
// This allows us to decouple the client from the database implementation and
// make adding new databases easier.
type databaseQuerier interface {
	TableStructure(table TableRef) (string, []any, error)
	Constraints(table TableRef) (string, []any, error)
	Indexes(table TableRef) (string, []any, error)
	Catalog(context.Context) (*DBNode, error)
}

// Client is used to store the pool of db connection.
type Client struct {
	db                *sqlx.DB
	dbName            string
	databaseQuerier   databaseQuerier
	driver, schema    string
	host              string
	paginationManager *pagination.Manager
	limit             uint
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
		dbName: opts.DBName,
		host:   opts.Host,
		driver: opts.Driver,
		limit:  opts.Limit,
	}

	if opts.Schema != "" {
		c.schema = opts.Schema
	}

	// This is where an implementation of databaseQuerier is getting picked up.
	switch c.driver {
	case drivers.Postgres, drivers.PostgreSQL, drivers.PostgresSSH:
		c.databaseQuerier = newPostgres(c.dbName, c.schema, c.db)
	case drivers.MySQL:
		c.databaseQuerier = newMySQL(c.dbName, c.db)
	case drivers.SQLite:
		c.databaseQuerier = newSQLite(c.dbName, c.db)
	case drivers.Oracle:
		c.databaseQuerier = newOracle(c.dbName, c.schema, c.db)
	case drivers.SQLServer:
		c.databaseQuerier = newMSSQL(c.dbName, c.schema, c.db)
	default:
		return nil, fmt.Errorf("%s driver not supported", c.driver)
	}

	switch c.driver {
	case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH:
		if _, err = db.Exec(fmt.Sprintf("set search_path = '%s'", c.schema)); err != nil {
			return nil, err
		}
	case drivers.Oracle:
		if c.schema != "" {
			if _, err = db.Exec(fmt.Sprintf("ALTER SESSION SET CURRENT_SCHEMA = %s", c.schema)); err != nil {
				return nil, err
			}
		}
	}

	pm, err := pagination.New(c.limit, 0, "")
	if err != nil {
		return nil, err
	}

	c.paginationManager = pm

	return c, nil
}

// DB Return the db attribute.
func (c *Client) DB() *sqlx.DB {
	return c.db
}

// Driver returns the driver of the database.
func (c *Client) Driver() string {
	return c.driver
}

func (c *Client) Host() string {
	return c.host
}

// Query returns performs the query and returns the result set and the column names.
func (c *Client) Query(q string, args ...any) ([][]string, []string, error) {
	var (
		resultSet = [][]string{}
	)

	// Runs the query extracting the content of the view calling the Buffer method.
	rows, err := c.db.Queryx(q, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Gets the names of the columns of the result set.
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		// cols is an []any of all of the column results.
		cols, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}

		// Convert []any into []string.
		s := make([]string, len(cols))
		for i, v := range cols {
			switch v.(type) {
			case string, rune, []byte:
				s[i] = fmt.Sprintf("%s", v)
			case nil:
				s[i] = fmt.Sprint(v)
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
func (c *Client) Metadata(table TableRef) (*Metadata, error) {
	tcRows, tcColumns, err := c.tableContent(table)
	if err != nil {
		return nil, err
	}

	sRows, sColumns, err := c.tableStructure(table)
	if err != nil {
		return nil, err
	}

	cRows, cColumns, err := c.constraints(table)
	if err != nil {
		return nil, err
	}

	iRows, iColumns, err := c.indexes(table)
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

// TableContent returns all the rows of a table.
func (c *Client) tableContent(table TableRef) ([][]string, []string, error) {
	var query string

	switch c.driver {
	case drivers.Postgres, drivers.PostgreSQL, drivers.PostgresSSH:
		query = fmt.Sprintf(
			"SELECT * FROM %s.%s LIMIT %d OFFSET %d;",
			table.Schema,
			table.Name,
			c.paginationManager.Limit(),
			c.paginationManager.Offset(),
		)
	case drivers.Oracle:
		query = fmt.Sprintf(
			"SELECT * FROM %s.%s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			strings.ToUpper(table.Schema),
			strings.ToUpper(table.Name),
			c.paginationManager.Offset(),
			c.paginationManager.Limit(),
		)
	case drivers.SQLServer:
		query = fmt.Sprintf(
			"SELECT * FROM %s ORDER BY (SELECT NULL) OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			table.Name,
			c.paginationManager.Offset(),
			c.paginationManager.Limit(),
		)
	default:
		query = fmt.Sprintf(
			"SELECT * FROM %s LIMIT %d OFFSET %d;",
			table.Name,
			c.paginationManager.Limit(),
			c.paginationManager.Offset(),
		)
	}

	return c.Query(query)
}

// tableStructure returns the structure of the table columns.
func (c *Client) tableStructure(table TableRef) ([][]string, []string, error) {
	var (
		query string
		err   error
		args  []any
	)

	query, args, err = c.databaseQuerier.TableStructure(table)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(query, args...)
}

// constraints returns the resultet of from information_schema.table_constraints.
func (c *Client) constraints(table TableRef) ([][]string, []string, error) {
	sql, args, err := c.databaseQuerier.Constraints(table)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(sql, args...)
}

// indexes returns a resulset with the information of the indexes given a table name.
func (c *Client) indexes(table TableRef) ([][]string, []string, error) {
	query, args, err := c.databaseQuerier.Indexes(table)
	if err != nil {
		return nil, nil, err
	}

	return c.Query(query, args...)
}

func (c *Client) Catalog(ctx context.Context) (*DBNode, error) {
	return c.databaseQuerier.Catalog(ctx)
}
