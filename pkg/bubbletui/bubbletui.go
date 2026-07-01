package bubbletui

import (
	"cmp"
	"context"
	"io"
	"os"
	"slices"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Digital-Shane/treeview/v2"
	"github.com/common-nighthawk/go-figure"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/davecgh/go-spew/spew"
)

// MaxQueries limits the total number of queries executed per batch.
const MaxQueries = 5

type focusState int

var (
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
)

const (
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

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(1, 2)
)

// metadataSucessMsg struct used to retrieve a given table's metadata asynchronously.
type metadataSuccessMsg struct {
	metadata *client.Metadata
	isTable  bool
}

// metadataErrMsg struct used to report error to user at the time to retrieve metadata.
type metadataErrMsg struct{ err error }

// querySuccessMsg struct used to get result sets from executed queries asynchronously.
// Sometimes, tables can be created, altered of deleted, so the this returns a refreshed list of tables.
type querySuccessMsg struct {
	reloadCatalog bool
	queriesResult []client.QueryResult
}

// queryErrMsg struct used to report when the query execution fails.
type queryErrMsg struct{ err error }

// updateGraphMsg struct used to refresh the database graph from executed queries asynchronously.
// It's triggered when the user either submits a DDL (Data Definition Language) query with a drop, create, alter, etc.
type updateGraphMsg struct {
	tree *treeview.TuiTreeModel[*client.DBNode]
}

// queryErrMsg struct used to report when the grap update fails.
type updateGraphErrMsg struct{ err error }

type Model struct {
	// database client.
	c *client.Client

	// models.
	editor          Editor
	sidebarViewport SidebarViewport
	resulstset      ResultSet

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

	// constant text on the client.
	footer        string
	renderedTitle string

	dump io.Writer

	// Stores the kill switch.
	cancelQuery context.CancelFunc
}

