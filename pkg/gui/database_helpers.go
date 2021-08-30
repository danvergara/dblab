package gui

import (
	"fmt"

	"github.com/danvergara/gocui"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

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

// render renders a the resultset from a given function using a table name as input.
// This method receives a function as parameter to call it in the body passing the
// table name as its parameter. The output is a result set as a product of an
// expected query.
func (gui *Gui) render(
	g *gocui.Gui,
	v *gocui.View,
	viewName string,
	query func(string) ([][]string, []string, error),
) error {
	v.Rewind()

	_, cy := v.Cursor()

	t, err := v.Line(cy)
	if err != nil {
		return err
	}

	// f is the function to be executed to get a result set.
	// This gives flexibility to use whaterever query method
	// from the client, without the need to know what query is
	// being executed.
	resultSet, columnNames, err := query(t)
	if err != nil {
		return err
	}

	ov, err := gui.g.View(viewName)
	if err != nil {
		return err
	}

	// Cleans the view.
	ov.Clear()

	if err := ov.SetCursor(0, 0); err != nil {
		return err
	}

	renderTable(ov, columnNames, resultSet)

	return nil
}

// metadata get the metadata from a given table name.
func (gui *Gui) metadata(g *gocui.Gui, v *gocui.View) error {
	var viewQueries = []struct {
		name  string
		query func(string) ([][]string, []string, error)
	}{
		{"rows", gui.c.TableContent},
		{"structure", gui.c.TableStructure},
		{"constraints", gui.c.Constraints},
		{"indexes", gui.c.Indexes},
	}

	for _, vq := range viewQueries {
		if err := gui.render(g, v, vq.name, vq.query); err != nil {
			return err
		}
	}

	return nil
}
