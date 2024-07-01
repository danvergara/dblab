package command

import "os"

// Options is a struct that stores the provided commands by the user.
type Options struct {
	Driver string
	URL    string
	Host   string
	Port   string
	User   string
	Pass   string
	DBName string
	// PostgreSQL only.
	Schema string
	Limit  uint
	Socket string
	SSL    string
	// SSL connection params.
	SSLCert     string
	SSLKey      string
	SSLPassword string
	SSLRootcert string
	// oracle specific.
	TraceFile string
	SSLVerify string
	Wallet    string
	// sql server.
	Encrypt                string
	TrustServerCertificate string
	ConnectionTimeout      string
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
