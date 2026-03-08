package layout

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderFooter renders the footer with consistent controls
func RenderFooter(controls string, width int) string {
	footerStyle := styles.Footer
	if width > 0 {
		footerStyle = footerStyle.Width(width).AlignHorizontal(lipgloss.Center)
	}

	return "\n" + footerStyle.Render(controls)
}

// GetFooterHeight calculates footer height
func GetFooterHeight(footer string) int {
	return lipgloss.Height(footer)
}
