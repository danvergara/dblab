package bubbletui

import (
	"context"
	"io"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Digital-Shane/treeview/v2"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/davecgh/go-spew/spew"
)

func dbObjectHasType(nodeType string) func(*treeview.Node[*client.DBNode]) bool {
	return func(n *treeview.Node[*client.DBNode]) bool {
		return (*n.Data()).Type == nodeType
	}
}

type selectTableMsg struct {
	Schema string
	Table  string
}

type SidebarViewport struct {
	c        *client.Client
	bindings *command.TUIKeyMap

	sidebarViewport viewport.Model
	dbTree          *treeview.TuiTreeModel[*client.DBNode]
	width, height   int

	selected bool
	dump     io.Writer
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
	var dump *os.File
	if _, ok := os.LookupEnv("DBLAB_DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("sidebar_messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	svp := SidebarViewport{
		c:        c,
		bindings: kb,
		dump:     dump,
	}

	svp.sidebarViewport = viewport.New(viewport.WithHeight(0), viewport.WithWidth(0))
	svp.sidebarViewport.KeyMap = viewport.KeyMap{}

	root, err := c.Catalog(ctx)
	if err != nil {
		return SidebarViewport{}, err
	}

	tree, err := treeview.NewTreeFromNestedData[*client.DBNode](
		ctx,
		[]*client.DBNode{root},
		&DBGraphTreeBuilderProvider{},
		treeview.WithProvider(createCyberpunkProvider()),
	)
	if err != nil {
		return svp, err
	}

	svp.dbTree = svp.newTuiTreeModel(tree, 0, 80)

	return svp, nil
}

func (s *SidebarViewport) SetSize(w, h int) {
	s.width = w - 4
	s.height = h - 2

	s.sidebarViewport.SetWidth(s.width)
	s.sidebarViewport.SetHeight(s.height)

	if s.dbTree != nil {
		s.dbTree = s.newTuiTreeModel(s.dbTree.Tree, 0, s.height-2)
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
	if s.dump != nil {
		spew.Fdump(s.dump, msg)
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEnter:
			selectedNode := s.dbTree.GetFocusedNode()
			if selectedNode != nil && selectedNode.Data() != nil {
				switch (*selectedNode.Data()).Type {
				case "table":
					selectTableCmd := func() tea.Msg {
						stm := selectTableMsg{Table: selectedNode.Name()}
						switch s.c.Driver() {
						case drivers.PostgreSQL, drivers.Postgres, drivers.PostgresSSH, drivers.Oracle:
							stm.Schema = (*selectedNode.Data()).ParentName
						}
						return stm
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

		s.sidebarViewport.SetContent(s.dbTree.View().Content)

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

		return s, cmd
	case querySuccessMsg:
		var cmd tea.Cmd
		if msg.reloadCatalog {
			cmd = s.updateGraph()
		}
		return s, cmd
	case updateGraphMsg:
		s.dbTree = msg.tree
		return s, nil
	case updateGraphErrMsg:
		return s, nil
	}

	return s, cmd
}

func (s SidebarViewport) View() string {
	s.sidebarViewport.SetContent(s.dbTree.View().Content)
	sideViewContent := s.sidebarViewport.View()

	listBorder := darkPurple
	if s.selected {
		listBorder = neonPurple
	}
	sideViewContent = tablesListStyle.BorderForeground(listBorder).Height(s.height).Render(sideViewContent)

	return sideViewContent
}

// updateGraph method refreshes the database catalog asynchronously, so it does not block the bubbletea execution.
// If it succeeds, returns a updateGraphMsg with the a new database tree. Otherwise, it returns  updateGraphErrMsg with the error.
func (s *SidebarViewport) updateGraph() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		root, err := s.c.Catalog(ctx)
		if err != nil {
			return updateGraphErrMsg{err}
		}

		tree, err := treeview.NewTreeFromNestedData[*client.DBNode](
			ctx,
			[]*client.DBNode{root},
			&DBGraphTreeBuilderProvider{},
			treeview.WithProvider(createCyberpunkProvider()),
		)
		if err != nil {
			return updateGraphErrMsg{err}
		}

		dbTree := s.newTuiTreeModel(tree, 0, s.height-2)
		_, _ = dbTree.SetFocusedID(ctx, root.ID)

		return updateGraphMsg{tree: dbTree}
	}
}

func createCyberpunkProvider() *treeview.DefaultNodeProvider[*client.DBNode] {
	// Icons for database objects.
	databaseIconRule := treeview.WithIconRule(dbObjectHasType("database"), "⛃")
	schemaIconRule := treeview.WithIconRule(dbObjectHasType("schema"), "📁")
	tableIconRule := treeview.WithIconRule(dbObjectHasType("table"), "📋")

	return treeview.NewDefaultNodeProvider[*client.DBNode](
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
