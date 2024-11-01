package tui

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
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
	errorPage       = "error"

	// Tables' metadata page titles.
	columnsPageTitle     = "Columns"
	structurePageTitle   = "Structure"
	indexesPageTitle     = "Indexes"
	constraintsPageTitle = "Constraints"
	errorPageTitle       = "error"

	// Titles.
	queriesAreaTitle = "SQL query"
	tablesListTitle  = "tables"
)

type Tui struct {
	app *tview.Application
	c   *client.Client
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
	errorView := tview.NewTextView().SetDynamicColors(true)
	errorView.SetTitle(errorPageTitle).SetBorder(true)

	banner := tview.NewTextView().SetDynamicColors(true)
	banner.SetBorderColor(tcell.ColorGreen)
	banner.SetBorder(true)

	dblabFigure := figure.NewFigure("dblab", "", true)
	coloredBannerContent := fmt.Sprintf("[purple]%s", dblabFigure.String())
	fmt.Fprintln(banner, coloredBannerContent)

	tables := tview.NewList()
	tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			tables.Clear()
			structure.Clear()
		})
	tables.SetBorder(true).SetTitle(tablesListTitle).SetBorderColor(tcell.ColorPurple)

	tableMetadata := tview.NewPages().
		AddPage(columnsPage, columns, true, true).
		AddPage(structurePage, structure, true, false).
		AddPage(constraintsPage, constraints, true, false).
		AddPage(indexesPage, indexes, true, false).
		AddPage(errorPage, errorView, true, false)

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(banner, 0, 1, false).
		AddItem(tables, 0, 5, true)

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(queries, 0, 1, false).
		AddItem(tableMetadata, 0, 3, false)

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(leftFlex, 0, 1, true).AddItem(rightFlex, 0, 4, false)

	pagination := tview.NewTextView().SetText(" 1 / 10")

	prevButton := tview.NewButton("<").SetSelectedFunc(func() {
		t.app.SetFocus(tables)

		currentTable, page, err := t.c.PreviousPage()
		if err != nil {
			return
		}

		columns.Clear()

		totalPages := t.c.TotalPages()
		pagination.SetText(fmt.Sprintf("%4d / %4d", page, totalPages))

		for i, tc := range currentTable.Columns {
			columns.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range currentTable.Rows {
			for j, sc := range sr {
				if i == 0 {
					columns.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					columns.SetCellSimple(i+1, j, sc)
				}
			}
		}
	})

	nextButton := tview.NewButton(">").SetSelectedFunc(func() {
		t.app.SetFocus(tables)

		currentTable, page, err := t.c.NextPage()
		if err != nil {
			return
		}

		columns.Clear()

		totalPages := t.c.TotalPages()
		pagination.SetText(fmt.Sprintf("%4d / %4d", page, totalPages))

		for i, tc := range currentTable.Columns {
			columns.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range currentTable.Rows {
			for j, sc := range sr {
				if i == 0 {
					columns.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					columns.SetCellSimple(i+1, j, sc)
				}
			}
		}
	})

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(prevButton, 10, 1, false). // Adjust width as needed
		AddItem(nil, 2, 0, false).         // Spacer with width 2
		AddItem(pagination, 15, 1, false). // Adjust width as needed
		// AddItem(nil, 2, 0, false).         // Spacer with width 2
		AddItem(nextButton, 10, 1, false)

	helpInfo := tview.NewTextView().SetDynamicColors(true)
	helpStr := fmt.Sprintf(
		"[green]%s",
		"Press Ctrl-C to exit, press Ctrl-(vim motions) to move between panels, press Ctrl-L to execute queries, press j/k to scroll on tables panel",
	)
	fmt.Fprintln(helpInfo, helpStr)

	footer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(buttonFlex, 0, 1, false).
				AddItem(helpInfo, 0, 4, false),
			1, 1, false,
		)

	appFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(flex, 0, 20, true).
		AddItem(footer, 0, 1, false)

	tables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL:
			t.app.SetFocus(queries)
		}

		switch event.Rune() {
		// Use 'j' to move down.
		case 'j':
			tables.SetCurrentItem(tables.GetCurrentItem() + 1)
			return nil
			// Use 'k' to move up.
		case 'k':
			tables.SetCurrentItem(tables.GetCurrentItem() - 1)
			return nil
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

	tables.SetFocusFunc(func() {
		columns.Clear()
		structure.Clear()
		indexes.Clear()
		constraints.Clear()
		tableName, _ := tables.GetItemText(tables.GetCurrentItem())
		m, err := t.c.Metadata(tableName)
		if err != nil {
			errorView.Clear()
			errorMsg := fmt.Sprintf("[red]%s", err.Error())
			fmt.Fprintln(errorView, errorMsg)
			tableMetadata.SwitchToPage(errorPage)
			return
		}

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

		pagination.SetText(fmt.Sprintf("%4d / %4d", 1, m.TotalPages))
	})

	tables.SetChangedFunc(func(i int, tableName string, st string, s rune) {
		columns.Clear()
		structure.Clear()
		indexes.Clear()
		constraints.Clear()

		m, err := t.c.Metadata(tableName)
		if err != nil {
			errorView.Clear()
			errorMsg := fmt.Sprintf("[red]%s", err.Error())
			fmt.Fprintln(errorView, errorMsg)
			tableMetadata.SwitchToPage(errorPage)
			return
		}

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

		pagination.SetText(fmt.Sprintf("%4d / %4d", 1, m.TotalPages))
	})

	queries.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			columns.Clear()
			structure.Clear()
			constraints.Clear()
			indexes.Clear()

			tableMetadata.SwitchToPage(columnsPage)

			if err := t.c.ResetPagination(); err != nil {
				errorView.Clear()
				errorMsg := fmt.Sprintf("[red]%s", err.Error())
				fmt.Fprintln(errorView, errorMsg)

				tableMetadata.SwitchToPage(errorPage)
				return event
			}

			query := queries.GetText()

			resultSet, columnNames, err := t.c.Query(query)
			if err != nil {
				errorView.Clear()
				errorMsg := fmt.Sprintf("[red]%s", err.Error())
				fmt.Fprintln(errorView, errorMsg)

				tableMetadata.SwitchToPage(errorPage)

				return event
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

	t.app.SetRoot(appFlex, true).EnableMouse(true).SetFocus(tables)

	return nil
}

func (t *Tui) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}

	return nil
}
