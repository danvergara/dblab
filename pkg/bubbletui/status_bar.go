package bubbletui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
)

const (
	primaryColor   = "#8A2BE2" // Morado
	secondaryColor = "#008080" // Teal
	endArrow       = ""
	startArrow     = ""
)

type StatusBar struct {
	mode     Mode
	width    int
	bindings *command.TUIKeyMap
	fixed    string
	client   *client.Client
	focus    focusState
}

func NewStatusBar(mode Mode, kb *command.TUIKeyMap, client *client.Client) StatusBar {
	var statusKb = lipgloss.NewStyle().
		Background(KbOddBg).
		Foreground(KbOddText).
		Render(fmt.Sprintf(" %s %s ", kb.Quit.Help().Key, kb.Quit.Help().Desc)) +
		lipgloss.NewStyle().
			Background(KbEvenBg).
			Foreground(KbOddBg).
			Render(endArrow) +
		lipgloss.NewStyle().
			Background(KbEvenBg).
			Foreground(KbEvenText).
			Render(fmt.Sprintf(" %s %s ", kb.Help.Help().Key, kb.Help.Help().Desc)) +
		lipgloss.NewStyle().
			Foreground(KbEvenBg).
			Render(endArrow) +
		lipgloss.NewStyle().
			Foreground(KbEvenText).
			Render("  "+client.Driver()+": "+client.Conn())
	return StatusBar{mode: mode, bindings: kb, fixed: statusKb, client: client, focus: focusEditor}
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
	return tea.NewView(f.view())
}

func (f *StatusBar) view() string {

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

	return lipgloss.JoinHorizontal(lipgloss.Left, leftBlock, spacer, rightBlock) + "\n"
}
