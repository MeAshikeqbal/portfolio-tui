package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Runtime string `yaml:"runtime"`
	} `yaml:"app"`

	Owner struct {
		FullName string `yaml:"full_name"`
		Username string `yaml:"username"`
		Tagline  string `yaml:"tagline"`
		Role     string `yaml:"role"`
		Bio      string `yaml:"bio"`
	} `yaml:"owner"`

	Network struct {
		Host string `yaml:"host"`
	} `yaml:"network"`

	Contact struct {
		Email string `yaml:"email"`
		Phone string `yaml:"phone"`
	} `yaml:"contact"`

	Social struct {
		GitHub    string `yaml:"github"`
		LinkedIn  string `yaml:"linkedin"`
		Twitter   string `yaml:"twitter"`
		Website   string `yaml:"website"`
		YouTube   string `yaml:"youtube"`
		Instagram string `yaml:"instagram"`
	} `yaml:"social"`

	Branding struct {
		SidebarTitle string `yaml:"sidebar_title"`
		AsciiLogo    string `yaml:"ascii_logo"`
		SidebarLogo  string `yaml:"sidebar_logo"`
	} `yaml:"branding"`

	SSH struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		HostKeyPath string `yaml:"host_key_path"`
	} `yaml:"ssh"`
}

var defaultConfig = Config{
	App: struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Runtime string `yaml:"runtime"`
	}{
		Name:    "portfolio-tui",
		Version: "v1.0",
		Runtime: "Bubble Tea",
	},
	Owner: struct {
		FullName string `yaml:"full_name"`
		Username string `yaml:"username"`
		Tagline  string `yaml:"tagline"`
		Role     string `yaml:"role"`
		Bio      string `yaml:"bio"`
	}{
		FullName: "Ashik Eqbal",
		Username: "ashikeqbal",
		Tagline:  "Portfolio",
		Role:     "Developer • DevOps • Homelab Builder",
		Bio:      "I build backend systems, homelab infrastructure,\nand developer tools using Go, TypeScript and Linux.",
	},
	Network: struct {
		Host string `yaml:"host"`
	}{
		Host: "localhost",
	},
	Contact: struct {
		Email string `yaml:"email"`
		Phone string `yaml:"phone"`
	}{
		Email: "ashik@example.com",
		Phone: "",
	},
	Social: struct {
		GitHub    string `yaml:"github"`
		LinkedIn  string `yaml:"linkedin"`
		Twitter   string `yaml:"twitter"`
		Website   string `yaml:"website"`
		YouTube   string `yaml:"youtube"`
		Instagram string `yaml:"instagram"`
	}{
		GitHub:    "https://github.com/MeAshikeqbal",
		LinkedIn:  "",
		Twitter:   "",
		Website:   "",
		YouTube:   "",
		Instagram: "",
	},
	Branding: struct {
		SidebarTitle string `yaml:"sidebar_title"`
		AsciiLogo    string `yaml:"ascii_logo"`
		SidebarLogo  string `yaml:"sidebar_logo"`
	}{
		SidebarTitle: "PORTFOLIO",
		AsciiLogo: strings.Join([]string{
			"        /\\",
			"       /  \\",
			"      / /\\ \\",
			"     / ____ \\",
			"    /_/    \\_\\",
		}, "\n"),
		SidebarLogo: strings.Join([]string{
			"      /\\",
			"     /  \\",
			"    / /\\ \\",
			"   / ____ \\",
			"  /_/    \\_\\",
		}, "\n"),
	},
	SSH: struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		HostKeyPath string `yaml:"host_key_path"`
	}{
		Host:        "0.0.0.0",
		Port:        "23234",
		HostKeyPath: ".ssh/term_key",
	},
}

type SSHRuntimeConfig struct {
	Host        string
	Port        string
	HostKeyPath string
}

type SanityRuntimeConfig struct {
	ProjectID  string
	Dataset    string
	APIVersion string
}

var (
	globalConfig *Config
	envLoadOnce  sync.Once
)

// Load loads configuration from file or uses defaults
func Load(configPath string) (*Config, error) {
	loadEnvFile()

	// Try to load from file
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err == nil {
			var cfg Config
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return nil, err
			}
			// Apply defaults for any missing values
			applyDefaults(&cfg)
			globalConfig = &cfg
			return globalConfig, nil
		}
	}

	// Try loading from default locations
	defaultPaths := []string{
		"config.yaml",
		"config.yml",
		".portfolio-tui.yaml",
		".portfolio-tui.yml",
	}

	for _, path := range defaultPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			var cfg Config
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				continue
			}
			applyDefaults(&cfg)
			globalConfig = &cfg
			return globalConfig, nil
		}
	}

	// Use defaults if no config file found
	cfg := defaultConfig
	globalConfig = &cfg
	return globalConfig, nil
}

