package history

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type QueryHistory struct {
	QueryText string
	Timestamp time.Time
	Success   bool
}

func (q QueryHistory) Title() string       { return q.QueryText }
func (q QueryHistory) Description() string { return q.Timestamp.Format(time.RFC1123) }
func (q QueryHistory) FilterValue() string { return q.QueryText }

func ReadHistory(baseDir string) ([]QueryHistory, error) {
	filePath := filepath.Join(baseDir, "dblab", "dblab.gob")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	var history []QueryHistory

	err = decoder.Decode(&history)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func SaveHistory(baseDir string, newQuery ...QueryHistory) error {
	fullPath := filepath.Join(baseDir, "dblab", "dblab.gob")
	dirOnly := filepath.Dir(fullPath)

	if err := os.MkdirAll(dirOnly, 0755); err != nil {
		return fmt.Errorf("error at creating the dblab app-specific subdirectory, if it does not exist: %w", err)
	}

	history, err := ReadHistory(baseDir)
	if err != nil {
		history = []QueryHistory{}
	}

	history = append(history, newQuery...)

	return saveHistory(fullPath, history)
}

func saveHistory(filename string, history []QueryHistory) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//  Create a new Gob encoder attached to the file.
	encoder := gob.NewEncoder(file)

	// Encode the data
	err = encoder.Encode(history)
	if err != nil {
		return err
	}

	return nil
}
