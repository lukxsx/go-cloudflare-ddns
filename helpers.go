package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	httpClient http.Client
)

// Check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// Helper function to make a HTTP GET request with custom headers
func httpGETRequest(url string) (*http.Response, error) {
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfApiToken)

	// Create client and make request
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Debug(fmt.Sprintf("%d GET %s", res.StatusCode, url))

	return res, nil
}

// Helper function to parse JSON from the request body
func parseJSON(res *http.Response, v any) error {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &v)
	if err != nil {
		return err
	}

	return nil
}
