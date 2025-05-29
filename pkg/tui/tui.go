package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/common-nighthawk/go-figure"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
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
	queriesAreaTitle     = "SQL query"
	tablesListTitle      = "Tables"
	databaseCatalogTitle = "Database Catalog"
)

// Tui struct is the main struct when it comes to manage the UI.
// It is composed by a pointer to a tview application and the database client responsible for making queries to the database.
// Tha last field is a AppWidgets instance, so every widget can be accessed through the Tui's reference.
type Tui struct {
	app      *tview.Application
	c        *client.Client
	aw       AppWidgets
	bindings *command.TUIKeyBindings
}

// Option is a functional option type that allows us to configure the Tui object.
type Option func(*Tui)

// WithClient adds and optional database client to the Tui.
func WithClient(c *client.Client) Option {
	return func(t *Tui) {
		t.c = c
	}
}

// WithKeyBinding sets a TUIKeyBindings to the Tui struct.
func WithKeyBinding(kb *command.TUIKeyBindings) Option {
	return func(t *Tui) {
		t.bindings = kb
	}
}

// AppWidgets struct holds the widgets needed to run the app and manage the multiple behaviors supported.
// This is done this way, because some widgets make refernce to others when they get focused, clicked, etc.
// Besides, tview always returns pointers from its constructor functions.
type AppWidgets struct {
	queries            *tview.TextArea
	structure          *tview.Table
	content            *tview.Table
	constraints        *tview.Table
	indexes            *tview.Table
	errorView          *tview.TextView
	banner             *tview.TextView
	tables             *tview.List
	tableMetadata      *tview.Pages
	catalogPage        *tview.Pages
	databaseCatalog    *tview.TreeView
	leftSideFlex       *tview.Flex
	rightSideFlex      *tview.Flex
	mainViewFlex       *tview.Flex
	activeDatabaseText *tview.TextView
	appFlex            *tview.Flex
}

// New is a constructor that returns a pointer to a Tui struct.
// The functions starts the app, initializes an AppWidgets and runs the prepare function,
// responsible for setting up the whole app and its multiple behaviors.
func New(options ...Option) (*Tui, error) {
	t := &Tui{}

	for _, opt := range options {
		opt(t)
	}

	app := tview.NewApplication()

	t.app = app
	t.aw = AppWidgets{}

	if err := t.prepare(); err != nil {
		return nil, err
	}

	return t, nil
}

// setupQueries function sets up the queries text area, the widget responsible for receiving the text input from the user,
// and call the Query method from the database client.
func (t *Tui) setupQueries() {
	t.aw.queries = tview.NewTextArea().SetPlaceholder("Enter your query here...")
	t.aw.queries.SetTitle(queriesAreaTitle).SetBorder(true)

	t.aw.queries.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// Runs the query when the designated Key to run queries gets pressed.
		case t.bindings.RunQuery:
			// Clears the views and switches to the content table so the user can see the result set.
			t.aw.content.Clear()
			t.aw.structure.Clear()
			t.aw.constraints.Clear()
			t.aw.indexes.Clear()

			t.aw.tableMetadata.SwitchToPage(contentPage)

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
					tview.NewTableCell(tc).
						SetAlign(tview.AlignCenter).
						SetTextColor(tcell.ColorYellow),
				)
			}

			for i, sr := range resultSet {
				for j, sc := range sr {
					if i == 0 {
						t.aw.content.SetCell(
							i+1,
							j,
							tview.NewTableCell(sc).SetTextColor(tcell.ColorRed),
						)
					} else {
						t.aw.content.SetCell(i+1, j, tview.NewTableCell(sc).SetMaxWidth(0))
					}
				}
			}

			t.aw.content.ScrollToBeginning()
			t.aw.content.Select(0, 0)
			t.app.SetFocus(t.aw.content)

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
		case t.bindings.Navigation.Down:
			// switch to the tableMetadata page if designated Navigation Down Key gets pressed.
			t.app.SetFocus(t.aw.tableMetadata)
		case t.bindings.Navigation.Left:
			// switch to the list of tables page if designated Navigation Left Key gets pressed.
			t.app.SetFocus(t.aw.catalogPage)
			return nil
		case t.bindings.ClearEditor:
			t.aw.queries.SetText("", true)
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
	t.aw.content.SetSelectable(true, true)
	t.aw.content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, col := t.aw.content.GetSelection()
			cell := t.aw.content.GetCell(row, col)
			if cell != nil {
				clipboard.WriteAll(cell.Text)
			}
		}
		return event
	})

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
		case t.bindings.Navigation.Left:
			t.app.SetFocus(t.aw.catalogPage)
		case t.bindings.Navigation.Up:
			t.app.SetFocus(t.aw.queries)
		case t.bindings.Structure:
			if name == contentPage {
				t.aw.tableMetadata.SwitchToPage(structurePage)
				t.aw.structure.ScrollToBeginning()
			} else {
				t.aw.tableMetadata.SwitchToPage(contentPage)
				t.aw.content.ScrollToBeginning()
			}
		case t.bindings.Indexes:
			if name == contentPage {
				t.aw.tableMetadata.SwitchToPage(indexesPage)
				t.aw.indexes.ScrollToBeginning()
			} else {
				t.aw.tableMetadata.SwitchToPage(contentPage)
				t.aw.content.ScrollToBeginning()
			}
		case t.bindings.Constraints:
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

