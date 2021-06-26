package connection

import (
	"errors"
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/assert"
)

func TestBuildConnectionFromOptsFromURL(t *testing.T) {
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
		// postgres
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
				err:      ErrInvalidPostgresURLFormat,
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
				err:      ErrInvalidPostgresURLFormat,
			},
		},
		// mysql
		{
			name: "valid mysql localhost",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(localhost:3306)/db",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(localhost:3306)/db",
			},
		},
		{
			name: "valid mysql localhost with params",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
			},
		},
		{
			name: "valid mysql remote url",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
			},
		},
		{
			name: "error misspelled mysql",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysq://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidMySQLURLFormat,
			},
		},
		{
			name: "error invalid url",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "not-a-url",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidMySQLURLFormat,
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

func TestBuildConnectionFromOptsUserData(t *testing.T) {
	type given struct {
		opts command.Options
	}
	type want struct {
		uri string
	}
	var cases = []struct {
		name  string
		given given
		want  want
	}{
		// postgres
		{
			name: "success - localhost with no explicit ssl mode - postgres",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					User:   "user",
					Pass:   "password",
					Host:   "localhost",
					Port:   "5432",
					DBName: "db",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslmode=disable",
			},
		},
		{
			name: "success - 127.0.0.1 with no explicit ssl mode - postgres",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					User:   "user",
					Pass:   "password",
					Host:   "127.0.0.1",
					Port:   "5432",
					DBName: "db",
				},
			},
			want: want{
				uri: "postgres://user:password@127.0.0.1:5432/db?sslmode=disable",
			},
		},
		{
			name: "success - remote host - postgres",
			given: given{
				opts: command.Options{
					Driver: "postgres",
					User:   "user",
					Pass:   "password",
					Host:   "your-amazonaws-uri.com",
					Port:   "5432",
					DBName: "db",
				},
			},
			want: want{
				uri: "postgres://user:password@your-amazonaws-uri.com:5432/db",
			},
		},
		// mysql
		{
			name: "success - localhost - mysql",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					Pass:   "password",
					Host:   "localhost",
					Port:   "3306",
					DBName: "db",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(localhost:3306)/db",
			},
		},
		{
			name: "success - 127.0.0.1 - mysql",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					Pass:   "password",
					Host:   "127.0.0.1",
					Port:   "3306",
					DBName: "db",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(127.0.0.1:3306)/db",
			},
		},
		{
			name: "success - remote host -mysql",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					Pass:   "password",
					Host:   "your-amazonaws-uri.com",
					Port:   "3306",
					DBName: "db",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/db",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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
				err:      ErrInvalidPostgresURLFormat,
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
				err:      ErrInvalidPostgresURLFormat,
			},
		},
	}

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
	}{
		{
			name: "valid mysql localhost",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(localhost:3306)/db",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(localhost:3306)/db",
			},
		},
		{
			name: "valid mysql localhost with params",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
			},
		},
		{
			name: "valid mysql remote url",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
				},
			},
			want: want{
				uri: "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
			},
		},
		{
			name: "error misspelled mysql",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "mysq://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidMySQLURLFormat,
			},
		},
		{
			name: "error invalid url",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					URL:    "not-a-url",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidMySQLURLFormat,
			},
		},
	}

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
