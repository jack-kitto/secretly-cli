package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	LANDING = "LANDING"
	LOGIN   = "LOGIN"
	PROJECT = "PROJECT"
)

type MainModel struct {
	state   string
	landing LandingModel
	login   LoginModel
	project ProjectModel
}

func (m MainModel) Init() tea.Cmd {
	// Initialize the spinner and other components
	if m.state == LOGIN {
		return m.login.Init()
	}
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case MsgSwitchToLogin:
		m.state = LOGIN
	case MsgExitApp:
		return m, tea.Quit
	case MsgSwitchToProject:
		m.state = PROJECT
	}

	switch m.state {
	case LANDING:
		updatedLanding, landingCmd := m.landing.Update(msg)
		m.landing = updatedLanding.(LandingModel)
		cmd = landingCmd
	case LOGIN:
		updatedLogin, loginCmd := m.login.Update(msg)
		m.login = updatedLogin.(LoginModel)
		cmd = loginCmd
	case PROJECT:
		updatedProject, projectCmd := m.project.Update(msg)
		m.project = updatedProject.(ProjectModel)
		cmd = projectCmd
	}

	return m, cmd
}

func (m MainModel) View() string {
	switch m.state {
	case LANDING:
		return m.landing.View()
	case LOGIN:
		return m.login.View()
	case PROJECT:
		return m.project.View()
	}
	return "Main Model View"
}

func MainModel_New() MainModel {
	return MainModel{
		state:   LANDING,
		landing: LandingModle_New(),
		login:   LoginModel_New(),
		project: ProjectModel_New(),
	}
}
