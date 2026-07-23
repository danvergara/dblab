package bubbletui

import (
	"fmt"
	"os"
	"slices"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/danvergara/dblab/internal/history"
)

// querySelectedMsg async message to let the bubbletea app know about a query selected from the confi file.
type querySelectedMsg struct {
	QueryText string
}

// queryHistoryLoadedMsg async message load the query history from the config file.
type queryHistoryLoadedMsg struct {
	items []list.Item
}

// backToNormalMsg asyn message to escape from the query history view.
type backToNormalMsg struct{}

// queryHistoryErrMsg message to communicate errors while trying to get query history from the config file.
type queryHistoryErrMsg struct{ err error }

// HistoryModel struct is the model use to display the query history.
type HistoryModel struct {
	// checks if the history is loading or ready to show.
	state state
	// spinner model to show while the data is loading.
	spinner spinner.Model
	// list model to show the query history.
	list list.Model
	// model size.
	width         int
	height        int
	loadingAction string
}

// NewHistoryModel returns a pointer to the HistoryModel, with the models initialized.
func NewHistoryModel() *HistoryModel {
	delegate := list.NewDefaultDelegate()

	// white text for the normal title.
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(whiteText)

	// Cyber green color for the selected items with a magent border foreground.
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(cyberGreen).
		BorderLeftForeground(hiMagenta)

	// Set up the list model.
	queryList := list.New([]list.Item{}, delegate, 0, 0)
	queryList.Title = "Select a query from the history"
	queryList.Styles.Title = queryList.Styles.Title.
		Background(darkPurple).
		Foreground(hiMagenta).
		Bold(true)
	queryList.SetShowStatusBar(false)
	queryList.SetShowHelp(true)

	// Set up the spinner model.
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &HistoryModel{
		state:         stateLoading,
		spinner:       s,
		list:          queryList,
		loadingAction: "Fetch queries from the history file...",
	}
}

// SetSize method is used to set the model size when the main tui model routes the size from the tea.WindowSizeMsg message to this model.
func (h *HistoryModel) SetSize(width, height int) {
	h.height = height
	h.width = width
	h.list.SetSize(h.width-6, 14)
}

// Init method is used to initialize the spinner Tick command and fetch the query history.
func (h *HistoryModel) Init() tea.Cmd {
	return tea.Batch(h.spinner.Tick, fetchQueryHistoryCmd())
}

func (h *HistoryModel) Update(msg tea.Msg) (*HistoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		h.list.SetSize(h.width-6, 14)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// Select the current query from the list.
			selectedQuery := h.list.SelectedItem().(history.QueryHistory).QueryText
			return h, func() tea.Msg {
				return querySelectedMsg{QueryText: selectedQuery}
			}
		case "esc":
			// Press esc to get back to the main app if a query was not selected.
			return h, func() tea.Msg {
				return backToNormalMsg{}
			}
		}
	// catch the queryHistoryLoadedMsg with the query history and reverse the content to show the history in descending order based on the timestamp.
	case queryHistoryLoadedMsg:
		slices.Reverse(msg.items)
		cmd := h.list.SetItems(msg.items)
		h.state = stateForm
		return h, cmd
	}

	// Manage the state of this model.
	// If it's loading, routes the messages to the spinner,
	// otherwise, it routes the messages to the list.
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

// View method renders the spinner or the list content.
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
		content = setModalContent(h.list.View(), h.width, h.height)
	}

	v.SetContent(content)
	return v
}

// fetchQueryHistoryCmd is used to fetch the query history when the model is loaded.
func fetchQueryHistoryCmd() tea.Cmd {
	return func() tea.Msg {
		// Read the config based dir path.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		// Retrieve the query history from the config file.
		queryHistory, err := history.ReadHistory(configDir)
		if err != nil {
			return queryHistoryErrMsg{err: err}
		}

		// Populate the items slice with queryHistory which is []history.QueryHistory that implements the item interface.
		var items []list.Item
		for _, query := range queryHistory {
			items = append(items, query)
		}
		return queryHistoryLoadedMsg{items: items}
	}
}
