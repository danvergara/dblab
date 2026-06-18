package bubbletui

import (
	"errors"
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
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

// profilesLoadedMsg struct is a message what will carry the profiles back to the update loop.
type profilesLoadedMsg struct {
	options  []huh.Option[string]
	profiles map[string]command.Options
}

type errMsg struct{ err error }

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

		var options []huh.Option[string]
		for name, profile := range profiles {
			var label string
			if profile.Driver == drivers.SQLite {
				label = fmt.Sprintf("%s (sqlite://%s)", name, profile.Host)
			} else {
				label = fmt.Sprintf("%s (%s://%s@%s)", name, profile.Driver, profile.User, profile.Host)
			}
			options = append(options, huh.NewOption(label, name))
		}

		return profilesLoadedMsg{options: options, profiles: profiles}
	}
}

type ConnectModel struct {
	state          state
	form           *huh.Form
	spinner        spinner.Model
	selectedOption string
	err            error
	profiles       map[string]command.Options
	width          int
	height         int
	aborted        bool
}

func initConnectModel() *ConnectModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return &ConnectModel{
		state:   stateLoading,
		spinner: s,
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
		if m.form != nil {
			totalWidth := min(50, m.width)
			m.form.WithWidth(totalWidth - 6)
		}
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			m.aborted = true
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case profilesLoadedMsg:
		totalWidth := min(50, m.width)
		formWidth := totalWidth - 6

		m.profiles = msg.profiles
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select the connection database profile").
					Options(msg.options...).
					Value(&m.selectedOption).WithWidth(formWidth),
			),
		)

		m.state = stateForm
		return m, m.form.Init()
	}

	switch m.state {
	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case stateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		if m.form.State == huh.StateCompleted {
			return m, tea.Quit
		}
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
		content = fmt.Sprintf("\n %s Fetch profiles from the config file...\n", m.spinner.View())
	case stateForm:
		formStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)
		boxedForm := formStyle.Render(m.form.View())
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
