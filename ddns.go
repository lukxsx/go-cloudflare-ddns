package main

import (
	"fmt"
	"net"
)

func checkAndUpdate() error {
	logger.Info("Checking for updates")
	// Get current IP
	currentIP, err := getIP()
	if err != nil {
		return err
	}

	// Get DNS records
	dnsRecords, err := getRecords()
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
			err = updateRecord(record, currentIP)
			if err != nil {
				return err
			}
		} else {
			logger.Debug(fmt.Sprintf("IPs match: %s: %s, local: %s", record.Name, record.Content, currentIP.String()))
		}
	}

	return nil
}
