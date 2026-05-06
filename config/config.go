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
		Timeout:       30 * time.Second,
		AllowedHosts:  []string{},
		BlockedHosts:  []string{},
		UserAgent:     "Mozilla/5.0 (compatible; Ladder/1.0)",
		XForwardedFor: false,
		MaxBodySize:   0,
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
			return nil, fmt.Errorf("invalid LADDER_PROXY_URL %q: %w", v, err)
		}
		cfg.ProxyURL = u
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port %d is out of valid range (1-65535)", c.Port)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %s", c.Timeout)
	}
	return nil
}

// Addr returns the host:port string for the server to listen on.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// splitTrimmed splits s by sep and trims whitespace from each element,
// omitting empty strings.
func splitTrimmed(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}
