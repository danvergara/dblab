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

// runQuery run the introduced query in the query panel.
func (gui *Gui) runQuery() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		resultSet := [][]string{}

		// Cleans the view.
		v.Rewind()

		// Runs the query extracting the content of the view calling the Buffer method.
		rows, err := gui.c.DB().Queryx(v.Buffer())
		if err != nil {
			return err
		}

		// Gets the names of the columns of the result set.
		columnNames, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			// cols is an []interface{} of all of the column results.
			cols, err := rows.SliceScan()
			if err != nil {
				return err
			}

			// Convert []interface{} into []string.
			s := make([]string, len(cols))
			for i, v := range cols {
				s[i] = fmt.Sprint(v)
			}

			resultSet = append(resultSet, s)
		}

		ov, err := gui.g.View("rows")
		if err != nil {
			return err
		}

		// Cleans the view.
		ov.Rewind()
		ov.Clear()

		// Setup the table.
		table := tablewriter.NewWriter(ov)
		table.SetCenterSeparator("|")
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetHeader(columnNames)
		// Add Bulk Data.
		table.AppendBulk(resultSet)
		table.Render()

		return nil
	}
}

func (gui *Gui) selectTable(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	t, err := v.Line(cy)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT * FROM %s;", t)

	rows, err := gui.c.DB().Queryx(query)
	if err != nil {
		return err
	}

	// Gets the names of the columns of the result set.
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	resultSet := [][]string{}

	for rows.Next() {
		// cols is an []interface{} of all of the column results.
		cols, err := rows.SliceScan()
		if err != nil {
			return err
		}

		// Convert []interface{} into []string.
		s := make([]string, len(cols))
		for i, v := range cols {
			s[i] = fmt.Sprint(v)
		}

		resultSet = append(resultSet, s)
	}

	ov, err := gui.g.View("rows")
	if err != nil {
		return err
	}

	// Cleans the view.
	ov.Rewind()
	ov.Clear()

	// Setup the table.
	table := tablewriter.NewWriter(ov)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	table.SetHeader(columnNames)
	// Add Bulk Data.
	table.AppendBulk(resultSet)
	table.Render()
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		l, err := v.Line(cy + 1)
		if err != nil {
			return err
		}
		if l != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
