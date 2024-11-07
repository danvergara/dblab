package tui

import (
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/danvergara/dblab/pkg/client"
)

const (
	// Tables' metadata pages.
	contentPage     = "content"
	structurePage   = "structure"
	indexesPage     = "indexes"
	constraintsPage = "constraints"
	errorPage       = "error"

	// Tables' metadata page titles.
	contentPageTitle     = "Table Content"
	structurePageTitle   = "Structure"
	indexesPageTitle     = "Indexes"
	constraintsPageTitle = "Constraints"
	errorPageTitle       = "error"

	// Titles.
	queriesAreaTitle = "SQL query"
	tablesListTitle  = "tables"
)

// Tui struct is the main struct when it comes to manage the UI.
// It is composed by a pointer to a tview application and the database client responsible for making queries to the database.
// Tha last field is a AppWidgets instance, so every widget can be accessed through the Tui's reference.
type Tui struct {
	app *tview.Application
	c   *client.Client
	aw  AppWidgets
}

// AppWidgets struct holds the widgets needed to run the app and manage the multiple behaviors supported.
// This is done this way, because some widgets make refernce to others when they get focused, clicked, etc.
// Besides, tview always returns pointers from its constructor functions.
type AppWidgets struct {
	queries       *tview.TextArea
	structure     *tview.Table
	content       *tview.Table
	constraints   *tview.Table
	indexes       *tview.Table
	errorView     *tview.TextView
	banner        *tview.TextView
	tables        *tview.List
	tableMetadata *tview.Pages
	leftSideFlex  *tview.Flex
	rightSideFlex *tview.Flex
	mainViewFlex  *tview.Flex
	pagination    *tview.TextView
	prevButton    *tview.Button
	nextButton    *tview.Button
	appFlex       *tview.Flex
}

// New is a constructor that returns a pointer to a Tui struct.
// The functions starts the app, initializes an AppWidgets and runs the prepare function,
// responsible for setting up the whole app and its multiple behaviors.
func New(c *client.Client) (*Tui, error) {
	app := tview.NewApplication()

	t := Tui{
		app: app,
		c:   c,
		aw:  AppWidgets{},
	}

	if err := t.prepare(); err != nil {
		return nil, err
	}

	return &t, nil
}

// setupQueries function sets up the queries text area, the widget responsible for receiving the text input from the user,
// and call the Query method from the database client.
func (t *Tui) setupQueries() {
	t.aw.queries = tview.NewTextArea().SetPlaceholder("Enter your query here...")
	t.aw.queries.SetTitle(queriesAreaTitle).SetBorder(true)

	t.aw.queries.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Runs the query when Ctrl+Space is pressed.
		case tcell.KeyCtrlSpace:
			// Clears the views and switches to the content table so the user can see the result set.
			t.aw.content.Clear()
			t.aw.structure.Clear()
			t.aw.constraints.Clear()
			t.aw.indexes.Clear()

			t.aw.tableMetadata.SwitchToPage(contentPage)

			// reset the pagination.
			if err := t.c.ResetPagination(); err != nil {
				// This is the way errors get handled and this pattern repeats multiple times accross the board.
				// Clear the error view.
				t.aw.errorView.Clear()
				// Print the error message on the error view.
				errorMsg := fmt.Sprintf("[red]%s", err.Error())
				fmt.Fprintln(t.aw.errorView, errorMsg)
				// Switch to the error view.
				t.aw.tableMetadata.SwitchToPage(errorPage)
				return event
			}

			t.aw.pagination.SetText(fmt.Sprintf("%4d / %4d", 1, 1))

			query := t.aw.queries.GetText()

			// Call the Query method from the client and populate the content page.
			resultSet, columnNames, err := t.c.Query(query)
			if err != nil {
				t.aw.errorView.Clear()
				errorMsg := fmt.Sprintf("[red]%s", err.Error())
				fmt.Fprintln(t.aw.errorView, errorMsg)

				t.aw.tableMetadata.SwitchToPage(errorPage)

				return event
			}

			for i, tc := range columnNames {
				t.aw.content.SetCell(
					0,
					i,
					&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
				)
			}

			for i, sr := range resultSet {
				for j, sc := range sr {
					if i == 0 {
						t.aw.content.SetCell(
							i+1,
							j,
							&tview.TableCell{Text: sc, Color: tcell.ColorRed},
						)
					} else {
						t.aw.content.SetCellSimple(i+1, j, sc)
					}
				}
			}

			// Update the table list if the tables get updated somehow.
			switch {
			case strings.Contains(strings.ToLower(query), "alter table"):
				fallthrough
			case strings.Contains(strings.ToLower(query), "drop table"):
				fallthrough
			case strings.Contains(strings.ToLower(query), "create table"):
				ts, err := t.c.ShowTables()
				if err != nil {
					t.aw.errorView.Clear()
					errorMsg := fmt.Sprintf("[red]%s", err.Error())
					fmt.Fprintln(t.aw.errorView, errorMsg)

					t.aw.tableMetadata.SwitchToPage(errorPage)

					return event
				}

				t.aw.tables.Clear()

				for _, ta := range ts {
					t.aw.tables.AddItem(ta, "", 0, nil)
				}

				t.aw.tables.SetCurrentItem(0)
			}
		case tcell.KeyCtrlJ:
			// switch to the tableMetadata page if Ctrl+J gets pressed.
			t.app.SetFocus(t.aw.tableMetadata)
		case tcell.KeyCtrlH:
			// switch to the list of tables page if Ctrl+H gets pressed.
			t.app.SetFocus(t.aw.tables)
		}
		return event
	})
}

