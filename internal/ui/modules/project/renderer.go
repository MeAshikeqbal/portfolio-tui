package project

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
)

// RenderProjectContent builds the detailed project page content.
func RenderProjectContent(p *sanity.Project) string {
	var sb strings.Builder

	sb.WriteString(p.Title)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", len(p.Title)))
	sb.WriteString("\n\n")

	sb.WriteString("Description\n")
	sb.WriteString(strings.Repeat("-", len("Description")))
	sb.WriteString("\n")
	sb.WriteString(p.Description)
	sb.WriteString("\n\n")

	sb.WriteString("Technologies\n")
	sb.WriteString(strings.Repeat("-", len("Technologies")))
	sb.WriteString("\n")
	if len(p.Technologies) == 0 {
		sb.WriteString("- Not specified\n")
	} else {
		for _, tech := range p.Technologies {
			sb.WriteString("- ")
			sb.WriteString(tech)
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString("Project ID: ")
	sb.WriteString(p.ID)
	sb.WriteString("\n")

	return sb.String()
}
