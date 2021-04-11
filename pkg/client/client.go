package client

import (
	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	// postgres driver.
	_ "github.com/lib/pq"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/jmoiron/sqlx"
)

// Client is used to store the pool of db connection.
type Client struct {
	db *sqlx.DB
}

// New return an instance of the client.
func New(opts command.Options) (*Client, error) {
	db, err := sqlx.Open(opts.Driver, opts.URL)
	if err != nil {
		return nil, err
	}

	c := Client{
		db: db,
	}

	return &c, nil
}