// setupTablesMetadata sets up all the table's data related text views and the tableMetadata page.
func (t *Tui) setupTablesMetadata() {
	t.aw.structure = tview.NewTable().SetBorders(true)
	t.aw.structure.SetBorder(true).SetTitle(structurePageTitle)

	t.aw.content = tview.NewTable().SetBorders(true)
	t.aw.content.SetBorder(true).SetTitle(contentPageTitle)

	t.aw.constraints = tview.NewTable().SetBorders(true)
	t.aw.constraints.SetBorder(true).SetTitle(constraintsPageTitle)

	t.aw.indexes = tview.NewTable().SetBorders(true)
	t.aw.indexes.SetBorder(true).SetTitle(indexesPageTitle)

	t.aw.errorView = tview.NewTextView().SetDynamicColors(true)
	t.aw.errorView.SetTitle(errorPageTitle).SetBorder(true)

	t.aw.tableMetadata = tview.NewPages().
		AddPage(contentPage, t.aw.content, true, true).
		AddPage(structurePage, t.aw.structure, true, false).
		AddPage(constraintsPage, t.aw.constraints, true, false).
		AddPage(indexesPage, t.aw.indexes, true, false).
		AddPage(errorPage, t.aw.errorView, true, false)

		// Define how to navigate between pages.
	t.aw.tableMetadata.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		name, _ := t.aw.tableMetadata.GetFrontPage()

		switch event.Key() {
		case tcell.KeyCtrlH:
			t.app.SetFocus(t.aw.tables)
		case tcell.KeyCtrlK:
			t.app.SetFocus(t.aw.queries)
		case tcell.KeyCtrlS:
			if name == contentPage {
				t.aw.tableMetadata.SwitchToPage(structurePage)
				t.aw.structure.ScrollToBeginning()
			} else {
				t.aw.tableMetadata.SwitchToPage(contentPage)
				t.aw.content.ScrollToBeginning()
			}
		case tcell.KeyCtrlI:
			if name == contentPage {
				t.aw.tableMetadata.SwitchToPage(indexesPage)
				t.aw.indexes.ScrollToBeginning()
			} else {
				t.aw.tableMetadata.SwitchToPage(contentPage)
				t.aw.content.ScrollToBeginning()
			}
		case tcell.KeyCtrlT:
			if name == contentPage {
				t.aw.tableMetadata.SwitchToPage(constraintsPage)
				t.aw.constraints.ScrollToBeginning()
			} else {
				t.aw.tableMetadata.SwitchToPage(contentPage)
				t.aw.content.ScrollToBeginning()
			}
		}

		return event
	})
}

