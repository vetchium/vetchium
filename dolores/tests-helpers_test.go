package dolores

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"sync"
	"time"

	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
	"github.com/vetchium/vetchium/typespec/hub"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	serverURL  = "http://localhost:8080"
	mailPitURL = "http://localhost:8025"
)

func setupTestDB() *pgxpool.Pool {
	connStr := os.Getenv("POSTGRES_URI")
	if connStr == "" {
		log.Fatal("POSTGRES_URI environment variable is required")
	}
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return pool
}

func seedDatabase(db *pgxpool.Pool, fileName string) {
	seed, err := os.ReadFile(fileName)
	Expect(err).ShouldNot(HaveOccurred())
	_, err = db.Exec(context.Background(), string(seed))
	Expect(err).ShouldNot(HaveOccurred())
}

// Message corresponds to the message object in mailpit
type Message struct {
	ID string `json:"ID"`
}

type MailPitResponse struct {
	Messages []Message `json:"messages"`
}

type MailPitDeleteRequest struct {
	IDs []string `json:"IDs"`
}

type SigninError struct {
	StatusCode int
}

func (e SigninError) Error() string {
	return fmt.Sprintf("signin failed with status code: %d", e.StatusCode)
}

// employerSigninAsync performs the signin operation asynchronously and returns
// the session token in the given token pointer. The wg will be decremented by
// one when the signin operation is completed, irrespective of whether it
// succeeds or fails.
func employerSigninAsync(
	clientID, email, password string,
	token *string,
	wg *sync.WaitGroup,
) {
	go func(clientID, email, password string, token *string, wg *sync.WaitGroup) {
		defer GinkgoRecover()
		defer wg.Done()
		var err error
		gotToken, err := employerSignin(clientID, email, password)
		Expect(err).ShouldNot(HaveOccurred())
		*token = gotToken
		fmt.Fprintf(
			GinkgoWriter,
			"email: %s, password: %s, gotToken: %s\n",
			email,
			password,
			gotToken,
		)
	}(
		clientID,
		email,
		password,
		token,
		wg,
	)
}

