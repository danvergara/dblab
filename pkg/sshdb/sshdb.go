package sshdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/danvergara/dblab/pkg/drivers"
)

// default path to the known_hosts file.
var defaultKnownHostsPath = filepath.Join(os.Getenv("HOME"), ".ssh")

// createKnownHosts function creates known_hosts if does not exist.
// It uses the os package which has an OpenFile function, this function accepts 3 arguments:
// 1. the file path
// 2. the flag (e.g. os.O_CREATE|os.O_APPEND creates the file if not exists, if exists, appends to the file)
// 3. the last argument is the permission.
func createKnownHosts(knownHostsPath string) (err error) {
	f, err := os.OpenFile(
		filepath.Join(knownHostsPath, "known_hosts"),
		os.O_CREATE,
		0600,
	)
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	if err != nil {
		return err
	}

	return nil
}

// checkKnownHosts fucntion creates a know_hosts callback function with the New function.
// This callback function can be used to check if the host exists in the known_hosts file.
func checkKnownHosts(knownHostsPath string) (ssh.HostKeyCallback, error) {
	if knownHostsPath == "" {
		knownHostsPath = defaultKnownHostsPath
	}

	if err := createKnownHosts(knownHostsPath); err != nil {
		return nil, err
	}

	kh, err := knownhosts.New(filepath.Join(knownHostsPath, "known_hosts"))
	if err != nil {
		return nil, err
	}

	return kh, nil
}

// keyString create human-readable SSH-key strings.
func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(
		k.Marshal(),
	) // e.g. "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY...."
}

