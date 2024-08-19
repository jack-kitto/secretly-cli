package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	Secret struct {
		name          string
		value         string
		projectId     string
		environmentId string
	}
	Environment struct {
		id        string
		name      string
		projectId string
		secrets   []Secret
	}
	Project struct {
		environments []Environment
		name         string
		id           string
	}
)

type ProjectModel struct {
	project Project
	table   table.Model
}

func Secret_New(projectId string, environnentId string, name string, value string) Secret {
	return Secret{
		projectId:     projectId,
		environmentId: environnentId,
		name:          name,
		value:         name,
	}
}

func Environment_New(id string, name string, projectId string, secrets []Secret) Environment {
	return Environment{
		id:        id,
		name:      name,
		projectId: projectId,
		secrets:   secrets,
	}
}

func Project_New(environments []Environment, name string, id string) Project {
	return Project{
		environments: environments,
		name:         name,
		id:           id,
	}
}

func ProjectModel_New() ProjectModel {
	projectId := "123"
	devId := "1"
	stagingId := "2"
	productionId := "2"
	devSecret := Secret_New(projectId, devId, "DEV_SECRET_TEST", "HELLO_WORLD_DEV")
	prodSecret := Secret_New(projectId, productionId, "PROD_SECRET_TEST", "HELLO_WORLD_PROD")
	stagSecret := Secret_New(projectId, stagingId, "STAG_SECRET_TEST", "HELLO_WORLD_STAG")
	dev := Environment_New(devId, "development", projectId, []Secret{devSecret})
	stag := Environment_New(stagingId, "staging", projectId, []Secret{stagSecret})
	prod := Environment_New(productionId, "production", projectId, []Secret{prodSecret})
	project := Project_New([]Environment{dev, stag, prod}, "Project Foo", projectId)
	columns := []table.Column{
		{Title: "Name", Width: 50},
		{Title: "Environment", Width: 20},
	}

	rows := []table.Row{
		{devSecret.name, "development"},
		{stagSecret.name, "staging"},
		{prodSecret.name, "production"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return ProjectModel{
		project: project,
		table:   t,
	}
}

func (m ProjectModel) Init() tea.Cmd {
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m ProjectModel) View() string {
	s := fmt.Sprintf("\n\nProject: %s\n\n", m.project.name)
	for _, environment := range m.project.environments {
		s += fmt.Sprintf("\n- %s (%d Secrets)", environment.name, len(environment.secrets))
	}
	s += fmt.Sprintf("\n\n")
	s += baseStyle.Render(m.table.View()) + "\n"
	s += fmt.Sprintf("\n\n")
	s += "Controls:\n"
	s += "N - Create New Secret \n"
	s += "E - Edit Secret\n"
	s += "D - Delete Secret\n"

	return s
}
