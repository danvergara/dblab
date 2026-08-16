package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danvergara/dblab/pkg/config"
)

func TestInit(t *testing.T) {
	type want struct {
		host        string
		port        string
		dbname      string
		user        string
		pass        string
		driver      string
		schema      string
		limit       uint
		ssl         string
		sslcert     string
		sslkey      string
		sslpassword string
		sslrootcert string
		traceFile   string
		sslVerify   string
		wallet      string
		sshHost     string
		sshPort     string
		sshUser     string
		sshPass     string
		sshKeyFile  string
		sshKeyPass  string
		readOnly    bool
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
				host:     "localhost",
				port:     "5432",
				dbname:   "users",
				user:     "postgres",
				pass:     "password",
				driver:   "postgres",
				schema:   "public",
				ssl:      "disable",
				limit:    50,
				readOnly: true,
			},
		},
		{
			name:  "test config",
			input: "test",
			want: want{
				host:     "localhost",
				port:     "5432",
				dbname:   "users",
				user:     "postgres",
				pass:     "password",
				driver:   "postgres",
				schema:   "public",
				ssl:      "disable",
				limit:    50,
				readOnly: true,
			},
		},
		{
			name:  "production config",
			input: "prod",
			want: want{
				host:        "mydb.123456789012.us-east-1.rds.amazonaws.com",
				port:        "5432",
				dbname:      "users",
				user:        "postgres",
				pass:        "password",
				driver:      "postgres",
				schema:      "public",
				ssl:         "require",
				sslrootcert: "~/.postgresql/root.crt.",
				limit:       50,
			},
		},
		{
			name:  "ssh tunnel",
			input: "ssh-tunnel",
			want: want{
				host:    "localhost",
				port:    "5432",
				dbname:  "users",
				user:    "postgres",
				pass:    "password",
				driver:  "postgres",
				schema:  "public",
				ssl:     "disable",
				sshHost: "example.com",
				sshPort: "22",
				sshUser: "ssh-user",
				sshPass: "password",
				limit:   50,
			},
		},
		{
			name:  "realistic example",
			input: "realistic-ssh-example",
			want: want{
				host:       "rds-endpoint.region.rds.amazonaws.com",
				port:       "5432",
				dbname:     "database_name",
				user:       "db_user",
				pass:       "password",
				driver:     "postgres",
				schema:     "schema_name",
				ssl:        "require",
				sshHost:    "bastion.host.ip",
				sshPort:    "22",
				sshUser:    "ec2-user",
				sshKeyFile: "/path/to/ssh/key.pem",
				sshKeyPass: "hiuwiewnc092",
				limit:      50,
			},
		},
		{
			name:  "oracle",
			input: "oracle",
			want: want{
				host:      "localhost",
				port:      "1521",
				dbname:    "FREEPDB1 ",
				user:      "system",
				pass:      "password",
				driver:    "oracle",
				ssl:       "enable",
				sslVerify: "true",
				wallet:    "path/to/wallet",
				traceFile: "trace.log",
				limit:     50,
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

			// SSH validations.
			assert.Equal(t, tt.want.sshHost, opts.SSHHost)
			assert.Equal(t, tt.want.sshPort, opts.SSHPort)
			assert.Equal(t, tt.want.sshUser, opts.SSHUser)
			assert.Equal(t, tt.want.sshPass, opts.SSHPass)
			assert.Equal(t, tt.want.sshPass, opts.SSHPass)
			assert.Equal(t, tt.want.sshKeyFile, opts.SSHKeyFile)
			// Read Only mode.
			assert.Equal(t, tt.want.readOnly, opts.ReadOnly)
		})
	}
}

func TestSetupKeybindings(t *testing.T) {
	kb, err := config.GetKeyMap()
	assert.NoError(t, err)
	km := kb.KeyBindings

	assert.Contains(t, km.Help, "?")
	assert.Contains(t, km.Quit, "ctrl+c")
	assert.Contains(t, km.History, "f8")

	// Navigation.
	assert.Contains(t, km.Navigation.Up, "ctrl+k")
	assert.Contains(t, km.Navigation.Down, "ctrl+j")
	assert.Contains(t, km.Navigation.Right, "ctrl+l")
	assert.Contains(t, km.Navigation.Left, "ctrl+h")

	// Editor.
	assert.Contains(t, km.Editor.Up, "k")
	assert.Contains(t, km.Editor.Down, "j")
	assert.Contains(t, km.Editor.Right, "l")
	assert.Contains(t, km.Editor.Left, "h")
	assert.Contains(t, km.Editor.LineStart, "0")
	assert.Contains(t, km.Editor.LineEnd, "$")
	assert.Contains(t, km.Editor.GoToTop, "g")
	assert.Contains(t, km.Editor.GoToBottom, "G")
	assert.Contains(t, km.Editor.Insert, "i")
	assert.Contains(t, km.Editor.Normal, "esc")
	assert.Contains(t, km.Editor.ExecuteQuery, "ctrl+e")
	assert.Contains(t, km.Editor.ExecuteSingleQuery, "ctrl+r")

	// Sidebar.
	assert.Contains(t, km.Sidebar.GoToTop, "foo") // should be alt+k, but it tests the value comes from the config file
	assert.Contains(t, km.Sidebar.GoToBottom, "alt+j")

	// Result set.
	assert.Contains(t, km.ResultSet.NextTab, "bar") // should be tab, but it tests the value comes from the config file
	assert.Contains(t, km.ResultSet.PrevTab, "shift+tab")
	assert.Contains(t, km.ResultSet.LineStart, "0")
	assert.Contains(t, km.ResultSet.LineEnd, "$")
	assert.Contains(t, km.ResultSet.GoToTop, "g")
	assert.Contains(t, km.ResultSet.GoToBottom, "G")
}
