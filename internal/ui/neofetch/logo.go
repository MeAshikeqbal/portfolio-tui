package neofetch

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderSidebarLogo renders the compact sidebar logo
func RenderSidebarLogo(width int, portfolioOwner string) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	// Load branding and identity from config
	cfg := config.Get()
	logo := cfg.Branding.SidebarLogo

	username := strings.TrimSpace(cfg.Owner.Username)
	if username == "" {
		username = strings.ToLower(strings.ReplaceAll(portfolioOwner, " ", ""))
	}
	if username == "" {
		username = "user"
	}

	host := strings.TrimSpace(cfg.Network.Host)
	if host == "" {
		host = "localhost"
	}

	// Truncate if too long
	userHost := username + "@" + host
	if len(userHost) > 24 {
		userHost = userHost[:21] + "..."
	}

	innerWidth := max(1, width-2)
	identityWidth := min(lipgloss.Width(userHost), innerWidth)
	if identityWidth < 1 {
		identityWidth = 1
	}

	meta := styles.SidebarMeta.Render(userHost)
	line := styles.NeoSeparator.Render(strings.Repeat("─", identityWidth))

	logoBlock := lipgloss.NewStyle().
		Width(innerWidth).
		AlignHorizontal(lipgloss.Center).
		Render(logoStyle.Render(logo))
	metaBlock := lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, meta)
	lineBlock := lipgloss.PlaceHorizontal(innerWidth, lipgloss.Center, line)

	return lipgloss.JoinVertical(lipgloss.Left, logoBlock, metaBlock, lineBlock)
}
