package menu

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
)

// MenuItem represents a menu item configuration
type MenuItem struct {
	Label string
	Icon  string
}

// MenuSection represents a section in the menu
type MenuSection struct {
	Title string
	Items []MenuItem
}

// RenderMenu renders the main menu with sections
func RenderMenu(menu []string, selected int) string {
	var s strings.Builder
	sectionBreak := map[string]string{
		"Home":   "CORE",
		"Skills": "PROFILE",
		"Exit":   "SYSTEM",
	}

	for i, item := range menu {
		if section, ok := sectionBreak[item]; ok {
			if i > 0 {
				s.WriteString("\n")
			}
			s.WriteString(styles.SidebarSection.Render("[" + section + "]"))
			s.WriteString("\n")
		}

		entry := fmt.Sprintf("%s %s", GetMenuIcon(item), item)
		if i == selected {
			s.WriteString(styles.SidebarActiveItem.Render("> " + entry))
		} else {
			s.WriteString(styles.SidebarItem.Render("  " + entry))
		}
		s.WriteString("\n")
	}

	return s.String()
}

// GetMenuIcon returns the icon for a menu item
func GetMenuIcon(item string) string {
	switch item {
	case "Home":
		return "H"
	case "Projects":
		return "P"
	case "Skills":
		return "S"
	case "Experience":
		return "E"
	case "Education":
		return "D"
	case "Blogs":
		return "B"
	case "Contact Me":
		return "C"
	case "Exit":
		return "X"
	default:
		return ">"
	}
}
