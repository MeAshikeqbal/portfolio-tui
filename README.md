# Portfolio TUI

A beautiful terminal-based portfolio application built with Go and Bubble Tea, with dynamic content from Sanity CMS.

## Features

- 📱 **Interactive TUI** - Navigate with keyboard (vim-style shortcuts supported)
- 🎨 **Styled with Lip Gloss** - Beautiful color schemes and layouts
- 📜 **Scrollable Content** - Viewport for long content with smooth scrolling
- 📋 **Interactive Lists** - Browse projects and blog posts with List component from Bubbles
- 📖 **Blog Reader** - Individual blog post pages with smooth scrolling and progress indicator
- 🖱️ **Mouse Support** - Full mouse wheel scrolling support across all content views
- ⏳ **Loading Spinner** - Animated spinner while fetching from Sanity
- 🔍 **Filtering** - Built-in search/filter functionality for projects and blog posts
- ⌨️ **Vim Keybindings** - Full vim-style navigation (hjkl, gg, G, etc.)
- ❓ **Built-in Help** - Press `?` for keybindings
- 🌐 **Dynamic Content** - Fetches data from Sanity CMS
- 📝 **Blog Integration** - Displays your blog posts from Sanity

## Installation

```bash
# Clone the repository
git clone https://github.com/MeAshikeqbal/portfolio-tui
cd portfolio-tui

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your Sanity credentials

# Run the application
go run .

# Or build and run
go build -o portfolio-tui
./portfolio-tui
```

## Environment Variables

Create a `.env` file with your Sanity project details and personal information:

```env
# Sanity CMS Configuration
SANITY_PROJECT_ID=your_project_id
SANITY_DATASET=production
SANITY_API_VERSION=2024-12-21

# Personal Information
FULL_NAME=Your Name
TAGLINE=Your Professional Title

# Contact Information
CONTACT_EMAIL=your.email@example.com
CONTACT_PHONE=+1234567890

# Social Links (leave blank to hide)
SOCIAL_GITHUB=https://github.com/yourusername
SOCIAL_LINKEDIN=https://linkedin.com/in/yourprofile
SOCIAL_TWITTER=https://twitter.com/yourhandle
SOCIAL_WEBSITE=https://yourwebsite.com
SOCIAL_YOUTUBE=
SOCIAL_INSTAGRAM=
```

**Note:** Contact information is managed via `.env` file - no need to create separate Sanity documents for social links!

## Sanity Schema Setup

Your Sanity studio should have the following document types:

### 1. Project
```javascript
{
  _type: "project",
  title: string,
  description: string,
  technologies: array of strings,
  order: number
}
```

### 2. Skill
```javascript
{
  _type: "skill",
  category: string,
  items: array of strings
}
```

### 3. About
```javascript
{
  _type: "about",
  content: text,
  background: text
}
```

### 4. Contact
```javascript
{
  _type: "contact",
  platform: string,
  value: string
}
```

### 5. Post (Blog)
```javascript
{
  _type: "post",
  title: string,
  slug: slug,
  author: reference to author,
  publishedAt: datetime,
  categories: array of references to category,
  mainImage: image,
  body: blockContent
}
```

### 6. Author
```javascript
{
  _type: "author",
  name: string,
  slug: slug,
  image: image,
  bio: array of blocks
}
```

### 7. Category
```javascript
{
  _type: "category",
  title: string,
  slug: slug,
  description: text
}
```

## Usage

### Navigation

**Menu View:**
- `↑`/`k` - Move up
- `↓`/`j` - Move down
- `Enter` - Select item
- `?` - Toggle help
- `q` - Quit

**Content View:**
- `↑`/`k` - Scroll up / Navigate list up
- `↓`/`j` - Scroll down / Navigate list down
- `PgUp`/`b` - Page up
- `PgDn`/`f` - Page down
- `/` - Filter list (Projects and Blog only)
- `Enter` - Open blog post (in Blog list)
- `Esc` - Back to menu / Clear filter
- `?` - Toggle help
- `q` - Quit

**Blog Detail View:**
- `↑`/`k` - Scroll up line by line
- `↓`/`j` - Scroll down line by line
- `PgUp`/`b` - Page up
- `PgDn`/`f` - Page down
- `Home`/`g` - Jump to top of post
- `End`/`G` - Jump to bottom of post
- `Esc` - Back to blog list
- `?` - Toggle help
- `q` - Quit
- Mouse wheel scrolling supported

### Menu Items

1. **Projects** - Browse your projects in an interactive list with filtering
2. **Skills** - Display your technical skills
3. **Blog** - Browse your blog posts in an interactive list with filtering
4. **About** - Your bio and background
5. **Contact** - Contact information
6. **Exit** - Quit the application

## Development

```bash
# Run with live reload (requires air)
air

# Format code
go fmt ./...

# Run tests
go test ./...

# Build for production
go build -o portfolio-tui -ldflags="-s -w" .
```

## Project Structure

```
portfolio-tui/
├── internal/
│   ├── sanity/         # Sanity CMS client
│   │   └── client.go
│   ├── styles/         # Lip Gloss styles
│   │   └── styles.go
│   └── ui/             # Bubble Tea UI components
│       ├── components/
│       │   ├── keymap/
│       │   │   └── keymap.go
│       │   └── listitem/
│       │       └── listitem.go
│       ├── modules/
│       │   ├── blog/
│       │   │   └── renderer.go
│       │   └── content/
│       │       └── fetcher.go
│       ├── model.go
│       ├── update.go
│       ├── content.go
│       └── views.go
├── main.go             # Entry point
├── .env                # Environment variables (gitignored)
├── .env.example        # Environment template
├── .gitignore
├── go.mod
└── README.md
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [godotenv](https://github.com/joho/godotenv) - Environment variables

## License

MIT

## Author

Ashik Eqbal
- GitHub: [@MeAshikeqbal](https://github.com/MeAshikeqbal)
