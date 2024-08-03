package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	logger     *slog.Logger
	cfApiToken string
	cfZoneId   string
	domains    []string
)

func main() {
	godotenv.Load()

	// Setup logger
	logger = setupLogger()
	logger.Info("Starting Cloudflare DDNS client")

	// Read configuration values from environment variables
	err := configure()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	checkAndUpdate("google.com")
}

// Setup logger
func setupLogger() *slog.Logger {
	var logLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))
	if os.Getenv("DEBUG") == "true" {
		logLevel.Set(slog.LevelDebug)
	}
	return slog.New(h)
}

// Read configuration values from environment variables
func configure() error {
	logger.Debug("Reading configuration from environment variables")

	// Parse Cloudflare API credentials
	cfApiToken = os.Getenv("CF_API_TOKEN")
	cfZoneId = os.Getenv("CF_ZONE_ID")
	if cfApiToken == "" || cfZoneId == "" {
		return errors.New("missing Cloudflare API credentials")
	}

	// Parse domains from comma-separated list
	domainStr := os.Getenv("DOMAINS")
	if domainStr == "" {
		return errors.New("missing domain list")
	}
	domains = strings.Split(domainStr, ",")
	if len(domains) == 0 {
		return errors.New("missing domain list")
	}

	// TODO: validate all values

	logger.Debug("CF_API_TOKEN: " + cfApiToken)
	logger.Debug("CF_ZONE_ID: " + cfZoneId)
	logger.Debug("DOMAINS: " + strings.Join(domains, ", "))

	return nil
}
