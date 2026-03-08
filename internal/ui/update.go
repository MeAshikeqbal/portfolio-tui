package ui

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/listitem"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/layout"
	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"
	projectmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/project"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case introStageMsg:
		if msg.stage > m.introStage {
			m.introStage = msg.stage
		}

	case introDoneMsg:
		m.introTimerDone = true
		maybeCompleteIntro(&m)

	case ContentMsg:
		if msg.err != nil {
			m.loadingState = failed
			m.error = msg.err.Error()
		} else {
			m.content = msg.data
			m.projects = msg.projects
			m.posts = msg.posts
			m.loadingState = loaded

			if len(msg.projects) > 0 {
				items := make([]list.Item, len(msg.projects))
				for i, p := range msg.projects {
					items[i] = listitem.ProjectItem{Data: p}
				}
				m.projectsList.SetItems(items)
			}

			if len(msg.posts) > 0 {
				items := make([]list.Item, len(msg.posts))
				for i, p := range msg.posts {
					items[i] = listitem.PostItem{Data: p}
				}
				m.blogList.SetItems(items)
			}
		}

		maybeCompleteIntro(&m)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView(msg.Width))
		verticalMarginHeight := headerHeight + footerHeight
		availableHeight := max(1, msg.Height-verticalMarginHeight)

		// Keep resize logic aligned with shell layout (25% sidebar / 75% content)
		_, contentWidth := layout.CalculateShellWidths(msg.Width)

		// Constrain list width to content pane
		listWidth := max(20, contentWidth-4) // account for padding

		if !m.ready {
			m.viewport = viewport.New(contentWidth, availableHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.MouseWheelEnabled = true
			m.viewport.MouseWheelDelta = 3

			m.blogDetailViewport = viewport.New(msg.Width, availableHeight)
			m.blogDetailViewport.YPosition = headerHeight
			m.blogDetailViewport.MouseWheelEnabled = true
			m.blogDetailViewport.MouseWheelDelta = 3

			m.projectDetailViewport = viewport.New(msg.Width, availableHeight)
			m.projectDetailViewport.YPosition = headerHeight
			m.projectDetailViewport.MouseWheelEnabled = true
			m.projectDetailViewport.MouseWheelDelta = 3
			m.ready = true
		} else {
			m.viewport.Width = contentWidth
			m.viewport.Height = availableHeight
			m.blogDetailViewport.Width = msg.Width
			m.blogDetailViewport.Height = availableHeight
			m.projectDetailViewport.Width = msg.Width
			m.projectDetailViewport.Height = availableHeight
		}

		m.projectsList.SetWidth(listWidth)
		m.projectsList.SetHeight(availableHeight)
		m.blogList.SetWidth(listWidth)
		m.blogList.SetHeight(availableHeight)
		m.help.Width = msg.Width
		m.termHeight = msg.Height

	case tea.MouseMsg:
		if m.state == blogDetailView {
			m.blogDetailViewport, cmd = m.blogDetailViewport.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.state == projectDetailView {
			m.projectDetailViewport, cmd = m.projectDetailViewport.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.state == menuView {
			// In menuView, allow viewport scrolling for Home and other text content
			selectedItem := m.menu[m.selected]
			if selectedItem != "Projects" && selectedItem != "Blogs" {
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else if m.state == contentView {
			// In contentView, allow full interaction with lists and viewports
			selectedItem := m.menu[m.selected]
			if selectedItem == "Projects" {
				m.projectsList, cmd = m.projectsList.Update(msg)
				cmds = append(cmds, cmd)
			} else if selectedItem == "Blogs" {
				m.blogList, cmd = m.blogList.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.KeyMsg:
		// Handle help modal close on Esc or ?
		if m.showHelpModal {
			if msg.String() == "esc" || key.Matches(msg, m.keys.Help) {
				m.showHelpModal = false
				return m, nil
			}
			// Allow quit even while modal is open
			if key.Matches(msg, m.keys.Quit) {
				return m, tea.Quit
			}
			// Ignore other keys while modal is open
			return m, nil
		}

		switch m.state {
		case menuView:
			return m.updateMenuView(msg)
		case contentView:
			return m.updateContentView(msg)
		case blogDetailView:
			return m.updateBlogDetailView(msg)
		case projectDetailView:
			return m.updateProjectDetailView(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// In menuView, only handle sidebar navigation - lists are visual only

	// Letter shortcuts - jump directly to the section
	if idx := menuShortcutIndex(msg.String()); idx >= 0 && idx < len(m.menu) {
		m.selected = idx
		m.viewport.GotoTop()
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Up):
		if m.selected > 0 {
			m.selected--
			// Reset viewport when changing selection
			m.viewport.GotoTop()
		}
	case key.Matches(msg, m.keys.Down):
		if m.selected < len(m.menu)-1 {
			m.selected++
			// Reset viewport when changing selection
			m.viewport.GotoTop()
		}
	case key.Matches(msg, m.keys.Help):
		m.showHelpModal = !m.showHelpModal
	case msg.String() == "enter":
		selectedItem := m.menu[m.selected]
		if selectedItem == "Exit" {
			return m, tea.Quit
		}
		// Transition to contentView for interaction
		m.state = contentView
		if selectedItem != "Projects" && selectedItem != "Blogs" {
			m.viewport.SetContent(m.content[selectedItem])
			m.viewport.GotoTop()
		}
	}

	return m, nil
}

// menuShortcutIndex maps a pressed key to its menu index, or -1 if not a shortcut.
func menuShortcutIndex(key string) int {
	switch key {
	case "h", "H":
		return 0 // Home
	case "p", "P":
		return 1 // Projects
	case "s", "S":
		return 2 // Skills
	case "e", "E":
		return 3 // Experience
	case "d", "D":
		return 4 // Education
	case "b", "B":
		return 5 // Blogs
	case "c", "C":
		return 6 // Contact Me
	case "x", "X":
		return 7 // Exit
	}
	return -1
}

func (m Model) updateContentView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selectedItem := m.menu[m.selected]

	if selectedItem == "Projects" {
		if m.projectsList.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, m.keys.Back):
				m.state = menuView
				return m, nil
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.showHelpModal = !m.showHelpModal
				return m, nil
			case msg.String() == "enter":
				if selected := m.projectsList.SelectedItem(); selected != nil {
					if projectItem, ok := selected.(listitem.ProjectItem); ok {
						m.selectedProject = &projectItem.Data
						m.projectDetailViewport.SetContent(projectmodule.RenderProjectContent(m.selectedProject))
						m.projectDetailViewport.GotoTop()
						m.state = projectDetailView
						return m, nil
					}
				}
			}
		}

		var cmd tea.Cmd
		m.projectsList, cmd = m.projectsList.Update(msg)
		return m, cmd
	}

	if selectedItem == "Blogs" {
		if m.blogList.FilterState() != list.Filtering {
			switch {
			case key.Matches(msg, m.keys.Back):
				m.state = menuView
				return m, nil
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.showHelpModal = !m.showHelpModal
				return m, nil
			case msg.String() == "enter":
				if selected := m.blogList.SelectedItem(); selected != nil {
					if postItem, ok := selected.(listitem.PostItem); ok {
						m.selectedPost = &postItem.Data
						m.blogDetailViewport.SetContent(blogmodule.RenderPostContent(m.selectedPost))
						m.blogDetailViewport.GotoTop()
						m.state = blogDetailView
						return m, nil
					}
				}
			}
		}

		var cmd tea.Cmd
		m.blogList, cmd = m.blogList.Update(msg)
		return m, cmd
	}

	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = menuView
		return m, nil
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.showHelpModal = !m.showHelpModal
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) updateBlogDetailView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = contentView
		m.selectedPost = nil
		return m, nil
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.showHelpModal = !m.showHelpModal
		return m, nil
	case key.Matches(msg, m.keys.Home):
		m.blogDetailViewport.GotoTop()
		return m, nil
	case key.Matches(msg, m.keys.End):
		m.blogDetailViewport.GotoBottom()
		return m, nil
	}

	var cmd tea.Cmd
	m.blogDetailViewport, cmd = m.blogDetailViewport.Update(msg)
	return m, cmd
}

func (m Model) updateProjectDetailView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = contentView
		m.selectedProject = nil
		return m, nil
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.showHelpModal = !m.showHelpModal
		return m, nil
	case key.Matches(msg, m.keys.Home):
		m.projectDetailViewport.GotoTop()
		return m, nil
	case key.Matches(msg, m.keys.End):
		m.projectDetailViewport.GotoBottom()
		return m, nil
	}

	var cmd tea.Cmd
	m.projectDetailViewport, cmd = m.projectDetailViewport.Update(msg)
	return m, cmd
}
