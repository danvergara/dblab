package history

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"charm.land/lipgloss/v2"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Muted Green
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Muted Red
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
)

// QueryHistory struct used to serialized the queriesa as the gob (encoding/gob) format.
// It's also used to be shown as list item.
type QueryHistory struct {
	QueryText string
	Timestamp time.Time
	Success   bool
	RowCount  int
	Duration  time.Duration
}

// Title returns the query text.
func (q QueryHistory) Title() string { return q.QueryText }

// Description shows if the query succeeded or not ant the time when the query was executed.
func (q QueryHistory) Description() string {
	var statusIndicator string
	if q.Success {
		statusIndicator = successStyle.Render(fmt.Sprintf("✔ %d rows", q.RowCount))
	} else {
		statusIndicator = errorStyle.Render("✘ error")
	}
	return fmt.Sprintf(
		"%s  %s  %s",
		mutedStyle.Render(q.Timestamp.Format(time.RFC1123)),
		mutedStyle.Render(q.Duration.String()),
		statusIndicator,
	)
}
func (q QueryHistory) FilterValue() string { return q.QueryText }

// ReadHistory reads the history file and returns the content as slice of QueryHistory.
func ReadHistory(baseDir string) ([]QueryHistory, error) {
	// Create the full filepath.
	// The base directory is usually the content of $XDG_CONFIG_HOME.
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

// SaveHistory function saves one or more queries to the dblab.gob file.
func SaveHistory(baseDir string, newQuery ...QueryHistory) error {
	fullPath := filepath.Join(baseDir, "dblab", "dblab.gob")
	// Get the path before the file name.
	dirOnly := filepath.Dir(fullPath)
	// the dirOnly contentn is used to created the official dblab directory for config file, if it does not exist.
	if err := os.MkdirAll(dirOnly, 0755); err != nil {
		return fmt.Errorf("error at creating the dblab app-specific subdirectory, if it does not exist: %w", err)
	}

	// Call ReadHistory to get the content of the config file to append new queries to it.
	history, err := ReadHistory(baseDir)
	if err != nil {
		history = []QueryHistory{}
	}

	// append new queries to the current content of the gob config file.
	history = append(history, newQuery...)

	return saveHistory(fullPath, history)
}

// saveHistory method does the heavy lifting by opening the config file or creating it if it does not exist.
func saveHistory(filename string, history []QueryHistory) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//  Create a new Gob encoder attached to the file.
	encoder := gob.NewEncoder(file)

	// Encode the data.
	err = encoder.Encode(history)
	if err != nil {
		return err
	}

	return nil
}
