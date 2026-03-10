package ui

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type introDoneMsg struct{}
type introStageMsg struct {
	stage int
}

func introDelayCmd() tea.Cmd {
	return tea.Tick(1500*time.Millisecond, func(time.Time) tea.Msg {
		return introDoneMsg{}
	})
}

func introStageCmd(stage int, delay time.Duration) tea.Cmd {
	return tea.Tick(delay, func(time.Time) tea.Msg {
		return introStageMsg{stage: stage}
	})
}

func maybeCompleteIntro(m *Model) {
	if m.introComplete {
		return
	}

	if m.introTimerDone && m.loadingState != loading {
		m.introComplete = true
	}
}

func resolvePortfolioOwner() string {
	name := config.Get().Owner.FullName
	if strings.TrimSpace(name) == "" {
		return "Ashik Eqbal"
	}
	return name
}

func resolveTerminalName() string {
	term := os.Getenv("TERM")
	if strings.TrimSpace(term) == "" {
		return "xterm-256color"
	}
	return term
}

func resolveUserIP() string {
	sshConn := os.Getenv("SSH_CONNECTION")
	if strings.TrimSpace(sshConn) == "" {
		return "127.0.0.1"
	}

	parts := strings.Fields(sshConn)
	if len(parts) == 0 {
		return "127.0.0.1"
	}

	return strings.TrimSpace(parts[0])
}

// resolveSessionField returns the appropriate session field, preferring
// SSH-provided values over local env fallbacks.
func resolveSessionField(session *SessionInfo, field string) string {
	switch field {
	case "ip":
		if session != nil && strings.TrimSpace(session.UserIP) != "" {
			return session.UserIP
		}
		return resolveUserIP()
	case "terminal":
		if session != nil && strings.TrimSpace(session.Terminal) != "" {
			return session.Terminal
		}
		return resolveTerminalName()
	default:
		return ""
	}
}

func generateSessionID() string {
	buf := make([]byte, 3)
	if _, err := rand.Read(buf); err != nil {
		return "00000"
	}

	id := hex.EncodeToString(buf)
	if len(id) > 5 {
		return id[:5]
	}
	return id
}
