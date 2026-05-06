package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Port == 0 {
		t.Error("expected non-zero default port")
	}

	if cfg.Timeout == 0 {
		t.Error("expected non-zero default timeout")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Run("uses defaults when no env vars set", func(t *testing.T) {
		clearEnv()
		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defaults := DefaultConfig()
		if cfg.Port != defaults.Port {
			t.Errorf("expected port %d, got %d", defaults.Port, cfg.Port)
		}
	})

	t.Run("reads PORT from env", func(t *testing.T) {
		clearEnv()
		os.Setenv("PORT", "9090")
		defer os.Unsetenv("PORT")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
	})

	t.Run("invalid PORT returns error", func(t *testing.T) {
		clearEnv()
		os.Setenv("PORT", "not-a-number")
		defer os.Unsetenv("PORT")

		_, err := LoadFromEnv()
		if err == nil {
			t.Error("expected error for invalid PORT, got nil")
		}
	})

	t.Run("reads ALLOWED_HOSTS from env", func(t *testing.T) {
		clearEnv()
		os.Setenv("ALLOWED_HOSTS", "example.com, foo.org, bar.net")
		defer os.Unsetenv("ALLOWED_HOSTS")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.AllowedHosts) != 3 {
			t.Errorf("expected 3 allowed hosts, got %d", len(cfg.AllowedHosts))
		}
		if cfg.AllowedHosts[0] != "example.com" {
			t.Errorf("expected first host to be 'example.com', got '%s'", cfg.AllowedHosts[0])
		}
	})

	t.Run("reads BLOCKED_HOSTS from env", func(t *testing.T) {
		clearEnv()
		os.Setenv("BLOCKED_HOSTS", "evil.com,bad.org")
		defer os.Unsetenv("BLOCKED_HOSTS")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cfg.BlockedHosts) != 2 {
			t.Errorf("expected 2 blocked hosts, got %d", len(cfg.BlockedHosts))
		}
	})

	// Verify that USER_AGENT env var is respected when set.
	t.Run("reads USER_AGENT from env", func(t *testing.T) {
		clearEnv()
		os.Setenv("USER_AGENT", "MyCustomAgent/1.0")
		defer os.Unsetenv("USER_AGENT")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.UserAgent != "MyCustomAgent/1.0" {
			t.Errorf("expected user agent 'MyCustomAgent/1.0', got '%s'", cfg.UserAgent)
		}
	})

	// PORT=0 should be treated as invalid since binding to port 0 is not useful here.
	t.Run("PORT=0 returns error", func(t *testing.T) {
		clearEnv()
		os.Setenv("PORT", "0")
		defer os.Unsetenv("PORT")

		_, err := LoadFromEnv()
		if err == nil {
			t.Error("expected error for PORT=0, got nil")
		}
	})
}

func TestSplitTrimmed(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"a", []string{"a"}},
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b , c ", []string{"a", "b", "c"}},
		{"a,,b", []string{"a", "b"}},
		{",,,", []string{}},
	}

	for _, tc := range tests {
		result := splitTrimmed(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("splitTrimmed(%q): expected len %d, got %d", tc.input, len(tc.expected), len(result))
			continue
		}
		for i, v := range result {
			if v != tc.expected[i] {
				t.Errorf("splitTrimmed(%q)[%d]: expected %q, got %q", tc.input, i, tc.expected[i], v)
			}
		}
	}
}
