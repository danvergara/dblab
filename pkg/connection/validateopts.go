package connection

import (
	"fmt"

	"github.com/danvergara/dblab/pkg/command"
)

// ValidateOpts make sure the important fields used to open a connection aren't empty.
func ValidateOpts(opts command.Options) error {
	if opts.Host == "" && opts.Port == "" && opts.User == "" && opts.Pass == "" &&
		opts.DBName == "" &&
		opts.Driver == "" &&
		opts.URL == "" {
		return fmt.Errorf("non-empty values required to open a session with a database")
	}
	return nil
}
