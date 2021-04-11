package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/assert"
)

var (
	driver   string
	user     string
	password string
	host     string
	port     string
	name     string
)

func TestMain(m *testing.M) {
	driver = os.Getenv("DB_DRIVER")
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	name = os.Getenv("DB_NAME")

	os.Exit(m.Run())
}

func generateURL(driver string) string {
	switch driver {
	case "postgres":
		return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=mode", driver, user, password, host, port, name)
	case "mysql":
		return fmt.Sprintf("%s://%s:%s@tcp(%s:%s)/%s?sslmode=mode", driver, user, password, host, port, name)
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
	}

	c, err := New(opts)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, c)
}
