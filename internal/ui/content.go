package ui

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	contentmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/content"
	tea "github.com/charmbracelet/bubbletea"
)

// fetchContentCmd initiates fetching content from Sanity.
func fetchContentCmd() tea.Cmd {
	return func() tea.Msg {
		payload := contentmodule.Fetch(sanity.NewClient())
		return ContentMsg{
			data:     payload.Data,
			projects: payload.Projects,
			posts:    payload.Posts,
			err:      payload.Err,
		}
	}
}
