package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Back     key.Binding
	Quit     key.Binding
	Help     key.Binding
	PageDown key.Binding
	PageUp   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "b"),
		key.WithHelp("pgup/b", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "f"),
		key.WithHelp("pgdn/f", "page down"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to menu"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type viewState int

const (
	menuView viewState = iota
	contentView
)

type loadingState int

const (
	idle loadingState = iota
	loading
	loaded
	failed
)

// ContentMsg carries fetched content from Sanity
type ContentMsg struct {
	data  map[string]string
	err   error
}

type Model struct {
	menu         []string
	selected     int
	viewport     viewport.Model
	help         help.Model
	keys         keyMap
	ready        bool
	state        viewState
	content      map[string]string
	loadingState loadingState
	error        string
	sanityClient *sanity.Client
}

func InitialModel() Model {
	m := Model{
		menu: []string{
			"Projects",
			"Skills",
			"Blog",
			"About",
			"Contact",
			"Exit",
		},
		selected:     0,
		help:         help.New(),
		keys:         keys,
		state:        menuView,
		content:      make(map[string]string),
		loadingState: loading,
		sanityClient: sanity.NewClient(),
	}
	return m
}

// fetchContentCmd initiates fetching content from Sanity
func fetchContentCmd() tea.Cmd {
	return func() tea.Msg {
		client := sanity.NewClient()
		content := make(map[string]string)

		// Fetch Projects
		if projects, err := client.GetProjects(); err == nil {
			var sb strings.Builder
			sb.WriteString("🚀 Projects\n\n")
			for i, p := range projects {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, p.Title))
				sb.WriteString(fmt.Sprintf("   %s\n", p.Description))
				if len(p.Technologies) > 0 {
					sb.WriteString(fmt.Sprintf("   Technologies: %s\n", strings.Join(p.Technologies, ", ")))
				}
				sb.WriteString("\n")
			}
			content["Projects"] = sb.String()
		} else {
			content["Projects"] = "📡 Fetching projects...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultProjects
		}

		// Fetch Skills
		if skills, err := client.GetSkills(); err == nil {
			var sb strings.Builder
			sb.WriteString("💻 Technical Skills\n\n")
			for _, s := range skills {
				sb.WriteString(fmt.Sprintf("%s:\n", s.Category))
				for _, item := range s.Items {
					sb.WriteString(fmt.Sprintf("  • %s\n", item))
				}
				sb.WriteString("\n")
			}
			content["Skills"] = sb.String()
		} else {
			content["Skills"] = "📡 Fetching skills...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultSkills
		}

		// Fetch About
		if about, err := client.GetAbout(); err == nil {
			var sb strings.Builder
			sb.WriteString("👨‍💻 About Me\n\n")
			sb.WriteString(about.Content)
			sb.WriteString("\n\n")
			sb.WriteString(about.Background)
			content["About"] = sb.String()
		} else {
			content["About"] = "📡 Fetching about...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultAbout
		}

		// Fetch Contact
		if contacts, err := client.GetContacts(); err == nil {
			var sb strings.Builder
			sb.WriteString("📫 Contact Information\n\n")
			for _, c := range contacts {
				sb.WriteString(fmt.Sprintf("%s:\n  %s\n\n", c.Platform, c.Value))
			}
			content["Contact"] = sb.String()
		} else {
			content["Contact"] = "📡 Fetching contact...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultContact
		}

		// Fetch Blog Posts
		if posts, err := client.GetPosts(); err == nil {
			var sb strings.Builder
			sb.WriteString("📝 Blog Posts\n\n")
			if len(posts) == 0 {
				sb.WriteString("No posts published yet.\n")
			} else {
				for i, post := range posts {
					sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, post.Title))
					
					// Format published date
					if post.PublishedAt != "" {
						sb.WriteString(fmt.Sprintf("   Published: %s\n", formatDate(post.PublishedAt)))
					}
					
					// Show author if available (from dereferenced query)
					if post.Author != nil {
						if author, ok := post.Author.(string); ok && author != "" {
							sb.WriteString(fmt.Sprintf("   Author: %s\n", author))
						}
					}
					
					// Show categories if available (from dereferenced query)
					if len(post.Categories) > 0 {
						var cats []string
						for _, cat := range post.Categories {
							if catStr, ok := cat.(string); ok {
								cats = append(cats, catStr)
							}
						}
						if len(cats) > 0 {
							sb.WriteString(fmt.Sprintf("   Categories: %s\n", strings.Join(cats, ", ")))
						}
					}
					
					sb.WriteString("\n")
				}
			}
			content["Blog"] = sb.String()
		} else {
			content["Blog"] = "📡 Fetching blog posts...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultBlog
		}

		return ContentMsg{data: content, err: nil}
	}
}

