package ui

import (
	"time"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/keymap"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InitialModel creates and initializes the UI model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	projectDelegate := list.NewDefaultDelegate()
	projectsList := list.New([]list.Item{}, projectDelegate, 0, 0)
	projectsList.Title = "🚀 Projects"
	projectsList.SetShowStatusBar(false)
	projectsList.SetShowHelp(false)
	projectsList.SetFilteringEnabled(true)
	projectsList.Styles.Title = styles.Title
	projectsList.Styles.TitleBar = lipgloss.NewStyle()

	blogDelegate := list.NewDefaultDelegate()
	blogList := list.New([]list.Item{}, blogDelegate, 0, 0)
	blogList.Title = "📝 Blog Posts"
	blogList.SetShowStatusBar(false)
	blogList.SetShowHelp(false)
	blogList.SetFilteringEnabled(true)
	blogList.Styles.Title = styles.Title
	blogList.Styles.TitleBar = lipgloss.NewStyle()

	return Model{
		menu: []string{
			"Home",
			"Projects",
			"Skills",
			"Experience",
			"Education",
			"Blogs",
			"Contact Me",
			"Exit",
		},
		selected:        0,
		projectsList:    projectsList,
		blogList:        blogList,
		help:            help.New(),
		keys:            keymap.Default(),
		introComplete:   false,
		introStage:      0,
		introTimerDone:  false,
		state:           menuView,
		content:         make(map[string]string),
		portfolioOwner:  resolvePortfolioOwner(),
		sessionUserIP:   resolveMaskedUserIP(),
		sessionTerminal: resolveTerminalName(),
		sessionID:       generateSessionID(),
		loadingState:    loading,
		sanityClient:    sanity.NewClient(),
		spinner:         s,
	}
}

// Init initializes the model and returns initial commands
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchContentCmd(),
		m.spinner.Tick,
		introStageCmd(1, 350*time.Millisecond),
		introStageCmd(2, 750*time.Millisecond),
		introDelayCmd(),
	)
}
