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

	// Blog detail view has a special header
	if m.state == blogDetailView && m.selectedPost != nil {
		title := styles.BlogTitleStyle.Render(m.selectedPost.Title)
		line := strings.Repeat("─", max(0, m.blogDetailViewport.Width-lipgloss.Width(title)))
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

	// Show scroll percentage in blog detail view
	if m.state == blogDetailView {
		scrollInfo := styles.BlogScrollInfo.Render(fmt.Sprintf(" %3.f%% ", m.blogDetailViewport.ScrollPercent()*100))
		line := strings.Repeat("─", max(0, m.blogDetailViewport.Width-lipgloss.Width(scrollInfo)))
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
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.loadingState == loading {
		loadingText := styles.LoadingStyle.Render("Loading content from Sanity...")
		return fmt.Sprintf("\n\n  %s %s\n\n  %s\n",
			m.spinner.View(),
			loadingText,
			styles.Footer.Render("Fetching projects, skills, blog posts & more..."))
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
