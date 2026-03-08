package ui

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
)

// ProjectItem implements list.Item for projects
type ProjectItem struct {
	project sanity.Project
}

func (i ProjectItem) FilterValue() string { return i.project.Title }
func (i ProjectItem) Title() string       { return i.project.Title }
func (i ProjectItem) Description() string {
	desc := i.project.Description
	if len(i.project.Technologies) > 0 {
		desc += " • " + strings.Join(i.project.Technologies, ", ")
	}
	return desc
}

// PostItem implements list.Item for blog posts
type PostItem struct {
	post sanity.Post
}

func (i PostItem) FilterValue() string { return i.post.Title }
func (i PostItem) Title() string       { return i.post.Title }
func (i PostItem) Description() string {
	desc := ""
	if i.post.PublishedAt != "" {
		desc = "Published: " + formatDate(i.post.PublishedAt)
	}

	// Add author if available
	if i.post.Author != nil {
		if author, ok := i.post.Author.(string); ok && author != "" {
			if desc != "" {
				desc += " • "
			}
			desc += "By " + author
		}
	}

	// Add categories if available
	if len(i.post.Categories) > 0 {
		var cats []string
		for _, cat := range i.post.Categories {
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
