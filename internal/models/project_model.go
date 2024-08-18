package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
	return ProjectModel{
		project: project,
	}
}

func (m ProjectModel) Init() tea.Cmd {
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m ProjectModel) View() string {
	s := fmt.Sprintf("\n\nProject: %s\n\n", m.project.name)
	for _, environment := range m.project.environments {
		s += fmt.Sprintf("\n- %s (%d Secrets)", environment.name, len(environment.secrets))
	}
	return s
}
