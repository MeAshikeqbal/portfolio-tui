package neofetch

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/utils"
	"github.com/charmbracelet/lipgloss"
)

// RenderHome renders the neofetch-style home screen
func RenderHome(availableWidth int, portfolioOwner string, projects []sanity.Project, posts []sanity.Post, sessionTerminal, sessionID string) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	logo := strings.Join([]string{
		"        /\\",
		"       /  \\",
		"      / /\\ \\",
		"     / ____ \\",
		"    /_/    \\_\\",
	}, "\n")

	owner := strings.ToLower(strings.ReplaceAll(portfolioOwner, " ", ""))
	if owner == "" {
		owner = "user"
	}

	title := owner + "@portfolio-tui"
	if len(title) > 40 {
		title = title[:37] + "..."
	}

	projectsCount := fmt.Sprintf("%d", len(projects))
	postsCount := fmt.Sprintf("%d", len(posts))

	// Get host from environment
	host := utils.GetHost()

	// Truncate session terminal name if too long
	termName := sessionTerminal
	if len(termName) > 30 {
		termName = termName[:27] + "..."
	}

	// Role/tagline
	tagline := "Developer • DevOps • Homelab Builder"

	infoLines := []string{
		styles.NeoTitle.Render(title),
		styles.NeoSeparator.Render(strings.Repeat("─", min(len(title), 40))),
		styles.NeoLabel.Render(tagline),
		"",
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("App"), "portfolio-tui"),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Version"), "v1.0"),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Runtime"), "Bubble Tea"),
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
