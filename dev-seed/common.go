package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	serverURL  = "http://localhost:8081"
	mailPitURL = "http://localhost:8025"
)

// makeRequest is a helper function to make HTTP requests
func makeRequest(
	method, path string,
	token string,
	reqBody interface{},
	respBody interface{},
) {
	var body io.Reader
	if reqBody != nil {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatalf("failed to marshal request body: %v", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, serverURL+path, body)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to read error response body: %v", err)
		}
		log.Fatalf(
			"request failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			log.Fatalf("failed to decode response body: %v", err)
		}
	}
}
