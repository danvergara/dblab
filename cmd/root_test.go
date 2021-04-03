package cmd

import "testing"

func TestRootcmd(t *testing.T) {
	cmd := NewRootCmd()

	err := cmd.Execute()

	if err != nil {
		t.Error(err)
	}
}
