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

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

func addOfficialEmails() {
	var wg sync.WaitGroup
	for _, user := range hubUsers {
		wg.Add(1)
		go func(user HubSeedUser) {
			defer wg.Done()
			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			for _, job := range user.Jobs {
				addOfficialEmail(user, authToken, job.Website)
			}
		}(user)
	}
	wg.Wait()
}

func addOfficialEmail(user HubSeedUser, authToken string, domain string) {
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

	if resp.StatusCode == http.StatusConflict {
		color.Yellow("skipping %q because it already exists", email)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error adding official email: %q %v", email, resp.Status)
	}
	color.Green("added official email: %q for %q\n", email, user.Handle)

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
	sleepInterval := 5 * time.Second
	for i := 0; i < 5; i++ {
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
		time.Sleep(sleepInterval)
		sleepInterval *= 3
	}

	if messageID == "" {
		log.Fatalf("no verification email found in mailpit for %q", email)
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

	color.Cyan("verified email: %s\n", email)

	// Delete the email from mailpit
	deleteBody := struct {
		IDs []string `json:"IDs"`
	}{
		IDs: []string{messageID},
	}

	deleteJSON, err := json.Marshal(deleteBody)
	if err != nil {
		log.Fatalf("failed to marshal delete request: %v", err)
	}

	deleteReq, err := http.NewRequest(
		"DELETE",
		mailPitURL+"/api/v1/messages",
		bytes.NewBuffer(deleteJSON),
	)
	if err != nil {
		log.Fatalf("failed to create delete request: %v", err)
	}
	deleteReq.Header.Set("Accept", "application/json")
	deleteReq.Header.Set("Content-Type", "application/json")

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		log.Fatalf("failed to delete email: %v", err)
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK {
		log.Fatalf("failed to delete email. status: %d", deleteResp.StatusCode)
	}
}
