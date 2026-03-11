package tocmodal

import (
	"fmt"
	"strings"

	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  lipgloss.Style
	numStyle    lipgloss.Style
	entryStyle  lipgloss.Style
	hintStyle   lipgloss.Style
	borderStyle lipgloss.Style
)

func init() {
	SetRenderer(nil)
}

func SetRenderer(r *lipgloss.Renderer) {
	newStyle := lipgloss.NewStyle
	if r != nil {
		newStyle = r.NewStyle
	}

	titleStyle = newStyle().Bold(true).Foreground(lipgloss.Color("226"))
	numStyle = newStyle().Bold(true).Foreground(lipgloss.Color("51"))
	entryStyle = newStyle().Foreground(lipgloss.Color("252"))
	hintStyle = newStyle().Foreground(lipgloss.Color("244"))
	borderStyle = newStyle().Foreground(lipgloss.Color("238"))
}

// TOCModal renders a centered table-of-contents popup.
type TOCModal struct {
	width   int
	height  int
	entries []blogmodule.TOCEntry
}

// New creates a new TOC modal.
func New(width, height int, entries []blogmodule.TOCEntry) *TOCModal {
	return &TOCModal{
		width:   width,
		height:  height,
		entries: entries,
	}
}

// View renders the modal.
func (t *TOCModal) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("TABLE OF CONTENTS"))
	sb.WriteString("\n")
	sb.WriteString(borderStyle.Render(strings.Repeat("─", 40)))
	sb.WriteString("\n\n")

	for i, entry := range t.entries {
		if i >= 9 {
			break // max 9 entries for single-digit jump
		}
		sb.WriteString(numStyle.Render(fmt.Sprintf("  %d", i+1)))
		sb.WriteString(entryStyle.Render("  " + entry.Title))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(hintStyle.Render("Press number to jump • Esc/t to close"))

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2, 3)

	modalBox := modalStyle.Render(sb.String())

	horizontalCentered := lipgloss.PlaceHorizontal(t.width, lipgloss.Center, modalBox)
	verticalCentered := lipgloss.PlaceVertical(t.height, lipgloss.Center, horizontalCentered)

	return verticalCentered
}
