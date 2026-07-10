package bubbletui

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/internal/history"
)

type querySelectedMsg struct {
	QueryText string
}

type queryHistoryLoadedMsg struct {
	items []list.Item
}

type backToNormalMsg struct{}

type queryHistoryErrMsg struct{ err error }

type HistoryModel struct {
	state         state
	spinner       spinner.Model
	list          list.Model
	width         int
	height        int
	loadingAction string
}

func NewHistoryModel() *HistoryModel {
	delegate := list.NewDefaultDelegate()
	// delegate.ShowDescription = false

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(whiteText)

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(cyberGreen).
		BorderLeftForeground(hiMagenta)

	queryList := list.New([]list.Item{}, delegate, 0, 0)
	queryList.Title = "Select a query from the history"
	queryList.Styles.Title = queryList.Styles.Title.
		Background(darkPurple).
		Foreground(hiMagenta).
		Bold(true)
	queryList.SetShowStatusBar(false)
	queryList.SetShowHelp(true)
	queryList.KeyMap.Quit.Unbind()

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &HistoryModel{
		state:         stateLoading,
		spinner:       s,
		list:          queryList,
		loadingAction: "Fetch queries from the history file...",
	}
}

func (h *HistoryModel) SetSize(width, height int) {
	h.height = height
	h.width = width
	h.list.SetSize(min(50, h.width)-6, 14)
}

func (h *HistoryModel) Init() tea.Cmd {
	return tea.Batch(h.spinner.Tick, fetchQueryHistoryCmd())
}

func (h *HistoryModel) Update(msg tea.Msg) (*HistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		h.list.SetSize(min(50, h.width)-6, 14)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			selectedQuery := h.list.SelectedItem().(history.QueryHistory).QueryText

			return h, func() tea.Msg {
				return querySelectedMsg{QueryText: selectedQuery}
			}
		case "esc":
			return h, func() tea.Msg {
				return backToNormalMsg{}
			}
		}
	case queryHistoryLoadedMsg:
		cmd := h.list.SetItems(msg.items)
		h.state = stateForm
		return h, cmd
	}

	switch h.state {
	case stateLoading:
		var cmd tea.Cmd
		h.spinner, cmd = h.spinner.Update(msg)
		return h, cmd

	case stateForm:
		var cmd tea.Cmd
		h.list, cmd = h.list.Update(msg)
		return h, cmd
	}

	return h, nil
}

func (h *HistoryModel) View() tea.View {
	var (
		v       tea.View
		content string
	)

	v.AltScreen = true

	switch h.state {
	case stateLoading:
		spinnerView := fmt.Sprintf("\n %s %s\n", h.spinner.View(), h.loadingAction)
		content = lipgloss.Place(h.width, h.height, lipgloss.Center, lipgloss.Center, spinnerView)
	case stateForm:
		formStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(neonPurple).
			Padding(1, 2)
		boxedForm := formStyle.Render(h.list.View())
		content = lipgloss.Place(
			h.width,
			h.height,
			lipgloss.Center,
			lipgloss.Center,
			boxedForm,
		)
	}

	v.SetContent(content)
	return v
}

func fetchQueryHistoryCmd() tea.Cmd {
	return func() tea.Msg {
		// Read the config based dir path.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		var items []list.Item
		queryHistory, err := history.ReadHistory(configDir)
		if err != nil {
			return queryHistoryErrMsg{err: err}
		}

		for _, query := range queryHistory {
			items = append(items, query)
		}
		return queryHistoryLoadedMsg{items: items}
	}
}
