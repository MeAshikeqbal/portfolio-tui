package helpmodal

import (
	"github.com/charmbracelet/lipgloss"
)

// HelpModal renders a centered help popup modal
type HelpModal struct {
	width  int
	height int
	help   string
}

// New creates a new help modal
func New(width, height int, helpText string) *HelpModal {
	return &HelpModal{
		width:  width,
		height: height,
		help:   helpText,
	}
}

// View renders the modal
func (h *HelpModal) View() string {
	// Create modal style with enhanced styling
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2, 3)

	// Create the modal box with the help content
	modalBox := modalStyle.Render(h.help)

	// Center the modal using lipgloss for symmetric alignment
	horizontalCentered := lipgloss.PlaceHorizontal(h.width, lipgloss.Center, modalBox)
	verticalCentered := lipgloss.PlaceVertical(h.height, lipgloss.Center, horizontalCentered)

	return verticalCentered
}

// SetSize updates the modal dimensions
func (h *HelpModal) SetSize(width, height int) {
	h.width = width
	h.height = height
}

