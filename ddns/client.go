package ddns

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Client struct
type Client struct {
	// Cloudflare API credentials
	CFApiToken string
	CFZoneID   string
	Domains    []string

	currentIP net.IP
}

// Read configuration values from environment variables
func (c *Client) Configure() error {
	// Check for .env file
	godotenv.Load()

	// Parse Cloudflare API credentials
	c.CFApiToken = os.Getenv("CF_API_TOKEN")
	c.CFZoneID = os.Getenv("CF_ZONE_ID")
	if c.CFApiToken == "" || c.CFZoneID == "" {
		return errors.New("missing Cloudflare API credentials")
	}

	// Parse domains from comma-separated list
	domains := os.Getenv("DOMAINS")
	if domains == "" {
		return errors.New("missing domain list")
	}
	split_domains := strings.Split(domains, ",")
	if len(split_domains) == 0 {
		return errors.New("missing domain list")
	}
	c.Domains = split_domains

	// TODO: validate all values

	return nil
}
