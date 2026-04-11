package bubbletui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/davecgh/go-spew/spew"
)

type ResultSet struct {
	focused       bool
	tabs          []string
	activeTab     int
	width, height int
	tabStyles     *tabStyles

	bindings *command.TUIKeyBindings

	viewport       viewport.Model
	tablesMetadata []table.Model
	dump           io.Writer
}

func NewResultSet(kb *command.TUIKeyBindings) ResultSet {
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
		viewport: viewport.New(0, 0),
		dump:     dump,
	}

	rs.tabStyles = newTabStyles()
	rs.setupTable()

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

	r.viewport.Width = w - 4
	r.viewport.Height = h
}

func (r *ResultSet) setupTable() {
	columns := setupTable(r.height - 2)
	data := setupTable(r.height - 2)
	constraints := setupTable(r.height)
	indexes := setupTable(r.height)
	r.tablesMetadata = []table.Model{
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, r.bindings.NextTab):
			r.activeTab = min(r.activeTab+1, len(r.tabs)-1)
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View())
			return r, nil
		case key.Matches(msg, r.bindings.PrevTab):
			r.activeTab = max(r.activeTab-1, 0)
			r.viewport.SetContent(r.tablesMetadata[r.activeTab].View())
			return r, nil
		case key.Matches(msg, r.bindings.BeginningOfLine):
			r.viewport.SetXOffset(0)
			return r, nil
		case key.Matches(msg, r.bindings.EndOfLine):
			maxWidth := 0
			for line := range strings.SplitSeq(r.tablesMetadata[r.activeTab].View(), "\n") {
				w := lipgloss.Width(line)
				if w > maxWidth {
					maxWidth = w
				}
			}

			maxOffset := maxWidth - r.viewport.Width

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
		r.viewport.SetContent(r.tablesMetadata[r.activeTab].View())
		cmds = append(cmds, cmd)
	case queryErrMsg:
		errorText := fmt.Sprintf("❌ QUERY FAILED\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	case querySuccessMsg:
		r.clearTables()
		tableContentColumns, tableContentRows := populateTable(msg.columns, msg.rows)
		r.tablesMetadata[0].SetColumns(tableContentColumns)
		r.tablesMetadata[0].SetRows(tableContentRows)
		r.viewport.SetContent(r.tablesMetadata[0].View())
		r.viewport.GotoTop()
		return r, nil
	case metadataSuccessMsg:
		r.updateTableMetadataOnChange(msg.metadata)
		r.viewport.SetContent(r.tablesMetadata[r.activeTab].View())
		r.viewport.GotoTop()
		return r, nil
	case metadataErrMsg:
		errorText := fmt.Sprintf("❌ failed to get table metadata\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	case tablesFetchError:
		errorText := fmt.Sprintf("❌ failed to retrieve the current database tables\n\n%s", msg.err.Error())
		styledError := errorStyle.Render(errorText)
		r.viewport.SetContent(styledError)
		r.viewport.GotoTop()
		return r, nil
	}

	return r, tea.Batch(cmds...)
}

func (r ResultSet) View() string {
	var renderedTabs []string

	tableBorder := darkPurple
	if r.focused {
		tableBorder = neonPurple
	}

	doc := strings.Builder{}
	s := r.tabStyles
	numTabs := len(r.tabs)
	viewportWidth := r.width - 6

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
	return doc.String()
}

func (r *ResultSet) clearTables() {
	for i := range r.tablesMetadata {
		r.tablesMetadata[i] = setupTable(r.height)
	}
}

// updateTableMetadataOnChange method is used to print the table metadata retrieved asynchronously.
func (r *ResultSet) updateTableMetadataOnChange(metadata *client.Metadata) {
	if metadata != nil {
		r.clearTables()

		// table data.
		tableContentColumns, tableContentRows := populateTable(metadata.TableContent.Columns, metadata.TableContent.Rows)
		r.tablesMetadata[0].SetColumns(tableContentColumns)
		r.tablesMetadata[0].SetRows(tableContentRows)

		// table columns.
		tableStructureColumns, tableStructureRows := populateTable(metadata.Structure.Columns, metadata.Structure.Rows)
		r.tablesMetadata[1].SetColumns(tableStructureColumns)
		r.tablesMetadata[1].SetRows(tableStructureRows)

		// table indexes.
		tableIndexColumns, tableIndexRows := populateTable(metadata.Indexes.Columns, metadata.Indexes.Rows)
		r.tablesMetadata[2].SetColumns(tableIndexColumns)
		r.tablesMetadata[2].SetRows(tableIndexRows)

		// table constraints.
		tableConstraintsColumns, tableConstraintsRows := populateTable(metadata.Constraints.Columns, metadata.Constraints.Rows)
		r.tablesMetadata[3].SetColumns(tableConstraintsColumns)
		r.tablesMetadata[3].SetRows(tableConstraintsRows)
	}
}
