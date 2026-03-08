package ui

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"time"

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
	name := os.Getenv("FULL_NAME")
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

func resolveMaskedUserIP() string {
	sshConn := os.Getenv("SSH_CONNECTION")
	if strings.TrimSpace(sshConn) == "" {
		return "127.0.xx.xx"
	}

	parts := strings.Fields(sshConn)
	if len(parts) == 0 {
		return "127.0.xx.xx"
	}

	ip := strings.TrimSpace(parts[0])
	octets := strings.Split(ip, ".")
	if len(octets) == 4 {
		return octets[0] + "." + octets[1] + ".xx.xx"
	}

	return "private-network"
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
