package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarWidth     = 30
	minContentWidth  = 50
	maxContentWidth  = 120
	minTerminalWidth = sidebarWidth + minContentWidth + 6
	minHeight        = 20
)

// headerView renders the header based on current state
func (m Model) headerView() string {
	name := os.Getenv("FULL_NAME")
	if name == "" {
		name = "Ashik Eqbal"
	}

	tagline := os.Getenv("TAGLINE")
	if tagline == "" {
		tagline = "Portfolio"
	}

	// Blog detail view breadcrumb
	if m.state == blogDetailView && m.selectedPost != nil {
		title := styles.Title.Render(name + " – " + tagline + " > Blogs > " + m.selectedPost.Title)
		return title + "\n"
	}

	// Project detail view breadcrumb
	if m.state == projectDetailView && m.selectedProject != nil {
		title := styles.Title.Render(name + " – " + tagline + " > Projects > " + m.selectedProject.Title)
		return title + "\n"
	}

	// For menuView and contentView, show consistent header
	if m.state == menuView || m.state == contentView {
		selectedItem := m.menu[m.selected]
		if selectedItem == "Home" {
			// Home page doesn't show header (neofetch style)
			return ""
		}
		// All other pages show: Name - Portfolio > Section
		title := styles.Title.Render(name + " – " + tagline + " > " + selectedItem)
		return title + "\n"
	}

	return ""
}

// footerView renders the footer with consistent controls
func (m Model) footerView() string {
	var controls string

	// Build context-specific help text
	if m.help.ShowAll {
		controls = m.help.View(m.keys)
	} else {
		// Show concise controls based on current state
		switch m.state {
		case menuView:
			selectedItem := m.menu[m.selected]
			if selectedItem == "Projects" || selectedItem == "Blogs" {
				controls = "↑/↓ navigate • enter activate • ? help • q quit"
			} else if selectedItem == "Exit" {
				controls = "↑/↓ navigate • enter exit • ? help • q quit"
			} else {
				controls = "↑/↓ navigate • enter view • ? help • q quit"
			}
		case contentView:
			selectedItem := m.menu[m.selected]
			if selectedItem == "Projects" || selectedItem == "Blogs" {
				controls = "↑/↓ navigate • / filter • enter open • esc back • ? help • q quit"
			} else {
				controls = "↑/↓ scroll • pgup/pgdn page • home/end top/bottom • esc back • ? help • q quit"
			}
		case blogDetailView:
			scrollPct := fmt.Sprintf("%3.0f%%", m.blogDetailViewport.ScrollPercent()*100)
			controls = "↑/↓ scroll • pgup/pgdn page • home/end top/bottom • esc back • ? help • q quit • " + scrollPct
		case projectDetailView:
			scrollPct := fmt.Sprintf("%3.0f%%", m.projectDetailViewport.ScrollPercent()*100)
			controls = "↑/↓ scroll • pgup/pgdn page • home/end top/bottom • esc back • ? help • q quit • " + scrollPct
		}
	}

	return "\n" + styles.Footer.Render(controls)
}

// renderMenu renders the main menu
func (m Model) renderMenu() string {
	var s strings.Builder
	sectionBreak := map[string]string{
		"Home":   "CORE",
		"Skills": "PROFILE",
		"Exit":   "SYSTEM",
	}

	for i, item := range m.menu {
		if section, ok := sectionBreak[item]; ok {
			if i > 0 {
				s.WriteString("\n")
			}
			s.WriteString(styles.SidebarSection.Render("[" + section + "]"))
			s.WriteString("\n")
		}

		entry := fmt.Sprintf("%s %s", menuIcon(item), item)
		if i == m.selected {
			s.WriteString(styles.SidebarActiveItem.Render("> " + entry))
		} else {
			s.WriteString(styles.SidebarItem.Render("  " + entry))
		}
		s.WriteString("\n")
	}

	return s.String()
}

func menuIcon(item string) string {
	switch item {
	case "Home":
		return "H"
	case "Projects":
		return "P"
	case "Skills":
		return "S"
	case "Experience":
		return "E"
	case "Education":
		return "D"
	case "Blogs":
		return "B"
	case "Contact Me":
		return "C"
	case "Exit":
		return "X"
	default:
		return ">"
	}
}