// Get returns the global config instance
func Get() *Config {
	loadEnvFile()

	if globalConfig == nil {
		globalConfig = &defaultConfig
	}
	return globalConfig
}

func loadEnvFile() {
	envLoadOnce.Do(func() {
		// Ignore missing .env files because container/runtime env may already be set.
		_ = godotenv.Load()
	})
}

// applyDefaults fills in missing values with defaults
func applyDefaults(cfg *Config) {
	if cfg.App.Name == "" {
		cfg.App.Name = defaultConfig.App.Name
	}
	if cfg.App.Version == "" {
		cfg.App.Version = defaultConfig.App.Version
	}
	if cfg.App.Runtime == "" {
		cfg.App.Runtime = defaultConfig.App.Runtime
	}
	if cfg.Owner.FullName == "" {
		cfg.Owner.FullName = defaultConfig.Owner.FullName
	}
	if cfg.Owner.Username == "" {
		cfg.Owner.Username = defaultConfig.Owner.Username
	}
	if cfg.Owner.Tagline == "" {
		cfg.Owner.Tagline = defaultConfig.Owner.Tagline
	}
	if cfg.Owner.Role == "" {
		cfg.Owner.Role = defaultConfig.Owner.Role
	}
	if cfg.Network.Host == "" {
		cfg.Network.Host = defaultConfig.Network.Host
	}
	if cfg.Contact.Email == "" {
		cfg.Contact.Email = defaultConfig.Contact.Email
	}
	if cfg.Contact.Phone == "" {
		cfg.Contact.Phone = defaultConfig.Contact.Phone
	}
	if cfg.Social.GitHub == "" {
		cfg.Social.GitHub = defaultConfig.Social.GitHub
	}
	if cfg.Social.LinkedIn == "" {
		cfg.Social.LinkedIn = defaultConfig.Social.LinkedIn
	}
	if cfg.Social.Twitter == "" {
		cfg.Social.Twitter = defaultConfig.Social.Twitter
	}
	if cfg.Social.Website == "" {
		cfg.Social.Website = defaultConfig.Social.Website
	}
	if cfg.Social.YouTube == "" {
		cfg.Social.YouTube = defaultConfig.Social.YouTube
	}
	if cfg.Social.Instagram == "" {
		cfg.Social.Instagram = defaultConfig.Social.Instagram
	}
	if cfg.Branding.SidebarTitle == "" {
		cfg.Branding.SidebarTitle = defaultConfig.Branding.SidebarTitle
	}
	if cfg.Branding.AsciiLogo == "" {
		cfg.Branding.AsciiLogo = defaultConfig.Branding.AsciiLogo
	}
	if cfg.Branding.SidebarLogo == "" {
		cfg.Branding.SidebarLogo = defaultConfig.Branding.SidebarLogo
	}
	if cfg.SSH.Host == "" {
		cfg.SSH.Host = defaultConfig.SSH.Host
	}
	if cfg.SSH.Port == "" {
		cfg.SSH.Port = defaultConfig.SSH.Port
	}
	if cfg.SSH.HostKeyPath == "" {
		cfg.SSH.HostKeyPath = defaultConfig.SSH.HostKeyPath
	}
}

// getEnvOrDefault returns the value of an environment variable or a fallback.
func getEnvOrDefault(env, fallback string) string {
	loadEnvFile()

	v := strings.TrimSpace(os.Getenv(env))
	if v != "" {
		return trimOptionalQuotes(v)
	}
	return fallback
}

func trimOptionalQuotes(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}

	return v
}

// GetSSHConfig returns SSH config values, preferring environment variables.
func GetSSHConfig(cfg *Config) SSHRuntimeConfig {
	if cfg == nil {
		cfg = Get()
	}

	return SSHRuntimeConfig{
		Host:        getEnvOrDefault("PORTFOLIO_SSH_HOST", cfg.SSH.Host),
		Port:        getEnvOrDefault("PORTFOLIO_SSH_PORT", cfg.SSH.Port),
		HostKeyPath: getEnvOrDefault("PORTFOLIO_SSH_KEY", cfg.SSH.HostKeyPath),
	}
}

// GetSanityConfig returns Sanity config values from the shared env loading path.
func GetSanityConfig() SanityRuntimeConfig {
	return SanityRuntimeConfig{
		ProjectID:  getEnvOrDefault("SANITY_PROJECT_ID", "0oqq6f9z"),
		Dataset:    getEnvOrDefault("SANITY_DATASET", "production"),
		APIVersion: getEnvOrDefault("SANITY_API_VERSION", "2024-12-21"),
	}
}
