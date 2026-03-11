package blog

import (
	"fmt"
	"strings"

	"github.com/MeAshikeqbal/portfolio-tui/internal/sanity"
	"github.com/charmbracelet/lipgloss"
)

// styles used throughout the renderer
var (
	titleStyle      lipgloss.Style
	metaLabelStyle  lipgloss.Style
	metaValueStyle  lipgloss.Style
	metaStyle       lipgloss.Style
	tagSepStyle     lipgloss.Style
	dividerStyle    lipgloss.Style
	h1Style         lipgloss.Style
	h2Style         lipgloss.Style
	h2RulerStyle    lipgloss.Style
	h3Style         lipgloss.Style
	h4Style         lipgloss.Style
	bodyStyle       lipgloss.Style
	quoteBarStyle   lipgloss.Style
	quoteTextStyle  lipgloss.Style
	codeBlockStyle  lipgloss.Style
	codeBorderStyle lipgloss.Style
	inlineCodeStyle lipgloss.Style
	bulletStyle     lipgloss.Style
)

func init() {
	SetRenderer(nil)
}

func SetRenderer(r *lipgloss.Renderer) {
	newStyle := lipgloss.NewStyle
	if r != nil {
		newStyle = r.NewStyle
	}

	titleStyle = newStyle().Bold(true).Foreground(lipgloss.Color("205"))
	metaLabelStyle = newStyle().Bold(true).Foreground(lipgloss.Color("81"))
	metaValueStyle = newStyle().Foreground(lipgloss.Color("252"))
	metaStyle = newStyle().Foreground(lipgloss.Color("244"))
	tagSepStyle = newStyle().Foreground(lipgloss.Color("62")).Bold(true)
	dividerStyle = newStyle().Foreground(lipgloss.Color("238"))
	h1Style = newStyle().Bold(true).Foreground(lipgloss.Color("205"))
	h2Style = newStyle().Bold(true).Foreground(lipgloss.Color("81"))
	h2RulerStyle = newStyle().Foreground(lipgloss.Color("238"))
	h3Style = newStyle().Bold(true).Foreground(lipgloss.Color("220"))
	h4Style = newStyle().Bold(true).Foreground(lipgloss.Color("248"))
	bodyStyle = newStyle().Foreground(lipgloss.Color("252"))
	quoteBarStyle = newStyle().Foreground(lipgloss.Color("62"))
	quoteTextStyle = newStyle().Foreground(lipgloss.Color("244")).Italic(true)
	codeBlockStyle = newStyle().Foreground(lipgloss.Color("108"))
	codeBorderStyle = newStyle().Foreground(lipgloss.Color("238"))
	inlineCodeStyle = newStyle().Foreground(lipgloss.Color("108"))
	bulletStyle = newStyle().Foreground(lipgloss.Color("62")).Bold(true)
}

// TOCEntry holds the title and line offset of an h2 heading.
type TOCEntry struct {
	Title string
	Line  int
}

func FormatDate(isoDate string) string {
	if len(isoDate) >= 10 {
		return isoDate[:10]
	}
	return isoDate
}

