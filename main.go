package main

import (
	"fmt"
	"os"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	_, err := config.Load("")
	if err != nil {
		fmt.Println("Warning: Failed to load config, using defaults:", err)
	}

	p := tea.NewProgram(
		ui.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}