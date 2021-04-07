package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetEnvValues(t *testing.T) {
	if err := os.Setenv("DATABASE_URL", "postgres://user:password@host:port/database?sslmode=disable"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DB_USER", "user"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DB_PASS", "password"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DB_NAME", "database"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		t.Fatal(err)
	}

	opts := Options{}

	result := SetDefault(opts)

	assert.Equal(t, "postgres://user:password@host:port/database?sslmode=disable", result.URL)
	assert.Equal(t, "localhost", result.Host)
	assert.Equal(t, "user", result.User)
	assert.Equal(t, "password", result.Pass)
	assert.Equal(t, "database", result.DBName)
	assert.Equal(t, "5432", result.Port)
}
