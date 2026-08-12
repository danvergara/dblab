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
	mode   Mode
	width  int
	keyMap keys.KeyMap
	fixed  string
	focus  focusState
}

func NewStatusBar(mode Mode, km keys.KeyMap, driver, conn string) StatusBar {
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
	return StatusBar{mode: mode, keyMap: km, fixed: statusKb, focus: focusEditor}
}

func (f StatusBar) Init() tea.Cmd {
	return nil
}

func (f StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case modeChangeMsg:
		f.mode = msg.mode
	}
	return f, nil
}

func (f *StatusBar) ShowFocus(focus focusState) {
	f.focus = focus
}

func (f *StatusBar) SetWidth(width int) {
	f.width = width - 4
}

func (f StatusBar) View() tea.View {
	modeColorBg := NormalModeBg
	modeColorText := NormalModeText

	if f.mode == InsertMode {
		modeColorBg = InsertModeBg
		modeColorText = InsertModeText
	}

	leftBlock := lipgloss.NewStyle().
		Bold(true).
		Background(modeColorBg).
		Foreground(modeColorText).
		Render("  "+f.mode.String()+"  ") +
		lipgloss.NewStyle().
			Background(KbOddBg).
			Foreground(modeColorBg).
			Render(endArrow) +
		f.fixed

	rightBlock := lipgloss.NewStyle().
		Foreground(FocusBg).
		Render(startArrow) +
		lipgloss.NewStyle().
			Bold(true).
			Background(FocusBg).
			Foreground(FocusText).
			Render(" "+f.focus.String()+" ")

	spacerSize := f.width - lipgloss.Width(leftBlock) - lipgloss.Width(rightBlock)

	spacer := lipgloss.NewStyle().
		Width(spacerSize).
		Render("")

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, leftBlock, spacer, rightBlock))
}
