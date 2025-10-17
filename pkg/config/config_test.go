package config_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
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
				host:   "localhost",
				port:   "5432",
				dbname: "users",
				user:   "postgres",
				pass:   "password",
				driver: "postgres",
				schema: "public",
				ssl:    "disable",
				limit:  50,
			},
		},
		{
			name:  "test config",
			input: "test",
			want: want{
				host:   "localhost",
				port:   "5432",
				dbname: "users",
				user:   "postgres",
				pass:   "password",
				driver: "postgres",
				schema: "public",
				ssl:    "disable",
				limit:  50,
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
		})
	}
}

func TestSetupKeybindings(t *testing.T) {
	kb, err := config.SetupKeybindings()
	assert.NoError(t, err)

	assert.Equal(t, tcell.Key(tcell.KeyCtrlSpace), kb.RunQuery)
	assert.Equal(t, tcell.Key(tcell.KeyCtrlK), kb.Navigation.Up)
	assert.Equal(t, tcell.Key(tcell.KeyCtrlJ), kb.Navigation.Down)
	assert.Equal(
		t,
		tcell.Key(tcell.KeyCtrlL),
		kb.Navigation.Right,
	)
	assert.Equal(t, tcell.Key(tcell.KeyCtrlH), kb.Navigation.Left)
	assert.Equal(
		t,
		tcell.Key(tcell.KeyCtrlT),
		kb.Constraints,
	)
	assert.Equal(
		t,
		tcell.Key(tcell.KeyCtrlI),
		kb.Indexes,
	)
	assert.Equal(
		t,
		tcell.Key(tcell.KeyCtrlS),
		kb.Structure,
	)
	assert.Equal(
		t,
		tcell.Key(tcell.KeyCtrlD),
		kb.ClearEditor,
	)
}
