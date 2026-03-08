package styles

import "github.com/charmbracelet/lipgloss"

var Title = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205"))

var MenuItem = lipgloss.NewStyle().
	PaddingLeft(2)

var SelectedItem = lipgloss.NewStyle().
	PaddingLeft(2).
	Foreground(lipgloss.Color("170")).
	Bold(true)

var Footer = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

var ViewportStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	PaddingLeft(2).
	PaddingRight(2)