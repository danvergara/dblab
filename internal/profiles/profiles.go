package profiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/danvergara/dblab/pkg/command"

	"github.com/zalando/go-keyring"
)

// ProfileNotFound is a custom error, that implements the error interface.
// It holds the profile name to erich the error message.
type ProfileNotFound struct {
	Name string
}

func (e *ProfileNotFound) Error() string {
	return fmt.Sprintf("profile %s not found", e.Name)
}

// Config struct represents the profile configuration content.
// The configuration file is machine-driven, which means,
// is not meant to be manipulated by the user.
type Config struct {
	Profiles map[string]command.Options `json:"profiles"`
}

// FetcProfile function reads the config file and returns a profile, given the file path and the profile name.
func FetcProfile(filePath, name string) (command.Options, error) {
	// Reads the file and always returns the error.
	data, err := os.ReadFile(filePath)
	if err != nil {
		return command.Options{}, err
	}

	var cfg Config
	if len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return command.Options{}, nil
		}
	}

	// return the profile from the Profiles map.
	if cfg.Profiles != nil {
		return cfg.Profiles[name], nil
	}

	// The profile is not found.
	// Returns custom error with the profile name.
	return command.Options{}, &ProfileNotFound{Name: name}
}

// addProfileToConfig functions adds a profile to the configuration file.
func addProfileToConfig(filePath string, name string, profile command.Options) error {
	// read from the file.
	data, err := os.ReadFile(filePath)
	// Only returns the error if the files exists, but it could not be opened.
	// If the file does not exists, it will be created later on.
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("the file could not be opened: %w", err)
	}

	// Unmarshal the file content into the Config instance.
	var cfg Config
	if len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
	}

	// if the Profiles slice is nil, allocate a new map of profiles.
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]command.Options)
	}

	// store the profile in the Config.
	// If the profile exists, will be over-written.
	cfg.Profiles[name] = profile

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Creates a temporary file and writes the json back into the configuration file.
	// If the files does not exists, the file will be created.
	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, out, 0644); err != nil {
		return fmt.Errorf("failed to write to the file: %w", err)
	}

	// renames the temporary file back to the original path.
	// If the file does not exists, the new one will be renamed to the desired path.
	return os.Rename(tempFile, filePath)
}

// SaveProfile function creates the offical config path to the configuration file, if it does not exist.
// Then, saves the password in the OS keyring system.
// Finally, saves the profile in the dblab's config file.
func SaveProfile(basedDir, name string, profile command.Options) error {
	// Build the full path to your intended FILE.
	fullPath := filepath.Join(basedDir, "dblab", "dblab.json")
	// Extract just the directory portion.
	// This changes "~/.config/dblab/dblab.json" to "~/.config/dblab"
	dirOnly := filepath.Dir(fullPath)

	// Create the application-specific subdirectory for dblab, if it does not exist.
	// Create the directory (and any necessary parents).
	// We use 0755 for directories (executable bit required to open folders).
	if err := os.MkdirAll(dirOnly, 0755); err != nil {
		return fmt.Errorf("error at creating the dblab app-specific subdirectory, if it does not exist: %w", err)
	}

	// Safely create the file.
	// We use 0666 for files (read/write, no execute).
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("failed to create the dblab.json file, if it does not exist: %w", err)
		}
	}
	defer file.Close()

	// Save the password providing the profile name as reference/service and the database user.
	// go-keyring uses the OS keyring interface.
	if err := keyring.Set(name, profile.User, profile.Pass); err != nil {
		return err
	}

	// Save the profile in the config file. It ignores the password contained in the profile object, which is an instance of command.Option.
	if err := addProfileToConfig(fullPath, name, profile); err != nil {
		return err
	}

	return nil
}
