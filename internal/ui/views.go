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

	for i, item := range m.menu {
		if i == m.selected {
			s.WriteString(styles.SelectedItem.Render("> " + item))
		} else {
			s.WriteString(styles.MenuItem.Render(item))
		}
		s.WriteString("\n")
	}

	return s.String()
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
	if m.state == menuView {
		content = m.renderMenu()
	} else if m.state == blogDetailView {
		content = m.blogDetailViewport.View()
	} else if m.state == projectDetailView {
		content = m.projectDetailViewport.View()
	} else {
		// Render appropriate content based on selected menu item
		selectedItem := m.menu[m.selected]
		if selectedItem == "Projects" {
			content = m.projectsList.View()
		} else if selectedItem == "Blog" {
			content = m.blogList.View()
		} else {
			content = m.viewport.View()
		}
	}

	return header + "\n" + content + footer
}
