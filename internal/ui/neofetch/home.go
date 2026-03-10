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
func RenderHome(availableWidth int, portfolioOwner string, projects []sanity.Project, posts []sanity.Post, skillsContent string) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
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

	roleTagline := cfg.Owner.Role

	groupLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true)
	groupSep := styles.NeoSeparator.Render("───────────")

	infoLines := []string{
		styles.NeoTitle.Render(title),
		styles.NeoSeparator.Render(strings.Repeat("─", min(len(title), 40))),
		styles.NeoLabel.Render(roleTagline),
		"",
		groupLabel.Render("Environment"),
		groupSep,
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("App"), cfg.App.Name),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Runtime"), cfg.App.Runtime),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Host"), host),
		"",
		groupLabel.Render("Portfolio"),
		groupSep,
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

	neofetchBlock := lipgloss.JoinHorizontal(lipgloss.Center, renderedLogo, renderedInfo)

	// ── Styles ──
	dimText := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	sectionHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	projectTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	projectDesc := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	blogTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	blogDate := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	bioStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Italic(true).
		PaddingLeft(1).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("62"))

	// Full width for rulers, content indented
	rulerWidth := max(10, availableWidth)
	indent := "  "

	ruler := func(text string) string {
		tl := len(text)
		pl := max(1, (rulerWidth-tl-2)/2)
		pr := max(1, rulerWidth-tl-2-pl)
		return sectionHeader.Render(strings.Repeat("─", pl) + " " + text + " " + strings.Repeat("─", pr))
	}

	truncate := func(s string, maxW int) string {
		if lipgloss.Width(s) <= maxW {
			return s
		}
		for len(s) > 0 && lipgloss.Width(s) > maxW-3 {
			s = s[:len(s)-1]
		}
		return s + "..."
	}

	contentWidth := max(20, availableWidth-2)

	var below strings.Builder

	// ── Bio ──
	bio := cfg.Owner.Bio
	if strings.TrimSpace(bio) != "" {
		below.WriteString("\n")
		below.WriteString(bioStyle.Render(bio))
		below.WriteString("\n")
	}

	// ── Skills ──
	skillLines := parseSkillNames(skillsContent)
	if len(skillLines) > 0 {
		below.WriteString("\n")
		below.WriteString(ruler("Skills"))
		below.WriteString("\n")
		skillText := strings.Join(skillLines, " • ")
		below.WriteString(indent + dimText.Render(truncate(skillText, contentWidth)))
	}

	// ── Featured Projects ──
	if len(projects) > 0 {
		below.WriteString("\n\n")
		below.WriteString(ruler("Featured Projects"))
		showProjects := min(3, len(projects))
		// Find max title width for alignment
		maxNameW := 0
		for i := 0; i < showProjects; i++ {
			if len(projects[i].Title) > maxNameW {
				maxNameW = len(projects[i].Title)
			}
		}
		for i := 0; i < showProjects; i++ {
			p := projects[i]
			below.WriteString("\n")
			paddedName := p.Title + strings.Repeat(" ", maxNameW-len(p.Title))
			if p.Description != "" {
				descMaxW := contentWidth - maxNameW - 5 // " – " + indent
				desc := truncate(p.Description, descMaxW)
				below.WriteString(indent + projectTitle.Render(paddedName) + projectDesc.Render("  – "+desc))
			} else {
				below.WriteString(indent + projectTitle.Render(paddedName))
			}
		}
	}

	// ── Latest Blog Posts ──
	if len(posts) > 0 {
		below.WriteString("\n\n")
		below.WriteString(ruler("Latest Blog Posts"))
		showCount := min(3, len(posts))
		for i := 0; i < showCount; i++ {
			post := posts[i]
			below.WriteString("\n")
			date := ""
			if post.PublishedAt != "" && len(post.PublishedAt) >= 10 {
				date = post.PublishedAt[:10]
			}
			dateW := lipgloss.Width(date)
			maxTitleW := contentWidth - dateW - 2
			title := truncate(post.Title, maxTitleW)
			if date != "" {
				titleW := lipgloss.Width(blogTitle.Render(title))
				gap := contentWidth - titleW - dateW
				if gap < 1 {
					gap = 1
				}
				below.WriteString(indent + blogTitle.Render(title) + strings.Repeat(" ", gap) + blogDate.Render(date))
			} else {
				below.WriteString(indent + blogTitle.Render(title))
			}
		}
	}

	below.WriteString("\n")

	return neofetchBlock + "\n" + below.String()
}

// parseSkillNames extracts skill names from the fetched skills content string.
func parseSkillNames(content string) []string {
	var skills []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "•") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "•"))
			if name != "" {
				skills = append(skills, name)
			}
		}
	}
	return skills
}
