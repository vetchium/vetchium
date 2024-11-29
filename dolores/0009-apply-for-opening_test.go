package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Apply For Opening", Ordered, func() {
	var db *pgxpool.Pool
	var activeHubUserToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0009-apply-for-opening-up.pgsql")

		// Login as active hub user
		loginReqBody, err := json.Marshal(vetchi.LoginRequest{
			Email:    "active@applyopening.example",
			Password: "NewPassword123$",
		})
		Expect(err).ShouldNot(HaveOccurred())

		loginResp, err := http.Post(
			serverURL+"/hub/login",
			"application/json",
			bytes.NewBuffer(loginReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))

		var loginRespObj vetchi.LoginResponse
		err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
		Expect(err).ShouldNot(HaveOccurred())
		tfaToken := loginRespObj.Token

		// Get TFA code from email
		baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
		Expect(err).ShouldNot(HaveOccurred())
		query := url.Values{}
		query.Add(
			"query",
			"to:active@applyopening.example subject:Vetchi Two Factor Authentication",
		)
		baseURL.RawQuery = query.Encode()

		var messageID string
		for i := 0; i < 3; i++ {
			<-time.After(10 * time.Second)
			mailPitResp, err := http.Get(baseURL.String())
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

		// Get the email content
		mailResp, err := http.Get(mailPitURL + "/api/v1/message/" + messageID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

		body, err := io.ReadAll(mailResp.Body)
		Expect(err).ShouldNot(HaveOccurred())

		re := regexp.MustCompile(
			`Your Two Factor authentication code is:\s*([0-9]+)`,
		)
		matches := re.FindStringSubmatch(string(body))
		Expect(len(matches)).Should(BeNumerically(">=", 2))
		tfaCode := matches[1]

		// Clean up the email
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

		// Complete TFA flow
		tfaReqBody, err := json.Marshal(vetchi.HubTFARequest{
			TFAToken:   tfaToken,
			TFACode:    tfaCode,
			RememberMe: false,
		})
		Expect(err).ShouldNot(HaveOccurred())

		tfaResp, err := http.Post(
			serverURL+"/hub/tfa",
			"application/json",
			bytes.NewBuffer(tfaReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tfaResp.StatusCode).Should(Equal(http.StatusOK))

		var tfaRespObj vetchi.HubTFAResponse
		err = json.NewDecoder(tfaResp.Body).Decode(&tfaRespObj)
		Expect(err).ShouldNot(HaveOccurred())
		activeHubUserToken = tfaRespObj.SessionToken
	})

	AfterAll(func() {
		seedDatabase(db, "0009-apply-for-opening-down.pgsql")
		db.Close()
	})

	Describe("Apply For Opening", func() {
		type applyForOpeningTestCase struct {
			description string
			token       string
			request     vetchi.ApplyForOpeningRequest
			wantStatus  int
		}

		It("should handle application requests correctly", func() {
			testCases := []applyForOpeningTestCase{
				{
					description: "valid application",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "applyopening.example",
						Resume:                 "base64encodedresume",
						CoverLetter:            "I am interested in this position",
						Filename:               "resume.pdf",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "without auth token",
					token:       "",
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "applyopening.example",
						Resume:                 "base64encodedresume",
						Filename:               "resume.pdf",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid opening ID",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "invalid-id",
						CompanyDomain:          "applyopening.example",
						Resume:                 "base64encodedresume",
						Filename:               "resume.pdf",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "with invalid company domain",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "invalid.example",
						Resume:                 "base64encodedresume",
						Filename:               "resume.pdf",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "with empty resume",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "applyopening.example",
						Resume:                 "",
						Filename:               "resume.pdf",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with empty filename",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "applyopening.example",
						Resume:                 "base64encodedresume",
						Filename:               "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with too long cover letter",
					token:       activeHubUserToken,
					request: vetchi.ApplyForOpeningRequest{
						OpeningIDWithinCompany: "2024-Mar-09-001",
						CompanyDomain:          "applyopening.example",
						Resume:                 "base64encodedresume",
						CoverLetter: string(
							make([]byte, 4097),
						), // Max is 4096
						Filename: "resume.pdf",
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				reqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/hub/apply-for-opening",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())

				if tc.token != "" {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if resp.StatusCode == http.StatusBadRequest {
					var validationErrors vetchi.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(validationErrors.Errors).ShouldNot(BeEmpty())
				}
			}
		})
	})
})
