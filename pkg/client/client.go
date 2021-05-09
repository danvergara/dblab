package client

import (
	"fmt"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	// postgres driver.
	_ "github.com/lib/pq"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/connection"
	"github.com/jmoiron/sqlx"
)

// Client is used to store the pool of db connection.
type Client struct {
	db     *sqlx.DB
	driver string
}

// New return an instance of the client.
func New(opts command.Options) (*Client, error) {
	conn, err := connection.BuildConnectionFromOpts(opts)
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
	}

	return &c, nil
}

// Query returns performs the query and returns the result set and the colum names.
func (c *Client) Query(q string) ([][]string, []string, error) {
	resultSet := [][]string{}

	// Runs the query extracting the content of the view calling the Buffer method.
	rows, err := c.db.Queryx(q)
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

// TableContent returns all the rows of a table.
func (c *Client) TableContent(tableName string) ([][]string, []string, error) {
	query := fmt.Sprintf("SELECT * FROM %s;", tableName)

	return c.Query(query)
}

// DB Return the db attribute.
func (c *Client) DB() *sqlx.DB {
	return c.db
}

// Driver returns the driver of the database.
func (c *Client) Driver() string {
	return c.driver
}
