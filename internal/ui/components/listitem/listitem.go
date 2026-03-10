package listitem

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
)

// ProjectItem adapts a project to list.Item.
type ProjectItem struct {
	Data sanity.Project
}

func (i ProjectItem) FilterValue() string {
	// Return searchable content: title, description, technologies, github
	searchable := i.Data.Title + " " + i.Data.Description
	for _, tech := range i.Data.Technologies {
		searchable += " " + tech
	}
	searchable += " " + i.Data.GitHub + " " + i.Data.URL
	return searchable
}
func (i ProjectItem) Title() string       { return i.Data.Title }
func (i ProjectItem) Description() string {
	desc := i.Data.Description
	if i.Data.PublishedAt != "" {
		desc += " • " + formatDate(i.Data.PublishedAt)
	}
	if i.Data.GitHub != "" {
		desc += " • GitHub"
	}
	if i.Data.URL != "" {
		desc += " • Live"
	}
	if len(i.Data.Technologies) > 0 {
		desc += " • " + strings.Join(i.Data.Technologies, ", ")
	}
	return desc
}

// PostItem adapts a post to list.Item.
type PostItem struct {
	Data sanity.Post
}

func (i PostItem) FilterValue() string {
	// Return searchable content: title, author, categories
	searchable := i.Data.Title
	if i.Data.Author != nil {
		if author, ok := i.Data.Author.(string); ok {
			searchable += " " + author
		}
	}
	for _, cat := range i.Data.Categories {
		if catStr, ok := cat.(string); ok {
			searchable += " " + catStr
		}
	}
	return searchable
}
func (i PostItem) Title() string       { return i.Data.Title }
func (i PostItem) Description() string {
	desc := ""
	if i.Data.PublishedAt != "" {
		desc = "Published: " + formatDate(i.Data.PublishedAt)
	}
	if i.Data.Author != nil {
		if author, ok := i.Data.Author.(string); ok && author != "" {
			if desc != "" {
				desc += " • "
			}
			desc += "By " + author
		}
	}
	if len(i.Data.Categories) > 0 {
		var cats []string
		for _, cat := range i.Data.Categories {
			if catStr, ok := cat.(string); ok {
				cats = append(cats, catStr)
			}
		}
		if len(cats) > 0 {
			if desc != "" {
				desc += " • "
			}
			desc += strings.Join(cats, ", ")
		}
	}
	return desc
}

func formatDate(isoDate string) string {
	if len(isoDate) >= 10 {
		return isoDate[:10]
	}
	return isoDate
}
