package utils

import "os"

// GetFullName returns the full name from environment or default
func GetFullName() string {
	name := os.Getenv("FULL_NAME")
	if name == "" {
		return "Ashik Eqbal"
	}
	return name
}

// GetTagline returns the tagline from environment or default
func GetTagline() string {
	tagline := os.Getenv("TAGLINE")
	if tagline == "" {
		return "Portfolio"
	}
	return tagline
}

// GetHost returns the host from environment or default
func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		return "localhost"
	}
	return host
}
