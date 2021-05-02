package client

import (
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

// DB Return the db attribute.
func (c *Client) DB() *sqlx.DB {
	return c.db
}

// Driver returns the driver of the database.
func (c *Client) Driver() string {
	return c.driver
}
