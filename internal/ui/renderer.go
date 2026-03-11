package ui

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui/components/tocmodal"
	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"
	"github.com/charmbracelet/lipgloss"
)

// SetRenderer rebinds global Lip Gloss styles to a session-aware renderer.
func SetRenderer(r *lipgloss.Renderer) {
	if r == nil {
		return
	}

	lipgloss.SetDefaultRenderer(r)
	styles.SetRenderer(r)
	blogmodule.SetRenderer(r)
	tocmodal.SetRenderer(r)
}
