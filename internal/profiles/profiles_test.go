package profiles

import (
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func saveTestProfile(basedDir, name string, profile command.Options) {
	_ = SaveProfile(basedDir, name, profile)
}

func TestSaveProfile(t *testing.T) {
	keyring.MockInit()
	sandboxDir := t.TempDir()

	duplicateProfile := command.Options{
		Driver: "mysql",
		Host:   "localhost",
		Port:   "3306",
		User:   "sakila",
		Pass:   "12345",
		DBName: "sakila",
		Limit:  50,
	}

	saveTestProfile(sandboxDir, "sakila", duplicateProfile)

	type input struct {
		baseDir     string
		profileName string
		profile     command.Options
	}

	type expected struct {
		expectError bool
	}

	var tests = []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "Creates new file successfully",
			input: input{
				profileName: "pagila",
				baseDir:     sandboxDir,
				profile: command.Options{
					Driver: "postgres",
					Host:   "localhost",
					Port:   "5432",
					Pass:   "12345",
					User:   "postgres",
					DBName: "postgres",
					Schema: "public",
					Limit:  50,
					SSL:    "disable",
				},
			},
		},
		{
			name: "Handles existing file gracefully",
			input: input{
				profileName: "sakila",
				baseDir:     sandboxDir,
				profile:     duplicateProfile,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SaveProfile(test.input.baseDir, test.input.profileName, test.input.profile)
			if test.expected.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