// addHostKey adds the host key to known_hosts file by using Normalize and Line functions of knownhosts package.
// This functions implements the ssh.HostKeyCallback type wiich is a function type which signature goes like this:
// type HostKeyCallback func(hostname string, remote net.Addr, key PublicKey) error.
func addHostKey(_ string, remote net.Addr, pubKey ssh.PublicKey, knownHostsPath string) error {
	if knownHostsPath == "" {
		knownHostsPath = defaultKnownHostsPath
	}

	khFilePath := filepath.Join(knownHostsPath, "known_hosts")

	f, err := os.OpenFile(khFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	knownHosts := knownhosts.Normalize(remote.String())
	_, err = f.WriteString(knownhosts.Line([]string{knownHosts}, pubKey))
	return err
}

// PostgresViaSSHDialer implements the driver.Driver interface to register the connection to the database via the ssh tunnel.
type PostgresViaSSHDialer struct {
	client *ssh.Client
}

func (sd *PostgresViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(sd, s)
}

func (sd *PostgresViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return sd.client.Dial(network, address)
}

func (sd *PostgresViaSSHDialer) DialTimeout(
	network, address string,
	timeout time.Duration,
) (net.Conn, error) {
	return sd.client.Dial(network, address)
}

// MySQLViaSSHDialer used to register the database connection via the ssh tunnel.
type MySQLViaSSHDialer struct {
	client *ssh.Client
}

func (m *MySQLViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return m.client.Dial("tcp", addr)
}

// SSHConfig struct setup the ssh tunnel to connect with a given database.
type SSHConfig struct {
	sshUser        string
	sshPass        string
	sshKeyFile     string
	sshKeyPass     string
	sshHost        string
	sshPort        string
	sshClient      *ssh.Client
	dbDriver       string
	dbURL          string
	knownHostsPath string
}

type Option func(*SSHConfig)

func New(opts ...Option) *SSHConfig {
	c := &SSHConfig{}

	for _, o := range opts {
		o(c)
	}

	return c
}

func WithSSHUser(sshUser string) Option {
	return func(c *SSHConfig) {
		c.sshUser = sshUser
	}
}

func WithPass(sshPass string) Option {
	return func(c *SSHConfig) {
		c.sshPass = sshPass
	}
}

func WithSSHKeyFile(sshKeyFile string) Option {
	return func(c *SSHConfig) {
		c.sshKeyFile = sshKeyFile
	}
}

func WithSSHKeyPass(sshKeyPass string) Option {
	return func(c *SSHConfig) {
		c.sshKeyPass = sshKeyPass
	}
}

func WithSShHost(sshHost string) Option {
	return func(c *SSHConfig) {
		c.sshHost = sshHost
	}
}

func WithSShPort(sshPort string) Option {
	return func(c *SSHConfig) {
		c.sshPort = sshPort
	}
}

func WithDBDriver(driver string) Option {
	return func(c *SSHConfig) {
		c.dbDriver = driver
	}
}

func WithDBDURL(url string) Option {
	return func(c *SSHConfig) {
		c.dbURL = url
	}
}

func WithKnownHostsPath(knownHostsPath string) Option {
	return func(c *SSHConfig) {
		c.knownHostsPath = knownHostsPath
	}
}

// SSHTunnel method sets up the ssh tunnel and does a number of things:
// Create a ssh client config object that witht he user.
// Define a HostKeyCallback to ensures known ssh server is the actual server.
// If host key checking is ignore then any server that has the same FQDN or IP address can impersonate the actual ssh server.
// Define the authentication method to perform the ssh tunnel (passsword or private key).
// Register the ViaSSHDialer with the ssh connection as a parameter.
func (c *SSHConfig) SSHTunnel() error {
	// Reference: https://github.com/melbahja/goph/blob/master/client.go
	// Reference: https://github.com/melbahja/goph/blob/master/hosts.go
	// Study the client.go and hosts.go to understand how to write host key call back.
	var (
		keyErr      *knownhosts.KeyError
		signer      ssh.Signer
		parseKeyErr error
	)
	config := &ssh.ClientConfig{
		User: c.sshUser,
		HostKeyCallback: ssh.HostKeyCallback(
			func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
				kh, err := checkKnownHosts(c.knownHostsPath)
				if err != nil {
					return err
				}

				hErr := kh(host, remote, pubKey)
				if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
					// Reference: https://www.godoc.org/golang.org/x/crypto/ssh/knownhosts#KeyError
					// if keyErr.Want slice is empty then host is unknown, if keyErr.Want is not empty
					// and if host is known then there is key mismatch the connection is then rejected.
					log.Printf(
						"%v is not a key of %s, either a MiTM attack or %s has reconfigured the host pub key.",
						keyString(pubKey),
						host,
						host,
					)
					return keyErr
				} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
					// host key not found in known_hosts then give a warning and continue to connect.
					log.Printf("%s is not trusted, adding this key: %q to known_hosts file.", host, keyString(pubKey))
					return addHostKey(host, remote, pubKey, c.knownHostsPath)
				}

				log.Printf("pubkey exists for %s.", host)
				return nil
			},
		),
	}

	if c.sshPass != "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(c.sshPass)}
	} else if c.sshKeyFile != "" {
		// Load the private key for SSH authentication.
		key, err := os.ReadFile(c.sshKeyFile)
		if err != nil {
			return fmt.Errorf("error reading private key: %w", err)
		}

		// Parse the private using a passphrase if required.
		if c.sshKeyPass != "" {
			signer, parseKeyErr = ssh.ParsePrivateKeyWithPassphrase(key, []byte(c.sshKeyPass))
		} else {
			signer, parseKeyErr = ssh.ParsePrivateKey(key)
		}
		if parseKeyErr != nil {
			return fmt.Errorf("error parsing private key: %w", parseKeyErr)
		}

		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.sshHost, c.sshPort), config)
	if err != nil {
		return fmt.Errorf("failed to connect to the ssh server: %w", err)
	}

	c.sshClient = client

	switch c.dbDriver {
	case drivers.PostgreSQL, drivers.Postgres:
		sql.Register("postgres+ssh", &PostgresViaSSHDialer{c.sshClient})
	case drivers.MySQL:
		mysql.RegisterDialContext(
			"mysql+tcp",
			func(_ context.Context, addr string) (net.Conn, error) {
				dialer := &MySQLViaSSHDialer{c.sshClient}
				return dialer.Dial(addr)
			},
		)
	}

	if c.dbURL != "" {
		switch {
		case strings.Contains(c.dbURL, drivers.Postgres):
			fallthrough
		case strings.Contains(c.dbURL, drivers.PostgreSQL):
			sql.Register("postgres+ssh", &PostgresViaSSHDialer{c.sshClient})
		case strings.Contains(c.dbURL, drivers.MySQL):
			mysql.RegisterDialContext(
				"mysql+tcp",
				func(_ context.Context, addr string) (net.Conn, error) {
					dialer := &MySQLViaSSHDialer{c.sshClient}
					return dialer.Dial(addr)
				},
			)
		}

	}

	return nil
}

// Close method closes the tcp connection.
func (c *SSHConfig) Close() error {
	return c.sshClient.Close()
}
