package connection

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/danvergara/dblab/pkg/command"
)

var (
	// ErrInvalidUPostgresRLFormat is the error used to notify that the postgres given url is not valid.
	ErrInvalidUPostgresRLFormat = errors.New("Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode")
	// ErrInvalidUMySQLRLFormat is the error used to notify that the given mysql url is not valid.
	ErrInvalidUMySQLRLFormat = errors.New("Invalid URL. Valid format: mysql://user:password@tcp(host:port)/db")
	// ErrInvalidURLFormat is used to notify the url is invalid.
	ErrInvalidURLFormat = errors.New("invalid url")
)

// BuildConnectionFromOpts return the connection uri string given the options passed by the uses.
func BuildConnectionFromOpts(opts command.Options) (string, error) {
	if opts.URL != "" {
		switch opts.Driver {
		case "postgres":
			return formatPostgresURL(opts)
		case "mysql":
			return formatMySQLURL(opts)
		default:
			return "", fmt.Errorf("%s: %w", opts.URL, ErrInvalidURLFormat)
		}
	}

	return "", nil
}

// formatPostgresURL returns valid uri for postgres connection.
func formatPostgresURL(opts command.Options) (string, error) {
	if !hasValidPosgresPrefix(opts.URL) {
		return "", fmt.Errorf("invalid prefix %s : %w", opts.URL, ErrInvalidUPostgresRLFormat)
	}

	uri, err := url.Parse(opts.URL)
	if err != nil {
		return "", fmt.Errorf("%v : %w", err, ErrInvalidUPostgresRLFormat)
	}

	result := map[string]string{}
	for k, v := range uri.Query() {
		result[strings.ToLower(k)] = v[0]
	}

	if result["sslmode"] == "" {
		if opts.SSL == "" {
			if strings.Contains(uri.Host, "localhost") || strings.Contains(uri.Host, "127.0.0.1") {
				result["sslmode"] = "disable"
			}
		} else {
			result["sslmode"] = opts.SSL
		}
	}

	query := url.Values{}
	for k, v := range result {
		query.Add(k, v)
	}
	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

// formatMySQLURL returns valid uri for mysql connection.
func formatMySQLURL(opts command.Options) (string, error) {
	if !hasValidMySQLPrefix(opts.URL) {
		return "", ErrInvalidUMySQLRLFormat
	}

	uri, err := url.Parse(opts.URL)
	if err != nil {
		return "", ErrInvalidUMySQLRLFormat
	}

	result := map[string]string{}
	for k, v := range uri.Query() {
		result[strings.ToLower(k)] = v[0]
	}

	query := url.Values{}
	for k, v := range result {
		query.Add(k, v)
	}
	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

// hasValidPosgresPrefix checks if a given url has the driver name in it.
func hasValidPosgresPrefix(rawurl string) bool {
	return strings.HasPrefix(rawurl, "postgres://") || strings.HasPrefix(rawurl, "postgresql://")
}

// hasValidMySQLPrefix checks if a given url has the driver name in it.
func hasValidMySQLPrefix(rawurl string) bool {
	return strings.HasPrefix(rawurl, "mysql://")
}
