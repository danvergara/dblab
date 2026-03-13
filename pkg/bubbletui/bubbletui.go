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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type focusState int

const (
	// colors.
	green  = lipgloss.Color("#1fb009")
	purple = lipgloss.Color("#800080")

	cyberGreen = lipgloss.Color("#39FF14") // High-visibility neon green
	mutedGreen = lipgloss.Color("#2ECC71") // Softer green for standard text
	neonPurple = lipgloss.Color("#BF40BF") // Bright purple for highlights
	darkPurple = lipgloss.Color("#4B0082") // Deep violet for backgrounds
	whiteText  = lipgloss.Color("#E0E0E0") // Off-white for readability

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

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")). // Bright Red
			Bold(true).
			Padding(1, 2)
)

type metadataSucessMsg struct {
	metadata *client.Metadata
}

type metadataErrMsg struct{ err error }

type querySuccessMsg struct {
	columns []string
	rows    [][]string
}

type queryErrMsg struct{ err error }

type tabStyles struct {
	doc         lipgloss.Style
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
	window      lipgloss.Style
}

func newTabStyles() *tabStyles {
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	s := new(tabStyles)
	s.doc = lipgloss.NewStyle().
		Padding(1, 2, 1, 2)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(darkPurple).
		Padding(0, 0)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	s.window = lipgloss.NewStyle().
		BorderForeground(neonPurple).
		Padding(2, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
	return s
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

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
	c         *client.Client
	tabs      []string
	activeTab int
	// bindings  *command.TUIKeyBindings
	tables    []table.Writer
	viewport  viewport.Model
	input     textarea.Model
	list      list.Model
	focus     focusState
	width     int
	height    int
	styles    styles
	tabStyles *tabStyles
	ready     bool

	leftWidth  int
	rightWidth int

	titleHeight int
	titleWidth  int

	tableListHeight int
	tableListWidth  int

	resultSetHeight int
	resultSetWidth  int

	editorHeight int
	editorWidth  int
}

func NewModel(c *client.Client) (*Model, error) {
	m := &Model{
		focus: focusInput,
		c:     c,
		tabs:  []string{"Content", "Structure", "Indexes", "Constraints"},
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

		m.leftWidth = m.width / 5
		m.rightWidth = m.width - m.leftWidth

		m.titleHeight = availableHeight/4 - 2
		m.titleWidth = m.leftWidth - 2

		m.tableListHeight = availableHeight - m.titleHeight - 4
		m.tableListWidth = m.leftWidth - 2

		m.editorWidth = m.rightWidth - 4
		m.editorHeight = availableHeight/3 - 2

		m.resultSetHeight = availableHeight - m.editorHeight - 6
		m.resultSetWidth = m.rightWidth - 4

		m.list.SetSize(m.tableListWidth, m.tableListHeight)

		m.input.SetWidth(m.editorWidth - 4)
		m.input.SetHeight(m.editorHeight - 2)

		if !m.ready {
			m.viewport = viewport.New(m.resultSetWidth-4, m.resultSetHeight-2)
			m.ready = true
		} else {
			m.viewport.Width = m.resultSetWidth - 4
			m.viewport.Height = m.resultSetHeight - 2
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
		case tea.KeyEnter:
			if m.focus == focusList {
				return m, m.runTableMetadata("")
			}
		}
		switch msg.String() {
		case "ctrl+e":
			if m.focus == focusInput {
				query := m.input.Value()
				if strings.TrimSpace(query) == "" {
					return m, nil
				}
				return m, m.executeQueryCmd(query)
			}
		case "left", "h":
			if m.focus == focusTable {
				m.viewport.ScrollLeft(4)
			}
		case "right", "l":
			if m.focus == focusTable {
				m.viewport.ScrollRight(4)
			}
		case "n", "tab":
			if m.focus == focusTable {
				m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
				m.viewport.SetContent(m.tables[m.activeTab].Render())
				return m, nil
			}
		case "p", "shift+tab":
			if m.focus == focusTable {
				m.activeTab = max(m.activeTab-1, 0)
				m.viewport.SetContent(m.tables[m.activeTab].Render())
				return m, nil
			}
		}
	case querySuccessMsg:
		m.clearTables()
		m.tables[0].AppendHeader(populateTableHeaders(msg.columns))
		m.tables[0].AppendRows(populateTableRows(msg.rows))
		m.viewport.SetContent(m.tables[0].Render())
		m.viewport.GotoTop()
		return m, nil
	case queryErrMsg:
		errorText := fmt.Sprintf("❌ QUERY FAILED\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)

		m.viewport.SetContent(styledError)

		m.viewport.GotoTop()
	case metadataSucessMsg:
		m.updateTableMetadataOnChange(msg.metadata)
		m.viewport.SetContent(m.tables[m.activeTab].Render())
		m.viewport.GotoTop()
	case metadataErrMsg:
		errorText := fmt.Sprintf("❌ table metadata failed\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		m.viewport.SetContent(styledError)
		m.viewport.GotoTop()
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

	doc := strings.Builder{}
	s := m.tabStyles

	var renderedTabs []string

	footerView := footerStyle.Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)")

	dblabFigure := figure.NewFigure("dblab", "", true)

	titleBox := titleStyle.Width(m.titleWidth).Height(m.titleHeight).Render(dblabFigure.String())
	styledTableList := listStyle.BorderForeground(listBorder).Width(m.tableListWidth).Height(m.tableListHeight).Render(m.list.View())

	styledEditor := inputStyle.BorderForeground(textAreaBorder).Width(m.editorWidth).Height(m.editorHeight).Render(m.input.View())
	styledResultSet := tableStyle.BorderForeground(tableBorder).Width(m.resultSetWidth).Height(m.resultSetHeight).UnsetBorderTop()

	numTabs := len(m.tabs)
	viewportWidth := m.resultSetWidth - 6

	baseWidth := viewportWidth / numTabs
	remainder := viewportWidth % numTabs

	for i, t := range m.tabs {
		tabWidth := baseWidth

		if i < remainder {
			tabWidth++
		}

		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab

		if isActive {
			style = s.activeTab.Width(tabWidth)
			if m.focus == focusTable {
				style = style.BorderForeground(neonPurple)
			}
		} else {
			style = s.inactiveTab.Width(tabWidth)
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "│"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	lipgloss.JoinVertical(lipgloss.Left, row, m.viewport.View())
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(styledResultSet.Render(m.viewport.View()))

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, titleBox, styledTableList)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, styledEditor, doc.String())

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
	m.tabStyles = newTabStyles()
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
	structure := setupTable()
	content := setupTable()
	constraints := setupTable()
	indexes := setupTable()
	m.tables = []table.Writer{
		content,
		structure,
		indexes,
		constraints,
	}
}

func (m *Model) clearTables() {
	for i := range m.tables {
		m.tables[i] = setupTable()
	}
}

func setupTable() table.Writer {
	t := table.NewWriter()

	cyberStyle := table.StyleRounded

	cyberStyle.Color = table.ColorOptions{
		Border: text.Colors{text.FgHiMagenta},
		Header: text.Colors{text.FgHiGreen, text.Bold},
		Row:    text.Colors{text.FgGreen},
		Footer: text.Colors{text.FgHiGreen},
	}

	t.SetStyle(cyberStyle)
	return t
}

func (m *Model) setupDatabaseCatalog() error {
	ts, err := m.c.ShowTables()
	if err != nil {
		return err
	}

	tables := make([]list.Item, 0)
	for _, ta := range ts {
		tables = append(tables, item(ta))
	}

	l := list.New(tables, itemDelegate{}, 0, 0)

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

func (m *Model) updateTableMetadataOnChange(metadata *client.Metadata) {
	if metadata != nil {
		m.clearTables()

		// table content.
		m.tables[0].AppendHeader(populateTableHeaders(metadata.TableContent.Columns))
		m.tables[0].AppendRows(populateTableRows(metadata.TableContent.Rows))

		// table structure.
		m.tables[1].AppendHeader(populateTableHeaders(metadata.Structure.Columns))
		m.tables[1].AppendRows(populateTableRows(metadata.Structure.Rows))

		// table indexes.
		m.tables[2].AppendHeader(populateTableHeaders(metadata.Indexes.Columns))
		m.tables[2].AppendRows(populateTableRows(metadata.Indexes.Rows))

		// table constraints.
		m.tables[3].AppendHeader(populateTableHeaders(metadata.Constraints.Columns))
		m.tables[3].AppendRows(populateTableRows(metadata.Constraints.Rows))
	}
}

func (m *Model) runTableMetadata(tableName string) tea.Cmd {
	return func() tea.Msg {
		if tableName == "" {
			if len(m.list.Items()) == 0 {
				return metadataErrMsg{fmt.Errorf("empty list of tables")}
			}
			tableItem := m.list.Items()[m.list.Index()]
			i, ok := tableItem.(item)
			if !ok {
				return metadataErrMsg{fmt.Errorf("not valid tables list item %d", m.list.Index())}
			}

			tableName = i.Title()
		}

		metadata, err := m.c.Metadata(tableName)
		if err != nil {
			return metadataErrMsg{err}
		}

		return metadataSucessMsg{metadata}
	}
}

func (m *Model) executeQueryCmd(query string) tea.Cmd {
	return func() tea.Msg {
		rows, columns, err := m.c.Query(query)
		if err != nil {
			return queryErrMsg{err}
		}

		return querySuccessMsg{columns, rows}
	}
}

func populateTableHeaders(headers []string) table.Row {
	headerRow := make(table.Row, len(headers))

	for i, h := range headers {
		headerRow[i] = h
	}

	return headerRow
}

func populateTableRows(data [][]string) []table.Row {
	var convertedRows []table.Row

	for _, stringRow := range data {
		newRow := make(table.Row, len(stringRow))

		for i, cellData := range stringRow {
			newRow[i] = cellData
		}

		convertedRows = append(convertedRows, newRow)
	}

	return convertedRows
}