// Returns the session token for the employer with the given credentials
func employerSignin(clientID, email, password string) (string, error) {
	fmt.Fprintf(
		GinkgoWriter,
		"clientID: %s, email: %s, password: %s\n",
		clientID,
		email,
		password,
	)
	signinReqBody, err := json.Marshal(
		employer.EmployerSignInRequest{
			ClientID: clientID,
			Email:    common.EmailAddress(email),
			Password: common.Password(password),
		})
	Expect(err).ShouldNot(HaveOccurred())

	signinReq, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/signin",
		bytes.NewBuffer(signinReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())

	signinReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(signinReq)
	Expect(err).ShouldNot(HaveOccurred())
	if resp.StatusCode != http.StatusOK {
		// Callers of employerSignin may expect a fail intentionally
		return "", SigninError{StatusCode: resp.StatusCode}
	}

	var signinResp employer.EmployerSignInResponse
	err = json.NewDecoder(resp.Body).Decode(&signinResp)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(signinResp.Token).ShouldNot(BeEmpty())

	baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
	Expect(err).ShouldNot(HaveOccurred())
	query := url.Values{}
	query.Add(
		"query",
		fmt.Sprintf("to:%s subject:Vetchium Two Factor Authentication", email),
	)
	baseURL.RawQuery = query.Encode()

	url2 := baseURL.String()
	fmt.Fprintf(GinkgoWriter, "URL2: %s\n", url2)

	var messageID string
	for i := 0; i < 3; i++ {
		// Get the tfa code from the email by querying mailpit
		fmt.Fprintf(GinkgoWriter, "Sleeping to allow granger to email\n")
		<-time.After(10 * time.Second)
		fmt.Fprintf(GinkgoWriter, "Wokeup\n")

		mailPitReq1, err := http.NewRequest("GET", url2, nil)
		Expect(err).ShouldNot(HaveOccurred())
		mailPitReq1.Header.Add("Content-Type", "application/json")

		mailPitResp1, err := http.DefaultClient.Do(mailPitReq1)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mailPitResp1.StatusCode).Should(Equal(http.StatusOK))

		body, err := io.ReadAll(mailPitResp1.Body)
		Expect(err).ShouldNot(HaveOccurred())

		fmt.Fprintf(GinkgoWriter, "Body: %s\n", string(body))

		var mailPitResp1Obj MailPitResponse
		err = json.Unmarshal(body, &mailPitResp1Obj)
		Expect(err).ShouldNot(HaveOccurred())

		if len(mailPitResp1Obj.Messages) == 0 {
			continue
		}
		Expect(len(mailPitResp1Obj.Messages)).Should(Equal(1))
		messageID = mailPitResp1Obj.Messages[0].ID
		break
	}

	mailURL := "http://localhost:8025/api/v1/message/" + messageID
	fmt.Fprintf(GinkgoWriter, "Mail URL: %s\n", mailURL)

	mailPitReq2, err := http.NewRequest("GET", mailURL, nil)
	Expect(err).ShouldNot(HaveOccurred())
	mailPitReq2.Header.Add("Content-Type", "application/json")

	mailPitResp2, err := http.DefaultClient.Do(mailPitReq2)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(mailPitResp2.StatusCode).Should(Equal(http.StatusOK))

	body, err := io.ReadAll(mailPitResp2.Body)
	Expect(err).ShouldNot(HaveOccurred())

	// Extracting the token from the mail body
	re := regexp.MustCompile(`Token:\s*([a-zA-Z0-9]+)\s*`)

	tokens := re.FindAllStringSubmatch(string(body), -1)
	Expect(len(tokens)).Should(BeNumerically(">=", 1))

	tfaCode := tokens[0][1] // The token is captured in the first group
	fmt.Fprintf(GinkgoWriter, "TFACode: %s\n", tfaCode)
	fmt.Fprintf(GinkgoWriter, "TGT: %s\n", signinResp.Token)

	// TFA with the two tokens
	tfaReqBody, err := json.Marshal(
		employer.EmployerTFARequest{
			TFAToken: signinResp.Token,
			TFACode:  tfaCode,
		},
	)
	Expect(err).ShouldNot(HaveOccurred())

	tfaReq, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/tfa",
		bytes.NewBuffer(tfaReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())
	tfaReq.Header.Set("Content-Type", "application/json")

	tfaResp, err := http.DefaultClient.Do(tfaReq)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(tfaResp.StatusCode).Should(Equal(http.StatusOK))

	var tfaRespObj employer.EmployerTFAResponse
	err = json.NewDecoder(tfaResp.Body).Decode(&tfaRespObj)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(tfaRespObj.SessionToken).ShouldNot(BeEmpty())

	// Delete the email from mailpit so that we can run the test multiple times
	mailPitDeleteReqBody, err := json.Marshal(MailPitDeleteRequest{
		IDs: []string{messageID},
	})
	Expect(err).ShouldNot(HaveOccurred())

	mailPitReq3, err := http.NewRequest(
		"DELETE",
		"http://localhost:8025/api/v1/messages",
		bytes.NewBuffer(mailPitDeleteReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())
	mailPitReq3.Header.Set("Accept", "application/json")
	mailPitReq3.Header.Add("Content-Type", "application/json")

	mailPitDeleteResp, err := http.DefaultClient.Do(mailPitReq3)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(mailPitDeleteResp.StatusCode).Should(Equal(http.StatusOK))

	return tfaRespObj.SessionToken, nil
}

// testPOST performs a POST request to the given endpoint with the given request
// body, token and expects the given status code. The response body is not read.
func testPOST(
	token string,
	reqBody interface{},
	endpoint string,
	wantStatus int,
) {
	doPOST(token, reqBody, endpoint, wantStatus, false)
}

func testPOSTGetResp(
	token string,
	reqBody interface{},
	endpoint string,
	wantStatus int,
) interface{} {
	return doPOST(token, reqBody, endpoint, wantStatus, true)
}

func doPOST(
	token string,
	reqBody interface{},
	endpoint string,
	wantStatus int,
	wantResp bool,
) interface{} {
	body, err := json.Marshal(reqBody)
	Expect(err).ShouldNot(HaveOccurred())

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+endpoint,
		bytes.NewBuffer(body),
	)
	Expect(err).ShouldNot(HaveOccurred())

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())

	if resp.StatusCode != wantStatus {
		debug.PrintStack()

		respBody, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Fprintf(GinkgoWriter, "Response Body: %s\n", string(respBody))
		}

		Fail(
			fmt.Sprintf(
				"Expected status %d, got %d",
				wantStatus,
				resp.StatusCode,
			),
		)
	}

	if !wantResp {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	Expect(err).ShouldNot(HaveOccurred())

	return respBody
}

