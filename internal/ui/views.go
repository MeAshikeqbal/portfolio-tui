package ui

import (
	"fmt"
	"strings"

	appconfig "github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/helpmodal"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/menu"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/tocmodal"
	uiconfig "github.com/MeAshikeqbal/portfolio-tui/internal/ui/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/layout"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/neofetch"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/utils"
	"github.com/charmbracelet/lipgloss"
)

// headerView renders the header based on current state
func (m Model) headerView() string {
	breadcrumb := m.getBreadcrumb()
	return layout.RenderHeader(utils.GetFullName(), utils.GetTagline(), breadcrumb)
}

// getBreadcrumb returns the breadcrumb string for the current state.
// Returns "" for Home (no header) and any unhandled state.
func (m Model) getBreadcrumb() string {
	switch m.state {
	case blogDetailView, projectDetailView:
		// Breadcrumb is embedded in the detail box border, not the header
		return ""
	case menuView, contentView:
		selectedItem := m.menu[m.selected]
		if selectedItem == "Home" {
			return ""
		}
		return selectedItem
	}
	return ""
}

// footerView renders the footer with consistent controls
func (m Model) footerView(width int) string {
	var controls string

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
		controls = "↑/↓ scroll • t toc • esc back • ? help • q quit • " + scrollPct
	case projectDetailView:
		scrollPct := fmt.Sprintf("%3.0f%%", m.projectDetailViewport.ScrollPercent()*100)
		controls = "↑/↓ scroll • pgup/pgdn page • home/end top/bottom • esc back • ? help • q quit • " + scrollPct
	}

	return layout.RenderFooter(controls, width)
}

func (m Model) renderSidebar(width, height int) string {
	cfg := appconfig.Get()
	title := styles.SidebarTitle.Render(cfg.Branding.SidebarTitle)
	logo := neofetch.RenderSidebarLogo(width, m.portfolioOwner)
	body := menu.RenderMenu(m.menu, m.selected)

	// Use available height to avoid clipping in smaller terminals.
	constrainedHeight := max(1, height)

	box := lipgloss.NewStyle().
		Width(width).
		Height(constrainedHeight).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", logo, "", body))
}

func (m Model) renderContentPane(width, height int) string {
	selectedItem := m.menu[m.selected]

	// Use available height to avoid clipping in smaller terminals.
	constrainedHeight := max(1, height)

	// Calculate inner dimensions (accounting for padding and borders)
	innerWidth := width - 4
	innerHeight := constrainedHeight - 2

	var content string
	if selectedItem == "Home" {
		// Make Home content scrollable in case terminal is too small
		homeContent := neofetch.RenderHome(innerWidth, m.portfolioOwner, m.projects, m.posts, m.sessionTerminal, m.sessionID)
		m.viewport.Width = innerWidth
		m.viewport.Height = innerHeight
		m.viewport.SetContent(homeContent)
		content = m.viewport.View()
	} else if selectedItem == "Projects" {
		content = m.projectsList.View()
	} else if selectedItem == "Blogs" {
		content = m.blogList.View()
	} else {
		// All other text-based sections: unified styled title + scrollable viewport
		titleStr := styles.Title.Render(sectionIcon(selectedItem) + " " + selectedItem)
		viewportHeight := max(1, innerHeight-2) // reserve 2 lines for title + blank line
		m.viewport.Width = innerWidth
		m.viewport.Height = viewportHeight
		if m.state == menuView {
			textContent := m.content[selectedItem]
			if strings.TrimSpace(textContent) == "" {
				textContent = "Select an item from the sidebar and press Enter."
			}
			m.viewport.SetContent(utils.WrapText(textContent, innerWidth))
		}
		content = lipgloss.JoinVertical(lipgloss.Left, titleStr, "", m.viewport.View())
	}

	// Gray out inner content when browsing (not yet entered) for non-Home, non-Exit sections
	if m.state == menuView && selectedItem != "Home" && selectedItem != "Exit" {
		content = utils.DimContent(content)
	}

	box := lipgloss.NewStyle().
		Width(width).
		Height(constrainedHeight).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return box.Render(content)
}

