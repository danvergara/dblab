package bubbletui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/pkg/bubbletui/keys"
)

const (
	endArrow   = ""
	startArrow = ""
)

type StatusBar struct {
	width  int
	keyMap keys.KeyMap
	fixed  string
	schema string
	focus  focusState
}

func NewStatusBar(km keys.KeyMap, driver, conn, schema string) StatusBar {
	var statusKb = lipgloss.NewStyle().
		Background(KbOddBg).
		Foreground(KbOddText).
		Render(fmt.Sprintf(" %s %s ", km.Quit.Help().Key, km.Quit.Help().Desc)) +
		lipgloss.NewStyle().
			Background(KbEvenBg).
			Foreground(KbOddBg).
			Render(endArrow) +
		lipgloss.NewStyle().
			Background(KbEvenBg).
			Foreground(KbEvenText).
			Render(fmt.Sprintf(" %s %s ", km.Help.Help().Key, km.Help.Help().Desc)) +
		lipgloss.NewStyle().
			Foreground(KbEvenBg).
			Render(endArrow) +
		lipgloss.NewStyle().
			Foreground(KbEvenText).
			Render("  "+driver+": "+conn)
	return StatusBar{keyMap: km, fixed: statusKb, focus: focusEditor, schema: schema}
}

func (f StatusBar) Init() tea.Cmd {
	return nil
}

func (f StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	return f, nil
}

func (f *StatusBar) ShowFocus(focus focusState) {
	f.focus = focus
}

func (f *StatusBar) SetSchema(schema string) {
	f.schema = schema
}

func (f *StatusBar) SetWidth(width int) {
	f.width = width
}

func (f StatusBar) View() tea.View {
	leftBlock := lipgloss.NewStyle().
		Bold(true).
		Background(FocusBg).
		Foreground(FocusText).
		Render("  "+f.focus.String()+"  ") +
		lipgloss.NewStyle().
			Foreground(FocusBg).
			Render(endArrow) +
		f.fixed

	var rightBlock string
	if f.schema != "" {
		rightBlock =
			lipgloss.NewStyle().
				Foreground(cyberGreen).
				Render("active schema:" + " ")
	}
	rightBlock +=
		lipgloss.NewStyle().
			Foreground(SchemaBg).
			Render(startArrow) +
			lipgloss.NewStyle().
				Bold(true).
				Background(SchemaBg).
				Foreground(SchemaText).
				Render(" "+f.schema+" ")

	spacerSize := f.width - lipgloss.Width(leftBlock) - lipgloss.Width(rightBlock)

	spacer := lipgloss.NewStyle().
		Width(spacerSize).
		Render("")

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, leftBlock, spacer, rightBlock))
}
