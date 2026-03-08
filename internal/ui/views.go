package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/charmbracelet/lipgloss"
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

	// Blog detail view has a special header with breadcrumb
	if m.state == blogDetailView && m.selectedPost != nil {
		title := styles.BlogTitleStyle.Render("Blogs > " + m.selectedPost.Title)
		line := strings.Repeat("─", max(0, m.blogDetailViewport.Width-lipgloss.Width(title)))
		return lipgloss.JoinHorizontal(lipgloss.Center, title, line) + "\n"
	}

	// Project detail view has a special header with breadcrumb
	if m.state == projectDetailView && m.selectedProject != nil {
		title := styles.BlogTitleStyle.Render("Projects > " + m.selectedProject.Title)
		line := strings.Repeat("─", max(0, m.projectDetailViewport.Width-lipgloss.Width(title)))
		return lipgloss.JoinHorizontal(lipgloss.Center, title, line) + "\n"
	}

	title := styles.Title.Render(name + " – " + tagline)
	if m.state == contentView {
		title += " > " + styles.SelectedItem.Render(m.menu[m.selected])
	}
	return title + "\n"
}

// footerView renders the footer based on current state
func (m Model) footerView() string {
	if m.help.ShowAll {
		return "\n" + m.help.View(m.keys)
	}

	// Show scroll percentage in detail views
	if m.state == blogDetailView {
		scrollInfo := styles.BlogScrollInfo.Render(fmt.Sprintf(" %3.f%% ", m.blogDetailViewport.ScrollPercent()*100))
		line := strings.Repeat("─", max(0, m.blogDetailViewport.Width-lipgloss.Width(scrollInfo)))
		return "\n" + lipgloss.JoinHorizontal(lipgloss.Center, line, scrollInfo)
	}

	if m.state == projectDetailView {
		scrollInfo := styles.BlogScrollInfo.Render(fmt.Sprintf(" %3.f%% ", m.projectDetailViewport.ScrollPercent()*100))
		line := strings.Repeat("─", max(0, m.projectDetailViewport.Width-lipgloss.Width(scrollInfo)))
		return "\n" + lipgloss.JoinHorizontal(lipgloss.Center, line, scrollInfo)
	}

	return "\n" + styles.Footer.Render(m.help.View(m.keys))
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

	meta := styles.SidebarMeta.Render(owner + "@portfolio")
	line := styles.NeoSeparator.Render(strings.Repeat("-", 20))
	return lipgloss.JoinVertical(lipgloss.Left, logoStyle.Render(logo), meta, line)
}

func renderNeoColorSwatches() string {
	palette := []string{"160", "166", "184", "114", "75", "69", "105", "250"}
	chips := make([]string, 0, len(palette))
	for _, c := range palette {
		chips = append(chips, lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("  "))
	}
	return strings.Join(chips, "")
}

func (m Model) renderNeoHome() string {
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
	projects := fmt.Sprintf("%d", len(m.projects))
	posts := fmt.Sprintf("%d", len(m.posts))

	info := strings.Join([]string{
		styles.NeoTitle.Render(title),
		styles.NeoSeparator.Render(strings.Repeat("-", len(title))),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("OS"), "Portfolio Linux"),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Host"), "BubbleTea Terminal"),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("TERM"), m.sessionTerminal),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Session"), m.sessionID),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Projects"), projects),
		fmt.Sprintf("%s: %s", styles.NeoLabel.Render("Posts"), posts),
		"",
		renderNeoColorSwatches(),
	}, "\n")

	return lipgloss.JoinHorizontal(lipgloss.Top, logoStyle.Render(logo), "  ", info)
}

func (m Model) renderSidebar(height int) string {
	title := styles.SidebarTitle.Render("NEOFETCH")
	logo := m.renderNeoSidebarLogo()
	body := m.renderMenu()

	box := lipgloss.NewStyle().
		Width(28).
		Height(height).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", logo, "", body))
}

func (m Model) renderContentPane(height int) string {
	selectedItem := m.menu[m.selected]

	var content string
	if selectedItem == "Home" {
		content = m.renderNeoHome()
	} else if m.state == menuView {
		content = m.content[selectedItem]
		if strings.TrimSpace(content) == "" {
			content = "Select an item from the sidebar and press Enter."
		}
	} else {
		if selectedItem == "Projects" {
			content = m.projectsList.View()
		} else if selectedItem == "Blogs" {
			content = m.blogList.View()
		} else {
			content = m.viewport.View()
		}
	}

	box := lipgloss.NewStyle().
		Width(max(40, m.help.Width-32)).
		Height(height).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(content)
}

func (m Model) renderShellLayout() string {
	height := max(10, m.viewport.Height)
	sidebar := m.renderSidebar(height)
	content := m.renderContentPane(height)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, " ", content)
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
