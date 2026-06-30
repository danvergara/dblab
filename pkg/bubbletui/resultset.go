package bubbletui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

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
		Padding(0, 1)
	s.activeTab = s.inactiveTab.
		Border(activeTabBorder, true)
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

	bindings *command.TUIKeyMap

	viewport       viewport.Model
	tablesMetadata []MetadataPanel
	dump           io.Writer
}

func NewResultSet(kb *command.TUIKeyMap) ResultSet {
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
		bindings: kb,
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
		case key.Matches(msg, r.bindings.NextTab):
			r.activeTab = min(r.activeTab+1, len(r.tabs)-1)
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			return r, nil
		case key.Matches(msg, r.bindings.PrevTab):
			r.activeTab = max(r.activeTab-1, 0)
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View().Content)
			return r, nil
		case key.Matches(msg, r.bindings.BeginningOfLine):
			r.viewport.SetXOffset(0)
			return r, nil
		case key.Matches(msg, r.bindings.EndOfLine):
			maxWidth := 0
			for line := range strings.SplitSeq(r.tablesMetadata[r.activeTab].View().Content, "\n") {
				w := lipgloss.Width(line)
				if w > maxWidth {
					maxWidth = w
				}
			}

			maxOffset := maxWidth - r.viewport.Width()

			if maxOffset < 0 {
				maxOffset = 0
			}

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
	case querySuccessMsg:
		r.clearTables()
		r.setupTables()
		r.tabs = []string{"Data", "Columns", "Indexes", "Constraints"}
		r.activeTab = 0
		tableContentColumns, tableContentRows := populateTable(msg.columns, msg.rows)
		if tablePanel, ok := r.tablesMetadata[0].(*TablePanel); ok {
			tablePanel.table.SetColumns(tableContentColumns)
			tablePanel.table.SetRows(tableContentRows)
			r.tablesMetadata[0] = tablePanel
			r.viewport.SetContent(r.tablesMetadata[0].View().Content)
			r.viewport.GotoTop()
		}
		return r, nil
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
	var renderedTabs []string

	tableBorder := darkPurple
	if r.focused {
		tableBorder = neonPurple
	}

	doc := strings.Builder{}
	s := r.tabStyles
	numTabs := len(r.tabs)
	viewportWidth := r.width

	baseWidth := viewportWidth / numTabs
	remainder := viewportWidth % numTabs

	for i, t := range r.tabs {
		tabWidth := baseWidth

		if i < remainder {
			tabWidth++
		}

		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(r.tabs)-1, i == r.activeTab

		if isActive {
			style = s.activeTab.Width(tabWidth)
			style = style.BorderForeground(neonPurple)
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

	lipgloss.JoinVertical(lipgloss.Left, row, r.viewport.View())

	styledResultSet := resultSetStyle.BorderForeground(tableBorder).Width(r.width).Height(r.height).UnsetBorderTop()

	doc.WriteString(row)
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

func newTablePanel(height, width int) *TablePanel {
	t := table.New(
		table.WithFocused(true),
		table.WithWidth(width-2),
		table.WithHeight(height-2),
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
