package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// MsgSwitchToLogin is a message to switch the view to the login screen.
type MsgSwitchToLogin struct{}

// MsgExitApp is a message to quit the application.
type MsgExitApp struct{}

type LandingModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func LandingModle_New() LandingModel {
	return LandingModel{
		choices:  []string{"Login", "About", "Exit"},
		selected: make(map[int]struct{}),
	}
}

func (m LandingModel) Init() tea.Cmd {
	return nil
}

func (m LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

			switch m.cursor {
			case 0:
				return m, func() tea.Msg { return MsgSwitchToLogin{} }
			case 2:
				return m, func() tea.Msg { return MsgExitApp{} }
			}
		}
	}

	return m, nil
}

func (m LandingModel) View() string {
	s := fmt.Sprintf("\n\nWelcome to Secretly CLI!\n\n")

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
