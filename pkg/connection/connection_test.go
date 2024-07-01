package connection

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
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
		// sql server.
		{
			name: "valid sql server url",
			given: given{
				opts: command.Options{
					URL: "sqlserver://SA:myStrong(!)Password@localhost:1433?database=tempdb&encrypt=true&trustservercertificate=false&connection+timeout=30",
				},
			},
			want: want{
				uri: "sqlserver://SA:myStrong%28%21%29Password@localhost:1433?connection+timeout=30&database=tempdb&encrypt=true&trustservercertificate=false",
			},
		},
		// oracle.
		{
			name: "valid oracle url",
			given: given{
				opts: command.Options{
					URL: "oracle://user:pass@server:1521/service_name",
				},
			},
			want: want{
				uri: "oracle://user:pass@server:1521/service_name",
			},
		},
		// postgres
		{
			name: "valid socket connection with postgres with no password or user",
			given: given{
				opts: command.Options{
					URL: "postgres:///db?host=/path/to/socket",
				},
			},
			want: want{
				uri: fmt.Sprintf(
					"postgres:///db?host=%s",
					url.QueryEscape("/path/to/socket"),
				),
			},
		},
		{
			name: "valid socket connection with postgres",
			given: given{
				opts: command.Options{
					URL: "postgres://user:password@/db?host=/path/to/socket",
				},
			},
			want: want{
				uri: fmt.Sprintf(
					"postgres://user:password@/db?host=%s",
					url.QueryEscape("/path/to/socket"),
				),
			},
		},
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
		{
			name: "valid postgres url with sslmode equals to require",
			given: given{
				opts: command.Options{
					// keep url params in alphetic orden
					// see: https://github.com/golang/go/issues/29985
					URL: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require&sslrootcert=server-ca.crt",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require&sslrootcert=server-ca.crt",
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
					URL: fmt.Sprintf(
						"mysql://user:password@unix(%s)/dbname?charset=utf8",
						socketFile.Name(),
					),
				},
			},
			want: want{
				uri: fmt.Sprintf("user:password@unix(%s)/dbname?charset=utf8", socketFile.Name()),
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
			name: "sqlite file dsn example",
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
		// sql server.
		{
			name: "success - sql server",
			given: given{
				opts: command.Options{
					Driver:                 drivers.SQLServer,
					User:                   "SA",
					Pass:                   "password",
					Host:                   "localhost",
					Port:                   "1433",
					DBName:                 "tempdb",
					Encrypt:                "true",
					TrustServerCertificate: "false",
					ConnectionTimeout:      "30",
				},
			},
			want: want{
				uri: "sqlserver://SA:password@localhost:1433?connection%2Btimeout=30&database=tempdb&encrypt=true&trustservercertificate=false",
			},
		},
		// oracle.
		{
			name: "success - oracle",
			given: given{
				opts: command.Options{
					Driver: drivers.Oracle,
					User:   "user",
					Pass:   "password",
					Host:   "localhost",
					Port:   "1521",
					DBName: "db",
				},
			},
			want: want{
				uri: "oracle://user:password@localhost:1521/db?",
			},
		},
		{
			name: "success - ssl enable - oracle",
			given: given{
				opts: command.Options{
					Driver: drivers.Oracle,
					User:   "user",
					Pass:   "password",
					Host:   "localhost",
					Port:   "1521",
					DBName: "db",
					SSL:    "enable",
				},
			},
			want: want{
				uri: "oracle://user:password@localhost:1521/db?SSL=enable",
			},
		},
		// postgres
		{
			name: "success - postgres socket connection without password",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
					User:   "user",
					DBName: "db",
					Socket: "/path/to/socket",
				},
			},
			want: want{
				uri: fmt.Sprintf(
					"postgres://user@/db?host=%s",
					url.QueryEscape("/path/to/socket"),
				),
			},
		},
		{
			name: "success - postgres socket connection",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
					User:   "user",
					Pass:   "password",
					DBName: "db",
					Socket: "/path/to/socket",
				},
			},
			want: want{
				uri: fmt.Sprintf(
					"postgres://user:password@/db?host=%s",
					url.QueryEscape("/path/to/socket"),
				),
			},
		},
		{
			name: "success - localhost with no explicit ssl mode - postgres",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
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
					Driver: drivers.Postgres,
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
			name: "success  ssl mode - require",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
					User:   "user",
					Pass:   "password",
					// fake instance.
					Host:        "db-postgresql-nyc1-12345-do-user-123456-0.b.db.ondigitalocean.com",
					Port:        "5432",
					DBName:      "db",
					SSL:         "require",
					SSLCert:     "client-cert.pem",
					SSLKey:      "client-key.pem",
					SSLRootcert: "server-ca.crt",
				},
			},
			want: want{
				uri: "postgres://user:password@db-postgresql-nyc1-12345-do-user-123456-0.b.db.ondigitalocean.com:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require&sslrootcert=server-ca.crt",
			},
		},
		{
			name: "success  ssl mode - require",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
					User:   "user",
					Pass:   "password",
					// fake instance.
					Host:        "db-postgresql-nyc1-12345-do-user-123456-0.b.db.ondigitalocean.com",
					Port:        "5432",
					DBName:      "db",
					SSL:         "require",
					SSLRootcert: "server-ca.crt",
				},
			},
			want: want{
				uri: "postgres://user:password@db-postgresql-nyc1-12345-do-user-123456-0.b.db.ondigitalocean.com:5432/db?sslmode=require&sslrootcert=server-ca.crt",
			},
		},
		{
			name: "success - remote host - postgres",
			given: given{
				opts: command.Options{
					Driver: drivers.Postgres,
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
					Driver: drivers.MySQL,
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
					Driver: drivers.MySQL,
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
					Driver: drivers.MySQL,
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
					Driver: drivers.MySQL,
					User:   "user",
					Pass:   "password",
					DBName: "db",
					Socket: socketFile.Name(),
				},
			},
			want: want{
				uri: fmt.Sprintf("user:password@unix(%s)/db?charset=utf8", socketFile.Name()),
			},
		},
		{
			name: "error - invalid socket file name",
			given: given{
				opts: command.Options{
					Driver: drivers.MySQL,
					User:   "user",
					Pass:   "password",
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
					Driver: drivers.MySQL,
					User:   "user",
					Pass:   "password",
					DBName: "db",
					Socket: "/path/to/not-existing-file.sock",
				},
			},
			want: want{
				hasError: true,
				err:      ErrSocketFileDoNotExist,
			},
		},
		// sqlite.
		{
			name: "success - sqlite db file extension",
			given: given{
				opts: command.Options{
					Driver: drivers.SQLite,
					User:   "user",
					DBName: "users.db",
				},
			},
			want: want{
				uri: "users.db",
			},
		},
		{
			name: "success - sqlite sqlite3 file extension",
			given: given{
				opts: command.Options{
					Driver: drivers.SQLite,
					DBName: "users.sqlite3",
				},
			},
			want: want{
				uri: "users.sqlite3",
			},
		},
		{
			name: "success - valid sqlite3 file extension",
			given: given{
				opts: command.Options{
					Driver: drivers.SQLite,
					DBName: "users.rsd",
				},
			},
			want: want{
				uri: "users.rsd",
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
			name: "valid postgres url ssl mode require all params",
			given: given{
				opts: command.Options{
					URL: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require&sslrootcert=server-ca.pem",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require&sslrootcert=server-ca.pem",
			},
		},
		{
			name: "valid postgres url ssl mode require missing params",
			given: given{
				opts: command.Options{
					URL: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require",
				},
			},
			want: want{
				uri: "postgres://user:password@localhost:5432/db?sslcert=client-cert.pem&sslkey=client-key.pem&sslmode=require",
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
					URL: "mysql://user:password@unix(/path/to/socket)/dbname?charset=utf8",
				},
			},
			want: want{
				uri: "user:password@unix(/path/to/socket)/dbname?charset=utf8",
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
				dsn: "mysql://user:password@unix(/path/to/socket)/dbname?charset=utf8",
			},
			want: want{
				uri: "mysql://user:password@unix(/path/to/socket)/dbname?charset=utf8",
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
