package bubbletui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Digital-Shane/treeview"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
	// "github.com/google/uuid"
)

type dbObjectType string

const (
	typeDatabase   dbObjectType = "database"
	typeSchema     dbObjectType = "schema"
	typeTable      dbObjectType = "table"
	typeConnection dbObjectType = "connection"
)

type DBObject struct {
	ID       string
	Name     string
	Type     dbObjectType
	Driver   string
	Children []DBObject
}

func dbObjectHasType(nodeType string) func(*treeview.Node[DBObject]) bool {
	return func(n *treeview.Node[DBObject]) bool {
		return string(n.Data().Type) == nodeType
	}
}

func dbObjectIsConnection(nodeType, driver string) func(*treeview.Node[DBObject]) bool {
	return func(n *treeview.Node[DBObject]) bool {
		return string(n.Data().Type) == nodeType && n.Data().Driver == driver
	}
}

// item implements the Item interface for required for the List Model from bubbles.
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

// itemDelegate is used to inject styling to the list items.
// Implements the ItemDelegate interface.
// It's important to highlight the selected item.
type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type selectDatabaseMsg struct {
	ActiveDatabase string
}

type selectTableMsg struct {
	Table string
}

type SidebarViewport struct {
	c *client.Client

	tablesList      list.Model
	sidebarViewport viewport.Model
	dbTree          *treeview.TuiTreeModel[DBObject]

	bindings *command.TUIKeyMap

	activeDatabase string
	width, height  int
}

type DBGraphTreeBuilderProvider struct{}

func (d DBGraphTreeBuilderProvider) ID(do DBObject) string {
	return do.ID
}

func (d *DBGraphTreeBuilderProvider) Name(do DBObject) string {
	return do.Name
}
func (p *DBGraphTreeBuilderProvider) Children(do DBObject) []DBObject {
	return do.Children
}

func NewSidebarViewport(ctx context.Context, c *client.Client, kb *command.TUIKeyMap) (SidebarViewport, error) {
	svp := SidebarViewport{c: c, bindings: kb}

	svp.sidebarViewport = viewport.New(0, 0)
	svp.sidebarViewport.KeyMap = viewport.KeyMap{}

	_ = DBObject{
		ID:     "databases",
		Name:   "localhost:5432",
		Type:   "connection",
		Driver: drivers.Postgres,
		Children: []DBObject{
			{
				ID:   "databases-postgres",
				Name: "postgres",
				Type: "database",
				Children: []DBObject{
					{
						ID:   "schemas-public",
						Name: "Public",
						Type: "schema",
						Children: []DBObject{
							{
								ID:   "schemas-public-users",
								Name: "users",
								Type: "table",
							},
							{
								ID:   "schemas-public-employees",
								Name: "employees",
								Type: "table",
							},
						},
					},
					{
						ID:   "schemas-products",
						Name: "Products",
						Type: "schema",
						Children: []DBObject{
							{
								ID:   "schemas-products-products",
								Name: "products",
								Type: "table",
							},
							{
								ID:   "schemas-products-prices",
								Name: "prices",
								Type: "table",
							},
						},
					},
				},
			},
			{
				ID:   "databases-users",
				Name: "users",
				Type: "database",
				Children: []DBObject{
					{
						ID:   "users-schemas-public",
						Name: "Public",
						Type: "schema",
						Children: []DBObject{
							{
								ID:   "users-schemas-public-employees",
								Name: "employees",
								Type: "table",
							},
							{
								ID:   "users-schemas-public-families",
								Name: "families",
								Type: "table",
							},
						},
					},
				},
			},
			{
				ID:   "databases-films",
				Name: "films",
				Type: "database",
			},
		},
	}

	if svp.c.ShowDataCatalog() {
		dbs, err := svp.c.ShowDatabases()
		if err != nil {
			return svp, err
		}
		rootID := uuid.New().String()

		root := DBObject{
			ID:       rootID,
			Name:     c.Host(),
			Driver:   c.Driver(),
			Type:     typeConnection,
			Children: make([]DBObject, 0, len(dbs)),
		}

		for _, db := range dbs {
			nodeID := uuid.New().String()
			dbNode := DBObject{
				ID:       fmt.Sprintf("%s-%s", db, nodeID),
				Name:     db,
				Type:     typeDatabase,
				Children: []DBObject{},
			}
			root.Children = append(root.Children, dbNode)
		}
		tree, err := treeview.NewTreeFromNestedData(
			ctx,
			[]DBObject{root},
			&DBGraphTreeBuilderProvider{},
			treeview.WithExpandFunc(func(n *treeview.Node[DBObject]) bool {
				return string(n.Data().Type) == string(typeConnection)
			}),
			treeview.WithProvider(createCyberpunkProvider()),
		)
		if err != nil {
			return svp, err
		}

		svp.dbTree = svp.newTuiTreeModel(tree, 0, 80)

		// If there are databases, choose the first one as the default active one.
		if len(dbs) > 0 {
			var i int

			for n, err := range svp.dbTree.AllVisible(ctx) {
				if err != nil {
					break
				}

				if i == 1 {
					svp.c.SetActiveDatabase(n.Node.Name())
					svp.activeDatabase = n.Node.Name()
					_, _ = svp.dbTree.SetFocusedID(ctx, n.Node.ID())
					break
				}
				i++
			}
		}

	} else {
		ts, err := svp.c.ShowTables()
		if err != nil {
			return svp, err
		}

		tables := make([]list.Item, 0)
		for _, ta := range ts {
			tables = append(tables, item(ta))
		}

		l := list.New(tables, itemDelegate{}, 0, 0)
		l.Title = "Tables"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.SetShowHelp(false)
		l.KeyMap.Quit.Unbind()
		svp.tablesList = l
		svp.updateStyles()
	}

	return svp, nil
}

