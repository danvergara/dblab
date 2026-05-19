package bubbletui

import (
	"context"

	"github.com/Digital-Shane/treeview"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
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

type selectTableMsg struct {
	Table string
}

type SidebarViewport struct {
	c        *client.Client
	bindings *command.TUIKeyMap

	sidebarViewport viewport.Model
	dbTree          *treeview.TuiTreeModel[*client.DBNode]
	width, height   int
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

		return s, cmd
	}

	return s, cmd
}

func (s SidebarViewport) View() string {
	s.sidebarViewport.SetContent(s.dbTree.View())
	sideViewContent := s.sidebarViewport.View()

	listBorder := darkPurple
	sideViewContent = tablesListStyle.BorderForeground(listBorder).Width(s.width).Height(s.height).Render(sideViewContent)

	return sideViewContent
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
