package bubbletui

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	databaseProfiles "github.com/danvergara/dblab/internal/profiles"
	"github.com/danvergara/dblab/pkg/command"
)

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
			label := fmt.Sprintf("🗄️ %s (%s@%s)", name, profile.User, profile.Host)
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
}

func initConnectModel() ConnectModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return ConnectModel{
		state:   stateLoading,
		spinner: s,
	}
}

func Run() (command.Options, error) {
	m := initConnectModel()
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		return command.Options{}, nil
	}

	return m.profiles[m.selectedOption], nil
}

func (m ConnectModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchProfilesCmd())
}

func (m ConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	case profilesLoadedMsg:
		m.profiles = m.profiles
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select the connection database profile").
					Options(msg.options...).
					Value(&m.selectedOption),
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

func (m ConnectModel) View() tea.View {
	var (
		v       tea.View
		content string
	)

	v.AltScreen = true

	switch m.state {
	case stateLoading:
		content = fmt.Sprintf("\n %s Fetch profiles from the config file...\n", m.spinner.View())
	case stateForm:
		content = m.form.View()
	}

	v.SetContent(content)
	return v
}
