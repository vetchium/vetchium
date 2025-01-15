package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Message corresponds to the message object in mailpit
type Message struct {
	ID string `json:"ID"`
}

type MailPitResponse struct {
	Messages []Message `json:"messages"`
}

// employerSignin handles the employer signin process including TFA
func employerSignin(email, password, clientID string) string {
	// Check if we already have a token
	if token, ok := sessionTokens.Load(email); ok {
		return token.(string)
	}

	// Step 1: Initial signin
	signinReq := struct {
		ClientID string `json:"client_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		ClientID: clientID,
		Email:    email,
		Password: password,
	}

	jsonBody, err := json.Marshal(signinReq)
	if err != nil {
		log.Fatalf("failed to marshal signin request: %v", err)
	}

	resp, err := http.Post(
		serverURL+"/employer/signin",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Fatalf("failed to make signin request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to read error response body: %v", err)
		}
		log.Fatalf(
			"signin failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var signinResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&signinResp); err != nil {
		log.Fatalf("failed to decode signin response: %v", err)
	}

	// Step 2: Wait for TFA email and get code from mailpit
	time.Sleep(2 * time.Second)

	// Query mailpit for the TFA email
	baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
	if err != nil {
		log.Fatalf("failed to parse mailpit URL: %v", err)
	}
	query := url.Values{}
	query.Add(
		"query",
		fmt.Sprintf("to:%s subject:Vetchi Two Factor Authentication", email),
	)
	baseURL.RawQuery = query.Encode()

	var messageID string
	for i := 0; i < 3; i++ {
		mailResp, err := http.Get(baseURL.String())
		if err != nil {
			log.Fatalf("failed to query mailpit: %v", err)
		}

		var mailPitResp MailPitResponse
		if err := json.NewDecoder(mailResp.Body).Decode(&mailPitResp); err != nil {
			log.Fatalf("failed to decode mailpit response: %v", err)
		}
		mailResp.Body.Close()

		if len(mailPitResp.Messages) > 0 {
			messageID = mailPitResp.Messages[0].ID
			break
		}
		time.Sleep(2 * time.Second)
	}

	if messageID == "" {
		log.Fatal("no TFA email found in mailpit")
	}

	// Get the email content
	mailResp, err := http.Get(mailPitURL + "/api/v1/message/" + messageID)
	if err != nil {
		log.Fatalf("failed to get email content: %v", err)
	}
	defer mailResp.Body.Close()

	body, err := io.ReadAll(mailResp.Body)
	if err != nil {
		log.Fatalf("failed to read email body: %v", err)
	}

	// Extract TFA code from email
	re := regexp.MustCompile(`Token:\s*([a-zA-Z0-9]+)\s*`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		log.Fatal("could not find TFA code in email")
	}
	tfaCode := matches[1]

	// Delete the email from mailpit
	deleteBody := struct {
		IDs []string `json:"IDs"`
	}{
		IDs: []string{messageID},
	}

	jsonBody, err = json.Marshal(deleteBody)
	if err != nil {
		log.Fatalf("failed to marshal delete request: %v", err)
	}

	req, err := http.NewRequest(
		"DELETE",
		mailPitURL+"/api/v1/messages",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Fatalf("failed to create delete request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	deleteResp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to delete email: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(deleteResp.Body)
		if err != nil {
			log.Fatalf("failed to read delete error response: %v", err)
		}
		log.Fatalf(
			"failed to delete email with status %d: %s",
			deleteResp.StatusCode,
			string(bodyBytes),
		)
	}

	// Step 3: Submit TFA code
	tfaReq := struct {
		TFAToken   string `json:"tfa_token"`
		TFACode    string `json:"tfa_code"`
		RememberMe bool   `json:"remember_me"`
	}{
		TFAToken:   signinResp.Token,
		TFACode:    tfaCode,
		RememberMe: true,
	}

	jsonBody, err = json.Marshal(tfaReq)
	if err != nil {
		log.Fatalf("failed to marshal TFA request: %v", err)
	}

	resp, err = http.Post(
		serverURL+"/employer/tfa",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Fatalf("failed to make TFA request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to read error response body: %v", err)
		}
		log.Fatalf(
			"TFA failed with status %d: %s",
			resp.StatusCode,
			string(bodyBytes),
		)
	}

	var tfaResp struct {
		SessionToken string `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tfaResp); err != nil {
		log.Fatalf("failed to decode TFA response: %v", err)
	}

	// Store the token for future use
	sessionTokens.Store(email, tfaResp.SessionToken)
	return tfaResp.SessionToken
}
