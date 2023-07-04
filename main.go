package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelice/lazypsql/app"
)

func main() {
	if _, err := tea.NewProgram(&app.Model{}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
