package models

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jack-kitto/secretly-sdk"
)

var (
	ADD_SECRET_FORM    = "ADD_SECRET_FORM"
	UPDATE_SECRET_FORM = "UPDATE_SECRET_FORM"
	PROJECT_VIEW       = "PROJECT_VIEW"
)

type ProjectModel struct {
	project               secretly.Project
	table                 table.Model
	state                 string
	addSecretFormModel    SecretFormModel
	updateSecretFormModel SecretFormModel
}

func (self *ProjectModel) UpdateTableRows() {
	columns := []table.Column{
		{Title: "ID", Width: 50},
		{Title: "Name", Width: 50},
		{Title: "Environment", Width: 20},
	}

	rows := []table.Row{}
	for _, environment := range self.project.Environments {
		for _, secret := range environment.Secrets {
			row := table.Row{
				secret.ID,
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
		project:            p,
		state:              PROJECT_VIEW,
		table:              t,
		addSecretFormModel: SecretFormModel_New(p, make(map[int]struct{}), nil),
	}
	projectModel.UpdateTableRows()
	return projectModel
}

func (p *ProjectModel) DeleteSecretById(s secretly.Secret) {
	for envIndex, environment := range p.project.Environments {
		for secretIndex, secret := range environment.Secrets {
			if secret.ID == s.ID {
				// Remove the secret from the environment's Secrets slice
				p.project.Environments[envIndex].Secrets = append(
					environment.Secrets[:secretIndex],
					environment.Secrets[secretIndex+1:]...,
				)
				return // Exit once the secret is found and deleted
			}
		}
	}
}

func (m ProjectModel) Init() tea.Cmd {
	if m.state == ADD_SECRET_FORM {
		return m.addSecretFormModel.Init()
	}
	return nil
}

func (m *ProjectModel) GetSecretAtCursor() (*secretly.Secret, error) {
	index := m.table.Cursor()
	var secretID string
	for i, row := range m.table.Rows() {
		if i == index {
			secretID = row[0]
		}
	}
	for _, environment := range m.project.Environments {
		for _, secret := range environment.Secrets {
			if secretID == secret.ID {
				return &secret, nil
			}
		}
	}
	return nil, errors.New("No secret at cursor")
}

func (m *ProjectModel) EditSecretAtCursor() {
	secret, error := m.GetSecretAtCursor()
	if error != nil {
		return
	}

	initialSelections := make(map[int]struct{})
	selectedSelections := make(map[int]struct{})

	for i, environment := range m.project.Environments {
		if secret.InEnvironment(environment) {
			initialSelections[i] = struct{}{}
			selectedSelections[i] = struct{}{}
		}
	}

	m.updateSecretFormModel = SecretFormModel_New(m.project, initialSelections, secret)
	m.updateSecretFormModel.valueInput.SetValue(secret.Value)
	m.updateSecretFormModel.nameInput.SetValue(secret.Name)
	m.updateSecretFormModel.selected = selectedSelections

	// Change the state to UPDATE_SECRET_FORM
	m.state = UPDATE_SECRET_FORM
}

func (m *ProjectModel) DeleteSecretAtCursor() {
	secret, error := m.GetSecretAtCursor()
	if error != nil {
		return
	}
	m.DeleteSecretById(*secret)
	m.UpdateTableRows()
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == ADD_SECRET_FORM {
		var addSecretCmd tea.Cmd
		var model tea.Model
		model, addSecretCmd = m.addSecretFormModel.Update(msg)
		m.addSecretFormModel = model.(SecretFormModel)
		if m.addSecretFormModel.submitted {
			secrets, _, _ := m.addSecretFormModel.BuildSecrets()
			m.state = PROJECT_VIEW
			secrets, _, _ = m.addSecretFormModel.BuildSecrets()
			m.project.DistributeSecrets(secrets)
			m.UpdateTableRows()
			m.addSecretFormModel = SecretFormModel_New(m.project, make(map[int]struct{}), nil)
		}
		return m, addSecretCmd
	}

	if m.state == UPDATE_SECRET_FORM {
		var updateSecretCmd tea.Cmd
		var model tea.Model
		model, updateSecretCmd = m.updateSecretFormModel.Update(msg)
		m.updateSecretFormModel = model.(SecretFormModel)
		if m.updateSecretFormModel.submitted {
			secrets, _, _ := m.updateSecretFormModel.BuildSecrets()
			for _, updatedSecret := range secrets {
				m.DeleteSecretById(updatedSecret)
			}
			m.project.DistributeSecrets(secrets)

			m.state = PROJECT_VIEW
			m.UpdateTableRows()
			m.updateSecretFormModel = SecretFormModel_New(m.project, nil, nil)
		}
		return m, updateSecretCmd
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
		case "d":
			m.DeleteSecretAtCursor()
			return m, nil
		case "e":
			m.EditSecretAtCursor()
			return m, m.updateSecretFormModel.Init()
		case "N", "n", "a", "A":
			if m.state == PROJECT_VIEW {
				m.state = ADD_SECRET_FORM
				return m, m.addSecretFormModel.Init()
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
	if m.state == ADD_SECRET_FORM {
		return m.addSecretFormModel.View()
	}
	if m.state == UPDATE_SECRET_FORM {
		return m.updateSecretFormModel.View()
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
