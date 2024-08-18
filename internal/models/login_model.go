package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginModel struct {
	spinner       spinner.Model
	completedAuth bool
}

func LoginModel_New() LoginModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return LoginModel{
		spinner: s,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m LoginModel) View() string {
	s := fmt.Sprintf("\n\nDevice Code Login\n\n")

	s += fmt.Sprintf("\n\n 1. Open the following URL in your web browser \n\n")
	s += fmt.Sprintf("       >  https://secretly.kitto.sh/device?code=XXXX-YYYY")
	s += fmt.Sprintf("\n\n 2. Follow the instructions on the webpage to complete the authentication process. \n\n")
	s += fmt.Sprintf("\n\n 3. After completing the authentication, press Enter here to finalize the login process. \n\n")
	s += fmt.Sprintf("\n\n%s  Waiting for authentication... \n\n", m.spinner.View())

	// Footer
	s += "\nPress q to quit.\n"

	return s
}
