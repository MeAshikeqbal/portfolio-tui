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

var LoadingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Bold(true)

var BlogTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	PaddingLeft(1).
	PaddingRight(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderRight(false)

var BlogScrollInfo = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	PaddingLeft(1).
	PaddingRight(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderLeft(false)

var LineNumberStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

var SidebarTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("81"))

var SidebarMeta = lipgloss.NewStyle().
	Foreground(lipgloss.Color("248"))

var SidebarSection = lipgloss.NewStyle().
	Foreground(lipgloss.Color("244")).
	Bold(true)

var SidebarItem = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250")).
	PaddingLeft(1)

var SidebarActiveItem = lipgloss.NewStyle().
	Foreground(lipgloss.Color("230")).
	Background(lipgloss.Color("63")).
	Bold(true).
	PaddingLeft(1).
	PaddingRight(1)

var NeoTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("86"))

var NeoLabel = lipgloss.NewStyle().
	Foreground(lipgloss.Color("111")).
	Bold(true)

var NeoSeparator = lipgloss.NewStyle().
	Foreground(lipgloss.Color("244"))
