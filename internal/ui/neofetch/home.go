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
func RenderHome(availableWidth int, portfolioOwner string, projects []sanity.Project, posts []sanity.Post, skillsContent string, sessionTerminal, sessionID string) string {
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

	termName := sessionTerminal
	if len(termName) > 30 {
		termName = termName[:27] + "..."
	}

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

	neofetchBlock := lipgloss.JoinHorizontal(lipgloss.Center, renderedLogo, renderedInfo)

	// ── Styles ──
	dimText := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	sectionHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("111")).Bold(true)
	projectTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	projectDesc := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	blogTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	blogDate := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	// Use the info section width for uniform layout below neofetch
	sectionWidth := max(30, infoWidth)
	rulerWidth := max(10, sectionWidth-2)
	indent := "  "

	// helper to build a centered ruler
	ruler := func(text string) string {
		tl := len(text)
		pl := max(1, (rulerWidth-tl-2)/2)
		pr := max(1, rulerWidth-tl-2-pl)
		return sectionHeader.Render(strings.Repeat("─", pl) + " " + text + " " + strings.Repeat("─", pr))
	}

	// truncate long text to fit within section width
	truncate := func(s string, maxW int) string {
		if lipgloss.Width(s) <= maxW {
			return s
		}
		// rough byte trim
		for len(s) > 0 && lipgloss.Width(s) > maxW-3 {
			s = s[:len(s)-1]
		}
		return s + "..."
	}

	// Build content below neofetch in a uniform-width block
	var below strings.Builder

	// ── Welcome ──
	bio := cfg.Owner.Bio
	if strings.TrimSpace(bio) != "" {
		below.WriteString(dimText.Render(bio))
		below.WriteString("\n")
	}

	// ── Portfolio Overview ──
	below.WriteString("\n")
	below.WriteString(ruler("Portfolio Overview"))
	below.WriteString("\n")
	below.WriteString(fmt.Sprintf("%s%s  %s\n", indent, labelStyle.Render("Projects:"), valueStyle.Render(projectsCount)))
	below.WriteString(fmt.Sprintf("%s%s  %s\n", indent, labelStyle.Render("Blog Posts:"), valueStyle.Render(postsCount)))
	below.WriteString(fmt.Sprintf("%s%s  %s", indent, labelStyle.Render("Experience:"), valueStyle.Render(roleTagline)))

	// ── Skills ──
	skillLines := parseSkillNames(skillsContent)
	if len(skillLines) > 0 {
		below.WriteString("\n\n")
		below.WriteString(ruler("Skills"))
		below.WriteString("\n")
		skillText := strings.Join(skillLines, " • ")
		below.WriteString(indent + dimText.Render(truncate(skillText, sectionWidth-4)))
	}

	// ── Featured Project ──
	if len(projects) > 0 {
		below.WriteString("\n\n")
		below.WriteString(ruler("Featured Project"))
		p := projects[0]
		below.WriteString("\n")
		below.WriteString(indent + projectTitle.Render(p.Title))
		if p.Description != "" {
			desc := truncate(p.Description, sectionWidth-4)
			below.WriteString("\n")
			below.WriteString(indent + projectDesc.Render(desc))
		}
		var meta []string
		if len(p.Technologies) > 0 {
			meta = append(meta, strings.Join(p.Technologies, ", "))
		}
		if p.GitHubData != nil && p.GitHubData.Stars > 0 {
			meta = append(meta, fmt.Sprintf("★ %d", p.GitHubData.Stars))
		}
		if len(meta) > 0 {
			below.WriteString("\n")
			below.WriteString(indent + dimText.Render(truncate(strings.Join(meta, " • "), sectionWidth-4)))
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
			maxTitleW := sectionWidth - 4 - len(date) - 2
			title := truncate(post.Title, maxTitleW)
			if date != "" {
				below.WriteString(fmt.Sprintf("%s%s  %s", indent, blogTitle.Render(title), blogDate.Render(date)))
			} else {
				below.WriteString(indent + blogTitle.Render(title))
			}
		}
	}

	below.WriteString("\n")

	// Wrap the bottom block to match the info column width and center below neofetch
	belowBlock := lipgloss.NewStyle().
		Width(sectionWidth).
		Render(below.String())

	// Center the below block to align under the neofetch info area
	centeredBelow := lipgloss.NewStyle().
		Width(availableWidth).
		AlignHorizontal(lipgloss.Center).
		Render(belowBlock)

	return neofetchBlock + "\n" + centeredBelow
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
