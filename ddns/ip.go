package ddns

import (
	"errors"
	"io"
	"net"
	"net/http"
)

// Fetches the public IP address using the ipify API
func GetMyIPAddress() (net.IP, error) {
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

	return ip, nil
}
