package main

import (
	"fmt"
	"os"
	"secretly-cli/internal/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "secretly-cli",
	Short: "Securely share environment variables with your team.",
	Long: `Secretly is a powerful tool for securely managing and sharing environment variables with your team. 
       Integrated with the Secretly website at https://secretly.kitto.sh, it ensures your secrets are handled safely and easily, 
       so you can focus on building your projects without the worry of exposing sensitive information. 
       Complete documentation is available at https://secretly.kitto.sh/docs.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(models.MainModel_New())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
