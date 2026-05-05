package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort    = "8080"
	defaultTimeout = 30 * time.Second
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	ProxyUserAgent string
	AllowedHosts   string
	BlockedHosts   string
}

// loadConfig reads configuration from environment variables with sensible defaults.
func loadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	userAgent := os.Getenv("PROXY_USER_AGENT")
	if userAgent == "" {
		// Mimic a real browser UA to reduce likelihood of being blocked by target sites
		userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	}

	return &Config{
		Port:           port,
		ReadTimeout:    defaultTimeout,
		WriteTimeout:   defaultTimeout,
		IdleTimeout:    120 * time.Second,
		ProxyUserAgent: userAgent,
		AllowedHosts:   os.Getenv("ALLOWED_HOSTS"),
		BlockedHosts:   os.Getenv("BLOCKED_HOSTS"),
	}
}

func main() {
	cfg := loadConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/", proxyHandler(cfg))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Channel to listen for OS signals for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Ladder listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}

// healthzHandler responds to health check requests.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

// proxyHandler returns an HTTP handler that proxies requests through ladder.
func proxyHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Target URL is expected as a query parameter: ?url=https://example.com
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		log.Printf("Proxying request for: %s", targetURL)

		// TODO: implement full proxy logic with host filtering and content rewriting
		http.Error(w, "proxy not yet implemented", http.StatusNotImplemented)
	}
}
