package neofetch

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderColorSwatches renders the color palette at bottom of neofetch display
func RenderColorSwatches() string {
	palette := []string{"160", "166", "184", "114", "75", "69", "105", "250"}
	chips := make([]string, 0, len(palette))
	for _, c := range palette {
		chips = append(chips, lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("  "))
	}
	return strings.Join(chips, "")
}
