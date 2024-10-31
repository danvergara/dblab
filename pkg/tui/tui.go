package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/danvergara/dblab/pkg/client"
)

const (
	// Tables' metadata pages.
	columnsPage     = "columns"
	structurePage   = "structure"
	indexesPage     = "indexes"
	constraintsPage = "constraints"

	// Tables' metadata page titles.
	columnsPageTitle     = "Columns"
	structurePageTitle   = "Structure"
	indexesPageTitle     = "Indexes"
	constraintsPageTitle = "Constraints"

	// Titles.
	queriesAreaTitle = "SQL query"
	tablesListTitle  = "tables"
)

type Tui struct {
	app   *tview.Application
	c     *client.Client
	pages *tview.Pages
}

func New(c *client.Client) (*Tui, error) {
	// Start the application.
	app := tview.NewApplication()

	t := Tui{
		app: app,
		c:   c,
	}

	if err := t.prepare(); err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Tui) prepare() error {
	queries := tview.NewTextArea().SetPlaceholder("Enter your query here...")
	queries.SetTitle(queriesAreaTitle).SetBorder(true)

	// Tables metadata.
	structure := tview.NewTable().SetBorders(true)
	structure.SetBorder(true).SetTitle(structurePageTitle)

	columns := tview.NewTable().SetBorders(true)
	columns.SetBorder(true).SetTitle(columnsPageTitle)

	constraints := tview.NewTable().SetBorders(true)
	constraints.SetBorder(true).SetTitle(constraintsPageTitle)

	indexes := tview.NewTable().SetBorders(true)
	indexes.SetBorder(true).SetTitle(indexesPageTitle)

	tables := tview.NewList()
	tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			tables.Clear()
			structure.Clear()
		})
	tables.SetBorder(true).SetTitle(tablesListTitle)

	tableMetadata := tview.NewPages().
		AddPage(columnsPage, columns, true, true).
		AddPage(structurePage, structure, true, false).
		AddPage(constraintsPage, constraints, true, false).
		AddPage(indexesPage, indexes, true, false)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(queries, 0, 1, false).
		AddItem(tableMetadata, 0, 3, false)

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(tables, 0, 1, true).AddItem(rightFlex, 0, 3, false)

	// Set up the pages and show the dblab flexbox.
	t.pages = tview.NewPages().
		AddPage("dblab", flex, true, true)

	tables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL:
			t.app.SetFocus(queries)
		}

		return event
	})

	tableMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		name, _ := tableMetadata.GetFrontPage()

		switch event.Key() {
		case tcell.KeyCtrlH:
			t.app.SetFocus(tables)
		case tcell.KeyCtrlK:
			t.app.SetFocus(queries)
		case tcell.KeyCtrlS:
			if name == columnsPage {
				tableMetadata.SwitchToPage(structurePage)
				structure.ScrollToBeginning()
			} else {
				tableMetadata.SwitchToPage(columnsPage)
				columns.ScrollToBeginning()
			}
		case tcell.KeyCtrlI:
			if name == columnsPage {
				tableMetadata.SwitchToPage(indexesPage)
				indexes.ScrollToBeginning()
			} else {
				tableMetadata.SwitchToPage(columnsPage)
				columns.ScrollToBeginning()
			}
		case tcell.KeyCtrlT:
			if name == columnsPage {
				tableMetadata.SwitchToPage(constraintsPage)
				constraints.ScrollToBeginning()
			} else {
				tableMetadata.SwitchToPage(columnsPage)
				columns.ScrollToBeginning()
			}
		}

		return event
	})

	// When the user navigates to a table, show its columns.
	tables.SetChangedFunc(func(i int, tableName string, st string, s rune) {
		m, err := t.c.Metadata(tableName)
		if err != nil {
			panic(err)
		}

		columns.Clear()
		structure.Clear()
		indexes.Clear()
		constraints.Clear()

		columns.ScrollToBeginning()
		structure.ScrollToBeginning()
		indexes.ScrollToBeginning()
		constraints.ScrollToBeginning()

		tableMetadata.SwitchToPage(columnsPage)
		for i, tc := range m.TableContent.Columns {
			columns.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range m.TableContent.Rows {
			for j, sc := range sr {
				if i == 0 {
					columns.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					columns.SetCellSimple(i+1, j, sc)
				}
			}
		}

		for i, tc := range m.Structure.Columns {
			structure.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range m.Structure.Rows {
			for j, sc := range sr {
				if i == 0 {
					structure.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					structure.SetCellSimple(i+1, j, sc)
				}
			}
		}

		for i, tc := range m.Indexes.Columns {
			indexes.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range m.Indexes.Rows {
			for j, sc := range sr {
				if i == 0 {
					indexes.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					indexes.SetCellSimple(i+1, j, sc)
				}
			}
		}

		for i, tc := range m.Constraints.Columns {
			constraints.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range m.Constraints.Rows {
			for j, sc := range sr {
				if i == 0 {
					constraints.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					constraints.SetCellSimple(i+1, j, sc)
				}
			}
		}

	})

	queries.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			columns.Clear()
			structure.Clear()
			constraints.Clear()
			indexes.Clear()

			if err := t.c.ResetPagination(); err != nil {
				panic(err)
			}

			query := queries.GetText()
			resultSet, columnNames, err := t.c.Query(query)
			if err != nil {
				panic(err)
			}

			for i, tc := range columnNames {
				columns.SetCell(
					0,
					i,
					&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
				)
			}

			for i, sr := range resultSet {
				for j, sc := range sr {
					if i == 0 {
						columns.SetCell(
							i+1,
							j,
							&tview.TableCell{Text: sc, Color: tcell.ColorRed},
						)
					} else {
						columns.SetCellSimple(i+1, j, sc)
					}
				}
			}
			// tableMetadata.SwitchToPage(columnsPage)
		case tcell.KeyCtrlJ:
			t.app.SetFocus(tableMetadata)
		case tcell.KeyCtrlH:
			t.app.SetFocus(tables)
		}
		return event
	})

	ts, err := t.showTables()
	if err != nil {
		return err
	}

	tables.Clear()

	for _, ta := range ts {
		tables.AddItem(ta, "", 0, nil)
	}

	// Trigger the initial selection.
	tables.SetCurrentItem(0)

	t.app.SetRoot(t.pages, true).EnableMouse(true).SetFocus(tables)

	return nil
}

func (t *Tui) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}

	return nil
}
