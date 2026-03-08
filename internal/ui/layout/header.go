package layout

import (
	"github.com/MeAshikeqbal/portfolio-tui/internal/styles"
)

// RenderHeader renders the header breadcrumb
func RenderHeader(name, tagline, breadcrumb string) string {
	if breadcrumb == "" {
		return ""
	}
	title := styles.Title.Render(name + " – " + tagline + " > " + breadcrumb)
	return title + "\n"
}
