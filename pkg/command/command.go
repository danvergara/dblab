package command

import (
	"os"
)

// Options is a struct that stores the provided commands by the user.
type Options struct {
	Driver string `json:"driver"`
	URL    string `json:"url"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"-"`
	DBName string `json:"db_name"`
	Schema string `json:"schema"`
	Limit  uint   `json:"limit"`
	Socket string `json:"socket"`
	SSL    string `json:"ssl"`
	// SSH.
	SSHHost          string `json:"ssh_host"`
	SSHPort          string `json:"ssh_port"`
	SSHUser          string `json:"ssh_user"`
	SSHPass          string `json:"-"`
	SSHKeyFile       string `json:"ssh_key_file"`
	SSHKeyPassphrase string `json:"ssh_key_passphrase"`
	// SSL connection params.
	SSLCert     string `json:"ssl_cert"`
	SSLKey      string `json:"ssl_key"`
	SSLPassword string `json:"-"`
	SSLRootcert string `json:"ssl_rootcert"`
	// oracle specific.
	TraceFile string `json:"trace_file"`
	SSLVerify string `json:"ssl_verify"`
	Wallet    string `json:"wallet"`
	// sql server.
	Encrypt                string `json:"encrypt"`
	TrustServerCertificate string `json:"trust_server_certificate"`
	ConnectionTimeout      string `json:"connection_timeout"`
	// Read Only mode.
	ReadOnly bool `json:"read_only"`
}

// SetDefault returns a Options struct and fills the empty
// values with environment variables if any.
func SetDefault(opts Options) Options {
	if opts.URL == "" {
		opts.URL = os.Getenv("DATABASE_URL")
	}

	if opts.Driver == "" {
		opts.Driver = os.Getenv("DB_DRIVER")
	}

	if opts.Host == "" {
		opts.Host = os.Getenv("DB_HOST")
	}

	if opts.User == "" {
		opts.User = os.Getenv("DB_USER")
	}

	if opts.Pass == "" {
		opts.Pass = os.Getenv("DB_PASSWORD")
	}

	if opts.DBName == "" {
		opts.DBName = os.Getenv("DB_NAME")
	}

	if opts.Port == "" {
		opts.Port = os.Getenv("DB_PORT")
	}

	if opts.Schema == "" {
		opts.Schema = os.Getenv("DB_SCHEMA")
	}

	return opts
}
