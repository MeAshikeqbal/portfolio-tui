package content

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	blogmodule "github.com/MeAshikeqbal/portfolio-tui/internal/ui/modules/blog"
)

type Payload struct {
	Data     map[string]string
	Projects []sanity.Project
	Posts    []sanity.Post
	Err      error
}

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

	defaultHome = `Welcome.

This is my terminal portfolio, inspired by classic Neofetch and SSH dashboards.

Use the sidebar to explore Projects, Skills, Experience, Education, Blogs, and Contact details.`

	defaultExperience = `Experience
• Software Developer - Building backend services and developer tools
• Open-source contributor - CLI/TUI experiences and automation
• Focus areas - Go, API design, performance, DX`

	defaultEducation = `Education
• Bachelor-level Computer Science track
• Continuous learning in systems design and cloud platforms
• Practical learning through projects and open source`
)

func Fetch(client *sanity.Client) Payload {
	content := make(map[string]string)
	var projects []sanity.Project
	var posts []sanity.Post

	if fetchedProjects, err := client.GetProjects(); err == nil {
		projects = fetchedProjects
		var sb strings.Builder
		sb.WriteString("🚀 Projects\n\n")
		for i, p := range fetchedProjects {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, p.Title))
			sb.WriteString(fmt.Sprintf("   %s\n", p.Description))
			if p.PublishedAt != "" {
				sb.WriteString(fmt.Sprintf("   Published: %s\n", blogmodule.FormatDate(p.PublishedAt)))
			}
			if p.GitHub != "" {
				sb.WriteString(fmt.Sprintf("   GitHub: %s\n", p.GitHub))
			}
			if p.URL != "" {
				sb.WriteString(fmt.Sprintf("   Live URL: %s\n", p.URL))
			}
			if p.GitHubData != nil {
				sb.WriteString(fmt.Sprintf("   Stats: %d stars, %d commits\n", p.GitHubData.Stars, p.GitHubData.Commits))
				if p.GitHubData.License != "" {
					sb.WriteString(fmt.Sprintf("   License: %s\n", p.GitHubData.License))
				}
			}
			if len(p.Technologies) > 0 {
				sb.WriteString(fmt.Sprintf("   Technologies: %s\n", strings.Join(p.Technologies, ", ")))
			}
			sb.WriteString("\n")
		}
		content["Projects"] = sb.String()
	} else {
		content["Projects"] = "📡 Fetching projects...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultProjects
	}

	if skills, err := client.GetSkills(); err == nil && len(skills) > 0 {
		var sb strings.Builder
		for _, s := range skills {
			if s.Name != "" {
				sb.WriteString(fmt.Sprintf("  • %s\n", s.Name))
			}
		}
		sb.WriteString("\n")
		content["Skills"] = sb.String()
	} else {
		content["Skills"] = "📡 Fetching skills...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultSkills
	}

	// Sidebar aliases and static sections
	content["Home"] = defaultHome

	if experiences, err := client.GetExperiences(); err == nil && len(experiences) > 0 {
		var sb strings.Builder
		for _, e := range experiences {
			if e.Year != "" {
				sb.WriteString(fmt.Sprintf("[%s]\n", e.Year))
			}
			for _, w := range e.Works {
				if w.Name != "" {
					sb.WriteString(fmt.Sprintf("  %s\n", w.Name))
				}
				if w.Company != "" {
					sb.WriteString(fmt.Sprintf("  @ %s\n", w.Company))
				}
				if w.Desc != "" {
					sb.WriteString(fmt.Sprintf("  %s\n", w.Desc))
				}
				sb.WriteString("\n")
			}
		}
		content["Experience"] = sb.String()
	} else {
		content["Experience"] = defaultExperience
	}

	if educations, err := client.GetEducation(); err == nil && len(educations) > 0 {
		var sb strings.Builder
		for _, e := range educations {
			if e.Year != "" {
				sb.WriteString(fmt.Sprintf("[%s]\n", e.Year))
			}
			if e.Degree != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", e.Degree))
			}
			if e.School != "" {
				sb.WriteString(fmt.Sprintf("  @ %s\n", e.School))
			}
			if e.Desc != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", e.Desc))
			}
			sb.WriteString("\n")
		}
		content["Education"] = sb.String()
	} else {
		content["Education"] = defaultEducation
	}

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

	if content["Home"] == defaultHome && content["About"] != "" {
		content["Home"] = "🏠 Home\n\n" + content["About"]
	}

	if contacts, err := client.GetContacts(); err == nil {
		var sb strings.Builder
		for _, c := range contacts {
			sb.WriteString(fmt.Sprintf("%s:\n  %s\n\n", c.Platform, c.Value))
		}
		content["Contact"] = sb.String()
		content["Contact Me"] = sb.String()
	} else {
		content["Contact"] = "📡 Fetching contact...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultContact
		content["Contact Me"] = content["Contact"]
	}

	if fetchedPosts, err := client.GetPosts(); err == nil {
		posts = fetchedPosts
		var sb strings.Builder
		sb.WriteString("📝 Blog Posts\n\n")
		if len(fetchedPosts) == 0 {
			sb.WriteString("No posts published yet.\n")
		} else {
			for i, post := range fetchedPosts {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, post.Title))
				if post.PublishedAt != "" {
					sb.WriteString(fmt.Sprintf("   Published: %s\n", blogmodule.FormatDate(post.PublishedAt)))
				}
				if post.Author != nil {
					if author, ok := post.Author.(string); ok && author != "" {
						sb.WriteString(fmt.Sprintf("   Author: %s\n", author))
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
						sb.WriteString(fmt.Sprintf("   Categories: %s\n", strings.Join(cats, ", ")))
					}
				}
				sb.WriteString("\n")
			}
		}
		content["Blog"] = sb.String()
		content["Blogs"] = sb.String()
	} else {
		content["Blog"] = "📡 Fetching blog posts...\n\nError connecting to Sanity. Using fallback content.\n\n" + defaultBlog
		content["Blogs"] = content["Blog"]
	}

	return Payload{Data: content, Projects: projects, Posts: posts, Err: nil}
}
