package neofetch

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/utils"
	"github.com/charmbracelet/lipgloss"
)

// RenderHome renders the neofetch-style home screen
func RenderHome(availableWidth int, portfolioOwner string, projects []sanity.Project, posts []sanity.Post, sessionTerminal, sessionID string) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	// Load logo and config
	cfg := config.Get()
	logo := cfg.Branding.AsciiLogo

	username := utils.GetUsername()
	if strings.TrimSpace(username) == "" {
		username = strings.ToLower(strings.ReplaceAll(portfolioOwner, " ", ""))
	}
	if strings.TrimSpace(username) == "" {
		username = "user"
	}

	host := utils.GetHost()
	if strings.TrimSpace(host) == "" {
		host = "localhost"
	}

	title := username + "@" + host
	if len(title) > 40 {
		title = title[:37] + "..."
	}

	projectsCount := fmt.Sprintf("%d", len(projects))
	postsCount := fmt.Sprintf("%d", len(posts))

	// Get host from config
	host = utils.GetHost()

	// Truncate session terminal name if too long
	termName := sessionTerminal
	if len(termName) > 30 {
		termName = termName[:27] + "..."
	}

	// Use role/tagline from config
	roleTagline := cfg.Owner.Role

	infoLines := []string{
		styles.NeoTitle.Render(title),
		styles.NeoSeparator.Render(strings.Repeat("─", min(len(title), 40))),
		styles.NeoLabel.Render(roleTagline),
		"",
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("App"), cfg.App.Name),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Version"), cfg.App.Version),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Runtime"), cfg.App.Runtime),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Host"), host),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("TERM"), termName),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Session"), sessionID),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Projects"), projectsCount),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Posts"), postsCount),
		"",
		RenderColorSwatches(),
	}
	info := strings.Join(infoLines, "\n")

	// 40% logo, 60% info
	logoWidth := availableWidth * 40 / 100
	infoWidth := availableWidth - logoWidth
	if infoWidth < 20 {
		infoWidth = 20
	}

	renderedLogo := lipgloss.NewStyle().
		Width(logoWidth).
		AlignHorizontal(lipgloss.Center).
		Render(logoStyle.Render(logo))

	renderedInfo := lipgloss.NewStyle().
		Width(infoWidth).
		Render(info)

	return lipgloss.JoinHorizontal(lipgloss.Center, renderedLogo, renderedInfo)
}
