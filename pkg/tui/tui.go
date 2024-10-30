package tui

import (
	"github.com/gdamore/tcell/v2"
	_ "github.com/lib/pq"
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
	columns := tview.NewTable().SetBorders(true)
	columns.SetBorder(true).SetTitle("Columns")

	tables := tview.NewList()
	tables.ShowSecondaryText(false)
	tables.SetDoneFunc(func() {
		tables.Clear()
		columns.Clear()
	})
	tables.SetBorder(true).SetTitle("Tables")

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(tables, 0, 1, true).AddItem(columns, 0, 3, false)

	// Set up the pages and show the dblab flexbox.
	t.pages = tview.NewPages().
		AddPage("dblab", flex, true, true)

	t.app.SetFocus(tables)

	// When the user navigates to a table, show its columns.
	tables.SetChangedFunc(func(i int, tableName string, st string, s rune) {
		columns.Clear()
		m, err := t.c.Metadata(tableName)
		if err != nil {
			panic(err)
		}

		for i, tc := range m.Structure.Columns {
			columns.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range m.Structure.Rows {
			for j, sc := range sr {
				if i == 0 {
					columns.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					columns.SetCellSimple(i+1, j, sc)
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

	tables.SetCurrentItem(0) // Trigger the initial selection.

	t.app.SetRoot(t.pages, true).EnableMouse(true)

	return nil
}

func (t *Tui) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}

	return nil
}