// getTFACode retrieves the TFA code from the email sent to the specified address
func getTFACode(email string) (string, string) {
	baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
	Expect(err).ShouldNot(HaveOccurred())
	query := url.Values{}
	query.Add(
		"query",
		fmt.Sprintf(
			"to:%s subject:Vetchium Two Factor Authentication",
			email,
		),
	)
	baseURL.RawQuery = query.Encode()
	mailURL := baseURL.String()

	var messageID string
	for i := 0; i < 3; i++ {
		<-time.After(10 * time.Second)
		mailPitResp, err := http.Get(mailURL)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mailPitResp.StatusCode).Should(Equal(http.StatusOK))

		body, err := io.ReadAll(mailPitResp.Body)
		Expect(err).ShouldNot(HaveOccurred())

		var mailPitRespObj MailPitResponse
		err = json.Unmarshal(body, &mailPitRespObj)
		Expect(err).ShouldNot(HaveOccurred())

		if len(mailPitRespObj.Messages) > 0 {
			messageID = mailPitRespObj.Messages[0].ID
			break
		}
	}
	Expect(messageID).ShouldNot(BeEmpty())

	mailResp, err := http.Get(
		mailPitURL + "/api/v1/message/" + messageID,
	)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

	body, err := io.ReadAll(mailResp.Body)
	Expect(err).ShouldNot(HaveOccurred())

	re := regexp.MustCompile(
		`Your Two Factor authentication code is:\s*([0-9]+)`,
	)
	matches := re.FindStringSubmatch(string(body))
	Expect(len(matches)).Should(BeNumerically(">=", 2))

	return matches[1], messageID
}

// getSessionToken completes the TFA flow and returns a session token
func getSessionToken(tfaToken, tfaCode string, rememberMe bool) string {
	tfaReqBody, err := json.Marshal(hub.HubTFARequest{
		TFAToken:   tfaToken,
		TFACode:    tfaCode,
		RememberMe: rememberMe,
	})
	Expect(err).ShouldNot(HaveOccurred())

	tfaResp, err := http.Post(
		serverURL+"/hub/tfa",
		"application/json",
		bytes.NewBuffer(tfaReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(tfaResp.StatusCode).Should(Equal(http.StatusOK))

	var tfaRespObj hub.HubTFAResponse
	err = json.NewDecoder(tfaResp.Body).Decode(&tfaRespObj)
	Expect(err).ShouldNot(HaveOccurred())

	return tfaRespObj.SessionToken
}

// cleanupEmail deletes the email with the given messageID from mailpit
func cleanupEmail(messageID string) {
	deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
		IDs: []string{messageID},
	})
	Expect(err).ShouldNot(HaveOccurred())

	req, err := http.NewRequest(
		"DELETE",
		mailPitURL+"/api/v1/messages",
		bytes.NewBuffer(deleteReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	deleteResp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(deleteResp.StatusCode).Should(Equal(http.StatusOK))
}

// Add this function after employerSigninAsync
func hubSignin(email, password string) string {
	loginReqBody, err := json.Marshal(hub.LoginRequest{
		Email:    common.EmailAddress(email),
		Password: common.Password(password),
	})
	Expect(err).ShouldNot(HaveOccurred())

	loginResp, err := http.Post(
		serverURL+"/hub/login",
		"application/json",
		bytes.NewBuffer(loginReqBody),
	)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))

	var loginRespObj hub.LoginResponse
	err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
	Expect(err).ShouldNot(HaveOccurred())

	tfaCode, messageID := getTFACode(email)
	defer cleanupEmail(messageID)

	return getSessionToken(loginRespObj.Token, tfaCode, false)
}

// Add this function after hubSignin
func hubSigninAsync(
	email, password string,
	token *string,
	wg *sync.WaitGroup,
) {
	go func(email, password string, token *string, wg *sync.WaitGroup) {
		defer GinkgoRecover()
		defer wg.Done()
		*token = hubSignin(email, password)
		fmt.Fprintf(
			GinkgoWriter,
			"Hub user email: %s, password: %s, gotToken: %s\n",
			email,
			password,
			*token,
		)
	}(email, password, token, wg)
}

func strptr(s string) *string {
	return &s
}
