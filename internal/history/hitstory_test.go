package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSaveHistory(t *testing.T) {
	sandboxDir := t.TempDir()

	type input struct {
		baseDir string
		queries []QueryHistory
	}

	type expected struct {
		totalAfterSave int
		expectError    bool
	}

	now := time.Now()

	var tests = []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "Save single query to a fresh file",
			input: input{
				baseDir: sandboxDir,
				queries: []QueryHistory{
					{
						QueryText: "SELECT * FROM users",
						Timestamp: now,
						Success:   true,
					},
				},
			},
			expected: expected{
				totalAfterSave: 1,
			},
		},
		{
			name: "Append multiple queries to existing file",
			input: input{
				baseDir: sandboxDir,
				queries: []QueryHistory{
					{
						QueryText: "INSERT INTO users (name) VALUES ('alice')",
						Timestamp: now,
						Success:   true,
					},
					{
						QueryText: "DELETE FROM users WHERE id = 1",
						Timestamp: now,
						Success:   false,
					},
				},
			},
			expected: expected{
				totalAfterSave: 3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SaveHistory(test.input.baseDir, test.input.queries...)
			if test.expected.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			history, err := ReadHistory(test.input.baseDir)
			require.NoError(t, err)
			require.Len(t, history, test.expected.totalAfterSave)
		})
	}
}

func TestReadHistory(t *testing.T) {
	t.Run("Read from a valid gob file", func(t *testing.T) {
		sandboxDir := t.TempDir()

		now := time.Now().Truncate(time.Second)

		queries := []QueryHistory{
			{
				QueryText: "SELECT 1",
				Timestamp: now,
				Success:   true,
			},
			{
				QueryText: "SELECT 2",
				Timestamp: now,
				Success:   false,
			},
		}

		err := SaveHistory(sandboxDir, queries...)
		require.NoError(t, err)

		history, err := ReadHistory(sandboxDir)
		require.NoError(t, err)
		require.Len(t, history, 2)

		require.Equal(t, "SELECT 1", history[0].QueryText)
		require.Equal(t, true, history[0].Success)
		require.Equal(t, now, history[0].Timestamp)

		require.Equal(t, "SELECT 2", history[1].QueryText)
		require.Equal(t, false, history[1].Success)
	})

	t.Run("Read from a non-existent path", func(t *testing.T) {
		sandboxDir := t.TempDir()

		_, err := ReadHistory(sandboxDir)
		require.Error(t, err)
	})

	t.Run("Read from a corrupted file", func(t *testing.T) {
		sandboxDir := t.TempDir()

		dirPath := filepath.Join(sandboxDir, "dblab")
		err := os.MkdirAll(dirPath, 0755)
		require.NoError(t, err)

		filePath := filepath.Join(dirPath, "dblab.gob")
		err = os.WriteFile(filePath, []byte("not valid gob data"), 0666)
		require.NoError(t, err)

		_, err = ReadHistory(sandboxDir)
		require.Error(t, err)
	})
}
