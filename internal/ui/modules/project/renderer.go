package project

import (
	"strconv"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"
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

	if p.Slug.Current != "" {
		sb.WriteString("Slug: ")
		sb.WriteString(p.Slug.Current)
		sb.WriteString("\n")
	}
	if p.PublishedAt != "" {
		sb.WriteString("Published: ")
		sb.WriteString(blogmodule.FormatDate(p.PublishedAt))
		sb.WriteString("\n")
	}
	if p.MainImage.Alt != "" {
		sb.WriteString("Image Alt: ")
		sb.WriteString(p.MainImage.Alt)
		sb.WriteString("\n")
	}
	if p.GitHub != "" {
		sb.WriteString("GitHub: ")
		sb.WriteString(p.GitHub)
		sb.WriteString("\n")
	}
	if p.URL != "" {
		sb.WriteString("Live URL: ")
		sb.WriteString(p.URL)
		sb.WriteString("\n")
	}
	if p.GitHubData != nil {
		sb.WriteString("GitHub Data\n")
		sb.WriteString(strings.Repeat("-", len("GitHub Data")))
		sb.WriteString("\n")
		sb.WriteString("Stars: ")
		sb.WriteString(strconv.Itoa(p.GitHubData.Stars))
		sb.WriteString("\n")
		sb.WriteString("Commits: ")
		sb.WriteString(strconv.Itoa(p.GitHubData.Commits))
		sb.WriteString("\n")
		if p.GitHubData.License != "" {
			sb.WriteString("License: ")
			sb.WriteString(p.GitHubData.License)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

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

	return sb.String()
}
