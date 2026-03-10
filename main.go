package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishbubbletea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		fmt.Println("Warning: Failed to load config, using defaults:", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "serve" {
		runSSHServer(cfg)
	} else {
		runLocal()
	}
}

func runLocal() {
	p := tea.NewProgram(
		ui.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Extract client IP from the SSH connection
	clientIP := ""
	if addr := s.RemoteAddr(); addr != nil {
		clientIP = addr.String()
		// Strip port from addr (e.g. "192.168.1.5:54321" -> "192.168.1.5")
		if host, _, err := net.SplitHostPort(clientIP); err == nil {
			clientIP = host
		}
	}

	// Extract terminal type from the SSH PTY
	termType := "ssh"
	if pty, _, ok := s.Pty(); ok {
		if pty.Term != "" {
			termType = pty.Term
		}
	}

	m := ui.InitialModelWithSession(ui.SessionInfo{
		UserIP:   clientIP,
		Terminal: termType,
	})

	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

func runSSHServer(cfg *config.Config) {
	host := "0.0.0.0"
	port := "23234"

	if cfg != nil && cfg.SSH.Host != "" {
		host = cfg.SSH.Host
	}
	if cfg != nil && cfg.SSH.Port != "" {
		port = cfg.SSH.Port
	}

	keyPath := ".ssh/term_key"
	if cfg != nil && cfg.SSH.HostKeyPath != "" {
		keyPath = cfg.SSH.HostKeyPath
	}

	addr := host + ":" + port

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			wishbubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln("Could not create SSH server:", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	log.Printf("Starting SSH portfolio server on %s", addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalln("SSH server error:", err)
		}
	}()

	<-done
	log.Println("Shutting down SSH server...")
	if err := s.Close(); err != nil {
		log.Fatalln("Could not close SSH server:", err)
	}
}
