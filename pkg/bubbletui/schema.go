package bubbletui

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/pkg/client"
)

type schemaLoadedMsg struct {
	items []list.Item
}

type schemaSelectedMsg struct {
	Name string
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
	// checks if the schema list is loading or ready to show.
	state state
	// spinner model to show while the data is loading.
	spinner spinner.Model
	// list model to show the query history.
	list list.Model
	// database client.
	c *client.Client
	// model size.
	width         int
	height        int
	loadingAction string
}

func NewSchemaModel(c *client.Client) *SchemaModel {
	delegate := list.NewDefaultDelegate()

	// white text for the normal title.
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(whiteText)

	// Cyber green color for the selected items with a magent border foreground.
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(cyberGreen).
		BorderLeftForeground(hiMagenta)

	schemaList := list.New([]list.Item{}, delegate, 0, 0)

	schemaList.Title = "Select the current schema"
	schemaList.Styles.Title = schemaList.Styles.Title.
		Background(darkPurple).
		Foreground(hiMagenta).
		Bold(true)
	schemaList.SetShowStatusBar(false)
	schemaList.SetShowHelp(true)
	schemaList.KeyMap.Quit.Unbind()
	schemaList.KeyMap.ForceQuit.Unbind()

	// Set up the spinner model.
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &SchemaModel{
		state:         stateLoading,
		spinner:       s,
		c:             c,
		list:          schemaList,
		loadingAction: "Fetch schemas from the database...",
	}
}

// SetSize method is used to set the model size when the main tui model routes the size from the tea.WindowSizeMsg message to this model.
func (s *SchemaModel) SetSize(width, height int) {
	s.height = height
	s.width = width
	s.list.SetSize(s.width-6, 14)
}

// Init method is used to initialize the spinner Tick command and fetch the query history.
func (s *SchemaModel) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.fetchSchemasCmd())
}

// Update method is used to keep elm cycle going for the schema model.
func (s *SchemaModel) Update(msg tea.Msg) (*SchemaModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(s.width-6, 14)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			selectedSchema := s.list.SelectedItem().(dbSchema).name
			return s, s.setSchema(selectedSchema)
		case "esc":
			// Press esc to get back to the main app if a query was not selected.
			return s, func() tea.Msg {
				return backToNormalMsg{}
			}
		}
	case schemaLoadedMsg:
		cmd := s.list.SetItems(msg.items)
		s.state = stateForm
		return s, cmd
	}

	switch s.state {
	case stateLoading:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd

	case stateForm:
		var cmd tea.Cmd
		s.list, cmd = s.list.Update(msg)
		return s, cmd
	}

	return s, nil
}

func (s *SchemaModel) View() tea.View {
	var (
		v       tea.View
		content string
	)
	v.AltScreen = true

	switch s.state {
	case stateLoading:
		spinnerView := fmt.Sprintf("\n %s %s\n", s.spinner.View(), s.loadingAction)
		content = lipgloss.Place(s.width, s.height, lipgloss.Center, lipgloss.Center, spinnerView)
	case stateForm:
		content = setModalContent(s.list.View(), s.width, s.height)
	}

	v.SetContent(content)
	return v
}

func (s *SchemaModel) setSchema(schema string) tea.Cmd {
	return func() tea.Msg {
		if err := s.c.SetActiveSchema(context.Background(), schema); err != nil {
			return errSchemasMsg{err: err}
		}

		return schemaSelectedMsg{Name: schema}
	}
}

func (s *SchemaModel) fetchSchemasCmd() tea.Cmd {
	return func() tea.Msg {
		schemas, err := s.c.Schemas(context.Background())
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
