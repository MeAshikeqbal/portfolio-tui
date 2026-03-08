package ui

import (
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
)

// formatDate formats an ISO date string to a more readable format
func formatDate(isoDate string) string {
	// Simple formatting - just take the date part (YYYY-MM-DD)
	if len(isoDate) >= 10 {
		return isoDate[:10]
	}
	return isoDate
}

// renderBlogPostContent converts Sanity block content to readable text
func renderBlogPostContent(post *sanity.Post) string {
	var sb strings.Builder

	// Title
	sb.WriteString(post.Title)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", len(post.Title)))
	sb.WriteString("\n\n")

	// Metadata
	if post.PublishedAt != "" {
		sb.WriteString("📅 Published: ")
		sb.WriteString(formatDate(post.PublishedAt))
		sb.WriteString("\n")
	}

	if post.Author != nil {
		if author, ok := post.Author.(string); ok && author != "" {
			sb.WriteString("✍️  Author: ")
			sb.WriteString(author)
			sb.WriteString("\n")
		}
	}

	if len(post.Categories) > 0 {
		var cats []string
		for _, cat := range post.Categories {
			if catStr, ok := cat.(string); ok {
				cats = append(cats, catStr)
			}
		}
		if len(cats) > 0 {
			sb.WriteString("🏷️  Categories: ")
			sb.WriteString(strings.Join(cats, ", "))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 80))
	sb.WriteString("\n\n")

	// Body content
	if len(post.Body) > 0 {
		for _, block := range post.Body {
			if blockMap, ok := block.(map[string]interface{}); ok {
				// Check block type
				blockType, _ := blockMap["_type"].(string)

				if blockType == "block" {
					// Text block
					if children, ok := blockMap["children"].([]interface{}); ok {
						for _, child := range children {
							if childMap, ok := child.(map[string]interface{}); ok {
								if text, ok := childMap["text"].(string); ok {
									sb.WriteString(text)
								}
							}
						}
						sb.WriteString("\n\n")
					}
				} else if blockType == "image" {
					sb.WriteString("[Image]\n\n")
				} else if blockType == "code" {
					if code, ok := blockMap["code"].(string); ok {
						sb.WriteString("```\n")
						sb.WriteString(code)
						sb.WriteString("\n```\n\n")
					}
				}
			}
		}
	} else {
		sb.WriteString("No content available for this post.\n\n")
		sb.WriteString("This post may be a draft or the content has not been published yet.\n")
	}

	return sb.String()
}
