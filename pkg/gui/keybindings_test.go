package gui

import (
	"testing"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/gocui"
	"github.com/stretchr/testify/assert"
)

func TestKeyBindings(t *testing.T) {
	g, _ := New(&gocui.Gui{}, &client.Client{})

	err := g.keybindings()

	assert.NoError(t, err)
}