// RenderPostContent renders a full blog post and returns the rendered string
// along with a TOC built from h2 headings (title + line number).
func RenderPostContent(post *sanity.Post, width int) (string, []TOCEntry) {
	if width < 20 {
		width = 80
	}
	inner := width
	if inner < 10 {
		inner = 10
	}

	var sb strings.Builder
	var toc []TOCEntry
	lineCount := 0

	// write helper that also tracks line count
	write := func(s string) {
		sb.WriteString(s)
		lineCount += strings.Count(s, "\n")
	}

	// ── Title ───────────────────────────────────────────────────────────────
	write(titleStyle.Render(post.Title))
	write("\n")
	write(dividerStyle.Render(strings.Repeat("─", inner)))
	write("\n")

	// ── Metadata ─────────────────────────────────────────────────────────────
	if post.PublishedAt != "" {
		write(metaLabelStyle.Render("Published") + ": " + metaValueStyle.Render(FormatDate(post.PublishedAt)) + "\n")
	}
	if post.Author != nil {
		if author, ok := post.Author.(string); ok && author != "" {
			write(metaLabelStyle.Render("Author") + ": " + metaValueStyle.Render(author) + "\n")
		}
	}
	if len(post.Categories) > 0 {
		var cats []string
		for _, cat := range post.Categories {
			if catStr, ok := cat.(string); ok && catStr != "" {
				cats = append(cats, catStr)
			}
		}
		if len(cats) > 0 {
			write(metaLabelStyle.Render("Tags") + ": " + metaValueStyle.Render(strings.Join(cats, tagSepStyle.Render(" • "))) + "\n")
		}
	}

	write("\n")
	write(dividerStyle.Render(strings.Repeat("─", inner)))
	write("\n\n")

	// ── Body ────────────────────────────────────────────────────────────────
	if len(post.Body) == 0 {
		write(metaStyle.Render("No content available for this post."))
		write("\n")
		return sb.String(), toc
	}

	listCounters := map[int]int{}
	prevWasList := false

	for _, rawBlock := range post.Body {
		blockMap, ok := rawBlock.(map[string]interface{})
		if !ok {
			continue
		}

		blockType, _ := blockMap["_type"].(string)

		switch blockType {
		case "block":
			style, _ := blockMap["style"].(string)
			listItem, _ := blockMap["listItem"].(string)
			rawLevel, _ := blockMap["level"].(float64)
			level := int(rawLevel)
			if level < 1 {
				level = 1
			}

			text := renderChildren(blockMap)

			if listItem != "" {
				indent := strings.Repeat("  ", level)
				switch listItem {
				case "bullet":
					prefix := bulletStyle.Render("•")
					wrapped := wrapSubsequent(text, inner-lipgloss.Width(indent)-2, indent+"  ")
					write(indent + prefix + " " + bodyStyle.Render(wrapped) + "\n")
					listCounters[level] = 0
				case "number":
					listCounters[level]++
					num := bulletStyle.Render(fmt.Sprintf("%d.", listCounters[level]))
					wrapped := wrapSubsequent(text, inner-lipgloss.Width(indent)-3, indent+"   ")
					write(indent + num + " " + bodyStyle.Render(wrapped) + "\n")
				}
				prevWasList = true
				continue
			}

			if prevWasList {
				write("\n")
			}
			prevWasList = false
			listCounters = map[int]int{}

			switch style {
			case "h1":
				write("\n")
				write(h1Style.Render(text) + "\n")
				write(h1Style.Render(strings.Repeat("═", min(lipgloss.Width(text), inner))) + "\n\n")
			case "h2":
				// Record TOC entry at the line where the ruler will appear
				plainText := renderChildrenPlain(blockMap)
				toc = append(toc, TOCEntry{Title: plainText, Line: lineCount + 1}) // +1 for the leading \n
				prefix := h2RulerStyle.Render("── ")
				styledText := h2Style.Render(text)
				suffix := " "
				used := lipgloss.Width(prefix) + lipgloss.Width(text) + lipgloss.Width(suffix)
				remainder := inner - used
				if remainder < 0 {
					remainder = 0
				}
				ruler := prefix + styledText + suffix + h2RulerStyle.Render(strings.Repeat("─", remainder))
				write("\n")
				write(ruler + "\n\n")
			case "h3":
				write("\n")
				write(h3Style.Render("▶  "+text) + "\n\n")
			case "h4":
				write(h4Style.Render("◆  "+text) + "\n\n")
			case "blockquote":
				for _, l := range wrapToLines(text, inner-4) {
					write(quoteBarStyle.Render("┃ ") + quoteTextStyle.Render(l) + "\n")
				}
				write("\n")
			default: // "normal"
				if strings.TrimSpace(text) == "" {
					write("\n")
				} else {
					lines := wrapToLines(text, inner)
					for i, l := range lines {
						// Justify all lines except the last one in a paragraph
						if i < len(lines)-1 {
							l = justifyLine(l, inner)
						}
						write(bodyStyle.Render(l) + "\n")
					}
					write("\n")
				}
			}

		case "code":
			code, _ := blockMap["code"].(string)
			lang, _ := blockMap["language"].(string)
			lineLen := inner - 2
			if lineLen < 4 {
				lineLen = 4
			}

			var topBar string
			if lang != "" {
				label := " " + strings.ToUpper(lang) + " "
				rest := lineLen - lipgloss.Width(label) - 2
				if rest < 0 {
					rest = 0
				}
				topBar = "╭─" + label + strings.Repeat("─", rest) + "─╮"
			} else {
				topBar = "╭" + strings.Repeat("─", lineLen) + "╮"
			}
			bottomBar := "╰" + strings.Repeat("─", lineLen) + "╯"

			write("\n")
			write(codeBorderStyle.Render(topBar) + "\n")
			for _, line := range strings.Split(code, "\n") {
				write(codeBorderStyle.Render("│") + " " + codeBlockStyle.Render(line) + "\n")
			}
			write(codeBorderStyle.Render(bottomBar) + "\n\n")

		case "image":
			write(metaStyle.Render("[ Image ]") + "\n\n")
		}
	}

	return sb.String(), toc
}