// NewModel returns a pointer to the main dblab bubbletea model.
// It also buids the sub-models, along with styling and the app title.
// If DBLAB_DEBUG is set, the constructor function will create a messages.log file to log bubbletui events.
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

	dblabTitle := figure.NewFigure("dblab", "", true).String()

	m := &Model{
		focus:           focusEditor,
		c:               c,
		bindings:        kb,
		editor:          NewEditor(kb),
		sidebarViewport: svp,
		resulstset:      NewResultSet(kb),
		footer:          footerStyle.Render("\n  (Press Ctrl-C to exit. Keybindings are configurable, please see the documentation for more information.)"),
		renderedTitle:   dblabTitle,
		titleHeight:     lipgloss.Height(dblabTitle),
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

		m.titleWidth = m.leftWidth

		m.sidebarViewportHeight = availableHeight - m.titleHeight - 2
		m.sidebarViewportWidth = m.leftWidth

		m.editorWidth = m.rightWidth - 4
		m.editorHeight = availableHeight/3 - 2

		m.resultSetHeight = availableHeight - m.editorHeight - 4
		m.resultSetWidth = m.rightWidth - 4

		m.editor.SetHeight(m.editorHeight)
		m.editor.SetWidth(m.editorWidth)

		m.sidebarViewport.SetSize(m.sidebarViewportWidth, m.sidebarViewportHeight)
		m.resulstset.SetSize(m.resultSetWidth, m.resultSetHeight)
		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.cancelQuery != nil {
				m.cancelQuery()
				m.cancelQuery = nil
			} else if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}
		switch {
		case key.Matches(msg, m.bindings.Navigation.Right):
			if m.focus == focusList {
				m.focus = focusEditor
				m.sidebarViewport.selected = false
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
				m.sidebarViewport.selected = true
				m.resulstset.Blur()
			}

			if m.focus == focusEditor {
				m.editor.Blur()
				m.focus = focusList
				m.sidebarViewport.selected = true
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
	case selectTableMsg:
		tableRef := client.TableRef{Name: msg.Table}
		switch m.c.Driver() {
		case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH, drivers.Oracle:
			tableRef.Schema = msg.Schema
		}
		return m, m.runTableMetadata(tableRef)
	case selectViewMsg:
		viewRef := client.ViewRef{Name: msg.View}
		switch m.c.Driver() {
		case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH, drivers.Oracle:
			viewRef.Schema = msg.Schema
		}
		return m, m.runViewMetadata(viewRef)
	case executeQueryMsg:
		ctx, cancel := context.WithCancel(context.Background())

		m.cancelQuery = cancel
		return m, m.runConcurrentlyCmd(ctx, msg.queriesToRun, 4)
	case metadataErrMsg, metadataSuccessMsg, queryErrMsg, querySuccessMsg:
		m.resulstset, cmd = m.resulstset.Update(msg)
		cmds = append(cmds, cmd)
		m.sidebarViewport, cmd = m.sidebarViewport.Update(msg)
		cmds = append(cmds, cmd)
	case updateGraphMsg, updateGraphErrMsg:
		m.sidebarViewport, cmd = m.sidebarViewport.Update(msg)
		cmds = append(cmds, cmd)
		m.resulstset, cmd = m.resulstset.Update(msg)
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

func (m Model) View() tea.View {
	var v tea.View
	v.AltScreen = true

	textAreaBorder := darkPurple

	switch m.focus {
	case focusEditor:
		textAreaBorder = neonPurple
	case focusTable:
	}

	fullFooter := m.footer
	lipgloss.JoinVertical(
		lipgloss.Left,
		fullFooter,
	)

	tightBlock := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Render(m.renderedTitle)

	centeredLogo := titleStyle.
		Width(m.titleWidth).
		MaxHeight(m.titleHeight + 2).
		Height(m.titleHeight).
		Align(lipgloss.Center).
		Render(tightBlock)

	leftColumn := lipgloss.JoinVertical(lipgloss.Left, centeredLogo, m.sidebarViewport.View())
	leftColumn = lipgloss.NewStyle().
		Width(m.leftWidth).
		MaxWidth(m.leftWidth).
		Height(m.height - lipgloss.Height(m.footer)).
		MaxHeight(m.height - lipgloss.Height(m.footer)).
		Render(leftColumn)

	styledEditor := editorStyle.BorderForeground(textAreaBorder).Width(m.editorWidth).Height(m.editorHeight).Render(m.editor.View().Content)
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, styledEditor, m.resulstset.View().Content)

	contentLayout := lipgloss.JoinHorizontal(lipgloss.Bottom, leftColumn, rightColumn)
	v.SetContent(lipgloss.JoinVertical(lipgloss.Left, contentLayout, fullFooter))
	return v
}

func (m *Model) Run() error {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

// runTableMetadata gets the given table's metadata asynchronously.
// If the query succeeds, it returns metadataSucessMsg with the metadata,
// otherwise it returns metadataErrMsg with the error.
func (m *Model) runTableMetadata(table client.TableRef) tea.Cmd {
	return func() tea.Msg {
		metadata, err := m.c.Metadata(table)
		if err != nil {
			return metadataErrMsg{err}
		}

		return metadataSuccessMsg{metadata: metadata, isTable: true}
	}
}

// runViewMetadata gets the given view's metadata asynchronously.
// If the query succeeds, it returns viewMetadataSucessMsg with the metadata,
// otherwise it returns viewMetadataErrMsg with the error.
func (m *Model) runViewMetadata(view client.ViewRef) tea.Cmd {
	return func() tea.Msg {
		metadata, err := m.c.ViewMetadata(view)
		if err != nil {
			return metadataErrMsg{err}
		}

		return metadataSuccessMsg{metadata: metadata}
	}
}

// runConcurrentlyCmd runs multiple queries concurrently by calling AsyncQuery.
// First off, it check if any query is about to alter the database graph shown in the UI.
// If so, then sets reloadCatalog to true.
// Then, it calls AsyncQuery to run multiple concurrently.
// Finally, it reads results from the resultChan channel, in a blocking way, but it does not matter,
// because this is an asynchronous function handled by the bubbletea runtime, so it does not freeze the app execution.
func (m *Model) runConcurrentlyCmd(ctx context.Context, queries []string, maxConcurrency int) tea.Cmd {
	return func() tea.Msg {
		qsMsg := querySuccessMsg{}

		for _, q := range queries {
			cleanQuery := strings.TrimSpace(q)
			firstWord := ""
			if parts := strings.Fields(cleanQuery); len(parts) > 0 {
				firstWord = strings.ToLower(parts[0])
			}
			// Check if any query runs a DDL commands.
			switch firstWord {
			case "create", "drop", "alter", "truncate", "rename":
				qsMsg.reloadCatalog = true
			}
		}

		resultChan := m.c.AsyncQuery(ctx, queries, maxConcurrency)

		var finalResults []client.QueryResult

		// Range over the channel to collect results.
		// NOTE: This blocks, but because it is inside a tea.Cmd,
		// Bubble Tea is running it in a background goroutine.
		for res := range resultChan {
			finalResults = append(finalResults, res)
		}

		// Sort the finalResults by the query index, because they ared added in a random order to the finalResults slice,
		// due to the concurrent nature of the AsyncQuery method.
		slices.SortFunc(finalResults, func(a, b client.QueryResult) int {
			return cmp.Compare(a.QueryIndex, b.QueryIndex)
		})

		qsMsg.queriesResult = finalResults
		return qsMsg
	}
}

// prepareQueriesForExecution functions splits the text coming from the text editor by ';',
// into multiple queries,
// then, it removes the leading and trailing white spaces from every query.
// To keep resources under control, the maximum numbers allowed are 5 (MaxQueries).
func prepareQueriesForExecution(rawText string) []string {
	rawQueries := strings.Split(rawText, ";")

	var validQueries []string

	for _, q := range rawQueries {
		cleanQ := strings.TrimSpace(q)
		if cleanQ != "" {
			validQueries = append(validQueries, cleanQ)
		}
	}

	if len(validQueries) > MaxQueries {
		validQueries = validQueries[:MaxQueries]
	}

	return validQueries
}
