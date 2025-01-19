package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/codeshaine/curlify/tui"
)

func main() {
	// Initial state of the application

	model := tui.NewModel()

	// Create and start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
	}
}
