package ui

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	tea "github.com/charmbracelet/bubbletea"
)

// Default fallback content when Sanity is unavailable
var (
	defaultProjects = `1. Portfolio TUI
   A beautiful terminal-based portfolio application built with Go and Bubble Tea.
   
2. Web Dashboard
   Modern web dashboard with real-time analytics.
   
3. API Gateway
   High-performance API gateway with authentication.
   
4. Mobile App
   Cross-platform mobile application.
   
5. CLI Tool
   Command-line tool for development workflows.`

	defaultSkills = `Programming Languages:
  • Go
  • JavaScript/TypeScript
  • Python
  
Frameworks & Libraries:
  • Bubble Tea
  • React
  • Node.js
  
Tools & Technologies:
  • Docker & Kubernetes
  • PostgreSQL & MongoDB
  • Redis`

	defaultAbout = `I'm a passionate software developer with a love for building
elegant and efficient solutions.

Background:
I specialize in backend development, terminal applications,
and developer tools.

Philosophy:
Good software should be both powerful and elegant.`

	defaultContact = `Email:
  ashik@example.com

GitHub:
  github.com/MeAshikeqbal

LinkedIn:
  linkedin.com/in/ashikeqbal

Twitter:
  @ashikeqbal`

	defaultBlog = `1. Welcome to My Blog
   Learn about my journey in software development
   Published: 2024-01-15
   
2. Building Terminal Applications with Go
   A deep dive into Bubble Tea framework
   Published: 2024-01-10
   
3. Why I Love Open Source
   Contributing to the community
   Published: 2024-01-05`
)

// fetchContentCmd initiates fetching content from Sanity
func fetchContentCmd() tea.Cmd {
	return func() tea.Msg {
		client := sanity.NewClient()
		content := make(map[string]string)
		var projects []sanity.Project
		var posts []sanity.Post

		// Fetch Projects
		if fetchedProjects, err := client.GetProjects(); err == nil {
			projects = fetchedProjects
			var sb strings.Builder
			sb.WriteString("🚀 Projects\n\n")
			for i, p := range fetchedProjects {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, p.Title))
				sb.WriteString(fmt.Sprintf("   %s\n", p.Description))
				if len(p.Technologies) > 0 {
					sb.WriteString(fmt.Sprintf("   Technologies: %s\n", strings.Join(p.Technologies, ", ")))
				}
				sb.WriteString("\n")
			}
			content["Projects"] = sb.String()
		} else {
			content["Projects"] = "📡 Fetching projects...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultProjects
		}

		// Fetch Skills
		if skills, err := client.GetSkills(); err == nil {
			var sb strings.Builder
			sb.WriteString("💻 Technical Skills\n\n")
			for _, s := range skills {
				sb.WriteString(fmt.Sprintf("%s:\n", s.Category))
				for _, item := range s.Items {
					sb.WriteString(fmt.Sprintf("  • %s\n", item))
				}
				sb.WriteString("\n")
			}
			content["Skills"] = sb.String()
		} else {
			content["Skills"] = "📡 Fetching skills...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultSkills
		}

		// Fetch About
		if about, err := client.GetAbout(); err == nil {
			var sb strings.Builder
			sb.WriteString("👨‍💻 About Me\n\n")
			sb.WriteString(about.Content)
			sb.WriteString("\n\n")
			sb.WriteString(about.Background)
			content["About"] = sb.String()
		} else {
			content["About"] = "📡 Fetching about...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultAbout
		}

		// Fetch Contact
		if contacts, err := client.GetContacts(); err == nil {
			var sb strings.Builder
			sb.WriteString("📫 Contact Information\n\n")
			for _, c := range contacts {
				sb.WriteString(fmt.Sprintf("%s:\n  %s\n\n", c.Platform, c.Value))
			}
			content["Contact"] = sb.String()
		} else {
			content["Contact"] = "📡 Fetching contact...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultContact
		}

		// Fetch Blog Posts
		if fetchedPosts, err := client.GetPosts(); err == nil {
			posts = fetchedPosts
			var sb strings.Builder
			sb.WriteString("📝 Blog Posts\n\n")
			if len(fetchedPosts) == 0 {
				sb.WriteString("No posts published yet.\n")
			} else {
				for i, post := range fetchedPosts {
					sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, post.Title))

					// Format published date
					if post.PublishedAt != "" {
						sb.WriteString(fmt.Sprintf("   Published: %s\n", formatDate(post.PublishedAt)))
					}

					// Show author if available (from dereferenced query)
					if post.Author != nil {
						if author, ok := post.Author.(string); ok && author != "" {
							sb.WriteString(fmt.Sprintf("   Author: %s\n", author))
						}
					}

					// Show categories if available (from dereferenced query)
					if len(post.Categories) > 0 {
						var cats []string
						for _, cat := range post.Categories {
							if catStr, ok := cat.(string); ok {
								cats = append(cats, catStr)
							}
						}
						if len(cats) > 0 {
							sb.WriteString(fmt.Sprintf("   Categories: %s\n", strings.Join(cats, ", ")))
						}
					}

					sb.WriteString("\n")
				}
			}
			content["Blog"] = sb.String()
		} else {
			content["Blog"] = "📡 Fetching blog posts...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultBlog
		}

		return ContentMsg{data: content, projects: projects, posts: posts, err: nil}
	}
}
