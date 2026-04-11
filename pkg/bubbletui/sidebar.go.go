package bubbletui

import (
	"context"
	"fmt"

	"github.com/Digital-Shane/treeview"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/google/uuid"
)

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
	dbTree          *treeview.TuiTreeModel[string]

	bindings *command.TUIKeyBindings

	activeDatabase string
	width, height  int
}

func NewSidebarViewport(ctx context.Context, c *client.Client, kb *command.TUIKeyBindings) (SidebarViewport, error) {
	svp := SidebarViewport{c: c, bindings: kb}

	svp.sidebarViewport = viewport.New(0, 0)
	svp.sidebarViewport.KeyMap = viewport.KeyMap{}

	if svp.c.ShowDataCatalog() {
		dbs, err := svp.c.ShowDatabases()
		if err != nil {
			return svp, err
		}
		rootID := uuid.New().String()
		root := treeview.NewNode(fmt.Sprintf("%s-%s", "db", rootID), "db", "root")
		for _, db := range dbs {
			nodeID := uuid.New().String()
			root.AddChild(treeview.NewNode(fmt.Sprintf("%s-%s", db, nodeID), db, "database"))
		}

		root.Expand()

		treeRoot := treeview.NewTree([]*treeview.Node[string]{root}, treeview.WithProvider(createCyberpunkProvider()))

		svp.dbTree = treeview.NewTuiTreeModel(treeRoot,
			treeview.WithTuiWidth[string](0),
			treeview.WithTuiHeight[string](80),
		)

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
	s.width = w
	s.height = h

	if s.c.ShowDataCatalog() {
		if s.dbTree != nil {
			s.dbTree = s.newTuiTreeModel(s.dbTree.Tree, w, h-2)
		}
	} else {
		s.tablesList.SetSize(w, h-2)
	}

	s.sidebarViewport.Height = h - 4
	s.sidebarViewport.Width = w - 4
}

func (s *SidebarViewport) newTuiTreeModel(tree *treeview.Tree[string], width, height int) *treeview.TuiTreeModel[string] {
	// Create custom key map to avoid key conflicts
	keyMap := treeview.DefaultKeyMap()
	keyMap.SearchStart = []string{"/"}
	keyMap.Up = []string{"up", "k", "w"}
	keyMap.Down = []string{"down", "j", "s"}

	return treeview.NewTuiTreeModel(tree,
		treeview.WithTuiWidth[string](width),
		treeview.WithTuiHeight[string](height),
		treeview.WithTuiKeyMap[string](keyMap),
		treeview.WithTuiDisableNavBar[string](true),
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
					switch *selectedNode.Data() {
					case "database":
						s.c.SetActiveDatabase(selectedNode.Name())
						if selectedNode.IsExpanded() {
							selectedNode.Collapse()
							return s, nil
						} else {
							selectedNode.Expand()
							selectDatabaseCmd := func() tea.Msg {
								s.activeDatabase = selectedNode.Name()
								return selectDatabaseMsg{ActiveDatabase: selectedNode.Name()}
							}
							return s, selectDatabaseCmd
						}
					case "table":
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
			if s.dbTree != nil {
				updatedModel, treeCmd := s.dbTree.Update(msg)
				if newTreeModel, ok := updatedModel.(*treeview.TuiTreeModel[string]); ok {
					s.dbTree = newTreeModel
				}
				cmd = treeCmd
			}
		} else {
			s.tablesList, cmd = s.tablesList.Update(msg)
			s.sidebarViewport.SetContent(s.tablesList.View())
		}

	case tablesFetchedMsg:
		selectedNode := s.dbTree.GetFocusedNode()
		if selectedNode != nil {
			tables := make([]*treeview.Node[string], len(msg.tables))
			for i, t := range msg.tables {
				nodeID := uuid.New().String()
				tables[i] = treeview.NewNode(fmt.Sprintf("%s-%s", t, nodeID), t, "table")
			}
			selectedNode.SetChildren(tables)
		}

		cmds = append(cmds, cmd)
		return s, tea.Batch(cmds...)
	}
	return s, nil
}

func (s SidebarViewport) View() string {
	s.sidebarViewport.SetContent(s.tablesList.View())
	sideViewContent := s.sidebarViewport.View()
	if s.c.ShowDataCatalog() {
		sideViewContent = s.dbTree.View()
	}
	return sideViewContent
}
