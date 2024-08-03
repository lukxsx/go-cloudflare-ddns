package main

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
)

// DNSResponse represents the JSON response from the CloudFlare DNS over HTTPS API
type DNSResponse struct {
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

// Verify API credentials
func verifyCFToken() error {
	logger.Debug("Verifying Cloudflare API credentials")
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	req.Header.Set("Authorization", "Bearer "+cfApiToken)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return errors.New("invalid Cloudflare API token")
	}

	return nil
}

// Make a DNS query using Cloudflare's DNS over HTTPS API
func dnsQuery(domain string) (net.IP, error) {
	logger.Debug("Making DNS A query for: " + domain)
	req, err := http.NewRequest("GET", "https://1.1.1.1/dns-query?name="+domain+"&type=A", nil)
	req.Header.Set("accept", "application/dns-json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var dnsResponse DNSResponse

	err = json.Unmarshal([]byte(body), &dnsResponse)
	if err != nil {
		return nil, err
	}

	if len(dnsResponse.Answer) < 1 {
		return nil, errors.New("DNS records not found")
	}

	firstAnswer := dnsResponse.Answer[0]
	ip := net.ParseIP(firstAnswer.Data)

	if ip == nil {
		return nil, errors.New("failed to parse IP")
	}

	logger.Debug("DNS A query result: " + ip.String())

	return ip, nil
}
