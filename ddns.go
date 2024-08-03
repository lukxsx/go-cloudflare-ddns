package main

// Steps:
// Query current IP of the commputer
// Make DNS query on the domain
// If different, update

func checkAndUpdate(domain string) error {
	// Get current IP
	currentIP, err := getIP()
	if err != nil {
		return err
	}

	// Query domain's IP
	domainIP, err := dnsQuery(domain)
	if err != nil {
		return err
	}

	// Check if they match
	if currentIP.Equal(domainIP) {
		logger.Debug("IP addresses match. Nothing to do.")
	} else {
		logger.Info("IP has changed")
		// updateDomain(domain)
	}

	return nil
}
