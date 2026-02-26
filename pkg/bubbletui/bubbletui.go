package bubbletui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusState int

const (
	// colors
	colorGreen     = lipgloss.Color("#1fb009")
	activeBorder   = lipgloss.Color("62")
	inactiveBorder = lipgloss.Color("240")

	focusInput focusState = iota
	focusList
	focusTable
)

var (
	leftColumnStyle = lipgloss.NewStyle().
			MarginRight(2).
			Border(lipgloss.RoundedBorder())

	rightTopStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Border(lipgloss.NormalBorder())

	rightBottomStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder())

	footerStyle = lipgloss.NewStyle().
			Foreground(colorGreen)
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type Model struct {
	input   textarea.Model
	list    list.Model
	table   table.Model
	focused focusState
	width   int
	height  int
}

func NewModel() Model {
	ti := textarea.New()
	ti.Placeholder = "Search or enter text..."
	ti.Focus() // Start with the input focused

	var items []list.Item = []list.Item{
		item("users"),
		item("products"),
		item("admins"),
		item("invoices"),
		item("relays"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 20, 10)
	l.Title = "Tables"
	l.SetShowHelp(false)

	cols := []table.Column{
		{Title: "url", Width: 20},
		{Title: "status", Width: 10},
		{Title: "latency", Width: 10},
	}
	rows := []table.Row{
		{"wss://relay.damus.io", "Online", "42"},
		{"wss://nos.lol", "Online", "85"},
		{"wss://nostr.mom", "Online", "120"},
	}
	mTable := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithHeight(6),
	)
	return Model{
		table: mTable,
		input: ti,
		list:  l,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	// 1. Get the raw string outputs from the bubbles
	listView := m.list.View()
	inputView := m.input.View()
	tableView := m.table.View()

	listColor := inactiveBorder
	inputColor := inactiveBorder
	tableColor := inactiveBorder

	switch m.focused {
	case focusList:
		listColor = activeBorder
	case focusInput:
		inputColor = activeBorder
	case focusTable:
		tableColor = activeBorder
	}

	styledList := leftColumnStyle.BorderForeground(listColor).Render(listView)
	styledInput := rightTopStyle.BorderForeground(inputColor).Render(inputView)
	styledTable := rightBottomStyle.BorderForeground(tableColor).Render(tableView)
	// 4. Compose the layout
	// Stack input and table vertically
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, styledInput, styledTable)

	// Join the list and the right column horizontally
	finalLayout := lipgloss.JoinHorizontal(lipgloss.Top, styledList, rightColumn)

	// Add a little instruction text at the very bottom
	helpText := lipgloss.NewStyle().Foreground(colorGreen).Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)")
	return "\n" + finalLayout + helpText + "\n"
}

func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
