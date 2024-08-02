package ddns

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

// Query the DNS A record for a domain
func DNSQuery(domain string) (net.IP, error) {
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

	return ip, nil
}
