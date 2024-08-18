package models

import tea "github.com/charmbracelet/bubbletea"

const (
	LANDING = "LANDING"
	LOGIN   = "LOGIN"
)

type MainModel struct {
	state   string
	landing LandingModel
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case LANDING:
		return m.landing.Update(msg)
	}
	return m, cmd
}

func (m MainModel) View() string {
	switch m.state {
	case LANDING:
		return m.landing.View()
	}
	return "Main Model View"
}

func MainModel_New() MainModel {
	return MainModel{
		state:   LANDING,
		landing: LandingModle_New(),
	}
}
