package bubbletui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/danvergara/dblab/internal/history"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/davecgh/go-spew/spew"
)

const (
	dblabJSONStyle = "dblab-cyberpunk"
)

// Register the dblab-cyberpunk style Highlight json inspects.
var _ = styles.Register(chroma.MustNewStyle(dblabJSONStyle, chroma.StyleEntries{
	chroma.NameTag:             "#FF00FF bold", // keys        → hiMagenta
	chroma.LiteralStringDouble: "#39FF14",      // strings     → cyberGreen
	chroma.LiteralString:       "#39FF14",
	chroma.LiteralNumber:       "#BF40BF",      // numbers     → neonPurple
	chroma.KeywordConstant:     "#BF40BF bold", // true/false/null
	chroma.Punctuation:         "#E0E0E0",      // {} [] : ,   → whiteText
	chroma.Error:               "#FF0000 bold",
	chroma.Background:          "#E0E0E0", // default fg, no bg set
}))

// tabStyles is for tab styling.
// The tabs are used to show table metadata.
type tabStyles struct {
	inactiveTab      lipgloss.Style
	activeTab        lipgloss.Style
	activeTabFocused lipgloss.Style
	border           lipgloss.Border
}

// newTabStyles function retuns a pointer to the tabStyles.
// It defines the highlighted background used for the active tab and the
// border character shared with the result set window below the tab row.
func newTabStyles() *tabStyles {
	s := new(tabStyles)
	s.inactiveTab = lipgloss.NewStyle().
		Padding(0, 1)
	s.activeTabFocused = s.inactiveTab.
		Background(neonPurple)
	s.activeTab = s.inactiveTab.
		Background(darkPurple)
	s.border = lipgloss.RoundedBorder()
	return s
}

type MetadataPanel interface {
	tea.Model
}

type TablePanel struct {
	table table.Model
}

func (t *TablePanel) Init() tea.Cmd { return nil }

func (t *TablePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedTable, cmd := t.table.Update(msg)
	t.table = updatedTable
	return t, cmd
}

func (t *TablePanel) View() tea.View {
	return tea.NewView(t.table.View())
}

func (t *TablePanel) GotoTop() { t.table.GotoTop() }

func (t *TablePanel) GotoBottom() { t.table.GotoBottom() }

type TextPanel struct {
	content string
}

func (t *TextPanel) Init() tea.Cmd {
	return nil
}

func (t *TextPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *TextPanel) View() tea.View {
	return tea.NewView(t.content)
}

func (t *TextPanel) SetContent(content string) {
	t.content = content
}

type ResultSet struct {
	focused       bool
	tabs          []string
	activeTab     int
	width, height int
	tabStyles     *tabStyles

	keyMap keys.ResultSetKeyMap

	viewport       viewport.Model
	tablesMetadata []MetadataPanel
	dump           io.Writer
}

