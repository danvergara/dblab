package bubbletui

import (
	"context"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/danvergara/dblab/pkg/client"
)

type schemaLoadedMsg struct {
	items []list.Item
}

type errSchemasMsg struct{ err error }

type dbSchema struct {
	name string
}

// Implement the list.Item interface.
func (d dbSchema) Title() string       { return d.name }
func (d dbSchema) Description() string { return "" }
func (d dbSchema) FilterValue() string { return d.name }

type SchemaModel struct {
	// spinner model to show while the data is loading.
	spinner spinner.Model
	// list model to show the query history.
	list list.Model
	c    *client.Client
}

// Init method is used to initialize the spinner Tick command and fetch the query history.
func (h *SchemaModel) Init() tea.Cmd {
	return tea.Batch(h.spinner.Tick, fetchSchemasCmd(h.c))
}

func fetchSchemasCmd(c *client.Client) tea.Cmd {
	return func() tea.Msg {
		schemas, err := c.Schemas(context.Background())
		if err != nil {
			return errSchemasMsg{err: err}
		}

		var items []list.Item
		for _, s := range schemas {
			items = append(items, dbSchema{name: s})
		}

		return schemaLoadedMsg{items: items}
	}
}