// updateStyle setup the styles across the client.
func (s *SidebarViewport) updateStyles() {
	styles := newStyles()
	s.tablesList.Styles.Title = styles.title
	s.tablesList.Styles.PaginationStyle = styles.pagination
	s.tablesList.Styles.HelpStyle = styles.help
	s.tablesList.SetDelegate(itemDelegate{styles: &styles})
}

func (s *SidebarViewport) ActiveDatabase() string {
	return s.activeDatabase
}

func (s *SidebarViewport) SetSize(w, h int) {
	s.width = w - 4
	s.height = h - 2

	s.sidebarViewport.Width = s.width
	s.sidebarViewport.Height = s.height

	if s.c.ShowDataCatalog() {
		if s.dbTree != nil {
			s.dbTree = s.newTuiTreeModel(s.dbTree.Tree, 0, s.height)
		}
	} else {
		s.tablesList.SetSize(w, s.height)
	}
}

func (s *SidebarViewport) newTuiTreeModel(tree *treeview.Tree[DBObject], width, height int) *treeview.TuiTreeModel[DBObject] {
	// Create custom key map to avoid key conflicts
	keyMap := treeview.DefaultKeyMap()
	keyMap.SearchStart = []string{"/"}
	keyMap.Up = []string{"up", "k", "w"}
	keyMap.Down = []string{"down", "j", "s"}
	keyMap.Toggle = []string{"enter"}

	return treeview.NewTuiTreeModel(tree,
		treeview.WithTuiWidth[DBObject](width),
		treeview.WithTuiHeight[DBObject](height),
		treeview.WithTuiKeyMap[DBObject](keyMap),
		treeview.WithTuiDisableNavBar[DBObject](true),
		treeview.WithTuiAllowResize[DBObject](false),
	)
}

func (s SidebarViewport) Init() tea.Cmd {
	return nil
}

