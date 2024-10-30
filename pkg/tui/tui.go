package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/danvergara/dblab/pkg/client"
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
	queries.SetTitle("SQL query").SetBorder(true)

	// Tables metadata.
	structure := tview.NewTable().SetBorders(true)
	structure.SetBorder(true).SetTitle("Structure")

	columns := tview.NewTable().SetBorders(true)
	columns.SetBorder(true).SetTitle("Columns")

	constraints := tview.NewTable().SetBorders(true)
	constraints.SetBorder(true).SetTitle("Constraints")

	indexes := tview.NewTable().SetBorders(true)
	indexes.SetBorder(true).SetTitle("Indexes")

	tables := tview.NewList()
	tables.ShowSecondaryText(false)
	tables.SetDoneFunc(func() {
		tables.Clear()
		structure.Clear()
	})
	tables.SetBorder(true).SetTitle("Tables")

	tableMetadata := tview.NewPages().
		AddPage("structure", structure, true, true).
		AddPage("columns", columns, true, false).
		AddPage("constraints", constraints, true, false).
		AddPage("indexes", indexes, true, false)

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
		case tcell.KeyRight:
			tableMetadata.SwitchToPage("structure")
			t.app.SetFocus(tableMetadata)
		}

		return event
	})

	tableMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			t.app.SetFocus(tables)
		case tcell.KeyCtrlS:
			tableMetadata.SwitchToPage("structure")
		case tcell.KeyCtrlI:
			tableMetadata.SwitchToPage("indexes")
		case tcell.KeyCtrlF:
			tableMetadata.SwitchToPage("constraints")
		}

		return event
	})

	// When the user navigates to a table, show its columns.
	tables.SetChangedFunc(func(i int, tableName string, st string, s rune) {
		structure.Clear()
		m, err := t.c.Metadata(tableName)
		if err != nil {
			panic(err)
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
