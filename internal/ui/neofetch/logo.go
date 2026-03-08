package neofetch

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderSidebarLogo renders the compact sidebar logo
func RenderSidebarLogo(width int, portfolioOwner string) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	logo := strings.Join([]string{
		"      /\\",
		"     /  \\",
		"    / /\\ \\",
		"   / ____ \\",
		"  /_/    \\_\\",
	}, "\n")

	owner := strings.ToLower(strings.ReplaceAll(portfolioOwner, " ", ""))
	if owner == "" {
		owner = "user"
	}

	// Truncate if too long
	userHost := owner + "@portfolio"
	if len(userHost) > 24 {
		userHost = userHost[:21] + "..."
	}

	meta := styles.SidebarMeta.Render(userHost)
	line := styles.NeoSeparator.Render(strings.Repeat("─", min(len(userHost), 20)))

	innerWidth := max(1, width-2)
	logoBlock := lipgloss.NewStyle().
		Width(innerWidth).
		AlignHorizontal(lipgloss.Center).
		Render(logoStyle.Render(logo))
	metaBlock := lipgloss.NewStyle().
		Width(innerWidth).
		AlignHorizontal(lipgloss.Center).
		Render(meta)
	lineBlock := lipgloss.NewStyle().
		Width(innerWidth).
		AlignHorizontal(lipgloss.Center).
		Render(line)

	return lipgloss.JoinVertical(lipgloss.Left, logoBlock, metaBlock, lineBlock)
}
