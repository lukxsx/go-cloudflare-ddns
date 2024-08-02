package ddns

import (
	"log/slog"
	"net"
)

// Client config struct
type ClientConfig struct {
	CFApiToken string
	CFZoneID   string
	Domains    []string
}

// Client struct
type Client struct {
	// Cloudflare API credentials
	config ClientConfig

	currentIP net.IP
	logger    *slog.Logger
}

func CreateClient(config *ClientConfig, logger *slog.Logger) (*Client, error) {
	client := &Client{config: *config, logger: logger}
	err := client.UpdateIP()
	if err != nil {
		return nil, err
	}

	return client, nil
}