func (m Model) renderNeoSidebarLogo() string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	logo := strings.Join([]string{
		"      /\\",
		"     /  \\",
		"    / /\\ \\",
		"   / ____ \\",
		"  /_/    \\_\\",
	}, "\n")

	owner := strings.ToLower(strings.ReplaceAll(m.portfolioOwner, " ", ""))
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
	return lipgloss.JoinVertical(lipgloss.Center, logoStyle.Render(logo), meta, line)
}

func renderNeoColorSwatches() string {
	palette := []string{"160", "166", "184", "114", "75", "69", "105", "250"}
	chips := make([]string, 0, len(palette))
	for _, c := range palette {
		chips = append(chips, lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("  "))
	}
	return strings.Join(chips, "")
}

func (m Model) renderNeoHome(availableWidth int) string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	logo := strings.Join([]string{
		"        /\\",
		"       /  \\",
		"      / /\\ \\",
		"     / ____ \\",
		"    /_/    \\_\\",
	}, "\n")

	owner := strings.ToLower(strings.ReplaceAll(m.portfolioOwner, " ", ""))
	if owner == "" {
		owner = "user"
	}

	title := owner + "@portfolio-tui"
	if len(title) > 40 {
		title = title[:37] + "..."
	}

	projects := fmt.Sprintf("%d", len(m.projects))
	posts := fmt.Sprintf("%d", len(m.posts))

	// Get host from environment
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	// Truncate session terminal name if too long
	termName := m.sessionTerminal
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
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Session"), m.sessionID),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Projects"), projects),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Posts"), posts),
		"",
		renderNeoColorSwatches(),
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

func (m Model) renderSidebar(height int) string {
	title := styles.SidebarTitle.Render("PORTFOLIO")
	logo := m.renderNeoSidebarLogo()
	body := m.renderMenu()

	// Constrain height to prevent overflow
	constrainedHeight := max(minHeight, min(height, 60))

	box := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(constrainedHeight).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", logo, "", body))
}

func (m Model) renderContentPane(width, height int) string {
	selectedItem := m.menu[m.selected]

	// Constrain height to prevent overflow
	constrainedHeight := max(minHeight, min(height, 60))

	// Calculate inner dimensions (accounting for padding and borders)
	innerWidth := width - 4
	innerHeight := constrainedHeight - 2

	var content string
	if selectedItem == "Home" {
		// Make Home content scrollable in case terminal is too small
		homeContent := m.renderNeoHome(innerWidth)
		m.viewport.Width = innerWidth
		m.viewport.Height = innerHeight
		m.viewport.SetContent(homeContent)
		content = m.viewport.View()
	} else if selectedItem == "Projects" {
		// Show list for Projects
		listContent := m.projectsList.View()
		// In menuView, gray out the list to show it's not interactive
		if m.state == menuView {
			content = dimContent(listContent)
		} else {
			content = listContent
		}
	} else if selectedItem == "Blogs" {
		// Show list for Blogs
		listContent := m.blogList.View()
		// In menuView, gray out the list to show it's not interactive
		if m.state == menuView {
			content = dimContent(listContent)
		} else {
			content = listContent
		}
	} else if m.state == menuView {
		// For other text content in menuView, use viewport for scrolling
		textContent := m.content[selectedItem]
		if strings.TrimSpace(textContent) == "" {
			textContent = "Select an item from the sidebar and press Enter."
		}
		// Wrap long lines in content
		wrappedContent := wrapText(textContent, innerWidth)

		// Update viewport with wrapped content
		m.viewport.Width = innerWidth
		m.viewport.Height = innerHeight
		m.viewport.SetContent(wrappedContent)

		content = m.viewport.View()
	} else {
		// contentView for non-list items
		content = m.viewport.View()
	}

	box := lipgloss.NewStyle().
		Width(width).
		Height(constrainedHeight).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(content)
}

