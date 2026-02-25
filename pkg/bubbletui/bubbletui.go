package bubbletui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type focusState int

const (
	focusInput focusState = iota
	focusList
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	input textarea.Model
	list  list.Model
	table table.Model
}

func NewModel() Model {
	var emptyItems []list.Item
	return Model{
		table: table.New(),
		input: textarea.New(),
		list:  list.New(emptyItems, list.NewDefaultDelegate(), 20, 10),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}
