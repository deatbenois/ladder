// Package config provides configuration loading and validation for ladder.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration.
type Config struct {
	// Port is the port the server listens on.
	Port int

	// Host is the host the server binds to.
	Host string

	// Timeout is the maximum duration for proxy requests.
	Timeout time.Duration

	// AllowedHosts is a list of allowed target hosts. Empty means all are allowed.
	AllowedHosts []string

	// BlockedHosts is a list of blocked target hosts.
	BlockedHosts []string

	// UserAgent is the User-Agent header sent with proxied requests.
	UserAgent string

	// XForwardedFor controls whether to set X-Forwarded-For header.
	XForwardedFor bool

	// MaxBodySize is the maximum response body size in bytes (0 = unlimited).
	MaxBodySize int64

	// ProxyURL is an optional upstream proxy URL.
	ProxyURL *url.URL
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:          8080,
		Host:          "",
		// Increased from 30s to 60s — some sites I proxy are slow to respond.
		Timeout:       60 * time.Second,
		AllowedHosts:  []string{},
		BlockedHosts:  []string{},
		// Using a more generic Chrome UA — gets better results on sites that block bots.
		UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		XForwardedFor: false,
		// Default to 10MB limit — prevents runaway downloads on my low-memory VPS.
		MaxBodySize:   10 * 1024 * 1024,
		ProxyURL:      nil,
	}
}

// LoadFromEnv populates the Config from environment variables.
// Environment variables take precedence over defaults.
func LoadFromEnv() (*Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv("LADDER_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid LADDER_PORT %q: %w", v, err)
		}
		cfg.Port = port
	}

	if v := os.Getenv("LADDER_HOST"); v != "" {
		cfg.Host = v
	}

	if v := os.Getenv("LADDER_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid LADDER_TIMEOUT %q: %w", v, err)
		}
		cfg.Timeout = d
	}

	if v := os.Getenv("LADDER_ALLOWED_HOSTS"); v != "" {
		cfg.AllowedHosts = splitTrimmed(v, ",")
	}

	if v := os.Getenv("LADDER_BLOCKED_HOSTS"); v != "" {
		cfg.BlockedHosts = splitTrimmed(v, ",")
	}

	if v := os.Getenv("LADDER_USER_AGENT"); v != "" {
		cfg.UserAgent = v
	}

	if v := os.Getenv("LADDER_X_FORWARDED_FOR"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("invalid LADDER_X_FORWARDED_FOR %q: %w", v, err)
		}
		cfg.XForwardedFor = b
	}

	if v := os.Getenv("LADDER_MAX_BODY_SIZE"); v != "" {
		size, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid LADDER_MAX_BODY_SIZE %q: %w", v, err)
		}
		cfg.MaxBodySize = size
	}

	if v := os.Getenv("LADDER_PROXY_URL"); v != "" {
		u, err := url.Parse(v)
		if err != nil {
			return nil, fmt.Err