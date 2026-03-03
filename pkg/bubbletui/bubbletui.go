package bubbletui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type focusState int

const (
	// colors
	green  = lipgloss.Color("#1fb009")
	purple = lipgloss.Color("#800080")

	cyberGreen = lipgloss.Color("#39FF14") // High-visibility neon green
	mutedGreen = lipgloss.Color("#2ECC71") // Softer green for standard text
	neonPurple = lipgloss.Color("#BF40BF") // Bright purple for highlights
	darkPurple = lipgloss.Color("#4B0082") // Deep violet for backgrounds
	whiteText  = lipgloss.Color("#E0E0E0") // Off-white for readability

	activeBorder   = lipgloss.Color("62")
	inactiveBorder = lipgloss.Color("240")

	listHeight = 14

	focusInput focusState = iota
	focusList
	focusTable
)

var (
	baseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	listStyle  = baseStyle
	inputStyle = baseStyle
	tableStyle = baseStyle

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(green).
			Foreground(purple).
			AlignVertical(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Foreground(green)
)

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func newStyles() styles {
	var s styles

	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(cyberGreen).
		Background(darkPurple).
		BorderLeftForeground(neonPurple)
	s.pagination = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	return s
}

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	c        *client.Client
	bindings *command.TUIKeyBindings
	t        table.Writer
	viewport viewport.Model
	input    textarea.Model
	list     list.Model
	focus    focusState
	width    int
	height   int
	styles   styles
	ready    bool
}

func NewModel(c *client.Client) (*Model, error) {
	m := &Model{
		focus: focusInput,
	}

	if err := m.prepare(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		fixedFooterHeight := 2
		availableHeight := m.height - fixedFooterHeight

		leftWidth := m.width / 5
		rightWidth := m.width - leftWidth

		titleHeight := availableHeight / 4
		listContainerHeight := availableHeight - titleHeight
		inputHeight := availableHeight / 3

		tableHeight := availableHeight - inputHeight

		m.list.SetSize(m.width-4, listContainerHeight-2)
		m.input.SetWidth(rightWidth - 4)
		m.input.SetHeight(inputHeight - 2)

		if !m.ready {
			m.viewport = viewport.New(rightWidth-4, tableHeight-2)
			m.viewport.SetContent(m.t.Render())
			m.ready = true
		} else {
			m.viewport.Width = rightWidth - 4
			m.viewport.Height = tableHeight - 2
		}
		if m.ready {
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlL:
			if m.focus == focusList {
				m.focus = focusInput
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		case tea.KeyCtrlJ:
			if m.focus == focusInput {
				m.focus = focusTable
				m.input.Blur()
			}
		case tea.KeyCtrlH:
			if m.focus == focusTable {
				m.focus = focusList
			}

			if m.focus == focusInput {
				m.input.Blur()
				m.focus = focusList
			}
		case tea.KeyCtrlK:
			if m.focus == focusTable {
				m.focus = focusInput
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}
		switch msg.String() {
		case "left", "h":
			if m.focus == focusTable {
				m.viewport.ScrollLeft(4)
			}
		case "right", "l":
			if m.focus == focusTable {
				m.viewport.ScrollRight(4)
			}
		}
	}

	switch m.focus {
	case focusList:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case focusInput:
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	case focusTable:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	listBorder := darkPurple
	textAreaBorder := darkPurple
	tableBorder := darkPurple

	switch m.focus {
	case focusList:
		listBorder = neonPurple
	case focusInput:
		textAreaBorder = neonPurple
	case focusTable:
		tableBorder = neonPurple
	}

	leftWidth := m.width / 5
	rightWidth := m.width - leftWidth

	footerView := footerStyle.Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)")
	footerHeight := lipgloss.Height(footerView)
	availableHeight := m.height - footerHeight

	titleHeight := availableHeight / 4
	listHeight := availableHeight - titleHeight

	inputHeight := availableHeight / 3
	tableHeight := availableHeight - inputHeight

	dblabFigure := figure.NewFigure("dblab", "", true)

	titleBox := titleStyle.Width(leftWidth - 2).Height(titleHeight - 2).Render(dblabFigure.String())
	styledList := listStyle.BorderForeground(listBorder).Width(leftWidth - 2).Height(listHeight - 2).Render(m.list.View())

	styledInput := inputStyle.BorderForeground(textAreaBorder).Width(rightWidth - 2).Height(inputHeight - 2).Render(m.input.View())
	styledTable := tableStyle.BorderForeground(tableBorder).Width(rightWidth - 2).Height(tableHeight - 2).Render(m.viewport.View())

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, titleBox, styledList)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, styledInput, styledTable)

	contentLayout := lipgloss.JoinHorizontal(lipgloss.Bottom, leftColumn, rightColumn)

	return lipgloss.JoinVertical(lipgloss.Left, contentLayout, footerView)
}

func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

func (m *Model) updateStyles() {
	m.styles = newStyles()
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m *Model) prepare() error {
	m.setupTable()
	m.setupQueries()
	if err := m.setupDatabaseCatalog(); err != nil {
		return err
	}
	m.updateStyles()
	return nil
}

func (m *Model) setupTable() {
	t := table.NewWriter()

	cyberStyle := table.StyleRounded

	// 2. Define the Colors using standard ANSI high-intensity variants
	cyberStyle.Color = table.ColorOptions{
		// Make the borders Neon Purple
		Border: text.Colors{text.FgHiMagenta},

		// Make the Header text Cyber Green and Bold
		Header: text.Colors{text.FgHiGreen, text.Bold},

		// Make the Data rows Muted Green
		Row: text.Colors{text.FgGreen},

		// Optional: If you use a footer, style it here
		Footer: text.Colors{text.FgHiGreen},
	}

	t.SetStyle(cyberStyle)
	m.t = t
}

func (m *Model) setupDatabaseCatalog() error {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)

	l.Title = "Tables"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	m.list = l

	return nil
}

func (m *Model) setupQueries() {
	ti := textarea.New()
	ti.Placeholder = "Search or enter text..."
	ti.FocusedStyle.Text = lipgloss.NewStyle().Foreground(mutedGreen)
	ti.BlurredStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	ti.Focus()
	m.input = ti
}
