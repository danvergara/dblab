package client

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/danvergara/dblab/pkg/pagination"
	"github.com/jmoiron/sqlx"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	// postgres driver.
	_ "github.com/lib/pq"
	// sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

// Client is used to store the pool of db connection.
type Client struct {
	db                *sqlx.DB
	driver            string
	paginationManager *pagination.Manager
	limit             int
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

	c := Client{
		db:     db,
		driver: opts.Driver,
		limit:  opts.Limit,
	}

	pm, err := pagination.New(c.limit, 0, "")
	if err != nil {
		return nil, err
	}

	c.paginationManager = pm

	return &c, nil
}

// Query returns performs the query and returns the result set and the column names.
func (c *Client) Query(q string, args ...interface{}) ([][]string, []string, error) {
	resultSet := [][]string{}

	// Runs the query extracting the content of the view calling the Buffer method.
	rows, err := c.db.Queryx(q, args...)
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

// Metadata retunrs the most relevant data from a given table.
func (c *Client) Metadata(tableName string) (*Metadata, error) {
	count, err := c.tableCount(tableName)
	if err != nil {
		return nil, err
	}

	pm, err := pagination.New(c.limit, count, tableName)
	if err != nil {
		return nil, err
	}

	c.paginationManager = pm

	pages := c.paginationManager.TotalPages()

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
		TotalPages: pages,
	}

	return &m, nil
}

func (c *Client) TotalPages() int {
	if c.paginationManager != nil {
		return c.paginationManager.TotalPages()
	}

	return 0
}

// ShowTables list all the tables in the database on the tables panel.
func (c *Client) ShowTables() ([]string, error) {
	var query string
	tables := make([]string, 0)

	switch c.driver {
	case "postgres":
		fallthrough
	case "postgresql":
		query = `
		SELECT
			table_name
		FROM
			information_schema.tables
		WHERE
			table_schema = 'public'
		ORDER BY
			table_name;`
	case "mysql":
		query = "SHOW TABLES;"
	case "sqlite3":
		query = `
		SELECT
			name
		FROM
			sqlite_schema
		WHERE
			type ='table' AND
			name NOT LIKE 'sqlite_%';`
	}

	rows, err := c.db.Queryx(query)
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

	return tables, nil
}

// NextPage returns the next page of the given table, based off the limit and the offsite.
func (c *Client) NextPage() (*Table, int, error) {
	if err := c.paginationManager.NextPage(); err != nil {
		return nil, 0, err
	}

	r, col, err := c.tableContent(c.paginationManager.CurrentTable())
	if err != nil {
		return nil, 0, err
	}

	t := Table{
		name:    c.paginationManager.CurrentTable(),
		Rows:    r,
		Columns: col,
	}

	page := c.paginationManager.CurrentPage()

	return &t, page, nil
}

// PreviousPage returns the next page of the given table, based off the limit and the offsite.
func (c *Client) PreviousPage() (*Table, int, error) {
	if err := c.paginationManager.PreviousPage(); err != nil {
		return nil, 0, err
	}

	r, col, err := c.tableContent(c.paginationManager.CurrentTable())
	if err != nil {
		return nil, 0, err
	}

	t := Table{
		name:    c.paginationManager.CurrentTable(),
		Rows:    r,
		Columns: col,
	}

	page := c.paginationManager.CurrentPage()

	return &t, page, nil
}

// ResetPagination resets the paginationManager field.
func (c *Client) ResetPagination() error {
	pm, err := pagination.New(c.limit, 0, "")
	if err != nil {
		return err
	}

	c.paginationManager = pm
	return nil
}

// DB Return the db attribute.
func (c *Client) DB() *sqlx.DB {
	return c.db
}

// Driver returns the driver of the database.
func (c *Client) Driver() string {
	return c.driver
}

// TableContent returns all the rows of a table.
func (c *Client) tableContent(tableName string) ([][]string, []string, error) {
	var query string

	if c.driver == "postgres" || c.driver == "postgresql" {
		query = fmt.Sprintf("SELECT * FROM public.%s LIMIT %d OFFSET %d;", tableName, c.paginationManager.Limit(), c.paginationManager.Offset())
	} else {
		query = fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d;", tableName, c.paginationManager.Limit(), c.paginationManager.Offset())
	}

	return c.Query(query)
}

// tableCount returns the count of a given table.
func (c *Client) tableCount(tableName string) (int, error) {
	var (
		query string
		count int
	)

	if c.driver == "postgres" || c.driver == "postgresql" {
		query = fmt.Sprintf("SELECT COUNT(*) FROM public.%s;", tableName)
	} else {
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s;", tableName)
	}

	if err := c.db.Get(&count, query); err != nil {
		return 0, err
	}

	return count, nil
}

// tableStructure returns the structure of the table columns.
func (c *Client) tableStructure(tableName string) ([][]string, []string, error) {
	var query string

	switch c.driver {
	case "postgres":
		fallthrough
	case "postgresql":
		query = `
        SELECT
			c.column_name,
			c.is_nullable,
			c.data_type,
			c.character_maximum_length,
			c.numeric_precision,
			c.numeric_scale,
			c.ordinal_position,
			tc.constraint_type pkey
		FROM
			information_schema.columns c
		LEFT JOIN
			information_schema.constraint_column_usage AS ccu
		ON
			c.table_schema = ccu.table_schema
			AND c.table_name = ccu.table_name
			AND c.column_name = ccu.column_name
		LEFT JOIN
			information_schema.table_constraints AS tc
		ON
			ccu.constraint_schema = tc.constraint_schema
			and ccu.constraint_name = tc.constraint_name
		WHERE
			c.table_schema = 'public'
			AND c.table_name = $1;`
		return c.Query(query, tableName)
	case "mysql":
		query = fmt.Sprintf("DESCRIBE %s;", tableName)
		return c.Query(query)
	case "sqlite3":
		query = fmt.Sprintf("PRAGMA table_info(%s);", tableName)
		return c.Query(query)
	default:
		return nil, nil, errors.New("not supported driver")
	}
}

// constraints returns the resultet of from information_schema.table_constraints.
func (c *Client) constraints(tableName string) ([][]string, []string, error) {
	var (
		query sq.SelectBuilder
		sql   string
	)

	query = sq.Select(
		`tc.constraint_name`,
		`tc.table_name`,
		`tc.constraint_type`,
	).
		From("information_schema.table_constraints AS tc").
		Where("tc.table_name = ?")

	switch c.driver {
	case "sqlite3":
		sql = `
		SELECT *
		FROM
			sqlite_master
		WHERE
			type='table' AND name = ?;`
		return c.Query(sql, tableName)
	case "postgres":
		fallthrough
	case "postgresql":
		query = query.Where("tc.table_schema = 'public'")
		query = query.PlaceholderFormat(sq.Dollar)
	}

	sql, _, err := query.ToSql()
	if err != nil {
		return nil, nil, err
	}

	return c.Query(sql, tableName)
}

// indexes returns a resulset with the information of the indexes given a table name.
func (c *Client) indexes(tableName string) ([][]string, []string, error) {
	var query string

	switch c.driver {
	case "postgres":
		fallthrough
	case "postgresql":
		query = "SELECT * FROM pg_indexes WHERE tablename = $1;"
		return c.Query(query, tableName)
	case "mysql":
		query = fmt.Sprintf("SHOW INDEX FROM %s", tableName)
		return c.Query(query)
	case "sqlite3":
		query = `PRAGMA index_list(%s);`
		query = fmt.Sprintf(query, tableName)
		return c.Query(query)
	default:
		return nil, nil, errors.New("not supported driver")
	}
}