func sectionIcon(item string) string {
	switch item {
	case "Skills":
		return "💻"
	case "Experience":
		return "💼"
	case "Education":
		return "🎓"
	case "Contact Me":
		return "📫"
	default:
		return "◆"
	}
}

func (m Model) renderShellLayout() string {
	// Calculate available dimensions
	termWidth := m.help.Width
	termHeight := m.viewport.Height

	// Ensure minimum terminal size
	if termWidth < uiconfig.MinTerminalWidth {
		return layout.SmallTerminalWidthWarning(termWidth)
	}

	if termHeight < uiconfig.MinContentHeight {
		return layout.SmallTerminalHeightWarning(termWidth)
	}

	sidebarWidth, contentWidth := layout.CalculateShellWidths(termWidth)

	height := max(1, termHeight)

	sidebar := m.renderSidebar(sidebarWidth, height)
	content := m.renderContentPane(contentWidth, height)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, " ", content)
}

// embedBreadcrumbInBorder replaces the top border line of a lipgloss-rendered
// box with one that contains the breadcrumb text: ╭─ Blogs > Title ──────╮
func embedBreadcrumbInBorder(rendered, breadcrumb string, borderStyle, textStyle lipgloss.Style) string {
	if breadcrumb == "" {
		return rendered
	}
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}
	topWidth := lipgloss.Width(lines[0])
	label := " " + breadcrumb + " "
	labelW := lipgloss.Width(label)
	rest := topWidth - 3 - labelW // "╭─" = 2, "╮" = 1
	if rest < 0 {
		rest = 0
	}
	lines[0] = borderStyle.Render("╭─") + textStyle.Render(label) + borderStyle.Render(strings.Repeat("─", rest)+"╮")
	return strings.Join(lines, "\n")
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

	if m.viewport.Height > 0 && m.viewport.Height < uiconfig.MinContentHeight {
		return layout.SmallTerminalHeightWarning(m.help.Width)
	}

	header := m.headerView()
	// Gray out header when browsing (not yet entered) for non-Home, non-Exit sections
	if m.state == menuView {
		selectedItem := m.menu[m.selected]
		if selectedItem != "Home" && selectedItem != "Exit" {
			header = utils.DimContent(header)
		}
	}
	footer := m.footerView(m.help.Width)

	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 2)

	borderColor := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	breadcrumbStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	var content string
	if m.state == blogDetailView {
		var bc string
		if m.selectedPost != nil {
			bc = "Blogs > " + m.selectedPost.Title
		}
		rendered := detailBox.Render(m.blogDetailViewport.View())
		content = embedBreadcrumbInBorder(rendered, bc, borderColor, breadcrumbStyle)
	} else if m.state == projectDetailView {
		var bc string
		if m.selectedProject != nil {
			bc = "Projects > " + m.selectedProject.Title
		}
		rendered := detailBox.Render(m.projectDetailViewport.View())
		content = embedBreadcrumbInBorder(rendered, bc, borderColor, breadcrumbStyle)
	} else {
		content = m.renderShellLayout()
	}

	mainView := header + "\n" + content + footer

	// Render help modal as a full-screen overlay (no appending — preserves layout)
	if m.showHelpModal {
		helpText := m.getHelpContent()
		modal := helpmodal.New(m.help.Width, m.termHeight, helpText)
		return modal.View()
	}

	// Render TOC modal as a full-screen overlay
	if m.showTOC && m.state == blogDetailView {
		modal := tocmodal.New(m.help.Width, m.termHeight, m.tocEntries)
		return modal.View()
	}

	return mainView
}

