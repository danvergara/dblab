package config_test

import (
	"testing"

	"github.com/danvergara/dblab/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	opts, err := config.Init()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", opts.Host)
	assert.Equal(t, "5432", opts.Port)
	assert.Equal(t, "users", opts.DBName)
	assert.Equal(t, "postgres", opts.User)
	assert.Equal(t, "password", opts.Pass)
	assert.Equal(t, "postgres", opts.Driver)
	assert.Equal(t, 50, opts.Limit)
}
