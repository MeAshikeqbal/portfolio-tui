package layout

import (
	"fmt"

	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/config"
	"github.com/charmbracelet/lipgloss"
)

// StyledWarningBox creates a styled warning box
func StyledWarningBox(message string, width int) string {
	boxWidth := min(width-4, 60)
	if boxWidth < 30 {
		boxWidth = 30
	}

	warningStyle := lipgloss.NewStyle().
		Width(boxWidth).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("202")).
		Foreground(lipgloss.Color("202")).
		Bold(true).
		AlignHorizontal(lipgloss.Center)

	containerStyle := lipgloss.NewStyle().
		Width(width).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	return containerStyle.Render(warningStyle.Render(message))
}

// SmallTerminalHeightWarning returns a styled warning for terminal too short
func SmallTerminalHeightWarning(width int) string {
	message := fmt.Sprintf("⚠ Terminal Too Short\n\nPlease resize to at least %d lines", config.MinContentHeight)
	return StyledWarningBox(message, width)
}

// SmallTerminalWidthWarning returns a styled warning for terminal too narrow
func SmallTerminalWidthWarning(width int) string {
	message := fmt.Sprintf("⚠ Terminal Too Narrow\n\nPlease resize to at least %d columns", config.MinTerminalWidth)
	return StyledWarningBox(message, width)
}
