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
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/api/pkg/vetchi"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	serverURL  = "http://localhost:8081"
	mailPitURL = "http://localhost:8025"
)

func setupTestDB() *pgxpool.Pool {
	connStr := "host=localhost port=5432 user=user dbname=vdb password=pass sslmode=disable"
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
	go func() {
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
	}()
}

// Returns the session token for the employer with the given credentials
func employerSignin(clientID, email, password string) (string, error) {
	signinReqBody, err := json.Marshal(
		vetchi.EmployerSignInRequest{
			ClientID: clientID,
			Email:    vetchi.EmailAddress(email),
			Password: vetchi.Password(password),
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

	var signinResp vetchi.EmployerSignInResponse
	err = json.NewDecoder(resp.Body).Decode(&signinResp)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(signinResp.Token).ShouldNot(BeEmpty())

	// Get the tfa code from the email by querying mailpit
	fmt.Fprintf(GinkgoWriter, "Sleeping to allow granger to email\n")
	<-time.After(15 * time.Second)
	fmt.Fprintf(GinkgoWriter, "Wokeup\n")

	baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
	Expect(err).ShouldNot(HaveOccurred())
	query := url.Values{}
	query.Add(
		"query",
		fmt.Sprintf("to:%s subject:Vetchi Two Factor Authentication", email),
	)
	baseURL.RawQuery = query.Encode()

	url1 := "http://localhost:8025/api/v1/search?query=to%3Aadmin%40domain-onboarded.example%20subject%3AVetchi%20Two%20Factor%20Authentication"
	url2 := baseURL.String()
	fmt.Fprintf(GinkgoWriter, "URL1: %s\n", url1)
	fmt.Fprintf(GinkgoWriter, "URL2: %s\n", url2)

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
	Expect(len(mailPitResp1Obj.Messages)).Should(Equal(1))

	messageID := mailPitResp1Obj.Messages[0].ID
	mailURL := "http://localhost:8025/api/v1/message/" + messageID
	fmt.Fprintf(GinkgoWriter, "Mail URL: %s\n", mailURL)

	mailPitReq2, err := http.NewRequest("GET", mailURL, nil)
	Expect(err).ShouldNot(HaveOccurred())
	mailPitReq2.Header.Add("Content-Type", "application/json")

	mailPitResp2, err := http.DefaultClient.Do(mailPitReq2)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(mailPitResp2.StatusCode).Should(Equal(http.StatusOK))

	body, err = io.ReadAll(mailPitResp2.Body)
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
		vetchi.EmployerTFARequest{
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

	var tfaRespObj vetchi.EmployerTFAResponse
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
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(resp.StatusCode).Should(Equal(wantStatus))

	if !wantResp {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	Expect(err).ShouldNot(HaveOccurred())

	return respBody
}
