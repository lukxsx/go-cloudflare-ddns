package ddns

import (
	"errors"
	"io"
	"net"
	"net/http"
)

// Fetches the public IP address using the ipify API
func (c *Client) getMyIPAddress() (net.IP, error) {
	c.logger.Debug("Fetching public IP address")
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(body))
	if ip == nil {
		return nil, errors.New("failed to parse IP address")
	}

	c.logger.Debug("Fetched IP: " + ip.String())

	return ip, nil
}

func (c *Client) UpdateIP() error {
	ip, err := c.getMyIPAddress()
	if err != nil {
		return err
	}

	c.currentIP = ip
	return nil
}