// renderChildrenPlain extracts plain text from a block's children (no styling).
func renderChildrenPlain(blockMap map[string]interface{}) string {
	children, ok := blockMap["children"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, child := range children {
		childMap, ok := child.(map[string]interface{})
		if !ok {
			continue
		}
		text, _ := childMap["text"].(string)
		sb.WriteString(text)
	}
	return sb.String()
}

// renderChildren extracts inline text from a Portable Text block's children,
// applying inline decorators (bold, italic, inline code).
func renderChildren(blockMap map[string]interface{}) string {
	children, ok := blockMap["children"].([]interface{})
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, child := range children {
		childMap, ok := child.(map[string]interface{})
		if !ok {
			continue
		}
		text, _ := childMap["text"].(string)
		marks, _ := childMap["marks"].([]interface{})

		var bold, em, code bool
		for _, mark := range marks {
			if m, ok := mark.(string); ok {
				switch m {
				case "strong":
					bold = true
				case "em":
					em = true
				case "code":
					code = true
				}
			}
		}

		switch {
		case code:
			sb.WriteString(inlineCodeStyle.Render("`" + text + "`"))
		case bold && em:
			sb.WriteString(lipgloss.NewStyle().Bold(true).Italic(true).Render(text))
		case bold:
			sb.WriteString(lipgloss.NewStyle().Bold(true).Render(text))
		case em:
			sb.WriteString(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("252")).Render(text))
		default:
			sb.WriteString(text)
		}
	}
	return sb.String()
}

// justifyLine pads spaces between words so the line fills the target width.
func justifyLine(line string, width int) string {
	words := strings.Fields(line)
	if len(words) <= 1 {
		return line
	}
	totalWordLen := 0
	for _, w := range words {
		totalWordLen += lipgloss.Width(w)
	}
	gaps := len(words) - 1
	totalSpaces := width - totalWordLen
	if totalSpaces <= gaps {
		return line
	}
	baseGap := totalSpaces / gaps
	extra := totalSpaces % gaps

	var sb strings.Builder
	for i, w := range words {
		sb.WriteString(w)
		if i < gaps {
			pad := baseGap
			if i < extra {
				pad++
			}
			sb.WriteString(strings.Repeat(" ", pad))
		}
	}
	return sb.String()
}

// wrapToLines wraps text respecting ANSI escape codes in width measurement.
func wrapToLines(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	current := ""
	currentW := 0
	for _, word := range words {
		ww := lipgloss.Width(word)
		if currentW == 0 {
			current = word
			currentW = ww
		} else if currentW+1+ww <= width {
			current += " " + word
			currentW += 1 + ww
		} else {
			lines = append(lines, current)
			current = word
			currentW = ww
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// wrapSubsequent wraps text so that continuation lines are indented by `contIndent`.
func wrapSubsequent(text string, firstWidth int, contIndent string) string {
	lines := wrapToLines(text, firstWidth)
	return strings.Join(lines, "\n"+contIndent)
}
