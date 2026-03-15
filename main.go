package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/MeAshikeqbal/portfolio-tui/internal/config"
	"github.com/MeAshikeqbal/portfolio-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/keygen"
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
	logFile, err := configureLogging("logs/app.log")
	if err != nil {
		fmt.Println("Warning: Failed to initialize logging:", err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	p := tea.NewProgram(
		ui.InitialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Println("Error:", err)
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	ui.SetRenderer(wishbubbletea.MakeRenderer(s))

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
	// Use env/config for SSH settings
	sshCfg := config.GetSSHConfig(cfg)

	addr := sshCfg.Host + ":" + sshCfg.Port

	logFile, err := configureLogging("logs/ssh-server.log")
	if err != nil {
		log.Printf("Could not open log file: %v", err)
	}
	if logFile != nil {
		defer logFile.Close()
	}

	if err := ensureSSHHostKey(sshCfg.HostKeyPath); err != nil {
		log.Printf("Could not prepare SSH host key: %v", err)
	}

	customLogger := log.New(log.Writer(), "", log.LstdFlags)

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(sshCfg.HostKeyPath),
		wish.WithMiddleware(
			wishbubbletea.Middleware(teaHandler),
			logging.MiddlewareWithLogger(customLogger),
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

func configureLogging(path string) (*os.File, error) {
	return configureLoggingWithConsole(path, os.Stderr)
}

func configureLoggingWithConsole(path string, console io.Writer) (*os.File, error) {
	if err := ensureParentDir(path); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	if console == nil {
		console = io.Discard
	}

	log.SetOutput(io.MultiWriter(console, logFile))
	return logFile, nil
}

func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}

func ensureSSHHostKey(path string) error {
	if err := ensureParentDir(path); err != nil {
		return err
	}

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if _, err := keygen.New(path, keygen.WithWrite()); err != nil {
		return err
	}

	log.Printf("Generated SSH host key at %s", path)
	return nil
}
