package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
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
	Comment string `json:"comment"`
}

type DNSRecordList struct {
	Result []DNSEntry `json:"result"`
}

type DNSUpdateResult struct {
	Result DNSEntry `json:"result"`
}

// Verify API credentials
func verifyCFToken() error {
	logger.Debug("Verifying Cloudflare API credentials")
	res, err := httpGETRequest("https://api.cloudflare.com/client/v4/user/tokens/verify")
	if err != nil {
		return errors.New("invalid Cloudflare API token")
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("request failed with error code %d", res.StatusCode)
	}

	defer res.Body.Close()

	return nil
}

// List DNS records
func getRecords() ([]DNSEntry, error) {
	res, err := httpGETRequest("https://api.cloudflare.com/client/v4/zones/" + cfZoneId + "/dns_records")
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with error code %d", res.StatusCode)
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

// Update DNS record
func updateRecord(record DNSEntry, currentIP net.IP) error {
	//
	// Update IP in the entry
	logger.Info(fmt.Sprintf("Updating DNS record %s", record.Name))
	record.Content = currentIP.String()
	record.Comment = fmt.Sprintf("Updated by go-cloudflare-ddns on %s", time.Now().Format("2006/01/02 15:04"))
	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	res, err := httpPATCHRequest("https://api.cloudflare.com/client/v4/zones/"+cfZoneId+"/dns_records/"+record.Id, jsonData)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with error code %d", res.StatusCode)
	}

	var dnsUpdateResult DNSUpdateResult

	err = parseJSON(res, &dnsUpdateResult)
	if err != nil {
		return err
	}

	if !net.ParseIP(dnsUpdateResult.Result.Content).Equal(currentIP) {
		return errors.New("IP address doesn't match")
	}

	return nil
}
