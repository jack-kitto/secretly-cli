package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jack-kitto/secretly-sdk"
)

var (
	ADDING_SECRET = "ADDING_SECRET"
	PROJECT_VIEW  = "PROJECT_VIEW"
)

type ProjectModel struct {
	project        secretly.Project
	table          table.Model
	state          string
	addSecretModel AddSecretModel
}

func (self *ProjectModel) UpdateTableRows() {
	columns := []table.Column{
		{Title: "Name", Width: 50},
		{Title: "Environment", Width: 20},
	}

	rows := []table.Row{}
	for _, environment := range self.project.Environments {
		for _, secret := range environment.Secrets {
			row := table.Row{
				secret.Name,
				environment.Name,
			}
			rows = append(rows, row)
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	self.table = t
}

func ProjectModel_New() ProjectModel {
	p := secretly.Project_fake()
	columns := []table.Column{
		{Title: "Name", Width: 50},
		{Title: "Environment", Width: 20},
	}

	rows := []table.Row{}
	for _, environment := range p.Environments {
		for _, secret := range environment.Secrets {
			row := table.Row{
				secret.Name,
				environment.Name,
			}
			rows = append(rows, row)
		}
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

	projectModel := ProjectModel{
		project:        p,
		state:          PROJECT_VIEW,
		table:          t,
		addSecretModel: AddSecretModel_New(p),
	}
	projectModel.UpdateTableRows()
	return projectModel
}

func (m ProjectModel) Init() tea.Cmd {
	if m.state == ADDING_SECRET {
		return m.addSecretModel.Init()
	}
	return nil
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == ADDING_SECRET {
		var addSecretCmd tea.Cmd
		var model tea.Model
		model, addSecretCmd = m.addSecretModel.Update(msg)
		m.addSecretModel = model.(AddSecretModel)
		if m.addSecretModel.submitted {
			secrets := m.addSecretModel.BuildSecrets()
			for _, s := range secrets {
				s.Print()
			}
			m.state = PROJECT_VIEW
			m.project.DistributeSecrets(m.addSecretModel.BuildSecrets())
			m.UpdateTableRows()
			m.addSecretModel = AddSecretModel_New(m.project)
		}
		return m, addSecretCmd
	}

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
		case "N", "n", "a", "A":
			if m.state == PROJECT_VIEW {
				m.state = ADDING_SECRET
				_, addSecretCmd := m.addSecretModel.Update(nil)
				return m, addSecretCmd
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m ProjectModel) View() string {
	if m.state == ADDING_SECRET {
		return m.addSecretModel.View()
	}
	s := fmt.Sprintf("\n\nProject: %s\n\n", m.project.Name)
	for _, environment := range m.project.Environments {
		s += fmt.Sprintf("\n- %s (%d Secrets)", environment.Name, len(environment.Secrets))
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