// setupBanner function sets up the banner that shows the name of the project at the top-left corner.
func (t *Tui) setupBanner() {
	t.aw.banner = tview.NewTextView().SetDynamicColors(true)
	t.aw.banner.SetBorderColor(tcell.ColorGreen)
	t.aw.banner.SetBorder(true)

	dblabFigure := figure.NewFigure("dblab", "", true)
	coloredBannerContent := fmt.Sprintf("[purple]%s", dblabFigure.String())
	fmt.Fprintln(t.aw.banner, coloredBannerContent)
}

// setupTablesList function sets up the table list page.
func (t *Tui) setupTablesList() error {
	t.aw.tables = tview.NewList()
	t.aw.tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			t.aw.tables.Clear()
			t.aw.structure.Clear()
		})
	t.aw.tables.SetBorder(true).SetTitle(tablesListTitle).SetBorderColor(tcell.ColorPurple)

	// Get the list of thables available for the current user.
	ts, err := t.c.ShowTables()
	if err != nil {
		return err
	}

	t.aw.tables.Clear()

	for _, ta := range ts {
		t.aw.tables.AddItem(ta, "", 0, nil)
	}

	// Trigger the initial selection.
	t.aw.tables.SetCurrentItem(0)

	// Default list navigation is done by the arrow keys, but this callback adds another one: using the 'j' and 'k'.
	// Similar to Vim.
	t.aw.tables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL:
			t.app.SetFocus(t.aw.queries)
		}

		switch event.Rune() {
		// Use 'j' to move down.
		case 'j':
			t.aw.tables.SetCurrentItem(t.aw.tables.GetCurrentItem() + 1)
			return nil
			// Use 'k' to move up.
		case 'k':
			t.aw.tables.SetCurrentItem(t.aw.tables.GetCurrentItem() - 1)
			return nil
		}

		return event
	})

	t.aw.tables.SetFocusFunc(func() {
		t.updateTableMetadataOnChange("")
	})

	t.aw.tables.SetChangedFunc(func(i int, tableName string, st string, s rune) {
		t.updateTableMetadataOnChange(tableName)
	})

	return nil
}

func (t *Tui) setupPagination() {
	t.aw.pagination = tview.NewTextView().SetText(fmt.Sprintf("%4d / %4d", 1, 1))
}

// updateTableMetadataOnChange functions updates tables' data related views on different events.
func (t *Tui) updateTableMetadataOnChange(tableName string) {
	t.aw.content.Clear()
	t.aw.structure.Clear()
	t.aw.indexes.Clear()
	t.aw.constraints.Clear()

	if tableName == "" {
		tableName, _ = t.aw.tables.GetItemText(t.aw.tables.GetCurrentItem())
	}

	// Get the constraints, indexes, schema and the content of a given table.
	m, err := t.c.Metadata(tableName)
	if err != nil {
		t.aw.errorView.Clear()
		errorMsg := fmt.Sprintf("[red]%s", err.Error())
		fmt.Fprintln(t.aw.errorView, errorMsg)
		t.aw.tableMetadata.SwitchToPage(errorPage)
		return
	}

	// Prepare the tables for the new content.
	t.aw.content.ScrollToBeginning()
	t.aw.structure.ScrollToBeginning()
	t.aw.indexes.ScrollToBeginning()
	t.aw.constraints.ScrollToBeginning()

	t.aw.tableMetadata.SwitchToPage(contentPage)

	// Populate the table.
	for i, tc := range m.TableContent.Columns {
		t.aw.content.SetCell(
			0,
			i,
			&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
		)
	}

	for i, sr := range m.TableContent.Rows {
		for j, sc := range sr {
			if i == 0 {
				t.aw.content.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
			} else {
				t.aw.content.SetCellSimple(i+1, j, sc)
			}
		}
	}

	for i, tc := range m.Structure.Columns {
		t.aw.structure.SetCell(
			0,
			i,
			&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
		)
	}

	for i, sr := range m.Structure.Rows {
		for j, sc := range sr {
			if i == 0 {
				t.aw.structure.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
			} else {
				t.aw.structure.SetCellSimple(i+1, j, sc)
			}
		}
	}

	for i, tc := range m.Indexes.Columns {
		t.aw.indexes.SetCell(
			0,
			i,
			&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
		)
	}

	for i, sr := range m.Indexes.Rows {
		for j, sc := range sr {
			if i == 0 {
				t.aw.indexes.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
			} else {
				t.aw.indexes.SetCellSimple(i+1, j, sc)
			}
		}
	}

	for i, tc := range m.Constraints.Columns {
		t.aw.constraints.SetCell(
			0,
			i,
			&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
		)
	}

	for i, sr := range m.Constraints.Rows {
		for j, sc := range sr {
			if i == 0 {
				t.aw.constraints.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
			} else {
				t.aw.constraints.SetCellSimple(i+1, j, sc)
			}
		}
	}

	// Update the paginantion text view.
	t.aw.pagination.SetText(fmt.Sprintf("%4d / %4d", 1, m.TotalPages))
}

