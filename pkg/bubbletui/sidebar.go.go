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
	// "github.com/google/uuid"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	// "github.com/danvergara/dblab/pkg/drivers"
	// "github.com/google/uuid"
)

type dbObjectType string

const (
	dbObjectDatabase dbObjectType = "database"
	dbOjbectSchema   dbObjectType = "schema"
	dbObjectTable    dbObjectType = "table"
	dbObjectHost     dbObjectType = "host"
)

func dbObjectHasType(nodeType string) func(*treeview.Node[*client.DBNode]) bool {
	return func(n *treeview.Node[*client.DBNode]) bool {
		return (*n.Data()).Type == nodeType
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

type getChildrenMsg struct {
	parent     string
	parentType string
}
type getChildrenErrMsg struct{ err error }

type SidebarViewport struct {
	c *client.Client

	tablesList      list.Model
	sidebarViewport viewport.Model
	dbTree          *treeview.TuiTreeModel[*client.DBNode]

	bindings *command.TUIKeyMap

	activeDatabase string
	width, height  int
}

type DBGraphTreeBuilderProvider struct{}

func (d DBGraphTreeBuilderProvider) ID(do *client.DBNode) string {
	return do.ID
}

func (d *DBGraphTreeBuilderProvider) Name(do *client.DBNode) string {
	return do.Name
}
func (p *DBGraphTreeBuilderProvider) Children(do *client.DBNode) []*client.DBNode {
	return do.Children
}

func NewSidebarViewport(ctx context.Context, c *client.Client, kb *command.TUIKeyMap) (SidebarViewport, error) {
	svp := SidebarViewport{c: c, bindings: kb}

	svp.sidebarViewport = viewport.New(0, 0)
	svp.sidebarViewport.KeyMap = viewport.KeyMap{}

	root, err := c.Catalog()
	if err != nil {
		return SidebarViewport{}, err
	}

	tree, err := treeview.NewTreeFromNestedData[*client.DBNode](
		ctx,
		[]*client.DBNode{root},
		&DBGraphTreeBuilderProvider{},
		// treeview.WithProvider(createCyberpunkProvider()),
	)
	if err != nil {
		return svp, err
	}

	svp.dbTree = svp.newTuiTreeModel(tree, 0, 80)

	// if svp.c.ShowDataCatalog() {
	// } else {
	// 	ts, err := svp.c.ShowTables()
	// 	if err != nil {
	// 		return svp, err
	// 	}
	//
	// 	tables := make([]list.Item, 0)
	// 	for _, ta := range ts {
	// 		tables = append(tables, item(ta))
	// 	}
	//
	// 	l := list.New(tables, itemDelegate{}, 0, 0)
	// 	l.Title = "Tables"
	// 	l.SetShowStatusBar(false)
	// 	l.SetFilteringEnabled(false)
	// 	l.SetShowHelp(false)
	// 	l.KeyMap.Quit.Unbind()
	// 	svp.tablesList = l
	// 	svp.updateStyles()
	// }

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

	if s.dbTree != nil {
		s.dbTree = s.newTuiTreeModel(s.dbTree.Tree, 0, s.height)
	}
}

func (s *SidebarViewport) newTuiTreeModel(tree *treeview.Tree[*client.DBNode], width, height int) *treeview.TuiTreeModel[*client.DBNode] {
	// Create custom key map to avoid key conflicts
	keyMap := treeview.DefaultKeyMap()
	keyMap.SearchStart = []string{"/"}
	keyMap.Up = []string{"up", "k", "w"}
	keyMap.Down = []string{"down", "j", "s"}
	keyMap.Toggle = []string{"enter"}

	return treeview.NewTuiTreeModel(tree,
		treeview.WithTuiWidth[*client.DBNode](width),
		treeview.WithTuiHeight[*client.DBNode](height),
		treeview.WithTuiKeyMap[*client.DBNode](keyMap),
		treeview.WithTuiDisableNavBar[*client.DBNode](true),
		treeview.WithTuiAllowResize[*client.DBNode](false),
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
			selectedNode := s.dbTree.GetFocusedNode()
			if selectedNode != nil && selectedNode.Data() != nil {
				switch (*selectedNode.Data()).Type {
				case "table":
					selectTableCmd := func() tea.Msg {
						return selectTableMsg{Table: selectedNode.Name()}
					}

					return s, selectTableCmd
				}
			}
		}

		switch {
		case key.Matches(msg, s.bindings.PageTop):
			ctx := context.Background()

			for nodeInfo, err := range s.dbTree.AllVisible(ctx) {
				if err != nil {
					break
				}

				_, _ = s.dbTree.SetFocusedID(ctx, nodeInfo.Node.ID())
				break
			}
			return s, nil
		case key.Matches(msg, s.bindings.PageBottom):
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
		}

		s.sidebarViewport.SetContent(s.dbTree.View())

		switch msg.String() {
		case "left", "h":
			s.sidebarViewport.ScrollLeft(4)
			return s, nil
		case "right", "l":
			s.sidebarViewport.ScrollRight(4)
			return s, nil
		}

		if s.dbTree != nil {
			updatedModel, treeCmd := s.dbTree.Update(msg)
			if newTreeModel, ok := updatedModel.(*treeview.TuiTreeModel[*client.DBNode]); ok {
				s.dbTree = newTreeModel
			}
			cmd = treeCmd
		}

		cmds = append(cmds, cmd)
		return s, nil
	}

	return s, tea.Batch(cmds...)
}

func (s SidebarViewport) View() string {
	s.sidebarViewport.SetContent(s.dbTree.View())
	sideViewContent := s.sidebarViewport.View()

	listBorder := darkPurple
	sideViewContent = tablesListStyle.BorderForeground(listBorder).Width(s.width).Height(s.height).Render(sideViewContent)

	return sideViewContent
}

func createCyberpunkProvider() *treeview.DefaultNodeProvider[*client.DBNode] {
	// Icons for database connections.
	// postgresIconRule := treeview.WithIconRule(dbObjectIsConnection("host", drivers.Postgres), "🐘")
	// mysqlIconRule := treeview.WithIconRule(dbObjectIsConnection("host", drivers.MySQL), "🐬")
	// sqliteIconRule := treeview.WithIconRule(dbObjectIsConnection("host", drivers.SQLite), "🪶")
	// oracleIconRule := treeview.WithIconRule(dbObjectIsConnection("host", drivers.Oracle), "☀ ")
	// sqlServerIconRule := treeview.WithIconRule(dbObjectIsConnection("host", drivers.SQLServer), "🔷")

	// Icons for database entities.
	databaseIconRule := treeview.WithIconRule(dbObjectHasType("database"), "⛃")
	schemaIconRule := treeview.WithIconRule(dbObjectHasType("schema"), "📁")
	tableIconRule := treeview.WithIconRule(dbObjectHasType("table"), "📄")

	return treeview.NewDefaultNodeProvider[*client.DBNode](
		// postgresIconRule,
		// mysqlIconRule,
		// sqliteIconRule,
		// oracleIconRule,
		// sqlServerIconRule,
		databaseIconRule,
		schemaIconRule,
		tableIconRule,
		treeview.WithStyleRule(
			func(n *treeview.Node[*client.DBNode]) bool { return true },
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
		treeview.WithFormatter(func(node *treeview.Node[*client.DBNode]) (string, bool) {
			return node.Name(), true
		}),
	)
}