// formatDate formats an ISO date string to a more readable format
func formatDate(isoDate string) string {
	// Simple formatting - just take the date part (YYYY-MM-DD)
	if len(isoDate) >= 10 {
		return isoDate[:10]
	}
	return isoDate
}

var (
	defaultProjects = `1. Portfolio TUI
   A beautiful terminal-based portfolio application built with Go and Bubble Tea.
   
2. Web Dashboard
   Modern web dashboard with real-time analytics.
   
3. API Gateway
   High-performance API gateway with authentication.
   
4. Mobile App
   Cross-platform mobile application.
   
5. CLI Tool
   Command-line tool for development workflows.`

	defaultSkills = `Programming Languages:
  • Go
  • JavaScript/TypeScript
  • Python
  
Frameworks & Libraries:
  • Bubble Tea
  • React
  • Node.js
  
Tools & Technologies:
  • Docker & Kubernetes
  • PostgreSQL & MongoDB
  • Redis`

	defaultAbout = `I'm a passionate software developer with a love for building
elegant and efficient solutions.

Background:
I specialize in backend development, terminal applications,
and developer tools.

Philosophy:
Good software should be both powerful and elegant.`

	defaultContact = `Email:
  ashik@example.com

GitHub:
  github.com/MeAshikeqbal

LinkedIn:
  linkedin.com/in/ashikeqbal

Twitter:
  @ashikeqbal`

	defaultBlog = `1. Welcome to My Blog
   Learn about my journey in software development
   Published: 2024-01-15
   
2. Building Terminal Applications with Go
   A deep dive into Bubble Tea framework
   Published: 2024-01-10
   
3. Why I Love Open Source
   Contributing to the community
   Published: 2024-01-05`
)

func (m Model) Init() tea.Cmd {
	return fetchContentCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case ContentMsg:
		// Handle fetched content from Sanity
		if msg.err != nil {
			m.loadingState = failed
			m.error = msg.err.Error()
		} else {
			m.content = msg.data
			m.loadingState = loaded
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.help.Width = msg.Width

	case tea.KeyMsg:
		if m.state == menuView {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
				if m.selected > 0 {
					m.selected--
				}

			case key.Matches(msg, m.keys.Down):
				if m.selected < len(m.menu)-1 {
					m.selected++
				}

			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

			case msg.String() == "enter":
				selectedItem := m.menu[m.selected]
				if selectedItem == "Exit" {
					return m, tea.Quit
				}
				// Switch to content view
				m.state = contentView
				content := m.content[selectedItem]
				m.viewport.SetContent(content)
				m.viewport.GotoTop()
			}
		} else {
			// Content view
			switch {
			case key.Matches(msg, m.keys.Back):
				m.state = menuView
				return m, nil

			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			}

			// Pass viewport key events
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) headerView() string {
	name := os.Getenv("FULL_NAME")
	if name == "" {
		name = "Ashik Eqbal"
	}
	
	tagline := os.Getenv("TAGLINE")
	if tagline == "" {
		tagline = "Portfolio"
	}
	
	title := styles.Title.Render(name + " – " + tagline)
	if m.state == contentView {
		title += " > " + styles.SelectedItem.Render(m.menu[m.selected])
	}
	return title + "\n"
}

func (m Model) footerView() string {
	if m.help.ShowAll {
		return "\n" + m.help.View(m.keys)
	}
	return "\n" + styles.Footer.Render(m.help.View(m.keys))
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.loadingState == loading {
		return "\n  📡 Loading content from Sanity...\n\n  Please wait..."
	}

	if m.loadingState == failed {
		return fmt.Sprintf("\n  ❌ Error: %s\n\n  Using fallback content...", m.error)
	}

	header := m.headerView()
	footer := m.footerView()

	var content string
	if m.state == menuView {
		content = m.renderMenu()
	} else {
		content = m.viewport.View()
	}

	return header + "\n" + content + footer
}

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