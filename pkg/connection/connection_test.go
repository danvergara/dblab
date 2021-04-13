package connection

import (
	"errors"
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestBuildConnectionFromOpts(t *testing.T) {
	type given struct {
		opts command.Options
	}
	type want struct {
		uri      string
		hasError bool
		err      error
	}
	var cases = []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "valid postgres localhost",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					URL:    "postgres://user:password@localhost:5432/db?sslmode=disable",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslmode=disable",
			},
		},
		{
			name: "valid postgres localhost but add sslmode",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					URL:    "postgres://user:password@localhost:5432/db",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslmode=disable",
			},
		},
		{
			name: "valid postgres localhost postgresql as protocol",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					URL:    "postgresql://user:password@localhost:5432/db",
				},
			},
			want: want{
				uri: "postgresql://user:password@localhost:5432/db?sslmode=disable",
			},
		},
		{
			name: "error misspelled postgres",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					URL:    "potgre://user:password@localhost:5432/db",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidUPostgresRLFormat,
			},
		},
		{
			name: "error invalid url",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					URL:    "not-a-url",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidUPostgresRLFormat,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uri, err := BuildConnectionFromOpts(tc.given.opts)

			if tc.want.hasError {
				assert.Error(t, err)

				if !errors.Is(err, tc.want.err) {
					t.Errorf("got %v, wanted %v", err, tc.want.err)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.uri, uri)
		})
	}
}

func TestFormatPostgresURL(t *testing.T) {
	type given struct {
		opts command.Options
	}
	type want struct {
		uri      string
		hasError bool
		err      error
	}
	var cases = []struct {
		name  string
		given given
		want  want
	}{}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uri, err := formatPostgresURL(tc.given.opts)

			if tc.want.hasError {
				assert.Error(t, err)

				if !errors.Is(err, tc.want.err) {
					t.Errorf("got %v, wanted %v", err, tc.want.err)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.uri, uri)
		})
	}
}

func TestFormatMySQLURL(t *testing.T) {
	type given struct {
		opts command.Options
	}
	type want struct {
		uri      string
		hasError bool
		err      error
	}
	var cases = []struct {
		name  string
		given given
		want  want
	}{}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uri, err := formatMySQLURL(tc.given.opts)

			if tc.want.hasError {
				assert.Error(t, err)

				if !errors.Is(err, tc.want.err) {
					t.Errorf("got %v, wanted %v", err, tc.want.err)
				}

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.uri, uri)
		})
	}
}
