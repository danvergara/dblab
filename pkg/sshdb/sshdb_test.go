package sshdb_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/danvergara/dblab/pkg/sshdb"
)

func mockSSHServer(t *testing.T, privateKeyPath, authorizedKeyPath string) (net.Listener, error) {
	t.Helper()

	// Load the server's private key
	privateBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load the client's public key
	authorizedKeyBytes, err := os.ReadFile(authorizedKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read authorized keys: %w", err)
	}
	authorizedKey, _, _, _, err := ssh.ParseAuthorizedKey(authorizedKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorized key: %w", err)
	}

	// Configure the server to require the correct public key for authentication
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if string(key.Marshal()) == string(authorizedKey.Marshal()) {
				return nil, nil // Authentication successful
			}
			return nil, fmt.Errorf("unauthorized key")
		},
	}
	config.AddHostKey(privateKey)

	// Start the server
	listener, err := net.Listen("tcp", "127.0.0.1:0") // Bind to a random port
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go func() {
				_, _, _, err := ssh.NewServerConn(conn, config)
				if err != nil {
					log.Printf("failed to establish server connection: %v", err)
				}
				conn.Close()
			}()
		}
	}()

	return listener, nil
}

func TestSSHKeyFileAuthentication(t *testing.T) {
	privateKeyPath := "testdata/test_host_key"
	authorizedKeyPath := "testdata/test_client_key.pub"

	// Start the mock server
	listener, err := mockSSHServer(t, privateKeyPath, authorizedKeyPath)
	if err != nil {
		t.Fatalf("Failed to start mock SSH server: %v", err)
	}
	defer listener.Close()
	host, port, _ := net.SplitHostPort(listener.Addr().String())

	sc := sshdb.New(
		sshdb.WithDBDriver(drivers.Postgres),
		sshdb.WithSShHost(host),
		sshdb.WithSShPort(port),
		sshdb.WithSSHUser("testuser"),
		sshdb.WithSSHKeyFile("testdata/test_client_key"),
		sshdb.WithKnownHostsPath("testdata"),
	)
	if err := sc.SSHTunnel(); err != nil {
		t.Fatalf("failed to connect to the ssh server %v", err)
	}
}
