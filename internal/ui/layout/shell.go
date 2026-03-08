package layout

import "github.com/MeAshikeqbal/portfolio-tui/internal/ui/config"

// CalculateShellWidths calculates a 25/75 sidebar/content split for the usable width
func CalculateShellWidths(termWidth int) (int, int) {
	usableWidth := termWidth - 6 // borders + spacing in shell layout
	sidebar := usableWidth * config.SidebarRatioPct / 100
	if sidebar < config.MinSidebarWidth {
		sidebar = config.MinSidebarWidth
	}

	content := usableWidth - sidebar
	if content < config.MinContentWidth {
		content = config.MinContentWidth
		sidebar = usableWidth - content
	}

	return sidebar, content
}
