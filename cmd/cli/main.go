package main

import (
	"fmt"
	"secretly-cli/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(models.MainModel_New())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
