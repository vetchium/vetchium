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
	"sync"
	"time"

	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func addOfficialEmails() {
	var wg sync.WaitGroup
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubUser) {
			defer wg.Done()
			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			for _, domain := range user.WorkHistoryDomains {
				addOfficialEmail(user, authToken, domain)
			}
		}(user)
	}
	wg.Wait()
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
		log.Fatalf("error adding official email: %q %v", email, resp.Status)
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

	// log.Printf("email body: %s", string(mailBody))

	// Parse the email JSON
	var emailContent struct {
		Text string `json:"Text"`
	}
	if err := json.Unmarshal(mailBody, &emailContent); err != nil {
		log.Fatalf("failed to parse email JSON: %v", err)
	}

	// Extract TFA code from email
	re := regexp.MustCompile(
		`(?m)^\s*([a-zA-Z0-9]{4})\s*$`,
	)
	matches := re.FindStringSubmatch(emailContent.Text)
	if len(matches) < 2 {
		log.Fatal("could not find TFA code in email")
	}
	emailConfirmationCode := matches[1]

	// log.Printf("Email Confirmation code: %s", emailConfirmationCode)

	// Update the user's email confirmation code
	mailConfirmBody := hub.VerifyOfficialEmailRequest{
		Email: common.EmailAddress(email),
		Code:  emailConfirmationCode,
	}

	mailConfirmJSON, err := json.Marshal(mailConfirmBody)
	if err != nil {
		log.Fatalf("failed to marshal mail confirm body: %v", err)
	}

	mailConfirmReq, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/verify-official-email",
		bytes.NewBuffer(mailConfirmJSON),
	)
	if err != nil {
		log.Fatalf("failed to create mail confirm request: %v", err)
	}
	mailConfirmReq.Header.Set("Authorization", "Bearer "+authToken)
	mailConfirmReq.Header.Set("Content-Type", "application/json")

	mailConfirmResp, err := http.DefaultClient.Do(mailConfirmReq)
	if err != nil {
		log.Fatalf("failed to send mail confirm request: %v", err)
	}
	defer mailConfirmResp.Body.Close()

	if mailConfirmResp.StatusCode != http.StatusOK {
		log.Fatalf("failed to verify email: %v", mailConfirmResp.Status)
	}

	log.Printf("verified email: %s", email)
}