func NewResultSet(keyMap keys.ResultSetKeyMap) ResultSet {
	var dump *os.File
	if _, ok := os.LookupEnv("DBLAB_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("results_messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	rs := ResultSet{
		tabs:     []string{"Data", "Columns", "Indexes", "Constraints"},
		keyMap:   keyMap,
		viewport: viewport.New(viewport.WithHeight(0), viewport.WithWidth(0)),
		dump:     dump,
	}

	rs.tabStyles = newTabStyles()
	rs.setupTables()

	return rs
}

func (r *ResultSet) Focus() {
	r.focused = true
}

func (r *ResultSet) Blur() {
	r.focused = false
}

func (r *ResultSet) SetSize(w, h int) {
	r.width = w
	r.height = h
	r.viewport.SetWidth(w - 4)
	r.viewport.SetHeight(h)
	for _, panel := range r.tablesMetadata {
		if tp, ok := panel.(*TablePanel); ok {
			tp.table.SetHeight(h - 2)
			tp.table.SetWidth(w - 2)
		}
	}
}

func (r *ResultSet) setupViews() {
	viewDef := newTextPanel()
	columns := newTablePanel(r.height, r.width)
	r.tablesMetadata = []MetadataPanel{
		viewDef,
		columns,
	}
}

func (r *ResultSet) setupTables() {
	columns := newTablePanel(r.height, r.width)
	data := newTablePanel(r.height, r.width)
	constraints := newTablePanel(r.height, r.width)
	indexes := newTablePanel(r.height, r.width)
	r.tablesMetadata = []MetadataPanel{
		data,
		columns,
		indexes,
		constraints,
	}
}

func (r ResultSet) Init() tea.Cmd {
	return nil
}

func (r ResultSet) Update(msg tea.Msg) (ResultSet, tea.Cmd) {
	if r.dump != nil {
		spew.Fdump(r.dump, msg)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, r.keyMap.GoToTop):
			if tablePanel, ok := r.tablesMetadata[r.activeTab].(*TablePanel); ok {
				tablePanel.GotoTop()
			}
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			r.viewport.GotoTop()
			return r, nil
		case key.Matches(msg, r.keyMap.GoToBottom):
			if tablePanel, ok := r.tablesMetadata[r.activeTab].(*TablePanel); ok {
				tablePanel.GotoBottom()
			}
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			r.viewport.GotoBottom()
			return r, nil
		case key.Matches(msg, r.keyMap.NextTab):
			if r.activeTab == len(r.tabs)-1 {
				r.activeTab = 0
			} else {
				r.activeTab = min(r.activeTab+1, len(r.tabs)-1)
			}
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			return r, nil
		case key.Matches(msg, r.keyMap.PrevTab):
			if r.activeTab == 0 {
				r.activeTab = len(r.tabs) - 1
			} else {
				r.activeTab = max(r.activeTab-1, 0)
			}
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			return r, nil
		case key.Matches(msg, r.keyMap.LineStart):
			r.viewport.SetXOffset(0)
			return r, nil
		case key.Matches(msg, r.keyMap.LineEnd):
			maxWidth := 0
			for line := range strings.SplitSeq(r.tablesMetadata[r.activeTab].View().Content, "\n") {
				w := lipgloss.Width(line)
				if w > maxWidth {
					maxWidth = w
				}
			}

			maxOffset := max(maxWidth-r.viewport.Width(), 0)

			r.viewport.SetXOffset(maxOffset)
			return r, nil
		}

		switch msg.String() {
		case "left", "h":
			r.viewport.ScrollLeft(4)
			return r, nil
		case "right", "l":
			r.viewport.ScrollRight(4)
			return r, nil
		}

		r.viewport, cmd = r.viewport.Update(msg)
		cmds = append(cmds, cmd)

		r.tablesMetadata[r.activeTab], cmd = r.tablesMetadata[r.activeTab].Update(msg)
		r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
		cmds = append(cmds, cmd)
	case queryErrMsg:
		errorText := fmt.Sprintf("❌ QUERY FAILED\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	case updateGraphErrMsg:
		errorText := fmt.Sprintf("❌ FAILED TO LOAD THE CATALOG\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	case queryHistoryErrMsg:
		errorText := fmt.Sprintf("❌ FAILED TO LOAD THE QUERY HISOTORY\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	case querySuccessMsg:
		r.clearTables()

		r.tabs = make([]string, len(msg.queriesResult))
		r.tablesMetadata = make([]MetadataPanel, len(msg.queriesResult))

		for i, qr := range msg.queriesResult {
			r.tabs[i] = fmt.Sprintf("query #%d", i+1)

			if qr.Error != nil {
				errPanel := newTextPanel()
				errorText := fmt.Sprintf("query #%d failed\n\n%s", i+1, qr.Error.Error())
				styledError := errorStyle.Render(errorText)
				errPanel.SetContent(styledError)
				r.tablesMetadata[i] = errPanel
				continue
			}

			switch qr.QueryType {
			case client.JSONQuery:
				jsonPanel := newTextPanel()
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, qr.JSONData, "", "  "); err != nil {
					errorText := fmt.Sprintf("query #%d failed\n\n%s", i+1, err)
					styledError := errorStyle.Render(errorText)
					jsonPanel.SetContent(styledError)
					r.tablesMetadata[i] = jsonPanel
					continue
				}

				var highlighted bytes.Buffer
				if err := quick.Highlight(&highlighted, prettyJSON.String(), "json", "terminal256", dblabJSONStyle); err != nil {
					errorText := fmt.Sprintf("query #%d failed\n\n%s", i+1, err)
					styledError := errorStyle.Render(errorText)
					jsonPanel.SetContent(styledError)
					r.tablesMetadata[i] = jsonPanel
					continue
				}

				jsonPanel.SetContent(highlighted.String())
				r.tablesMetadata[i] = jsonPanel
			case client.NormalQuery:
				panel := newTablePanel(r.height, r.width)
				tableContentColumns, tableContentRows := populateTable(qr.Headers, qr.ResultSet)
				panel.table.SetColumns(tableContentColumns)
				panel.table.SetRows(tableContentRows)
				r.tablesMetadata[i] = panel
			}
		}

		r.activeTab = 0
		r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
		r.viewport.GotoTop()

		return r, saveQueriesCmd(msg.queriesResult)
	case metadataSuccessMsg:
		r.updateMetadataOnChange(msg.metadata, msg.isTable)
		r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
		r.viewport.GotoTop()
		return r, nil
	case metadataErrMsg:
		errorText := fmt.Sprintf("❌ failed to get metadata\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	}

	return r, tea.Batch(cmds...)
}

func (r ResultSet) View() tea.View {
	tableBorder := darkPurple
	if r.focused {
		tableBorder = neonPurple
	}

	s := r.tabStyles
	borderCharStyle := lipgloss.NewStyle().Foreground(tableBorder)

	var renderedTabs []string
	for i, t := range r.tabs {
		style := s.inactiveTab
		if i == r.activeTab && r.focused {
			style = s.activeTabFocused
		}
		if i == r.activeTab && !r.focused {
			style = s.activeTab
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	separatorWidth := (max(r.width, 0) / len(r.tabs)) / 2

	separator := strings.Repeat(lipgloss.NewStyle().Foreground(tableBorder).Render("─"), separatorWidth)
	tabs := strings.Join(renderedTabs, separator)

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs)
	centered := lipgloss.PlaceHorizontal(r.width-2, lipgloss.Center, row,
		lipgloss.WithWhitespaceChars(s.border.Top),
		lipgloss.WithWhitespaceStyle(borderCharStyle))
	top := borderCharStyle.Render("╭") + centered + borderCharStyle.Render("╮")

	styledResultSet := resultSetStyle.BorderForeground(tableBorder).Width(r.width).Height(r.height).UnsetBorderTop()

	doc := strings.Builder{}
	doc.WriteString(top)
	doc.WriteString("\n")
	doc.WriteString(styledResultSet.Render(r.viewport.View()))
	return tea.NewView(doc.String())
}

func (r *ResultSet) clearTables() {
	for i := range r.tablesMetadata {
		r.tablesMetadata[i] = newTablePanel(r.height, r.width)
	}
}

// updateMetadataOnChange method is used to print the table metadata retrieved asynchronously.
func (r *ResultSet) updateMetadataOnChange(metadata *client.Metadata, isTable bool) {
	if metadata != nil {
		r.clearTables()
		if isTable {
			r.setupTables()

			r.tabs = []string{"Data", "Columns", "Indexes", "Constraints"}
			r.activeTab = 0

			// table data.
			tableContentColumns, tableContentRows := populateTable(metadata.TableContent.Columns, metadata.TableContent.Rows)
			if tablePanel, ok := r.tablesMetadata[0].(*TablePanel); ok {
				tablePanel.table.SetColumns(tableContentColumns)
				tablePanel.table.SetRows(tableContentRows)
			}

			// table columns.
			tableStructureColumns, tableStructureRows := populateTable(metadata.Structure.Columns, metadata.Structure.Rows)
			if tablePanel, ok := r.tablesMetadata[1].(*TablePanel); ok {
				tablePanel.table.SetColumns(tableStructureColumns)
				tablePanel.table.SetRows(tableStructureRows)
			}

			// table indexes.
			tableIndexColumns, tableIndexRows := populateTable(metadata.Indexes.Columns, metadata.Indexes.Rows)
			if tablePanel, ok := r.tablesMetadata[2].(*TablePanel); ok {
				tablePanel.table.SetColumns(tableIndexColumns)
				tablePanel.table.SetRows(tableIndexRows)
			}

			// table constraints.
			tableConstraintsColumns, tableConstraintsRows := populateTable(metadata.Constraints.Columns, metadata.Constraints.Rows)
			if tablePanel, ok := r.tablesMetadata[3].(*TablePanel); ok {
				tablePanel.table.SetColumns(tableConstraintsColumns)
				tablePanel.table.SetRows(tableConstraintsRows)
			}
		} else {
			r.setupViews()
			r.tabs = []string{"View Def", "Data"}
			r.activeTab = 0

			if textPanel, ok := r.tablesMetadata[0].(*TextPanel); ok {
				if len(metadata.ViewDef.Rows) > 0 {
					if len(metadata.ViewDef.Columns[0]) > 0 {
						textPanel.SetContent(metadata.ViewDef.Rows[0][0])
					}
				}
			}

			viewContentColumns, viewContentRows := populateTable(metadata.TableContent.Columns, metadata.TableContent.Rows)
			if tablePanel, ok := r.tablesMetadata[1].(*TablePanel); ok {
				tablePanel.table.SetColumns(viewContentColumns)
				tablePanel.table.SetRows(viewContentRows)
			}
		}
	}
}

func newTablePanel(height, width int) *TablePanel {
	t := table.New(
		table.WithFocused(true),
		table.WithWidth(max(width-2, 0)),
		table.WithHeight(max(height-2, 0)),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(hiMagenta).
		BorderBottom(true).
		Foreground(cyberGreen).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(black).
		Background(cyberGreen).
		Bold(true)

	t.SetStyles(s)

	return &TablePanel{
		table: t,
	}
}

func newTextPanel() *TextPanel {
	return &TextPanel{}
}

func populateTable(headers []string, data [][]string) ([]table.Column, []table.Row) {
	colWidths := make([]int, len(headers))
	var rows []table.Row

	for _, stringRow := range data {
		row := make(table.Row, len(stringRow))

		copy(row, stringRow)

		rows = append(rows, row)
	}

	for _, row := range rows {
		for i, cell := range row {
			cellWidth := lipgloss.Width(cell)
			if cellWidth > colWidths[i] {
				colWidths[i] = cellWidth
			}
		}
	}

	var columns []table.Column
	for i, header := range headers {
		finalWidth := colWidths[i]

		headerWidth := len(header) + 5
		if finalWidth < headerWidth {
			finalWidth = headerWidth
		}
		if finalWidth < 15 {
			finalWidth = 15
		}

		columns = append(columns, table.Column{
			Title: header,
			Width: finalWidth,
		})
	}

	return columns, rows
}

func saveQueriesCmd(queriesResult []client.QueryResult) tea.Cmd {
	return func() tea.Msg {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		queryHistory := make([]history.QueryHistory, 0, len(queriesResult))
		for _, qr := range queriesResult {
			qh := history.QueryHistory{
				QueryText: qr.Query,
				Timestamp: qr.Timestamp,
				Duration:  qr.Duration,
			}

			if qr.Error == nil {
				qh.Success = true
				qh.RowCount = qr.RowCount
			}

			queryHistory = append(queryHistory, qh)
		}

		if err := history.SaveHistory(configDir, queryHistory...); err != nil {
			return errMsg{err: err}
		}

		return nil
	}
}