// setupDatabaseCatalog function sets up the database catalog.
func (t *Tui) setupDatabaseCatalog() error {
	root := tview.NewTreeNode("db")
	t.aw.databaseCatalog = tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	t.aw.tables = tview.NewList()
	t.aw.tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			t.aw.tables.Clear()
			t.aw.structure.Clear()
		})

	t.aw.catalogPage = tview.NewPages().
		AddPage(tablesListTitle, t.aw.tables, true, !t.c.ShowDataCatalog()).
		AddPage(databaseCatalogTitle, t.aw.databaseCatalog, true, t.c.ShowDataCatalog())

	t.aw.catalogPage.SetBorder(true).SetBorderColor(tcell.ColorPurple)

	if t.c.ShowDataCatalog() {
		t.aw.catalogPage.SetTitle(databaseCatalogTitle)
	} else {
		t.aw.catalogPage.SetTitle(tablesListTitle)
	}

	// Get the list of thables available for the current user.
	if t.c.ShowDataCatalog() {
		dbs, err := t.c.ShowDatabases()
		if err != nil {
			return err
		}

		for _, db := range dbs {
			node := tview.NewTreeNode(db).
				SetReference("database").
				SetSelectable(true)
			root.AddChild(node)
		}

		t.aw.databaseCatalog.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case t.bindings.Navigation.Right:
				t.app.SetFocus(t.aw.queries)
			}
			return event
		})

		t.aw.databaseCatalog.SetSelectedFunc(func(node *tview.TreeNode) {
			reference := node.GetReference()
			if reference == nil {
				// Selecting the root node does nothing.
				return
			}

			kind := reference.(string)

			switch kind {
			case "table":
				t.updateTableMetadataOnChange(node.GetText())
			case "database":
				children := node.GetChildren()

				databaseName := node.GetText()
				t.c.SetActiveDatabase(databaseName)

				t.aw.activeDatabaseText.SetText(
					fmt.Sprintf("[purple]Active database: [orange]%s", databaseName),
				)

				if len(children) == 0 {
					tables, err := t.c.ShowTablesPerDB(databaseName)
					if err != nil {
						return
					}

					for _, t := range tables {
						node.AddChild(
							tview.NewTreeNode(t).SetReference("table").SetSelectable(true),
						)
					}
				} else {
					// Collapse if visible, expand if collapsed.
					node.SetExpanded(!node.IsExpanded())
				}
			}
		})

	} else {
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
			case t.bindings.Navigation.Right:
				t.app.SetFocus(t.aw.queries)
			case tcell.KeyEnter:
				t.updateTableMetadataOnChange("")
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
	}

	return nil
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
			&tview.TableCell{Text: tc, Align: tview.AlignCenter, Color: tcell.ColorOrange},
		)
	}

	for i, sr := range m.TableContent.Rows {
		for j, sc := range sr {
			t.aw.content.SetCellSimple(i+1, j, sc)
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

	t.app.SetFocus(t.aw.content)
	t.aw.content.ScrollToBeginning()
	t.aw.content.Select(0, 0)
	t.aw.structure.ScrollToBeginning()
	t.aw.indexes.ScrollToBeginning()
	t.aw.constraints.ScrollToBeginning()
}

// setUpFlexBoxes function sets up the flex boxes needed to compose the app.
func (t *Tui) setUpFlexBoxes() {
	t.aw.leftSideFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.banner, 0, 1, false).
		AddItem(t.aw.catalogPage, 0, 5, true)

	t.aw.rightSideFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.queries, 0, 1, false).
		AddItem(t.aw.tableMetadata, 0, 3, false)

	// Create the main layout.
	t.aw.mainViewFlex = tview.NewFlex().
		AddItem(t.aw.leftSideFlex, 0, 1, true).AddItem(t.aw.rightSideFlex, 0, 4, false)

	helpInfo := tview.NewTextView().SetDynamicColors(true)
	helpStr := fmt.Sprintf(
		"[green]%s",
		"Press Ctrl-C to exit, press Ctrl-(vim motions) to move between panels, press Ctrl-L to execute queries, press j/k to scroll on tables panel",
	)
	fmt.Fprintln(helpInfo, helpStr)

	t.aw.activeDatabaseText = tview.NewTextView().SetDynamicColors(true)

	footer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(helpInfo, 0, 4, false).
				AddItem(t.aw.activeDatabaseText, 0, 1, false),
			1, 1, false,
		)

	t.aw.appFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.aw.mainViewFlex, 0, 20, true).
		AddItem(footer, 0, 1, false)
}

func (t *Tui) prepare() error {
	t.setupQueries()
	t.setupTablesMetadata()
	t.setupBanner()
	if err := t.setupDatabaseCatalog(); err != nil {
		return err
	}
	t.setUpFlexBoxes()

	t.app.SetRoot(t.aw.appFlex, true).EnableMouse(true).SetFocus(t.aw.queries)

	return nil
}

func (t *Tui) Run() error {
	if err := t.app.Run(); err != nil {
		return err
	}

	return nil
}
