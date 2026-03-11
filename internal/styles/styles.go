package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title             lipgloss.Style
	MenuItem          lipgloss.Style
	SelectedItem      lipgloss.Style
	Footer            lipgloss.Style
	ViewportStyle     lipgloss.Style
	LoadingStyle      lipgloss.Style
	BlogTitleStyle    lipgloss.Style
	BlogScrollInfo    lipgloss.Style
	LineNumberStyle   lipgloss.Style
	SidebarTitle      lipgloss.Style
	SidebarMeta       lipgloss.Style
	SidebarSection    lipgloss.Style
	SidebarItem       lipgloss.Style
	SidebarActiveItem lipgloss.Style
	NeoTitle          lipgloss.Style
	NeoLabel          lipgloss.Style
	NeoSeparator      lipgloss.Style
)

func init() {
	SetRenderer(nil)
}

func SetRenderer(r *lipgloss.Renderer) {
	newStyle := lipgloss.NewStyle
	if r != nil {
		newStyle = r.NewStyle
	}

	Title = newStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	MenuItem = newStyle().
		PaddingLeft(2)

	SelectedItem = newStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("170")).
		Bold(true)

	Footer = newStyle().
		Foreground(lipgloss.Color("240"))

	ViewportStyle = newStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingLeft(2).
		PaddingRight(2)

	LoadingStyle = newStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	BlogTitleStyle = newStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderRight(false)

	BlogScrollInfo = newStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderLeft(false)

	LineNumberStyle = newStyle().
		Foreground(lipgloss.Color("240"))

	SidebarTitle = newStyle().
		Bold(true).
		Foreground(lipgloss.Color("81"))

	SidebarMeta = newStyle().
		Foreground(lipgloss.Color("248"))

	SidebarSection = newStyle().
		Foreground(lipgloss.Color("244")).
		Bold(true)

	SidebarItem = newStyle().
		Foreground(lipgloss.Color("250")).
		PaddingLeft(1)

	SidebarActiveItem = newStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("63")).
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1)

	NeoTitle = newStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	NeoLabel = newStyle().
		Foreground(lipgloss.Color("111")).
		Bold(true)

	NeoSeparator = newStyle().
		Foreground(lipgloss.Color("244"))
}
