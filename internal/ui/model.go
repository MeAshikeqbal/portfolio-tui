package ui

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	menuView viewState = iota
	contentView
	blogDetailView
)

type loadingState int

const (
	idle loadingState = iota
	loading
	loaded
	failed
)

// ContentMsg carries fetched content from Sanity.
type ContentMsg struct {
	data     map[string]string
	projects []sanity.Project
	posts    []sanity.Post
	err      error
}

type Model struct {
	menu               []string
	selected           int
	viewport           viewport.Model
	blogDetailViewport viewport.Model
	projectsList       list.Model
	blogList           list.Model
	help               help.Model
	keys               keyMap
	ready              bool
	state              viewState
	content            map[string]string
	projects           []sanity.Project
	posts              []sanity.Post
	selectedPost       *sanity.Post
	loadingState       loadingState
	error              string
	sanityClient       *sanity.Client
	spinner            spinner.Model
}

func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	projectDelegate := list.NewDefaultDelegate()
	projectsList := list.New([]list.Item{}, projectDelegate, 0, 0)
	projectsList.Title = "🚀 Projects"
	projectsList.SetShowStatusBar(false)
	projectsList.SetFilteringEnabled(true)
	projectsList.Styles.Title = styles.Title

	blogDelegate := list.NewDefaultDelegate()
	blogList := list.New([]list.Item{}, blogDelegate, 0, 0)
	blogList.Title = "📝 Blog Posts"
	blogList.SetShowStatusBar(false)
	blogList.SetFilteringEnabled(true)
	blogList.Styles.Title = styles.Title

	return Model{
		menu: []string{
			"Projects",
			"Skills",
			"Blog",
			"About",
			"Contact",
			"Exit",
		},
		selected:     0,
		projectsList: projectsList,
		blogList:     blogList,
		help:         help.New(),
		keys:         keys,
		state:        menuView,
		content:      make(map[string]string),
		loadingState: loading,
		sanityClient: sanity.NewClient(),
		spinner:      s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchContentCmd(), m.spinner.Tick)
}
