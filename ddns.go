package main

import (
	"fmt"
	"net"
)

// Steps:
// Query current IP of the commputer
// Make DNS query on the domain
// If different, update

func checkAndUpdate() error {
	// Get current IP
	currentIP, err := getIP()
	if err != nil {
		return err
	}

	// Get DNS records
	dnsRecords, err := getDNSRecords()
	if err != nil {
		return err
	}

	var recordsToCheck []DNSEntry
	for _, record := range dnsRecords {
		if contains(domains, record.Name) {
			logger.Debug(fmt.Sprintf("Found record %s: %s", record.Name, record.Content))
			recordsToCheck = append(recordsToCheck, record)
		}
	}

	for _, record := range recordsToCheck {
		recordIP := net.ParseIP(record.Content)
		if !recordIP.Equal(currentIP) {
			logger.Debug(fmt.Sprintf("IPs don't match: %s: %s, local: %s", record.Name, record.Content, currentIP.String()))
			logger.Info(fmt.Sprintf("Updating DNS record for %s", record.Name))
			// updateRecord(record, currentIp)
		} else {
			logger.Debug(fmt.Sprintf("IPs match: %s: %s, local: %s", record.Name, record.Content, currentIP.String()))
		}
	}

	return nil
}