func (s SidebarViewport) Update(msg tea.Msg) (SidebarViewport, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if s.c.ShowDataCatalog() {
				selectedNode := s.dbTree.GetFocusedNode()
				if selectedNode != nil && selectedNode.Data() != nil {
					switch *&selectedNode.Data().Type {
					case typeDatabase:
						s.c.SetActiveDatabase(selectedNode.Name())
						if selectedNode.IsExpanded() {
							selectedNode.Collapse()
						} else {
							selectedNode.Expand()
						}
						selectDatabaseCmd := func() tea.Msg {
							s.activeDatabase = selectedNode.Name()
							return selectDatabaseMsg{ActiveDatabase: selectedNode.Name()}
						}
						return s, selectDatabaseCmd
					case typeTable:
						selectTableCmd := func() tea.Msg {
							return selectTableMsg{Table: selectedNode.Name()}
						}

						return s, selectTableCmd
					}
				}
			} else {
				tableItem := s.tablesList.Items()[s.tablesList.Index()]

				i := tableItem.(item)
				tableName := i.Title()

				selectTableCmd := func() tea.Msg {
					return selectTableMsg{Table: tableName}
				}

				return s, selectTableCmd
			}
		}

		switch {
		case key.Matches(msg, s.bindings.PageTop):
			if s.c.ShowDataCatalog() {
				ctx := context.Background()

				for nodeInfo, err := range s.dbTree.AllVisible(ctx) {
					if err != nil {
						break
					}

					_, _ = s.dbTree.SetFocusedID(ctx, nodeInfo.Node.ID())
					break
				}
				return s, nil
			} else {
				s.tablesList.Select(0)
				s.sidebarViewport.SetContent(s.tablesList.View())
				return s, nil
			}
		case key.Matches(msg, s.bindings.PageBottom):
			if s.c.ShowDataCatalog() {
				ctx := context.Background()

				var bottomNodeID string
				var found bool

				for nodeInfo, err := range s.dbTree.AllVisible(ctx) {
					if err != nil {
						break
					}

					bottomNodeID = nodeInfo.Node.ID()
					found = true
				}

				if found {
					_, _ = s.dbTree.SetFocusedID(ctx, bottomNodeID)
				}
				return s, nil
			} else {
				totalItems := len(s.tablesList.Items())
				if totalItems > 0 {
					s.tablesList.Select(totalItems - 1)
				}
				s.sidebarViewport.SetContent(s.tablesList.View())
				return s, nil
			}
		}
		if s.c.ShowDataCatalog() {
			s.sidebarViewport.SetContent(s.dbTree.View())
		}

		switch msg.String() {
		case "left", "h":
			s.sidebarViewport.ScrollLeft(4)
			return s, nil
		case "right", "l":
			s.sidebarViewport.ScrollRight(4)
			return s, nil
		}

		if s.c.ShowDataCatalog() {
			if s.dbTree != nil {
				updatedModel, treeCmd := s.dbTree.Update(msg)
				if newTreeModel, ok := updatedModel.(*treeview.TuiTreeModel[DBObject]); ok {
					s.dbTree = newTreeModel
				}
				cmd = treeCmd
			}
		} else {
			s.tablesList, cmd = s.tablesList.Update(msg)
			s.sidebarViewport.SetContent(s.tablesList.View())
		}
		cmds = append(cmds, cmd)

		// case tablesFetchedMsg:
		// 	selectedNode := s.dbTree.GetFocusedNode()
		// 	if selectedNode != nil {
		// 		tables := make([]*treeview.Node[DBObject], len(msg.tables))
		// 		for i, t := range msg.tables {
		// 			nodeID := uuid.New().String()
		// 			tables[i] = treeview.NewNode(fmt.Sprintf("%s-%s", t, nodeID), t, "table")
		// 		}
		// 		selectedNode.SetChildren(tables)
		// 	}
		return s, nil
	case querySuccessMsg:
		if len(msg.tables) > 0 {
			tables := make([]list.Item, 0)
			for _, ta := range msg.tables {
				tables = append(tables, item(ta))
			}
			s.tablesList.SetItems(tables)
		}
		return s, nil
	}

	return s, tea.Batch(cmds...)
}

func (s SidebarViewport) View() string {
	s.sidebarViewport.SetContent(s.tablesList.View())
	if s.c.ShowDataCatalog() {
		s.sidebarViewport.SetContent(s.dbTree.View())
	}
	sideViewContent := s.sidebarViewport.View()

	listBorder := darkPurple
	sideViewContent = tablesListStyle.BorderForeground(listBorder).Width(s.width).Height(s.height).Render(sideViewContent)

	return sideViewContent
}

func createCyberpunkProvider() *treeview.DefaultNodeProvider[DBObject] {
	// Icons for database connections.
	postgresIconRule := treeview.WithIconRule(dbObjectIsConnection("connection", drivers.Postgres), "🐘")
	mysqlIconRule := treeview.WithIconRule(dbObjectIsConnection("connection", drivers.MySQL), "🐬")
	sqliteIconRule := treeview.WithIconRule(dbObjectIsConnection("connection", drivers.SQLite), "🪶")
	oracleIconRule := treeview.WithIconRule(dbObjectIsConnection("connection", drivers.Oracle), "☀ ")
	sqlServerIconRule := treeview.WithIconRule(dbObjectIsConnection("connection", drivers.SQLServer), "🔷")

	// Icons for database entities.
	databaseIconRule := treeview.WithIconRule(dbObjectHasType("database"), "⛃")
	schemaIconRule := treeview.WithIconRule(dbObjectHasType("schema"), "📁")
	tableIconRule := treeview.WithIconRule(dbObjectHasType("table"), "📄")

	return treeview.NewDefaultNodeProvider(
		postgresIconRule,
		mysqlIconRule,
		sqliteIconRule,
		oracleIconRule,
		sqlServerIconRule,
		databaseIconRule,
		schemaIconRule,
		tableIconRule,
		treeview.WithStyleRule(
			func(n *treeview.Node[DBObject]) bool { return true },
			lipgloss.NewStyle().
				Foreground(whiteText).
				PaddingLeft(2),

			lipgloss.NewStyle().
				Foreground(cyberGreen).
				Background(darkPurple).
				Bold(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderLeft(true).
				BorderForeground(hiMagenta).
				PaddingLeft(1),
		),
		treeview.WithFormatter(func(node *treeview.Node[DBObject]) (string, bool) {
			return node.Name(), true
		}),
	)
}
