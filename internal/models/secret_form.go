package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jack-kitto/secretly-sdk"
)

type (
	errMsg error
)

type SecretFormModel struct {
	nameInput       textinput.Model
	valueInput      textinput.Model
	err             error
	choices         []secretly.Environment // items on the to-do list
	cursor          int                    // which to-do list item our cursor is pointing at
	selected        map[int]struct{}       // which to-do items are selected
	initialSelected map[int]struct{}       // which to-do items are selected
	envFocused      bool
	submitted       bool
	project         secretly.Project
	initialSecret   secretly.Secret
}

type (
	ADD_SECRET_COMPLETE_MSG struct{}
)

func SecretFormModel_New(project secretly.Project, initialSelections map[int]struct{}, initialSecret *secretly.Secret) SecretFormModel {
	name_ti := textinput.New()
	name_ti.Placeholder = "Name"
	name_ti.Focus()
	name_ti.CharLimit = 156
	name_ti.Width = 20

	value_ti := textinput.New()
	value_ti.Placeholder = "Value"
	value_ti.CharLimit = 156
	value_ti.Width = 20

	var initialSecretCopy secretly.Secret
	if initialSecret != nil {
		initialSecretCopy = *initialSecret
	}

	return SecretFormModel{
		nameInput:       name_ti,
		valueInput:      value_ti,
		err:             nil,
		choices:         project.Environments,
		selected:        make(map[int]struct{}),
		initialSelected: initialSelections,
		envFocused:      false,
		project:         project,
		initialSecret:   initialSecretCopy,
	}
}

func (m SecretFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SecretFormModel) BuildSecrets() ([]secretly.Secret, []secretly.Environment, []secretly.Environment) {
	var secrets []secretly.Secret
	value := m.valueInput.Value()
	name := m.nameInput.Value()

	// Environments to which the secret was added
	var addedEnvironments []secretly.Environment

	// Environments from which the secret was removed
	var removedEnvironments []secretly.Environment

	for i := range m.selected {
		environment := m.choices[i]
		var secret secretly.Secret

		if m.initialSecret == (secretly.Secret{}) {
			// If the initial secret is an empty struct, treat it as if it's nil
			secret = secretly.Secret_build(name, value, m.project, environment)
		} else {
			secret = m.initialSecret
			secret.Name = name
			secret.Value = value
		}

		secrets = append(secrets, secret)
		if _, wasInitiallySelected := m.initialSelected[i]; !wasInitiallySelected {
			addedEnvironments = append(addedEnvironments, environment)
		}
	}

	for i := range m.initialSelected {
		if _, isSelectedNow := m.selected[i]; !isSelectedNow {
			removedEnvironments = append(removedEnvironments, m.choices[i])
		}
	}

	return secrets, addedEnvironments, removedEnvironments
}

func (m SecretFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyShiftTab.String():
			if m.valueInput.Focused() {
				m.valueInput.Blur()
				m.nameInput.Focus()
			} else if m.envFocused {
				m.envFocused = false
				m.valueInput.Focus()
			}
		case tea.KeyUp.String(), "k":
			if m.envFocused {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case tea.KeyDown.String(), "j":
			if m.envFocused {
				if m.cursor < len(m.choices) {
					m.cursor++
				}
			}
		case tea.KeyTab.String():
			if m.nameInput.Focused() {
				m.nameInput.Blur()
				m.valueInput.Focus()
			} else if m.valueInput.Focused() {
				m.valueInput.Blur()
				m.envFocused = true
			} else {
				m.nameInput.Focus()
				m.envFocused = false
			}
		case tea.KeyEnter.String():
			if m.nameInput.Focused() {
				m.nameInput.Blur()
				m.valueInput.Focus()
			} else if m.valueInput.Focused() {
				m.valueInput.Blur()
				m.envFocused = true
			} else {
				m.envFocused = false
				m.submitted = true
				return m, nil
			}
		case tea.KeySpace.String():
			if m.envFocused {
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.nameInput.Focused() {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.valueInput, cmd = m.valueInput.Update(msg)
	}
	return m, cmd
}

func (m SecretFormModel) View() string {
	s := fmt.Sprintf(
		"What's your secret?\n\n%s",
		m.nameInput.View(),
	) + "\n"
	s += m.valueInput.View() + "\n\n Environments:\n"
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name)
	}
	return s
}
