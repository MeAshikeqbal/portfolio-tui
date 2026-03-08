package ui

import (
	"fmt"
	"strings"

	appconfig "github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/helpmodal"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/menu"
	uiconfig "github.com/MeAshikeqbal/portfolio-tui/internal/ui/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/layout"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/neofetch"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/utils"
	"github.com/charmbracelet/lipgloss"
)

// headerView renders the header based on current state
func (m Model) headerView() string {
	// Blog detail view breadcrumb
	if m.state == blogDetailView && m.selectedPost != nil {
		return layout.RenderHeader(utils.GetFullName(), utils.GetTagline(), "Blogs > "+m.selectedPost.Title)
	}

	// Project detail view breadcrumb
	if m.state == projectDetailView && m.selectedProject != nil {
		return layout.RenderHeader(utils.GetFullName(), utils.GetTagline(), "Projects > "+m.selectedProject.Title)
	}

	// For menuView and contentView, show consistent header
	if m.state == menuView || m.state == contentView {
		selectedItem := m.menu[m.selected]
		if selectedItem == "Home" {
			// Home page doesn't show header (neofetch style)
			return ""
		}
		// All other pages show: Name - Portfolio > Section
		return layout.RenderHeader(utils.GetFullName(), utils.GetTagline(), selectedItem)
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
		controls = "↑/↓ scroll • pgup/pgdn page • home/end top/bottom • esc back • ? help • q quit • " + scrollPct
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
		// Show list for Projects
		listContent := m.projectsList.View()
		// In menuView, gray out the list to show it's not interactive
		if m.state == menuView {
			content = utils.DimContent(listContent)
		} else {
			content = listContent
		}
	} else if selectedItem == "Blogs" {
		// Show list for Blogs
		listContent := m.blogList.View()
		// In menuView, gray out the list to show it's not interactive
		if m.state == menuView {
			content = utils.DimContent(listContent)
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
		wrappedContent := utils.WrapText(textContent, innerWidth)

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
	footer := m.footerView(m.help.Width)

	var content string
	if m.state == blogDetailView {
		content = m.blogDetailViewport.View()
	} else if m.state == projectDetailView {
		content = m.projectDetailViewport.View()
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

	case blogDetailView, projectDetailView:
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
