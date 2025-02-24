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

	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func addOfficialEmails() {
	for _, user := range hubUsers {
		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Fatalf("no auth token found for %s", user.Email)
		}
		authToken := tokenI.(string)

		for _, domain := range user.WorkHistoryDomains {
			addOfficialEmail(user, authToken, domain)
		}
	}
}

func addOfficialEmail(user HubUser, authToken string, domain string) {
	email := user.Handle + "@" + domain
	body := hub.AddOfficialEmailRequest{
		Email: common.EmailAddress(email),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error marshalling body: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/add-official-email",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error adding official email: %v", resp.Status)
	}

	// Wait for the email to be sent
	<-time.After(2 * time.Second)

	// Query mailpit for the TFA email
	baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
	if err != nil {
		log.Fatalf("failed to parse mailpit URL: %v", err)
	}
	query := url.Values{}
	query.Add(
		"query",
		fmt.Sprintf("to:%s subject:Vetchi - Confirm Email Ownership", email),
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
		log.Fatal("no verification email found in mailpit")
	}

	// Get the email content
	mailResp, err := http.Get(mailPitURL + "/api/v1/message/" + messageID)
	if err != nil {
		log.Fatalf("failed to get email content: %v", err)
	}
	defer mailResp.Body.Close()

	mailBody, err := io.ReadAll(mailResp.Body)
	if err != nil {
		log.Fatalf("failed to read email body: %v", err)
	}

	// Extract TFA code from email
	re := regexp.MustCompile(
		`profile page of your vetchi account:\s*([a-zA-Z0-9]+)\s*`,
	)
	matches := re.FindStringSubmatch(string(mailBody))
	if len(matches) < 2 {
		log.Fatal("could not find TFA code in email")
	}
	tfaCode := matches[1]

	log.Printf("TFA code: %s", tfaCode)
}
