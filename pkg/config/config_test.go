package config_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"

	"github.com/danvergara/dblab/pkg/config"
)

func TestInit(t *testing.T) {
	type want struct {
		host            string
		port            string
		dbname          string
		user            string
		pass            string
		driver          string
		schema          string
		limit           uint
		ssl             string
		sslcert         string
		sslkey          string
		sslpassword     string
		sslrootcert     string
		traceFile       string
		sslVerify       string
		wallet          string
		sshHost         string
		sshPort         string
		sshUser         string
		sshPass         string
		runQuery        int16
		navigationUp    int16
		navigationDown  int16
		navigationRight int16
		navigationLeft  int16
		structure       int16
		indexes         int16
		constraints     int16
		clearEditor     int16
	}
	var tests = []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "empty config name",
			input: "",
			want: want{
				host:            "localhost",
				port:            "5432",
				dbname:          "users",
				user:            "postgres",
				pass:            "password",
				driver:          "postgres",
				schema:          "public",
				ssl:             "disable",
				limit:           50,
				runQuery:        int16(tcell.KeyCtrlSpace),
				navigationUp:    int16(tcell.KeyCtrlK),
				navigationDown:  int16(tcell.KeyCtrlJ),
				navigationLeft:  int16(tcell.KeyCtrlH),
				navigationRight: int16(tcell.KeyCtrlL),
				structure:       int16(tcell.KeyCtrlS),
				indexes:         int16(tcell.KeyCtrlI),
				constraints:     int16(tcell.KeyCtrlT),
				clearEditor:     int16(tcell.KeyCtrlD),
			},
		},
		{
			name:  "test config",
			input: "test",
			want: want{
				host:            "localhost",
				port:            "5432",
				dbname:          "users",
				user:            "postgres",
				pass:            "password",
				driver:          "postgres",
				schema:          "public",
				ssl:             "disable",
				limit:           50,
				runQuery:        int16(tcell.KeyCtrlSpace),
				navigationUp:    int16(tcell.KeyCtrlK),
				navigationDown:  int16(tcell.KeyCtrlJ),
				navigationLeft:  int16(tcell.KeyCtrlH),
				navigationRight: int16(tcell.KeyCtrlL),
				structure:       int16(tcell.KeyCtrlS),
				indexes:         int16(tcell.KeyCtrlI),
				constraints:     int16(tcell.KeyCtrlT),
				clearEditor:     int16(tcell.KeyCtrlD),
			},
		},
		{
			name:  "production config",
			input: "prod",
			want: want{
				host:            "mydb.123456789012.us-east-1.rds.amazonaws.com",
				port:            "5432",
				dbname:          "users",
				user:            "postgres",
				pass:            "password",
				driver:          "postgres",
				schema:          "public",
				ssl:             "require",
				sslrootcert:     "~/.postgresql/root.crt.",
				limit:           50,
				runQuery:        int16(tcell.KeyCtrlSpace),
				navigationUp:    int16(tcell.KeyCtrlK),
				navigationDown:  int16(tcell.KeyCtrlJ),
				navigationLeft:  int16(tcell.KeyCtrlH),
				navigationRight: int16(tcell.KeyCtrlL),
				structure:       int16(tcell.KeyCtrlS),
				indexes:         int16(tcell.KeyCtrlI),
				constraints:     int16(tcell.KeyCtrlT),
				clearEditor:     int16(tcell.KeyCtrlD),
			},
		},
		{
			name:  "ssh tunnel",
			input: "ssh-tunnel",
			want: want{
				host:            "localhost",
				port:            "5432",
				dbname:          "users",
				user:            "postgres",
				pass:            "password",
				driver:          "postgres",
				schema:          "public",
				ssl:             "disable",
				sshHost:         "example.com",
				sshPort:         "22",
				sshUser:         "ssh-user",
				sshPass:         "password",
				limit:           50,
				runQuery:        int16(tcell.KeyCtrlSpace),
				navigationUp:    int16(tcell.KeyCtrlK),
				navigationDown:  int16(tcell.KeyCtrlJ),
				navigationLeft:  int16(tcell.KeyCtrlH),
				navigationRight: int16(tcell.KeyCtrlL),
				structure:       int16(tcell.KeyCtrlS),
				indexes:         int16(tcell.KeyCtrlI),
				constraints:     int16(tcell.KeyCtrlT),
				clearEditor:     int16(tcell.KeyCtrlD),
			},
		},
		{
			name:  "oracle",
			input: "oracle",
			want: want{
				host:            "localhost",
				port:            "1521",
				dbname:          "FREEPDB1 ",
				user:            "system",
				pass:            "password",
				driver:          "oracle",
				ssl:             "enable",
				sslVerify:       "true",
				wallet:          "path/to/wallet",
				traceFile:       "trace.log",
				limit:           50,
				runQuery:        int16(tcell.KeyCtrlSpace),
				navigationUp:    int16(tcell.KeyCtrlK),
				navigationDown:  int16(tcell.KeyCtrlJ),
				navigationLeft:  int16(tcell.KeyCtrlH),
				navigationRight: int16(tcell.KeyCtrlL),
				structure:       int16(tcell.KeyCtrlS),
				indexes:         int16(tcell.KeyCtrlI),
				constraints:     int16(tcell.KeyCtrlT),
				clearEditor:     int16(tcell.KeyCtrlD),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := config.Init(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want.host, opts.Host)
			assert.Equal(t, tt.want.port, opts.Port)
			assert.Equal(t, tt.want.dbname, opts.DBName)
			assert.Equal(t, tt.want.user, opts.User)
			assert.Equal(t, tt.want.pass, opts.Pass)
			assert.Equal(t, tt.want.driver, opts.Driver)
			assert.Equal(t, tt.want.schema, opts.Schema)
			assert.Equal(t, tt.want.limit, opts.Limit)
			assert.Equal(t, tt.want.ssl, opts.SSL)
			assert.Equal(t, tt.want.sslcert, opts.SSLCert)
			assert.Equal(t, tt.want.sslkey, opts.SSLKey)
			assert.Equal(t, tt.want.sslpassword, opts.SSLPassword)
			assert.Equal(t, tt.want.sslrootcert, opts.SSLRootcert)
			assert.Equal(t, tcell.Key(tt.want.runQuery), opts.TUIKeyBindings.RunQuery)
			assert.Equal(t, tcell.Key(tt.want.navigationUp), opts.TUIKeyBindings.Navigation.Up)
			assert.Equal(t, tcell.Key(tt.want.navigationDown), opts.TUIKeyBindings.Navigation.Down)
			assert.Equal(
				t,
				tcell.Key(tt.want.navigationRight),
				opts.TUIKeyBindings.Navigation.Right,
			)
			assert.Equal(t, tcell.Key(tt.want.navigationLeft), opts.TUIKeyBindings.Navigation.Left)
			assert.Equal(
				t,
				tcell.Key(tt.want.constraints),
				opts.TUIKeyBindings.Constraints,
			)
			assert.Equal(
				t,
				tcell.Key(tt.want.indexes),
				opts.TUIKeyBindings.Indexes,
			)
			assert.Equal(
				t,
				tcell.Key(tt.want.structure),
				opts.TUIKeyBindings.Structure,
			)
			assert.Equal(
				t,
				tcell.Key(tt.want.clearEditor),
				opts.TUIKeyBindings.ClearEditor,
			)
		})
	}
}
