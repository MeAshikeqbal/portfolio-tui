package ui

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/keymap"
	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
)

type viewState int

const (
	menuView viewState = iota
	contentView
	blogDetailView
	projectDetailView
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
	menu                  []string
	selected              int
	viewport              viewport.Model
	blogDetailViewport    viewport.Model
	projectDetailViewport viewport.Model
	projectsList          list.Model
	blogList              list.Model
	help                  help.Model
	keys                  keymap.Map
	ready                 bool
	introComplete         bool
	introStage            int
	introTimerDone        bool
	state                 viewState
	content               map[string]string
	projects              []sanity.Project
	posts                 []sanity.Post
	selectedPost          *sanity.Post
	selectedProject       *sanity.Project
	portfolioOwner        string
	sessionUserIP         string
	sessionTerminal       string
	sessionID             string
	loadingState          loadingState
	error                 string
	sanityClient          *sanity.Client
	spinner               spinner.Model
	showHelpModal         bool
	showTOC               bool
	tocEntries            []blogmodule.TOCEntry
	termHeight            int
}
