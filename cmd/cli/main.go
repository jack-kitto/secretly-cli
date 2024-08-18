package main

import (
	"fmt"
	"secretly-cli/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(models.MainModel_New())

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
