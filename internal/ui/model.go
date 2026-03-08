package ui

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

type Model struct {
	menu      []string
	selected  int
	viewport  viewport.Model
	help      help.Model
	keys      keyMap
	ready     bool
	state     viewState
	content   map[string]string
}

func InitialModel() Model {
	return Model{
		menu: []string{
			"Projects",
			"Skills",
			"About",
			"Contact",
			"Exit",
		},
		selected: 0,
		help:     help.New(),
		keys:     keys,
		state:    menuView,
		content:  generateContent(),
	}
}

func generateContent() map[string]string {
	return map[string]string{
		"Projects": `🚀 Projects

1. Portfolio TUI
   A beautiful terminal-based portfolio application built with Go and Bubble Tea.
   Technologies: Go, Bubble Tea, Bubbles, Lip Gloss
   
2. Web Dashboard
   Modern web dashboard with real-time analytics and data visualization.
   Technologies: React, TypeScript, D3.js, Node.js
   
3. API Gateway
   High-performance API gateway with authentication and rate limiting.
   Technologies: Go, Redis, PostgreSQL, Docker
   
4. Mobile App
   Cross-platform mobile application for task management.
   Technologies: React Native, Firebase, Redux
   
5. CLI Tool
   Command-line tool for automating development workflows.
   Technologies: Go, Cobra, Viper

Each project demonstrates different aspects of software development,
from frontend to backend, from mobile to CLI applications.`,

		"Skills": `💻 Technical Skills

Programming Languages:
  • Go - Advanced proficiency in concurrent programming
  • JavaScript/TypeScript - Full-stack development
  • Python - Data processing and automation
  • Rust - Systems programming

Frameworks & Libraries:
  • Bubble Tea - Terminal user interfaces
  • React - Modern web applications
  • Node.js - Backend services
  • FastAPI - REST APIs

Tools & Technologies:
  • Docker & Kubernetes - Container orchestration
  • Git & GitHub - Version control
  • PostgreSQL & MongoDB - Databases
  • Redis - Caching and queuing
  • Linux - System administration

Soft Skills:
  • Problem solving and critical thinking
  • Team collaboration and communication
  • Agile/Scrum methodologies
  • Technical documentation`,

		"About": `👨‍💻 About Me

Hi! I'm Ashik Eqbal, a passionate software developer with a love for
building elegant and efficient solutions.

Background:
I specialize in backend development, terminal applications, and
developer tools. My journey in software development started with
curiosity about how things work under the hood, which led me to
explore systems programming and low-level optimizations.

Interests:
  • Open source contribution
  • Terminal-based applications
  • Developer experience and tooling
  • Performance optimization
  • Clean code and architecture

Philosophy:
I believe in writing code that is not only functional but also
maintainable and enjoyable to work with. Good software should be
both powerful and elegant.

Currently working on building developer tools that make programming
more productive and fun!`,

		"Contact": `📫 Contact Information

Feel free to reach out through any of these channels:

Email:
  ashik@example.com

GitHub:
  github.com/MeAshikeqbal

LinkedIn:
  linkedin.com/in/ashikeqbal

Twitter:
  @ashikeqbal

Website:
  ashikeqbal.dev

I'm always open to interesting conversations about technology,
collaboration opportunities, or just to connect with fellow developers.

Response time: Usually within 24-48 hours.

Looking forward to hearing from you!`,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
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
	title := styles.Title.Render("Ashik Eqbal – Portfolio")
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