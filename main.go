package main

import (
	"log/slog"
	"os"

	"github.com/lukxsx/go-cloudflare-ddns/ddns"
)

func main() {
	logger := setupLogger()
	logger.Info("Starting Cloudflare DDNS client")

	client := &ddns.Client{Logger: logger}
	err := client.Configure()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// Setup logger
func setupLogger() *slog.Logger {
	var programLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	if os.Getenv("DEBUG") == "true" {
		programLevel.Set(slog.LevelDebug)
	}
	return slog.New(h)
}