// setUpFlexBoxes function sets up the flex boxes needed to compose the app.
func (t *Tui) setUpFlexBoxes() {
	t.aw.leftSideFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.banner, 0, 1, false).
		AddItem(t.aw.tables, 0, 5, true)

	t.aw.rightSideFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.queries, 0, 1, false).
		AddItem(t.aw.tableMetadata, 0, 3, false)

	// Create the main layout.
	t.aw.mainViewFlex = tview.NewFlex().
		AddItem(t.aw.leftSideFlex, 0, 1, true).AddItem(t.aw.rightSideFlex, 0, 4, false)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.aw.prevButton, 10, 1, false).
		AddItem(nil, 2, 0, false).
		AddItem(t.aw.pagination, 15, 1, false).
		AddItem(t.aw.nextButton, 10, 1, false)

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

	t.aw.appFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.mainViewFlex, 0, 20, true).
		AddItem(footer, 0, 1, false)
}

// setUpButtons functions denfines the button used to move to the previous or the next page on the database's content table.
func (t *Tui) setUpButtons() {
	t.aw.prevButton = tview.NewButton("<").SetSelectedFunc(func() {
		t.app.SetFocus(t.aw.tables)

		currentTable, page, err := t.c.PreviousPage()
		if err != nil {
			return
		}

		t.aw.content.Clear()

		totalPages := t.c.TotalPages()
		t.aw.pagination.SetText(fmt.Sprintf("%4d / %4d", page, totalPages))

		for i, tc := range currentTable.Columns {
			t.aw.content.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range currentTable.Rows {
			for j, sc := range sr {
				if i == 0 {
					t.aw.content.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					t.aw.content.SetCellSimple(i+1, j, sc)
				}
			}
		}
	})

	t.aw.nextButton = tview.NewButton(">").SetSelectedFunc(func() {
		t.app.SetFocus(t.aw.tables)

		currentTable, page, err := t.c.NextPage()
		if err != nil {
			return
		}

		t.aw.content.Clear()

		totalPages := t.c.TotalPages()
		t.aw.pagination.SetText(fmt.Sprintf("%4d / %4d", page, totalPages))

		for i, tc := range currentTable.Columns {
			t.aw.content.SetCell(
				0,
				i,
				&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorYellow},
			)
		}

		for i, sr := range currentTable.Rows {
			for j, sc := range sr {
				if i == 0 {
					t.aw.content.SetCell(i+1, j, &tview.TableCell{Text: sc, Color: tcell.ColorRed})
				} else {
					t.aw.content.SetCellSimple(i+1, j, sc)
				}
			}
		}
	})
}

func (t *Tui) prepare() error {
	t.setupQueries()
	t.setupTablesMetadata()
	t.setupBanner()
	if err := t.setupTablesList(); err != nil {
		return err
	}
	t.setUpButtons()
	t.setupPagination()
	t.setUpFlexBoxes()

	t.app.SetRoot(t.aw.appFlex, true).EnableMouse(true).SetFocus(t.aw.tables)

	return nil
}

func (t *Tui) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}

	return nil
}
