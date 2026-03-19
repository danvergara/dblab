package bubbletui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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
	"github.com/savannahostrowski/tree-bubble"
)

type focusState int

const (
	// colors.
	green      = lipgloss.Color("#1fb009")
	purple     = lipgloss.Color("#800080")
	cyberGreen = lipgloss.Color("#39FF14") // High-visibility neon green
	mutedGreen = lipgloss.Color("#2ECC71") // Softer green for standard text
	neonPurple = lipgloss.Color("#BF40BF") // Bright purple for highlights
	darkPurple = lipgloss.Color("#4B0082") // Deep violet for backgrounds
	whiteText  = lipgloss.Color("#E0E0E0") // Off-white for readability

	// focus state management.
	focusEditor focusState = iota
	focusList
	focusTable
)

var (
	baseStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	tablesListStyle = baseStyle
	editorStyle     = baseStyle
	resultSetStyle  = baseStyle

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(green).
			Foreground(purple).
			AlignVertical(lipgloss.Center)

	footerStyle = lipgloss.NewStyle().
			Foreground(green)

	activeLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#B200FF")).
				Bold(true)

	dbNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9D00")).
			Bold(true).
			PaddingRight(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")). // Bright Red
			Bold(true).
			Padding(1, 2)
)

// metadataSucessMsg struct used to retrieve a given table's metadata asynchronously.
type metadataSucessMsg struct {
	metadata *client.Metadata
}

// metadataErrMsg struct used to report error to user at the time to retrieve metadata.
type metadataErrMsg struct{ err error }

// tablesFetchedMsg struct used to get a given database's tables asynchronously.
type tablesFetchedMsg struct {
	dbName string
	tables []string
}

// tablesFetchError struct used to report errors to the user at the time to get the list of tables.
type tablesFetchError struct{ err error }

// querySuccessMsg struct used to get result sets from executed queries asynchronously.
// Sometimes, tables can be created, altered of deleted, so the this returns a refreshed list of tables.
type querySuccessMsg struct {
	columns []string
	rows    [][]string
	tables  []string
}

// queryErrMsg struct used to report when the query execution fails.
type queryErrMsg struct{ err error }

// tabStyles is for tab styling.
// The tabs are used to show table metadata.
type tabStyles struct {
	inactiveTab lipgloss.Style
	activeTab   lipgloss.Style
}

// newTabStyles function retuns a pointer to the tabStyles.
// It basically defines the default borders for bot active and inactive tabs.
func newTabStyles() *tabStyles {
	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")
	s := new(tabStyles)
	s.inactiveTab = lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(darkPurple).
		Padding(0, 0)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
	return s
}

// tabBorderWithBottom function is used to define the tab borders.
// Borders changes whether the tabs is inacative or inactive.
// Active tab misses the bottom border.
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// styles struct is for generic styling.
type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

// newStyles function retunrs a styles with defaults.
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

// item implements the Item interface for required for the List Model from bubbles.
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

// itemDelegate is used to inject styling to the list items.
// Implements the ItemDelegate interface.
// It's important to highlight the selected item.
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

// Model struct implements the bubbletea's Model interface.
type Model struct {
	// database client.
	c *client.Client

	activeDatabase string

	// tab management.
	tabs      []string
	activeTab int

	// models.
	tablesMetadata  []table.Writer
	viewport        viewport.Model
	editor          textarea.Model
	tablesList      list.Model
	dbTree          tree.Model
	sidebarViewport viewport.Model

	// Manages the focus on the app.
	focus focusState

	// app dimensions.
	width  int
	height int

	// flag used to let the app know that the viewport is ready.
	ready bool

	// app styles.
	styles    styles
	tabStyles *tabStyles

	// widget dimensions.
	leftWidth             int
	rightWidth            int
	titleHeight           int
	titleWidth            int
	sidebarViewportHeight int
	sidebarViewportWidth  int
	resultSetHeight       int
	resultSetWidth        int
	editorHeight          int
	editorWidth           int

	// Key Bindings.
	bindings *command.TUIKeyBindings
}

func NewModel(c *client.Client, kb *command.TUIKeyBindings) (*Model, error) {
	m := &Model{
		focus:    focusEditor,
		c:        c,
		bindings: kb,
		tabs:     []string{"Content", "Structure", "Indexes", "Constraints"},
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

		m.sidebarViewportHeight = availableHeight - m.titleHeight - 4
		m.sidebarViewportWidth = m.leftWidth - 2

		m.editorWidth = m.rightWidth - 4
		m.editorHeight = availableHeight/3 - 2

		m.resultSetHeight = availableHeight - m.editorHeight - 6
		m.resultSetWidth = m.rightWidth - 4

		m.tablesList.SetSize(m.sidebarViewportWidth, m.sidebarViewportHeight)
		m.dbTree.SetSize(m.sidebarViewportWidth, m.sidebarViewportHeight)

		m.sidebarViewport.Height = m.sidebarViewportHeight - 4
		m.sidebarViewport.Width = m.sidebarViewportWidth - 4

		m.editor.SetWidth(m.editorWidth - 4)
		m.editor.SetHeight(m.editorHeight - 2)

		if m.c.ShowDataCatalog() {
			m.sidebarViewport.SetContent(m.dbTree.View())
		} else {
			m.sidebarViewport.SetContent(m.tablesList.View())
		}

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
		case tea.KeyEnter:
			if m.focus == focusList {
				if m.c.ShowDataCatalog() {
					cursor := m.dbTree.Cursor()
					nodes := m.dbTree.Nodes()
					path := getPathToCursor(nodes, cursor)
					if len(path) == 0 {
						return m, nil
					}
					targetNode := path[len(path)-1]

					switch targetNode.Desc {
					case "database":
						m.c.SetActiveDatabase(targetNode.Value)
						m.activeDatabase = targetNode.Value
						if len(targetNode.Children) == 0 {
							return m, m.fetchTablesCmd(targetNode.Value)
						} else {
							currentNodes := m.dbTree.Nodes()
							updatedNodes := injectTablesIntoTree(currentNodes, targetNode.Value, []string{})
							m.dbTree.SetNodes(updatedNodes)
						}
					case "table":
						return m, m.runTableMetadata(targetNode.Value)
					}
				} else {
					return m, m.runTableMetadata("")
				}
			}
		}
		switch {
		case key.Matches(msg, m.bindings.ExecuteQuery):
			if m.focus == focusEditor {
				query := m.editor.Value()
				if strings.TrimSpace(query) == "" {
					return m, nil
				}
				return m, m.executeQueryCmd(query)
			}
		case key.Matches(msg, m.bindings.NextTab):
			if m.focus == focusTable {
				m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
				m.viewport.SetContent(m.tablesMetadata[m.activeTab].Render())
				return m, nil
			}
		case key.Matches(msg, m.bindings.PrevTab):
			if m.focus == focusTable {
				m.activeTab = max(m.activeTab-1, 0)
				m.viewport.SetContent(m.tablesMetadata[m.activeTab].Render())
				return m, nil
			}
		case key.Matches(msg, m.bindings.Navigation.Right):
			if m.focus == focusList {
				m.focus = focusEditor
				cmd = m.editor.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.bindings.Navigation.Down):
			if m.focus == focusEditor {
				m.focus = focusTable
				m.editor.Blur()
			}
		case key.Matches(msg, m.bindings.Navigation.Left):
			if m.focus == focusTable {
				m.focus = focusList
			}

			if m.focus == focusEditor {
				m.editor.Blur()
				m.focus = focusList
			}
		case key.Matches(msg, m.bindings.Navigation.Up):
			if m.focus == focusTable {
				m.focus = focusEditor
				cmd = m.editor.Focus()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.bindings.PageTop):
			if m.c.ShowDataCatalog() {
				if m.focus == focusTable {
					m.viewport.GotoTop()
				}
				if m.focus == focusList {
					for m.dbTree.Cursor() > 0 {
						m.dbTree, _ = m.dbTree.Update(tea.KeyMsg{Type: tea.KeyUp})
					}

					m.syncTreeToViewport()
					return m, nil
				}
			} else {
				m.tablesList.Select(0)
				m.sidebarViewport.SetContent(m.tablesList.View())
			}
			return m, nil
		case key.Matches(msg, m.bindings.PageBottom):
			if m.focus == focusTable {
				m.viewport.GotoBottom()
			}
			if m.focus == focusList {
				if m.c.ShowDataCatalog() {
					totalNodes := m.dbTree.NumberOfNodes()
					if totalNodes > 0 {
						m.dbTree.SetCursor(totalNodes - 1)
					}

					m.syncTreeToViewport()
					return m, nil
				} else {
					totalItems := len(m.tablesList.Items())
					if totalItems > 0 {
						m.tablesList.Select(totalItems - 1)
					}
					m.sidebarViewport.SetContent(m.tablesList.View())
				}
			}
			return m, nil
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
		switch m.focus {
		case focusList:
			if m.c.ShowDataCatalog() {
				m.dbTree, cmd = m.dbTree.Update(msg)
				m.sidebarViewport.SetContent(m.dbTree.View())
				m.syncTreeToViewport()
			} else {
				m.tablesList, cmd = m.tablesList.Update(msg)
				m.sidebarViewport.SetContent(m.tablesList.View())
			}

			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case focusEditor:
			m.editor, cmd = m.editor.Update(msg)
			cmds = append(cmds, cmd)
		case focusTable:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}

	case querySuccessMsg:
		m.clearTables()
		m.tablesMetadata[0].AppendHeader(populateTableHeaders(msg.columns))
		m.tablesMetadata[0].AppendRows(populateTableRows(msg.rows))
		m.viewport.SetContent(m.tablesMetadata[0].Render())

		if len(msg.tables) > 0 {
			tables := make([]list.Item, 0)
			for _, ta := range msg.tables {
				tables = append(tables, item(ta))
			}
			m.tablesList.SetItems(tables)
		}

		m.viewport.GotoTop()

		return m, nil
	case queryErrMsg:
		errorText := fmt.Sprintf("❌ QUERY FAILED\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)

		m.viewport.SetContent(styledError)

		m.viewport.GotoTop()
	case metadataSucessMsg:
		m.updateTableMetadataOnChange(msg.metadata)
		m.viewport.SetContent(m.tablesMetadata[m.activeTab].Render())
		m.viewport.GotoTop()
	case metadataErrMsg:
		errorText := fmt.Sprintf("❌ failed to get table metadata\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		m.viewport.SetContent(styledError)
		m.viewport.GotoTop()
	case tablesFetchedMsg:
		currentNodes := m.dbTree.Nodes()
		updatedNodes := injectTablesIntoTree(currentNodes, msg.dbName, msg.tables)
		m.dbTree.SetNodes(updatedNodes)
		m.sidebarViewport.SetContent(m.dbTree.View())
	case tablesFetchError:
		errorText := fmt.Sprintf("❌ failed to retrieve the current database tables\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		m.viewport.SetContent(styledError)
		m.viewport.GotoTop()
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
	case focusEditor:
		textAreaBorder = neonPurple
	case focusTable:
		tableBorder = neonPurple
	}

	doc := strings.Builder{}
	s := m.tabStyles

	var renderedTabs []string

	footerView := footerStyle.Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)")

	var rightText string

	if m.c.ShowDataCatalog() && m.activeDatabase != "" {
		label := activeLabelStyle.Render("Active: ")
		dbName := dbNameStyle.Render(m.activeDatabase + " ")
		rightText = label + dbName
	}

	gapWidth := m.width - lipgloss.Width(footerView) - lipgloss.Width(rightText)

	if gapWidth < 0 {
		gapWidth = 0
	}

	spacer := strings.Repeat(" ", gapWidth)
	fullFooter := footerView + spacer + rightText
	lipgloss.JoinVertical(
		lipgloss.Left,
		fullFooter,
	)

	dblabFigure := figure.NewFigure("dblab", "", true)

	titleBox := titleStyle.Width(m.titleWidth).Height(m.titleHeight).Render(dblabFigure.String())

	styledTableList := tablesListStyle.BorderForeground(listBorder).Width(m.sidebarViewportWidth).Height(m.sidebarViewportHeight).Render(m.sidebarViewport.View())

	styledEditor := editorStyle.BorderForeground(textAreaBorder).Width(m.editorWidth).Height(m.editorHeight).Render(m.editor.View())
	styledResultSet := resultSetStyle.BorderForeground(tableBorder).Width(m.resultSetWidth).Height(m.resultSetHeight).UnsetBorderTop()

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

	return lipgloss.JoinVertical(lipgloss.Left, contentLayout, fullFooter)
}

func (m *Model) Run() error {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

// updateStyle setup the styles across the client.
func (m *Model) updateStyles() {
	m.tabStyles = newTabStyles()
	m.styles = newStyles()
	m.tablesList.Styles.Title = m.styles.title
	m.tablesList.Styles.PaginationStyle = m.styles.pagination
	m.tablesList.Styles.HelpStyle = m.styles.help
	m.tablesList.SetDelegate(itemDelegate{styles: &m.styles})
}

// prepare method sets up the client defaults, such as the tables, the editor, the initial queries to show the either the databases or tables the user has access to and the styles.
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
	m.tablesMetadata = []table.Writer{
		content,
		structure,
		indexes,
		constraints,
	}
}

func (m *Model) clearTables() {
	for i := range m.tablesMetadata {
		m.tablesMetadata[i] = setupTable()
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

// setupDatabaseCatalog method shows the initial database/tables catalog the user has access to.
// If the user wants to see the database catalog, the client will present a tree view, with a graph with databases and tables.
// Otherwise, the user will see a list of tables of the database connected.
func (m *Model) setupDatabaseCatalog() error {
	m.sidebarViewport = viewport.New(0, 0)
	m.sidebarViewport.KeyMap = viewport.KeyMap{}

	if m.c.ShowDataCatalog() {

		dbs, err := m.c.ShowDatabases()
		if err != nil {
			return err
		}
		nodes := make([]tree.Node, len(dbs))
		for i, db := range dbs {
			nodes[i] = tree.Node{Value: db, Desc: "database"}
		}

		m.dbTree = tree.New(nodes, m.sidebarViewportWidth, m.sidebarViewportHeight)

		// Override the Down binding.
		m.dbTree.KeyMap.Down = key.NewBinding(
			key.WithKeys("down", "j", "s"),
			key.WithHelp("↓/j/s", "move down"),
		)

		// Override the Up binding.
		m.dbTree.KeyMap.Up = key.NewBinding(
			key.WithKeys("up", "k", "w"),
			key.WithHelp("↑/k/w", "move up"),
		)
	} else {
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
		m.tablesList = l
	}

	return nil
}

func (m *Model) setupQueries() {
	ti := textarea.New()
	ti.Placeholder = "Search or enter text..."
	ti.FocusedStyle.Text = lipgloss.NewStyle().Foreground(mutedGreen)
	ti.BlurredStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	ti.Focus()
	m.editor = ti
}

// updateTableMetadataOnChange method is used to print the table metadata retrieved asynchronously.
func (m *Model) updateTableMetadataOnChange(metadata *client.Metadata) {
	if metadata != nil {
		m.clearTables()

		// table content.
		m.tablesMetadata[0].AppendHeader(populateTableHeaders(metadata.TableContent.Columns))
		m.tablesMetadata[0].AppendRows(populateTableRows(metadata.TableContent.Rows))

		// table structure.
		m.tablesMetadata[1].AppendHeader(populateTableHeaders(metadata.Structure.Columns))
		m.tablesMetadata[1].AppendRows(populateTableRows(metadata.Structure.Rows))

		// table indexes.
		m.tablesMetadata[2].AppendHeader(populateTableHeaders(metadata.Indexes.Columns))
		m.tablesMetadata[2].AppendRows(populateTableRows(metadata.Indexes.Rows))

		// table constraints.
		m.tablesMetadata[3].AppendHeader(populateTableHeaders(metadata.Constraints.Columns))
		m.tablesMetadata[3].AppendRows(populateTableRows(metadata.Constraints.Rows))
	}
}

// runTableMetadata gets the given table's metadata asynchronously.
// If the query succeeds, it returns metadataSucessMsg with the metadata, otherwise it returns metadataErrMsg with the error.
func (m *Model) runTableMetadata(tableName string) tea.Cmd {
	return func() tea.Msg {
		if tableName == "" {
			if len(m.tablesList.Items()) == 0 {
				return metadataErrMsg{fmt.Errorf("empty list of tables")}
			}
			tableItem := m.tablesList.Items()[m.tablesList.Index()]
			i, ok := tableItem.(item)
			if !ok {
				return metadataErrMsg{fmt.Errorf("not valid tables list item %d", m.tablesList.Index())}
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

// executeQueryCmd method executes queryes asynchronously, so it does not block the bubbletea execution.
// If it succeeds, returns a querySuccessMsg with the resultset. Otherwise, it returns queryErrMsg with the error.
func (m *Model) executeQueryCmd(query string) tea.Cmd {
	return func() tea.Msg {
		var ts []string
		rows, columns, err := m.c.Query(query)
		if err != nil {
			return queryErrMsg{err}
		}

		switch {
		case strings.Contains(strings.ToLower(query), "alter table"):
			fallthrough
		case strings.Contains(strings.ToLower(query), "drop table"):
			fallthrough
		case strings.Contains(strings.ToLower(query), "create table"):
			ts, err = m.c.ShowTables()
			if err != nil {
				return queryErrMsg{err}
			}
		}

		return querySuccessMsg{columns: columns, rows: rows, tables: ts}
	}
}

// fetchTablesCmd method gets the list from a given database asynchronously.
// If it succeeds, returns tablesFetchedMsg. Otherwise, it returns tablesFetchError with the error.
func (m *Model) fetchTablesCmd(dbName string) tea.Cmd {
	return func() tea.Msg {
		ts, err := m.c.ShowTablesPerDB(dbName)
		if err != nil {
			return tablesFetchError{err}
		}

		return tablesFetchedMsg{
			dbName: dbName,
			tables: ts,
		}
	}
}

func (m *Model) syncTreeToViewport() {
	m.sidebarViewport.SetContent(m.dbTree.View())

	cursor := m.dbTree.Cursor()
	if cursor >= m.sidebarViewport.YOffset+m.sidebarViewport.Height {
		m.sidebarViewport.SetYOffset(cursor - m.sidebarViewport.Height + 1)
	}
	if cursor < m.sidebarViewport.YOffset {
		m.sidebarViewport.SetYOffset(cursor)
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

// getPathToCursor returns the full path of nodes from the Root down to the Cursor position.
func getPathToCursor(nodes []tree.Node, targetCursor int) []tree.Node {
	var currentPath []tree.Node
	var currentIndex int

	var traverse func(currentNodes []tree.Node) bool

	traverse = func(currentNodes []tree.Node) bool {
		for _, node := range currentNodes {
			currentPath = append(currentPath, node)

			if currentIndex == targetCursor {
				return true
			}

			currentIndex++

			if len(node.Children) > 0 {
				if traverse(node.Children) {
					return true // The target was found deep inside this branch
				}
			}

			currentPath = currentPath[:len(currentPath)-1]
		}
		return false
	}

	traverse(nodes)
	return currentPath
}

// injectTablesIntoTree traverses the graph to find a not which is the database of interest to add tables to it.
func injectTablesIntoTree(nodes []tree.Node, targetDB string, tables []string) []tree.Node {
	for i := range nodes {

		if nodes[i].Desc == "database" && nodes[i].Value == targetDB {

			var newChildren []tree.Node
			for _, tableName := range tables {
				newChildren = append(newChildren, tree.Node{
					Value: tableName,
					Desc:  "table",
				})
			}

			nodes[i].Children = newChildren

			break
		}
	}

	return nodes
}
