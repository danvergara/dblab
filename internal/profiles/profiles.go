package profiles

import (
	"encoding/json"
	"os"

	"github.com/danvergara/dblab/pkg/command"
)

type Config struct {
	Profiles map[string]command.Options `json:"profiles"`
}

func AddProfileToConfig(filePath string, name string, profile command.Options) error {
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var cfg Config
	if len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]command.Options)
	}

	cfg.Profiles[name] = profile

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tempFile := filePath + ".tmp"
	if err := os.WriteFile(tempFile, out, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, filePath)
}
