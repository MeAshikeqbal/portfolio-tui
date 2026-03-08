package utils

import "github.com/MeAshikeqbal/portfolio-tui/internal/config"

// GetFullName returns the full name from config
func GetFullName() string {
	return config.Get().Owner.FullName
}

// GetTagline returns the tagline from config
func GetTagline() string {
	return config.Get().Owner.Tagline
}

// GetHost returns the host from config
func GetHost() string {
	return config.Get().Network.Host
}

// GetUsername returns the username from config
func GetUsername() string {
	return config.Get().Owner.Username
}
