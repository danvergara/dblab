package connection

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/assert"
)

func createTempSocketFile(t *testing.T) *os.File {
	t.Helper()

	socketFile, err := os.CreateTemp("", "mysqld*.sock")
	assert.NoError(t, err)

	return socketFile
}

func TestBuildConnectionFromOptsFromURL(t *testing.T) {
	socketFile := createTempSocketFile(t)

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
					URL: "postgres://user:password@localhost:5432/db?sslmode=disable",
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
					URL: "postgres://user:password@localhost:5432/db",
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
					URL: "postgresql://user:password@localhost:5432/db",
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
					URL: "potgre://user:password@localhost:5432/db",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidURLFormat,
			},
		},
		// mysql
		{
			name: "valid mysql localhost",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@tcp(localhost:3306)/db",
				},
			},
			want: want{
				uri: "user:password@tcp(localhost:3306)/db",
			},
		},
		{
			name: "valid mysql localhost with params",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				uri: "user:password@tcp(localhost:3306)/db?charset=utf8",
			},
		},
		{
			name: "valid mysql remote url",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
				},
			},
			want: want{
				uri: "user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
			},
		},
		{
			name: "valid socket connection",
			given: given{
				opts: command.Options{
					URL: fmt.Sprintf("mysql://user@unix(%s)/dbname?charset=utf8", socketFile.Name()),
				},
			},
			want: want{
				uri: fmt.Sprintf("user@unix(%s)/dbname?charset=utf8", socketFile.Name()),
			},
		},
		{
			name: "error misspelled mysql",
			given: given{
				opts: command.Options{
					URL: "mysq://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidURLFormat,
			},
		},
		{
			name: "error invalid url",
			given: given{
				opts: command.Options{
					URL: "not-a-url",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidURLFormat,
			},
		},
		{
			name: "sqlite3 file dsn example",
			given: given{
				opts: command.Options{
					URL: "file:test.db?cache=shared&mode=memory",
				},
			},
			want: want{
				uri: "file:test.db?cache=shared&mode=memory",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uri, _, err := BuildConnectionFromOpts(tc.given.opts)

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
	socketFile := createTempSocketFile(t)

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
				uri: "user:password@tcp(localhost:3306)/db",
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
				uri: "user:password@tcp(127.0.0.1:3306)/db",
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
				uri: "user:password@tcp(your-amazonaws-uri.com:3306)/db",
			},
		},
		{
			name: "success - sockets connection",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					DBName: "db",
					Socket: socketFile.Name(),
				},
			},
			want: want{
				uri: fmt.Sprintf("user@unix(%s)/db?charset=utf8", socketFile.Name()),
			},
		},
		{
			name: "error - invalid socket file name",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					DBName: "db",
					Socket: "/path/to/not-wrong-file",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidSocketFile,
			},
		},
		{
			name: "error - socket file do not exist",
			given: given{
				opts: command.Options{
					Driver: "mysql",
					User:   "user",
					DBName: "db",
					Socket: "/path/to/not-existing-file.sock",
				},
			},
			want: want{
				hasError: true,
				err:      ErrSocketFileDoNotExist,
			},
		},
		// sqlite3.
		{
			name: "success - sqlite3 db file extension",
			given: given{
				opts: command.Options{
					Driver: "sqlite3",
					DBName: "users.db",
				},
			},
			want: want{
				uri: "users.db",
			},
		},
		{
			name: "success - sqlite3 sqlite3 file extension",
			given: given{
				opts: command.Options{
					Driver: "sqlite3",
					DBName: "users.sqlite3",
				},
			},
			want: want{
				uri: "users.sqlite3",
			},
		},
		{
			name: "error - wrong sqlite3 file extension",
			given: given{
				opts: command.Options{
					Driver: "sqlite3",
					DBName: "users.wrong",
				},
			},
			want: want{
				uri:      "users.wrong",
				hasError: true,
				err:      ErrInvalidSqlite3Extension,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uri, _, err := BuildConnectionFromOpts(tc.given.opts)

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
	}{
		{
			name: "valid postgres localhost",
			given: given{
				opts: command.Options{
					URL: "postgres://user:password@localhost:5432/db?sslmode=disable",
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
					URL: "postgres://user:password@localhost:5432/db",
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
					URL: "postgresql://user:password@localhost:5432/db",
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
					URL: "potgre://user:password@localhost:5432/db",
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
					URL: "not-a-url",
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
					URL: "mysql://user:password@tcp(localhost:3306)/db",
				},
			},
			want: want{
				uri: "user:password@tcp(localhost:3306)/db",
			},
		},
		{
			name: "valid mysql localhost with params",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@tcp(localhost:3306)/db?charset=utf8",
				},
			},
			want: want{
				uri: "user:password@tcp(localhost:3306)/db?charset=utf8",
			},
		},
		{
			name: "valid mysql remote url",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
				},
			},
			want: want{
				uri: "user:password@tcp(your-amazonaws-uri.com:3306)/dbname",
			},
		},
		{
			name: "valid mysql traditional",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@/dbname?sql_mode=TRADITIONAL",
				},
			},
			want: want{
				uri: "user:password@/dbname?sql_mode=TRADITIONAL",
			},
		},
		{
			name: "valid mysql google cloud sql on app engine",
			given: given{
				opts: command.Options{
					URL: "mysql://user:password@unix(/cloudsql/project-id:region-name:instance-name)/dbname",
				},
			},
			want: want{
				uri: "user:password@unix(/cloudsql/project-id:region-name:instance-name)/dbname",
			},
		},
		{
			name: "valid socket connection",
			given: given{
				opts: command.Options{
					URL: "mysql://user@unix(/path/to/socket)/dbname?charset=utf8",
				},
			},
			want: want{
				uri: "user@unix(/path/to/socket)/dbname?charset=utf8",
			},
		},
		{
			name: "error misspelled mysql",
			given: given{
				opts: command.Options{
					URL: "mysq://user:password@tcp(localhost:3306)/db?charset=utf8",
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
					URL: "not-a-url",
				},
			},
			want: want{
				hasError: true,
				err:      ErrInvalidMySQLURLFormat,
			},
		},
		{
			name: "valid with strange characters",
			given: given{
				opts: command.Options{
					URL: "mysql://myuser:5@klkbN#ABC@@tcp(localhost:3306)/mydb",
				},
			},
			want: want{
				uri: "myuser:5@klkbN#ABC@@tcp(localhost:3306)/mydb",
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

func TestParseDSN(t *testing.T) {
	type given struct {
		dsn string
	}
	type want struct {
		uri      string
		hasError bool
	}
	var cases = []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "valid mysql traditional",
			given: given{
				dsn: "mysql://user:password@/dbname?sql_mode=TRADITIONAL",
			},
			want: want{
				uri: "mysql://user:password@/dbname?sql_mode=TRADITIONAL",
			},
		},
		{
			name: "valid mysql google cloud sql on app engine",
			given: given{
				dsn: "mysql://user:password@unix(/cloudsql/project-id:region-name:instance-name)/dbname",
			},
			want: want{
				uri: "mysql://user:password@unix(/cloudsql/project-id:region-name:instance-name)/dbname",
			},
		},
		{
			name: "valid socket connection",
			given: given{
				dsn: "mysql://user@unix(/path/to/socket)/dbname?charset=utf8",
			},
			want: want{
				uri: "mysql://user@unix(/path/to/socket)/dbname?charset=utf8",
			},
		},
		{
			name: "error invalid url",
			given: given{
				dsn: "not-a-url",
			},
			want: want{
				hasError: true,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			url, err := parseDSN(tc.given.dsn)

			if tc.want.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want.uri, url)
		})
	}
}
