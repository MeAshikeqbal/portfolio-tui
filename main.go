package main

import (
	"fmt"
	"os"

	"github.com/MeAshikeqbal/portfolio-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

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