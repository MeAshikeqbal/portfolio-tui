package utils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// WrapText wraps text to a specified width, preserving existing line breaks
func WrapText(text string, width int) string {
	if width < 10 {
		width = 10
	}

	lines := strings.Split(text, "\n")
	var wrapped []string

	for _, line := range lines {
		if len(line) <= width {
			wrapped = append(wrapped, line)
			continue
		}

		// Simple word wrapping
		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, line)
			continue
		}

		currentLine := ""
		for _, word := range words {
			if currentLine == "" {
				currentLine = word
			} else if len(currentLine)+1+len(word) <= width {
				currentLine += " " + word
			} else {
				wrapped = append(wrapped, currentLine)
				currentLine = word
			}
		}
		if currentLine != "" {
			wrapped = append(wrapped, currentLine)
		}
	}

	return strings.Join(wrapped, "\n")
}

// StripANSI removes ANSI escape codes from a string
func StripANSI(str string) string {
	var result strings.Builder
	inEscape := false

	for i := 0; i < len(str); i++ {
		if str[i] == '\x1b' && i+1 < len(str) && str[i+1] == '[' {
			inEscape = true
			i++ // skip the '['
			continue
		}

		if inEscape {
			if (str[i] >= 'A' && str[i] <= 'Z') || (str[i] >= 'a' && str[i] <= 'z') {
				inEscape = false
			}
			continue
		}

		result.WriteByte(str[i])
	}

	return result.String()
}

// DimContent applies a gray/dimmed style to content
func DimContent(content string) string {
	lines := strings.Split(content, "\n")
	dimmedLines := make([]string, len(lines))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	for i, line := range lines {
		dimmedLines[i] = dimStyle.Render(StripANSI(line))
	}

	return strings.Join(dimmedLines, "\n")
}
