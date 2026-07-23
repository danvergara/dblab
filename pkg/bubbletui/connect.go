package bubbletui

import (
	"errors"
	"fmt"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	databaseProfiles "github.com/danvergara/dblab/internal/profiles"
	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/zalando/go-keyring"
)

var ErrConnectionFormAborted = errors.New("connection form exited with ctrl+c")

// Connect model custom keys.
type customKeyMap struct {
	connect key.Binding
	delete  key.Binding
	quit    key.Binding
}

func newCustomKeyMap() customKeyMap {
	return customKeyMap{
		connect: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "connect"),
		),
		delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "delete"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

// state of the menu.
type state int

const (
	// state for loading is used during the connections loading.
	stateLoading state = iota
	// state used when the data is loaded and available.
	stateForm
)

// dbProfile struct is used as the bubble list model item.
// label field is used to show the data on the menu.
// name field matches the name of the database profile in the config file.
type dbProfile struct {
	name  string
	label string
}

// Implement the list.Item interface.
func (d dbProfile) Title() string       { return d.label }
func (d dbProfile) Description() string { return "" }
func (d dbProfile) FilterValue() string { return d.label }

// profilesLoadedMsg struct is a message what will carry the profiles back to the update loop.
type profilesLoadedMsg struct {
	items    []list.Item
	profiles map[string]command.Options
}

// errMsg used to communicate errors asynchronously when reading or deleting profiles from the config file.
type errMsg struct{ err error }

// itemDeletedMsg struct is used communicate the index of the profile on the list to be deleted.
type itemDeletedMsg struct {
	index int
}

// fetchProfilesCmd async function fetches a profiles from the config file.
// returns a tea.Cmd so it does not block the main bubbletea's thread.
func fetchProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		// Read the config based dir path.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		// Call the ReadProfiles function to get the profiles in the official config file.
		profiles, err := databaseProfiles.ReadProfiles(configDir)
		if err != nil {
			return errMsg{err: err}
		}

		// This part will come up with a list of items.
		var items []list.Item
		for name, profile := range profiles {
			var label string
			// SQLite is the only special case, since the database is a file.
			if profile.Driver == drivers.SQLite {
				label = fmt.Sprintf("%s (sqlite://%s)", name, profile.Host)
			} else {
				label = fmt.Sprintf("%s (%s://%s@%s)", name, profile.Driver, profile.User, profile.Host)
			}
			// append a new dbProfile, being the name the profile name and the label the data to be showed on the menu.
			items = append(items, dbProfile{name: name, label: label})
		}

		return profilesLoadedMsg{items: items, profiles: profiles}
	}
}

// deleteItemCmd asyc function to deleted a selected profile from the menu.
func deleteItemCmd(name string, index int) tea.Cmd {
	return func() tea.Msg {
		// Read the base config file path.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		// Delete a profile given the name.
		if err := databaseProfiles.DeleteProfile(configDir, name); err != nil {
			return errMsg{err: err}
		}

		// communicate the index being deleted to remove from the UI list.
		return itemDeletedMsg{index: index}
	}
}

// ConnectModel struct is the main model for the connect feature.
// Manages the list and the spinner sub-models,
// keeps track of the app state,
// and stores the profiles to be returned to the main dblab app.
type ConnectModel struct {
	state          state
	list           list.Model
	spinner        spinner.Model
	selectedOption string
	err            error
	profiles       map[string]command.Options
	width          int
	height         int
	aborted        bool
	loadingAction  string
}

func initConnectModel() *ConnectModel {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(whiteText)

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(cyberGreen).
		BorderLeftForeground(hiMagenta)

	profileList := list.New([]list.Item{}, delegate, 0, 0)
	profileList.Title = "Select your database"
	profileList.Styles.Title = profileList.Styles.Title.
		Background(darkPurple).
		Foreground(hiMagenta).
		Bold(true)
	profileList.SetShowStatusBar(false)
	profileList.SetShowHelp(true)
	profileList.KeyMap.Quit.Unbind()
	customKeys := newCustomKeyMap()
	profileList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.connect,
			customKeys.delete,
			customKeys.quit,
		}
	}

	profileList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			customKeys.connect,
			customKeys.delete,
			customKeys.quit,
		}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &ConnectModel{
		state:         stateLoading,
		spinner:       s,
		list:          profileList,
		loadingAction: "Fetch profiles from the config file...",
	}
}

// Run function is the main method called to run the database profiles menu.
func Run() (command.Options, error) {
	m := initConnectModel()
	model, err := tea.NewProgram(m).Run()
	if err != nil {
		return command.Options{}, nil
	}

	// Check if the connect method was ended by pressing ctrl+c.
	if cm, ok := model.(*ConnectModel); ok {
		if cm.aborted {
			return command.Options{}, ErrConnectionFormAborted
		}
	}

	// Get the selected profile to connect to.
	profile := m.profiles[m.selectedOption]
	// Get the password from the OS keyring.
	pass, err := keyring.Get(m.selectedOption, profile.User)
	if err != nil {
		return command.Options{}, nil
	}

	// add the password to the profile.
	profile.Pass = pass

	// if the profile contains ssh credentials.
	if profile.SSHUser != "" {
		// Get the ssh password, if any.
		sshPass, err := keyring.Get(m.selectedOption+"-ssh", profile.SSHUser)
		if err != nil {
			// If the error is different than ErrNotFound, return the error.
			if !errors.Is(err, keyring.ErrNotFound) {
				return command.Options{}, err
			}
		} else {
			profile.SSHPass = sshPass
		}
	}

	return profile, nil
}

// Init method starts with the spinner and fetching the profiles asynchronously.
func (m *ConnectModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchProfilesCmd())
}

// Update method  manages four different events:
// Windows re-size to re-calculate the size of the terminal,
// ctrl+c to exit the form,
// ctrl+d to delete a profile from the config file,
// enter to connect to a database profile.
// It also manages async messages and the menu state.
func (m *ConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(min(50, m.width)-6, 14)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "ctrl+d":
			if m.state == stateForm {
				if i, ok := m.list.SelectedItem().(dbProfile); ok {
					m.state = stateLoading
					index := m.list.Index()
					m.loadingAction = fmt.Sprintf("Deleting %s...", i.name)
					return m, tea.Batch(m.spinner.Tick, deleteItemCmd(i.name, index))
				}
			}
		case "enter":
			if m.state == stateForm {
				if i, ok := m.list.SelectedItem().(dbProfile); ok {
					m.selectedOption = i.name
					return m, tea.Quit
				}
			}
		}
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case itemDeletedMsg:
		m.list.RemoveItem(msg.index)
		m.state = stateForm
		return m, fetchProfilesCmd()
	case profilesLoadedMsg:
		m.profiles = msg.profiles
		cmd := m.list.SetItems(msg.items)
		m.state = stateForm
		return m, cmd
	}

	switch m.state {
	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stateForm:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View method shows the spinner if the data is loading,
// and the list of profiles if the data is already fetched from the config file.
func (m *ConnectModel) View() tea.View {
	var (
		v       tea.View
		content string
	)

	v.AltScreen = true

	switch m.state {
	case stateLoading:
		spinnerView := fmt.Sprintf("\n %s %s\n", m.spinner.View(), m.loadingAction)
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, spinnerView)
	case stateForm:
		content = setModalContent(m.list.View(), m.width, m.height)
	}

	v.SetContent(content)
	return v
}
