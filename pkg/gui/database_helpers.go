package gui

import (
	"fmt"

	"github.com/danvergara/gocui"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// runQuery run the introduced query in the query panel.
func (gui *Gui) inputQuery() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		v.Rewind()

		ov, err := gui.g.View("rows")
		if err != nil {
			return err
		}

		// Cleans up the rows view.
		ov.Clear()

		resultSet, columnNames, err := gui.c.Query(v.Buffer())
		if err != nil {
			// Prints the error in red on the rows view.
			red := color.New(color.FgRed)
			boldRed := red.Add(color.Bold)
			boldRed.Fprintf(ov, "%s\n", err)
		} else {
			renderTable(ov, columnNames, resultSet)
		}

		return nil
	}
}

// metadata get the metadata from a given table name.
func (gui *Gui) metadata(g *gocui.Gui, v *gocui.View) error {
	v.Rewind()

	_, cy := v.Cursor()

	t, err := v.Line(cy)
	if err != nil {
		return err
	}

	m, err := gui.c.Metadata(t)
	if err != nil {
		return err
	}

	var viewQueries = []struct {
		name    string
		columns []string
		rows    [][]string
	}{
		{"rows", m.TableContent.Columns, m.TableContent.Rows},
		{"structure", m.Structure.Columns, m.Structure.Rows},
		{"constraints", m.Constraints.Columns, m.Constraints.Rows},
		{"indexes", m.Indexes.Columns, m.Indexes.Rows},
	}

	for _, vq := range viewQueries {
		if err := gui.render(g, v, vq.name, vq.columns, vq.rows); err != nil {
			return err
		}
	}

	return nil
}

// showTables list all the tables in the database on the tables panel.
func (gui *Gui) showTables() error {
	tables, err := gui.c.ShowTables()
	if err != nil {
		return err
	}

	rv, err := gui.g.View("tables")
	if err != nil {
		return err
	}

	for _, table := range tables {
		fmt.Fprintf(rv, "%s\n", table)
	}

	return nil
}

// renderTable renders the result set as a table in the terminal output.
func renderTable(v *gocui.View, columns []string, resultSet [][]string) {
	table := tablewriter.NewWriter(v)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	table.SetHeader(columns)
	// Add Bulk Data.
	table.AppendBulk(resultSet)
	table.Render()
}

// render renders the resultset from a given function using a table name as an input.
// This method receives a function as parameter to call it in the body, passing along
// the table name as its parameter. The output is a result set as a product of the
// expected query.
func (gui *Gui) render(
	g *gocui.Gui,
	v *gocui.View,
	viewName string,
	columns []string,
	resultSet [][]string,
) error {
	v.Rewind()

	ov, err := gui.g.View(viewName)
	if err != nil {
		return err
	}

	// Cleans the view.
	ov.Clear()

	if err := ov.SetCursor(0, 0); err != nil {
		return err
	}

	renderTable(ov, columns, resultSet)

	return nil
}
