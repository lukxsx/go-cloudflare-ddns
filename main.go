package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	logger         *slog.Logger
	cfApiToken     string
	cfZoneId       string
	domains        []string
	updateInterval int
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

	// Check if CF token is valid
	err = verifyCFToken()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initial update check
	err = checkAndUpdate()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check for updates: %s", err.Error()))
	}

	// Set up goroutine with timer to check for updates
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// Start the scheduled function
	wg.Add(1)
	go runner(&wg, stopCh)

	// Set up channel to receive interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for an interrupt signal
	<-signalChan

	// Signal the asynchronous function to stop
	close(stopCh)

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("Exiting.")
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

	intervalStr := os.Getenv("UPDATE_INTERVAL")
	if intervalStr == "" {
		updateInterval = 10 // minutes
	} else {
		intervalInt, err := strconv.Atoi(intervalStr)
		if err != nil {
			return errors.New("invalid update interval value")
		}
		updateInterval = intervalInt
	}
	logger.Info(fmt.Sprintf("Update interval set to %d minutes", updateInterval))

	// TODO: validate all values

	logger.Debug("CF_API_TOKEN: " + cfApiToken)
	logger.Debug("CF_ZONE_ID: " + cfZoneId)
	logger.Debug("DOMAINS: " + strings.Join(domains, ", "))

	return nil
}

// Run the update routine periodically
func runner(wg *sync.WaitGroup, stopCh <-chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(updateInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			err := checkAndUpdate()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to check for updates: %s", err.Error()))
			}
		}
	}
}
