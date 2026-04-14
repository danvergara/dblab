package bubbletui

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

type focusState int

const (
	// colors.
	green      = lipgloss.Color("#1fb009") // Normal green
	purple     = lipgloss.Color("#800080")
	cyberGreen = lipgloss.Color("#39FF14") // High-visibility neon green
	hiMagenta  = lipgloss.Color("#FF00FF") // High-visibility Magenta
	mutedGreen = lipgloss.Color("#2ECC71") // Softer green for standard text
	neonPurple = lipgloss.Color("#BF40BF") // Bright purple for highlights
	darkPurple = lipgloss.Color("#4B0082") // Deep violet for backgrounds
	whiteText  = lipgloss.Color("#E0E0E0") // Off-white for readability
	black      = lipgloss.Color("#000000")

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
			AlignVertical(lipgloss.Center).
			Align(lipgloss.Center)

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
type metadataSuccessMsg struct {
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

type Model struct {
	// database client.
	c *client.Client

	// models.
	editor          Editor
	sidebarViewport SidebarViewport
	resulstset      ResultSet

	activeDatabase string

	// Manages the focus on the app.
	focus focusState

	// widget dimensions.
	width                 int
	height                int
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
	bindings *command.TUIKeyMap

	footer string

	dump io.Writer
}

func NewModel(c *client.Client, kb *command.TUIKeyMap) (*Model, error) {
	var dump *os.File
	if _, ok := os.LookupEnv("DBLAB_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	ctx := context.Background()
	svp, err := NewSidebarViewport(ctx, c, kb)
	if err != nil {
		return nil, err
	}

	m := &Model{
		focus:           focusEditor,
		c:               c,
		bindings:        kb,
		editor:          NewEditor(kb),
		sidebarViewport: svp,
		resulstset:      NewResultSet(kb),
		footer:          footerStyle.Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)"),
		dump:            dump,
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		availableHeight := m.height - lipgloss.Height(m.footer)

		m.leftWidth = m.width / 5
		m.rightWidth = m.width - m.leftWidth

		m.titleHeight = availableHeight / 6
		m.titleWidth = m.leftWidth - 2

		m.sidebarViewportHeight = availableHeight - m.titleHeight - 2
		m.sidebarViewportWidth = m.leftWidth - 2

		m.editorWidth = m.rightWidth - 4
		m.editorHeight = availableHeight/3 - 2

		m.resultSetHeight = availableHeight - m.editorHeight - 6
		m.resultSetWidth = m.rightWidth - 4

		m.editor.SetHeight(m.editorHeight)
		m.editor.SetWidth(m.editorWidth)

		m.sidebarViewport.SetSize(m.sidebarViewportWidth, m.sidebarViewportHeight)
		m.resulstset.SetSize(m.resultSetWidth, m.resultSetHeight)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		switch {
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
				m.resulstset.Focus()
			}
		case key.Matches(msg, m.bindings.Navigation.Left):
			if m.focus == focusTable {
				m.focus = focusList
				m.resulstset.Blur()
			}

			if m.focus == focusEditor {
				m.editor.Blur()
				m.focus = focusList
			}
		case key.Matches(msg, m.bindings.Navigation.Up):
			if m.focus == focusTable {
				m.focus = focusEditor
				cmd = m.editor.Focus()
				m.resulstset.Blur()
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}
	case selectDatabaseMsg:
		m.activeDatabase = msg.ActiveDatabase
		m.c.SetActiveDatabase(msg.ActiveDatabase)
		return m, m.fetchTablesCmd(msg.ActiveDatabase)
	case selectTableMsg:
		return m, m.runTableMetadata(msg.Table)
	case executeQueryMsg:
		return m, m.executeQueryCmd(msg.Query)
	case metadataErrMsg, metadataSuccessMsg, tablesFetchError, tablesFetchedMsg, queryErrMsg, querySuccessMsg:
		m.resulstset, cmd = m.resulstset.Update(msg)
		cmds = append(cmds, cmd)
		m.sidebarViewport, cmd = m.sidebarViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch m.focus {
	case focusEditor:
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	case focusList:
		m.sidebarViewport, cmd = m.sidebarViewport.Update(msg)
		cmds = append(cmds, cmd)
	case focusTable:
		m.resulstset, cmd = m.resulstset.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var rightText string

	listBorder := darkPurple
	textAreaBorder := darkPurple

	switch m.focus {
	case focusList:
		listBorder = neonPurple
	case focusEditor:
		textAreaBorder = neonPurple
	case focusTable:
	}

	if m.c.ShowDataCatalog() {
		if m.activeDatabase == "" {
			m.activeDatabase = m.sidebarViewport.ActiveDatabase()
		}
		label := activeLabelStyle.Render("Active: ")
		dbName := dbNameStyle.Render(m.activeDatabase + " ")
		rightText = label + dbName
	}

	gapWidth := m.width - lipgloss.Width(m.footer) - lipgloss.Width(rightText)

	if gapWidth < 0 {
		gapWidth = 0
	}

	spacer := strings.Repeat(" ", gapWidth)
	fullFooter := m.footer + spacer + rightText
	lipgloss.JoinVertical(
		lipgloss.Left,
		fullFooter,
	)
	dblabFigure := figure.NewFigure("dblab", "", true)

	tightBlock := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Render(dblabFigure.String())

	centeredLogo := titleStyle.
		Width(m.titleWidth).Height(m.titleHeight).
		Align(lipgloss.Center).
		Render(tightBlock)

	styledEditor := editorStyle.BorderForeground(textAreaBorder).Width(m.editorWidth).Height(m.editorHeight).Render(m.editor.View())
	styledTableList := tablesListStyle.BorderForeground(listBorder).Width(m.sidebarViewportWidth).Height(m.sidebarViewportHeight - 2).Render(m.sidebarViewport.View())

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, centeredLogo, styledTableList)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, styledEditor, m.resulstset.View())

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

// runTableMetadata gets the given table's metadata asynchronously.
// If the query succeeds, it returns metadataSucessMsg with the metadata, otherwise it returns metadataErrMsg with the error.
func (m *Model) runTableMetadata(tableName string) tea.Cmd {
	return func() tea.Msg {
		metadata, err := m.c.Metadata(tableName)
		if err != nil {
			return metadataErrMsg{err}
		}

		return metadataSuccessMsg{metadata}
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
