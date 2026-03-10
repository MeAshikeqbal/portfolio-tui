package ui

import (
	"strings"
	"time"
	"unicode"

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
// SessionInfo holds session metadata, optionally provided by the SSH server.
type SessionInfo struct {
	UserIP   string
	Terminal string
}

func InitialModel() Model {
	return initialModel(nil)
}

// InitialModelWithSession creates a model pre-populated with SSH session info.
func InitialModelWithSession(info SessionInfo) Model {
	return initialModel(&info)
}

func initialModel(session *SessionInfo) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	projectDelegate := list.NewDefaultDelegate()
	projectsList := list.New([]list.Item{}, projectDelegate, 0, 0)
	projectsList.Title = "🚀 Projects"
	projectsList.SetShowStatusBar(false)
	projectsList.SetShowHelp(false)
	projectsList.SetFilteringEnabled(true)
	projectsList.Filter = strictContainsFilter
	projectsList.Styles.Title = styles.Title
	projectsList.Styles.TitleBar = lipgloss.NewStyle()

	blogDelegate := list.NewDefaultDelegate()
	blogList := list.New([]list.Item{}, blogDelegate, 0, 0)
	blogList.Title = "📝 Blog Posts"
	blogList.SetShowStatusBar(false)
	blogList.SetShowHelp(false)
	blogList.SetFilteringEnabled(true)
	blogList.Filter = strictContainsFilter
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
		sessionUserIP:   resolveSessionField(session, "ip"),
		sessionTerminal: resolveSessionField(session, "terminal"),
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

// strictContainsFilter performs case-insensitive substring matching for all
// typed terms and preserves original item ordering.
func strictContainsFilter(term string, targets []string) []list.Rank {
	query := strings.TrimSpace(strings.ToLower(term))
	if query == "" {
		return nil
	}

	tokens := strings.Fields(query)
	if len(tokens) == 0 {
		return nil
	}

	ranks := make([]list.Rank, 0, len(targets))
	for i, target := range targets {
		lowerTarget := strings.ToLower(target)
		targetWords := tokenizeWords(lowerTarget)
		matchPos := -1
		matchToken := ""
		matchedAll := true

		for _, token := range tokens {
			if len(token) <= 2 {
				if !containsWord(targetWords, token) {
					matchedAll = false
					break
				}
				continue
			}

			pos := strings.Index(lowerTarget, token)
			if pos == -1 {
				matchedAll = false
				break
			}
			if matchPos == -1 || pos < matchPos {
				matchPos = pos
				matchToken = token
			}
		}

		if !matchedAll {
			continue
		}

		matchedIndexes := make([]int, 0, len(matchToken))
		if matchPos >= 0 {
			for j := 0; j < len(matchToken); j++ {
				matchedIndexes = append(matchedIndexes, matchPos+j)
			}
		}

		ranks = append(ranks, list.Rank{
			Index:          i,
			MatchedIndexes: matchedIndexes,
		})
	}

	return ranks
}

func tokenizeWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func containsWord(words []string, token string) bool {
	for _, w := range words {
		if w == token {
			return true
		}
	}
	return false
}
