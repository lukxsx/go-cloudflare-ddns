package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lukxsx/go-cloudflare-ddns/ddns"
)

var (
	logger *slog.Logger
)

func main() {
	godotenv.Load()

	// Setup logger
	logger = setupLogger()
	logger.Info("Starting Cloudflare DDNS client")

	// Read configuration values from environment variables
	clientConfig, err := Configure()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Create the client
	_, err = ddns.CreateClient(clientConfig, logger)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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
func Configure() (*ddns.ClientConfig, error) {
	logger.Debug("Reading configuration from environment variables")
	c := &ddns.ClientConfig{}

	// Parse Cloudflare API credentials
	c.CFApiToken = os.Getenv("CF_API_TOKEN")
	c.CFZoneID = os.Getenv("CF_ZONE_ID")
	if c.CFApiToken == "" || c.CFZoneID == "" {
		return nil, errors.New("missing Cloudflare API credentials")
	}

	// Parse domains from comma-separated list
	domains := os.Getenv("DOMAINS")
	if domains == "" {
		return nil, errors.New("missing domain list")
	}
	split_domains := strings.Split(domains, ",")
	if len(split_domains) == 0 {
		return nil, errors.New("missing domain list")
	}
	c.Domains = split_domains

	// TODO: validate all values

	logger.Debug("CF_API_TOKEN: " + c.CFApiToken)
	logger.Debug("CF_ZONE_ID: " + c.CFZoneID)
	logger.Debug("DOMAINS: " + strings.Join(c.Domains, ", "))

	return c, nil
}
