package main

import (
	"errors"
	"net"
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

type DNSEntry struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Id      string `json:"id"`
}

type DNSRecordList struct {
	Result []DNSEntry `json:"result"`
}

// Verify API credentials
func verifyCFToken() error {
	logger.Debug("Verifying Cloudflare API credentials")
	res, err := httpGETRequest("https://api.cloudflare.com/client/v4/user/tokens/verify")
	if err != nil || res.StatusCode != 200 {
		return errors.New("invalid Cloudflare API token")
	}

	defer res.Body.Close()

	return nil
}

// List DNS records
func listDNSRecords() ([]DNSEntry, error) {
	res, err := httpGETRequest("https://api.cloudflare.com/client/v4/zones/" + cfZoneId + "/dns_records")
	if err != nil {
		return nil, err
	}

	var dnsRecordList DNSRecordList

	err = parseJSON(res, &dnsRecordList)
	if err != nil {
		return nil, err
	}

	if len(dnsRecordList.Result) < 1 {
		return nil, errors.New("DNS records not found")
	}

	return dnsRecordList.Result, nil
}

// Make a DNS query using Cloudflare's DNS over HTTPS API
func dnsQuery(domain string) (net.IP, error) {
	logger.Debug("Making DNS A query for: " + domain)
	res, err := httpGETRequest("https://1.1.1.1/dns-query?name=" + domain + "&type=A")
	if err != nil {
		return nil, err
	}

	var dnsResponse DNSResponse
	err = parseJSON(res, &dnsResponse)
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