// dimContent applies a gray/dimmed style to content to show it's not interactive
func dimContent(content string) string {
	lines := strings.Split(content, "\n")
	dimmedLines := make([]string, len(lines))

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	for i, line := range lines {
		// Strip existing styles and apply dim color
		dimmedLines[i] = dimStyle.Render(stripANSI(line))
	}

	return strings.Join(dimmedLines, "\n")
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(str string) string {
	// Simple ANSI stripper - matches ESC sequences
	var result strings.Builder
	inEscape := false

	for i := 0; i < len(str); i++ {
		if str[i] == '\x1b' && i+1 < len(str) && str[i+1] == '[' {
			inEscape = true
			i++ // skip the '['
			continue
		}

		if inEscape {
			if (str[i] >= 'A' && str[i] <= 'Z') || (str[i] >= 'a' && str[i] <= 'z') {
				inEscape = false
			}
			continue
		}

		result.WriteByte(str[i])
	}

	return result.String()
}

func (m Model) renderShellLayout() string {
	// Calculate available dimensions
	termWidth := m.help.Width
	termHeight := m.viewport.Height

	// Ensure minimum terminal size
	if termWidth < minTerminalWidth {
		return "\n  Terminal too narrow. Please resize to at least " + fmt.Sprintf("%d", minTerminalWidth) + " columns."
	}

	// Calculate content width: total - sidebar - borders - spacing
	contentWidth := termWidth - sidebarWidth - 6
	if contentWidth < minContentWidth {
		contentWidth = minContentWidth
	} else if contentWidth > maxContentWidth {
		contentWidth = maxContentWidth
	}

	height := max(minHeight, min(termHeight, 60))

	sidebar := m.renderSidebar(height)
	content := m.renderContentPane(contentWidth, height)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, " ", content)
}

// wrapText wraps text to a specified width, preserving existing line breaks
func wrapText(text string, width int) string {
	if width < 10 {
		width = 10
	}

	lines := strings.Split(text, "\n")
	var wrapped []string

	for _, line := range lines {
		if len(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}

		// Simple word wrapping
		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, line)
			continue
		}

		currentLine := ""
		for _, word := range words {
			if currentLine == "" {
				currentLine = word
			} else if len(currentLine)+1+len(word) <= width {
				currentLine += " " + word
			} else {
				wrapped = append(wrapped, currentLine)
				currentLine = word
			}
		}
		if currentLine != "" {
			wrapped = append(wrapped, currentLine)
		}
	}

	return strings.Join(wrapped, "\n")
}

// View renders the entire view
func (m Model) View() string {
	if !m.introComplete {
		logo := "  _ _                 _     _ _         _\n" +
			" (_) |               | |   (_) |       | |\n" +
			"  _| |_ ___  __ _ ___| |__  _| | __  __| | _____   __\n" +
			" | | __/ __|/ _` / __| '_ \\| | |/ / / _` |/ _ \\ \\ / /\n" +
			" | | |_\\__ \\ (_| \\__ \\ | | | |   < | (_| |  __/\\ V /\n" +
			" |_|\\__|___/\\__,_|___/_| |_|_|_|\\_(_)__,_|\\___| \\_/"

		var sb strings.Builder
		sb.WriteString("\n")
		sb.WriteString(logo)

		if m.introStage >= 1 {
			sb.WriteString("\n\n")
			sb.WriteString("Welcome to ")
			sb.WriteString(m.portfolioOwner)
			sb.WriteString("'s Terminal Portfolio\n\n")
			sb.WriteString("Session Info\n")
			sb.WriteString("----------------------------\n")
			sb.WriteString("User IP: ")
			sb.WriteString(m.sessionUserIP)
			sb.WriteString("\n")
			sb.WriteString("Terminal: ")
			sb.WriteString(m.sessionTerminal)
			sb.WriteString("\n")
			sb.WriteString("Session ID: ")
			sb.WriteString(m.sessionID)
			sb.WriteString("\n")
		}

		if m.introStage >= 2 {
			sb.WriteString("\n")
			sb.WriteString("Loading portfolio...\n")
			sb.WriteString(m.spinner.View())

			if m.loadingState == loading {
				sb.WriteString(" Fetching content....\n")
			} else if m.loadingState == loaded {
				sb.WriteString(" Content loaded. Launching interface\n")
			} else if m.loadingState == failed {
				sb.WriteString(" Failed to fetch some content, using fallback data\n")
			} else {
				sb.WriteString(" Initializing interface\n")
			}
		}

		return sb.String()
	}

	if !m.ready {
		return "\n  Initializing..."
	}

	if m.loadingState == failed {
		return fmt.Sprintf("\n  ❌ Error: %s\n\n  Using fallback content...", m.error)
	}

	header := m.headerView()
	footer := m.footerView()

	var content string
	if m.state == blogDetailView {
		content = m.blogDetailViewport.View()
	} else if m.state == projectDetailView {
		content = m.projectDetailViewport.View()
	} else {
		content = m.renderShellLayout()
	}

	return header + "\n" + content + footer
}
