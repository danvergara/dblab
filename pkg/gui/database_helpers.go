package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
)

// showTables list all the tables in the database on the tables panel.
func (gui *Gui) showTables() error {
	var query string

	switch gui.c.Driver() {
	case "postgres":
		fallthrough
	case "postgresql":
		query = `
		SELECT
			table_name
		FROM
			information_schema.tables
		WHERE
			table_schema = 'public'
		ORDER BY
			table_name;`
	case "mysql":
		query = "SHOW TABLES;"
	}

	rows, err := gui.c.DB().Queryx(query)
	if err != nil {
		return err
	}

	rv, err := gui.g.View("tables")
	if err != nil {
		return err
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}

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

// query returns performs the query and returns the result set and the colum names.
func (gui *Gui) query(q string) ([][]string, []string, error) {
	resultSet := [][]string{}

	// Runs the query extracting the content of the view calling the Buffer method.
	rows, err := gui.c.DB().Queryx(q)
	if err != nil {
		return nil, nil, err
	}

	// Gets the names of the columns of the result set.
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		// cols is an []interface{} of all of the column results.
		cols, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}

		// Convert []interface{} into []string.
		s := make([]string, len(cols))
		for i, v := range cols {
			s[i] = fmt.Sprint(v)
		}

		resultSet = append(resultSet, s)
	}

	return resultSet, columnNames, nil
}

// runQuery run the introduced query in the query panel.
func (gui *Gui) inputQuery() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		// Cleans the view.
		v.Rewind()

		resultSet, columnNames, err := gui.query(v.Buffer())
		if err != nil {
			return err
		}

		ov, err := gui.g.View("rows")
		if err != nil {
			return err
		}

		// Cleans the view.
		ov.Rewind()
		ov.Clear()

		renderTable(ov, columnNames, resultSet)

		return nil
	}
}

// selectTable perfoms a select statement based on the selected table.
func (gui *Gui) selectTable(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	t, err := v.Line(cy)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM %s;", t)
	resultSet, columnNames, err := gui.query(query)
	if err != nil {
		return err
	}

	ov, err := gui.g.View("rows")
	if err != nil {
		return err
	}

	// Cleans the view.
	ov.Rewind()
	ov.Clear()

	renderTable(ov, columnNames, resultSet)

	return nil
}