// getHelpContent generates context-aware help text
func (m Model) getHelpContent() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))

	divider := strings.Repeat("─", 40)

	var help strings.Builder

	help.WriteString(titleStyle.Render("KEYBOARD SHORTCUTS"))
	help.WriteString("\n")
	help.WriteString(divider)
	help.WriteString("\n\n")

	switch m.state {
	case menuView:
		help.WriteString(titleStyle.Render("Navigation"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  ↑/↓"))
		help.WriteString("    " + descStyle.Render("Navigate menu"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  enter"))
		help.WriteString("   " + descStyle.Render("View selected section"))
		help.WriteString("\n\n")

		help.WriteString(titleStyle.Render("Quick Jump"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  H"))
		help.WriteString("  " + descStyle.Render("Home     "))
		help.WriteString("  " + keyStyle.Render("P"))
		help.WriteString("  " + descStyle.Render("Projects"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  S"))
		help.WriteString("  " + descStyle.Render("Skills   "))
		help.WriteString("  " + keyStyle.Render("E"))
		help.WriteString("  " + descStyle.Render("Experience"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  D"))
		help.WriteString("  " + descStyle.Render("Education"))
		help.WriteString("  " + keyStyle.Render("B"))
		help.WriteString("  " + descStyle.Render("Blogs"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  C"))
		help.WriteString("  " + descStyle.Render("Contact  "))
		help.WriteString("  " + keyStyle.Render("X"))
		help.WriteString("  " + descStyle.Render("Exit"))
		help.WriteString("\n\n")

	case contentView:
		selectedItem := m.menu[m.selected]
		if selectedItem == "Projects" || selectedItem == "Blogs" {
			help.WriteString(titleStyle.Render("Navigation"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  ↑/↓"))
			help.WriteString("    " + descStyle.Render("Navigate items"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  enter"))
			help.WriteString("   " + descStyle.Render("View details"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  /"))
			help.WriteString("      " + descStyle.Render("Filter list"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  esc"))
			help.WriteString("    " + descStyle.Render("Back to menu"))
			help.WriteString("\n\n")
		} else {
			help.WriteString(titleStyle.Render("Scrolling"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  ↑/↓"))
			help.WriteString("      " + descStyle.Render("Scroll up/down"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  pgup/pgdn"))
			help.WriteString("  " + descStyle.Render("Page up/down"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  home/end"))
			help.WriteString("   " + descStyle.Render("Jump to top/bottom"))
			help.WriteString("\n\n")

			help.WriteString(titleStyle.Render("Navigation"))
			help.WriteString("\n")
			help.WriteString(keyStyle.Render("  esc"))
			help.WriteString("    " + descStyle.Render("Back to menu"))
			help.WriteString("\n\n")
		}

	case blogDetailView:
		help.WriteString(titleStyle.Render("Scrolling"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  ↑/↓"))
		help.WriteString("      " + descStyle.Render("Scroll up/down"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  pgup/pgdn"))
		help.WriteString("  " + descStyle.Render("Page up/down"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  home/end"))
		help.WriteString("   " + descStyle.Render("Jump to top/bottom"))
		help.WriteString("\n\n")

		help.WriteString(titleStyle.Render("Navigation"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  t"))
		help.WriteString("        " + descStyle.Render("Table of contents"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  esc"))
		help.WriteString("    " + descStyle.Render("Back to view"))
		help.WriteString("\n\n")

	case projectDetailView:
		help.WriteString(titleStyle.Render("Scrolling"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  ↑/↓"))
		help.WriteString("      " + descStyle.Render("Scroll up/down"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  pgup/pgdn"))
		help.WriteString("  " + descStyle.Render("Page up/down"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  home/end"))
		help.WriteString("   " + descStyle.Render("Jump to top/bottom"))
		help.WriteString("\n\n")

		help.WriteString(titleStyle.Render("Navigation"))
		help.WriteString("\n")
		help.WriteString(keyStyle.Render("  esc"))
		help.WriteString("    " + descStyle.Render("Back to view"))
		help.WriteString("\n\n")
	}

	help.WriteString(divider)
	help.WriteString("\n")
	help.WriteString(titleStyle.Render("Global"))
	help.WriteString("\n")
	help.WriteString(keyStyle.Render("  ?"))
	help.WriteString("       " + descStyle.Render("Toggle help"))
	help.WriteString("\n")
	help.WriteString(keyStyle.Render("  q"))
	help.WriteString("       " + descStyle.Render("Quit"))
	help.WriteString("\n\n")
	help.WriteString(descStyle.Render("Press ? or Esc to close"))

	return help.String()
}
