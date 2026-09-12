package bubbletui

import (
	"context"
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"

	"github.com/stretchr/testify/assert"
)

func TestNavigateSidebar(t *testing.T) {
	ctx := context.Background()
	km := SidebarKeyMapTest()
	client, err := NewTestClient()
	assert.Nil(t, err)
	s, err := NewSidebarViewport(ctx, client, km)
	assert.Nil(t, err)
	updates := []struct {
		Key          tea.KeyPressMsg
		ExpectedNode string
	}{
		{tea.KeyPressMsg{Code: tea.KeyEnter}, ":memory:"}, // Press enter to load table tree
		{tea.KeyPressMsg{Mod: tea.ModCtrl, Code: 'b'}, "actor_name - v"},
		{tea.KeyPressMsg{Mod: tea.ModCtrl, Code: 't'}, ":memory:"},
		{tea.KeyPressMsg{Code: tea.KeyDown}, "actor - t"},
		{tea.KeyPressMsg{Code: tea.KeyDown}, "country - t"},
	}
	for _, update := range updates {
		s, _ = s.Update(update.Key)
		assert.Equal(t, s.dbTree.GetFocusedNode().Name(), update.ExpectedNode)
	}
}

func TestSelectSidebarTable(t *testing.T) {
	ctx := context.Background()
	km := SidebarKeyMapTest()
	client, err := NewTestClient()
	assert.Nil(t, err)
	s, err := NewSidebarViewport(ctx, client, km)
	assert.Nil(t, err)
	assert.Equal(t, s.dbTree.GetFocusedNode().Name(), ":memory:", "The database name is unexpected")

	// Navidate to first table
	s, _ = s.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	s, _ = s.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	// Select first table
	s, cmd := s.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	tableMsg, ok := cmd().(selectTableMsg)
	assert.True(t, ok, "Message should be a selectTableMsg")
	assert.Equal(t, tableMsg.Table, "actor", "Should emit a query to the first table")

	// Navigate to last view
	s, _ = s.Update(tea.KeyPressMsg{Mod: tea.ModCtrl, Code: 'b'})
	s, cmd = s.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	viewMsg, ok := cmd().(selectViewMsg)
	assert.True(t, ok, "Message should be a selectViewMsg")
	assert.Equal(t, viewMsg.View, "actor_name", "Should emit a query to the last view")
}

func NewTestClient() (*client.Client, error) {
	var opts command.Options
	opts.Driver = drivers.SQLite
	opts.Host = ":memory:"
	opts.DBName = ":memory:"
	c, err := client.New(opts)
	if err != nil {
		return nil, err
	}
	stmts := []string{
		"CREATE TABLE actor (actor_id INTEGER NOT NULL, first_name VARCHAR(45) NOT NULL, last_name VARCHAR(45) NOT NULL, PRIMARY KEY (actor_id));",
		"CREATE TABLE country (country_id INTEGER NOT NULL, country VARCHAR(50) NOT NULL, last_update TIMESTAMP, PRIMARY KEY (country_id));",
		"CREATE VIEW actor_name (actor_id, fullname) AS SELECT actor_id, concat(first_name, ' ', last_name) as fullname FROM actor;",
	}
	for _, stmt := range stmts {
		_, err = c.DB().Exec(stmt)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func SidebarKeyMapTest() keys.SidebarKeyMap {
	return keys.SidebarKeyMap{
		GoToTop: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "go to top (sidebar database graph)"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("ctrl+b"),
			key.WithHelp("ctrl+b", "go to bottom (sidebar database graph)"),
		),
	}
}
