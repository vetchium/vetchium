package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("Hub Signup", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0027-hubsignup-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0027-hubsignup-down.pgsql")
		db.Close()
	})

	Describe("/hub/signup endpoint", func() {
		type signupHubUserTestCase struct {
			description string
			request     hub.SignupHubUserRequest
			wantStatus  int
		}

		It("should handle successful signup", func() {
			testCase := signupHubUserTestCase{
				description: "valid email with approved domain",
				request: hub.SignupHubUserRequest{
					Email: common.EmailAddress("new@0027-example.com"),
				},
				wantStatus: http.StatusOK,
			}

			testPOST("", testCase.request, "/hub/signup", testCase.wantStatus)

			// Verify that an invite email was sent using mailpit
			baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
			Expect(err).ShouldNot(HaveOccurred())
			query := url.Values{}
			query.Add(
				"query",
				fmt.Sprintf(
					"to:%s subject:\"Vetchium user signup invite\"",
					"new@0027-example.com",
				),
			)
			baseURL.RawQuery = query.Encode()
			mailURL := baseURL.String()

			var messageID string
			var foundEmail bool
			// Try a few times with delay to allow email to be delivered
			for i := 0; i < 3; i++ {
				<-time.After(2 * time.Second)
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
					foundEmail = true
					break
				}
			}
			Expect(
				foundEmail,
			).Should(BeTrue(), "Should have found an invite email")

			// Verify the email content
			mailResp, err := http.Get(
				mailPitURL + "/api/v1/message/" + messageID,
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mailResp.StatusCode).Should(Equal(http.StatusOK))

			emailBody, err := io.ReadAll(mailResp.Body)
			Expect(err).ShouldNot(HaveOccurred())
			emailContent := string(emailBody)

			// Verify the email contains the signup link
			Expect(
				emailContent,
			).Should(ContainSubstring("https://vetchium.com/hub/signup?token="))

			// Clean up the email after verification
			cleanupEmail(messageID)
		})

		It("should handle another successful signup", func() {
			testCase := signupHubUserTestCase{
				description: "another valid email with approved domain",
				request: hub.SignupHubUserRequest{
					Email: common.EmailAddress("another@0027-example.com"),
				},
				wantStatus: http.StatusOK,
			}

			testPOST("", testCase.request, "/hub/signup", testCase.wantStatus)

			// Verify that an invite email was sent using mailpit
			baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
			Expect(err).ShouldNot(HaveOccurred())
			query := url.Values{}
			query.Add(
				"query",
				fmt.Sprintf(
					"to:%s subject:\"Vetchium user signup invite\"",
					"another@0027-example.com",
				),
			)
			baseURL.RawQuery = query.Encode()
			mailURL := baseURL.String()

			var messageID string
			var foundEmail bool
			// Try a few times with delay to allow email to be delivered
			for i := 0; i < 3; i++ {
				<-time.After(2 * time.Second)
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
					foundEmail = true
					break
				}
			}
			Expect(
				foundEmail,
			).Should(BeTrue(), "Should have found an invite email")

			// Clean up the email after verification
			cleanupEmail(messageID)
		})

		It("should handle multiple test cases correctly", func() {
			testCases := []signupHubUserTestCase{
				{
					description: "existing user",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress("existing@0027-example.com"),
					},
					wantStatus: 461, // Custom status code for already invited/member
				},
				{
					description: "already invited user",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress("invited@0027-example.com"),
					},
					wantStatus: 461, // Custom status code for already invited/member
				},
				{
					description: "unapproved domain",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress(
							"invalid@truly-unapproved-0027-example.com",
						),
					},
					wantStatus: 460, // Custom status code for unsupported domain
				},
				{
					description: "invalid email format",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress("invalid-email"),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "empty email",
					request: hub.SignupHubUserRequest{
						Email: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "email with more than one @",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress(
							"invalid@invalid@invalid.com",
						),
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "email with + symbol",
					request: hub.SignupHubUserRequest{
						Email: common.EmailAddress(
							"invalid+symbol@invalid.com",
						),
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				By(tc.description)
				testPOST("", tc.request, "/hub/signup", tc.wantStatus)
			}
		})

		It("should handle malformed JSON", func() {
			// Send malformed JSON to test error handling
			malformedJSON := []byte(`{"email": "malformed@0027-example.com`)

			req, err := http.NewRequest(
				http.MethodPost,
				serverURL+"/hub/signup",
				bytes.NewBuffer(malformedJSON),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
	})
})
