package bubbletui

import (
	"errors"
	"fmt"
	"os"

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

type state int

const (
	stateLoading state = iota
	stateForm
)

type dbProfile struct {
	name  string
	label string
}

// Implement the list.Item interface
func (d dbProfile) Title() string       { return d.label }
func (d dbProfile) Description() string { return "" }
func (d dbProfile) FilterValue() string { return d.label }

// profilesLoadedMsg struct is a message what will carry the profiles back to the update loop.
type profilesLoadedMsg struct {
	items    []list.Item
	profiles map[string]command.Options
}

type errMsg struct{ err error }

type itemDeletedMsg struct {
	index int
}

type deleteFailedMsg struct {
	err  error
	item dbProfile
}

func fetchProfilesCmd() tea.Cmd {
	return func() tea.Msg {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		profiles, err := databaseProfiles.ReadProfiles(configDir)
		if err != nil {
			return errMsg{err: err}
		}

		var items []list.Item
		for name, profile := range profiles {
			var label string
			if profile.Driver == drivers.SQLite {
				label = fmt.Sprintf("%s (sqlite://%s)", name, profile.Host)
			} else {
				label = fmt.Sprintf("%s (%s://%s@%s)", name, profile.Driver, profile.User, profile.Host)
			}
			items = append(items, dbProfile{name: name, label: label})
		}

		return profilesLoadedMsg{items: items, profiles: profiles}
	}
}

func deleteItemCmd(name string, index int) tea.Cmd {
	return func() tea.Msg {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return errMsg{err: err}
		}

		if err := databaseProfiles.DeleteProfile(configDir, name); err != nil {
			return errMsg{err: err}
		}

		return itemDeletedMsg{index: index}
	}
}

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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &ConnectModel{
		state:         stateLoading,
		spinner:       s,
		list:          profileList,
		loadingAction: "Fetch profiles from the config file...",
	}
}

func Run() (command.Options, error) {
	m := initConnectModel()
	model, err := tea.NewProgram(m).Run()
	if err != nil {
		return command.Options{}, nil
	}

	if cm, ok := model.(*ConnectModel); ok {
		if cm.aborted {
			return command.Options{}, ErrConnectionFormAborted
		}
	}

	profile := m.profiles[m.selectedOption]
	pass, err := keyring.Get(m.selectedOption, profile.User)
	if err != nil {
		return command.Options{}, nil
	}

	profile.Pass = pass

	if profile.SSHUser != "" {
		sshPass, err := keyring.Get(m.selectedOption+"-ssh", profile.SSHUser)
		if err != nil {
			return command.Options{}, nil
		}

		profile.SSHPass = sshPass
	}

	return profile, nil
}

func (m *ConnectModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchProfilesCmd())
}

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
		formStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(neonPurple).
			Padding(1, 2)
		boxedForm := formStyle.Render(m.list.View())
		content = lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			boxedForm,
		)
	}

	v.SetContent(content)
	return v
}
