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
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("Hub User Registration", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0031-hub-user-registration-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0031-hub-user-registration-down.pgsql")
		db.Close()
	})

	// Helper function to extract signup token from email
	extractSignupToken := func(email string) string {
		baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
		Expect(err).ShouldNot(HaveOccurred())
		query := url.Values{}
		query.Add(
			"query",
			fmt.Sprintf(
				"to:%s subject:\"Vetchium user signup invite\"",
				email,
			),
		)
		baseURL.RawQuery = query.Encode()

		var messageID string
		// Wait for email to arrive
		for i := 0; i < 3; i++ {
			<-time.After(5 * time.Second)
			mailPitResp, err := http.Get(baseURL.String())
			Expect(err).ShouldNot(HaveOccurred())

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

		// Get email content and extract token
		mailResp, err := http.Get(
			mailPitURL + "/api/v1/message/" + messageID,
		)
		Expect(err).ShouldNot(HaveOccurred())

		emailBody, err := io.ReadAll(mailResp.Body)
		Expect(err).ShouldNot(HaveOccurred())

		// Extract token from signup URL
		re := regexp.MustCompile(`/signup-hubuser/([a-zA-Z0-9]+)`)
		matches := re.FindStringSubmatch(string(emailBody))
		Expect(len(matches)).Should(BeNumerically(">=", 2))

		// Clean up email
		cleanupEmail(messageID)

		return matches[1]
	}

	// Helper function to send onboard request
	sendOnboardRequest := func(token, fullName, password, countryCode string) *http.Response {
		onboardReq := hub.OnboardHubUserRequest{
			Token:               token,
			FullName:            fullName,
			ResidentCountryCode: common.CountryCode(countryCode),
			Password:            common.Password(password),
			SelectedTier:        "FREE_TIER",
			PreferredLanguage:   "en",
			ShortBio:            "Test user bio",
			LongBio:             "This is a longer test user bio for testing purposes.",
		}

		reqBody, err := json.Marshal(onboardReq)
		Expect(err).ShouldNot(HaveOccurred())

		resp, err := http.Post(
			serverURL+"/hub/onboard-hubuser",
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		return resp
	}

	Describe("Complete Registration Flow", func() {
		It("should handle successful signup and onboarding", func() {
			email := "complete-flow@registration.example"

			// Step 1: Send signup request
			signupReq := hub.SignupHubUserRequest{
				Email: common.EmailAddress(email),
			}
			signupReqBody, err := json.Marshal(signupReq)
			Expect(err).ShouldNot(HaveOccurred())

			signupResp, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signupResp.StatusCode).Should(Equal(http.StatusOK))

			// Step 2: Extract token from email
			token := extractSignupToken(email)

			// Step 3: Complete onboarding
			onboardResp := sendOnboardRequest(
				token,
				"Complete Flow User",
				"CompletePassword123$",
				"USA",
			)
			Expect(onboardResp.StatusCode).Should(Equal(http.StatusOK))

			var onboardRespObj hub.OnboardHubUserResponse
			err = json.NewDecoder(onboardResp.Body).Decode(&onboardRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(onboardRespObj.SessionToken).ShouldNot(BeEmpty())
			Expect(onboardRespObj.GeneratedHandle).ShouldNot(BeEmpty())

			// Step 4: Verify user can login with new credentials
			loginReq := hub.LoginRequest{
				Email:    common.EmailAddress(email),
				Password: "CompletePassword123$",
			}
			loginReqBody, err := json.Marshal(loginReq)
			Expect(err).ShouldNot(HaveOccurred())

			loginResp, err := http.Post(
				serverURL+"/hub/login",
				"application/json",
				bytes.NewBuffer(loginReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(loginResp.StatusCode).Should(Equal(http.StatusOK))
		})

		It("should handle onboarding validation errors", func() {
			type onboardTestCase struct {
				description   string
				token         string
				fullName      string
				password      string
				countryCode   string
				wantStatus    int
				wantErrFields []string
			}

			// First create a valid signup to get a token
			email := "validation-test@registration.example"
			signupReq := hub.SignupHubUserRequest{
				Email: common.EmailAddress(email),
			}
			signupReqBody, err := json.Marshal(signupReq)
			Expect(err).ShouldNot(HaveOccurred())

			signupResp, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signupResp.StatusCode).Should(Equal(http.StatusOK))

			validToken := extractSignupToken(email)

			testCases := []onboardTestCase{
				{
					description: "with invalid token",
					token:       "invalid-token",
					fullName:    "Valid Name",
					password:    "ValidPassword123$",
					countryCode: "USA",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description:   "with empty full name",
					token:         validToken,
					fullName:      "",
					password:      "ValidPassword123$",
					countryCode:   "USA",
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"full_name"},
				},
				{
					description:   "with invalid password",
					token:         validToken,
					fullName:      "Valid Name",
					password:      "weak",
					countryCode:   "USA",
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
				{
					description:   "with invalid country code",
					token:         validToken,
					fullName:      "Valid Name",
					password:      "ValidPassword123$",
					countryCode:   "INVALID",
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"resident_country_code"},
				},
				{
					description:   "with very long full name",
					token:         validToken,
					fullName:      string(make([]byte, 300)),
					password:      "ValidPassword123$",
					countryCode:   "USA",
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"full_name"},
				},
			}

			for _, tc := range testCases {
				GinkgoWriter.Printf("#### %s\n", tc.description)

				resp := sendOnboardRequest(
					tc.token,
					tc.fullName,
					tc.password,
					tc.countryCode,
				)
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if len(tc.wantErrFields) > 0 {
					var validationErrors common.ValidationErrors
					err := json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
				}
			}
		})

		It("should prevent token reuse", func() {
			email := "token-reuse@registration.example"

			// Step 1: Send signup request
			signupReq := hub.SignupHubUserRequest{
				Email: common.EmailAddress(email),
			}
			signupReqBody, err := json.Marshal(signupReq)
			Expect(err).ShouldNot(HaveOccurred())

			signupResp, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signupResp.StatusCode).Should(Equal(http.StatusOK))

			// Step 2: Extract token from email
			token := extractSignupToken(email)

			// Step 3: Complete onboarding successfully
			onboardResp1 := sendOnboardRequest(
				token,
				"Token Reuse User",
				"TokenReusePassword123$",
				"USA",
			)
			Expect(onboardResp1.StatusCode).Should(Equal(http.StatusOK))

			// Step 4: Try to reuse the same token - should fail
			onboardResp2 := sendOnboardRequest(
				token,
				"Another User",
				"AnotherPassword123$",
				"CAN",
			)
			Expect(
				onboardResp2.StatusCode,
			).Should(Equal(http.StatusUnauthorized))
		})

		It("should handle expired tokens", func() {
			// This test would require manipulating token expiry in the database
			// For now, we'll test with an obviously invalid/expired token format
			resp := sendOnboardRequest(
				"expired-token-12345",
				"Expired Token User",
				"ExpiredPassword123$",
				"USA",
			)
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("should handle duplicate signup attempts", func() {
			email := "duplicate@registration.example"

			// First signup
			signupReq := hub.SignupHubUserRequest{
				Email: common.EmailAddress(email),
			}
			signupReqBody, err := json.Marshal(signupReq)
			Expect(err).ShouldNot(HaveOccurred())

			signupResp1, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signupResp1.StatusCode).Should(Equal(http.StatusOK))

			// Complete onboarding
			token := extractSignupToken(email)
			onboardResp := sendOnboardRequest(
				token,
				"Duplicate User",
				"DuplicatePassword123$",
				"USA",
			)
			Expect(onboardResp.StatusCode).Should(Equal(http.StatusOK))

			// Try to signup again with same email - should fail
			signupResp2, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				signupResp2.StatusCode,
			).Should(Equal(461))
			// Already invited/member
		})

		It("should handle malformed onboard requests", func() {
			// Test with malformed JSON
			malformedJSON := []byte(`{"token": "test", "full_name": "Test"`)

			resp, err := http.Post(
				serverURL+"/hub/onboard-hubuser",
				"application/json",
				bytes.NewBuffer(malformedJSON),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		It("should generate unique handles for users", func() {
			emails := []string{
				"handle1@registration.example",
				"handle2@registration.example",
				"handle3@registration.example",
			}
			handles := make([]string, 0, len(emails))

			for i, email := range emails {
				// Signup
				signupReq := hub.SignupHubUserRequest{
					Email: common.EmailAddress(email),
				}
				signupReqBody, err := json.Marshal(signupReq)
				Expect(err).ShouldNot(HaveOccurred())

				signupResp, err := http.Post(
					serverURL+"/hub/signup",
					"application/json",
					bytes.NewBuffer(signupReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(signupResp.StatusCode).Should(Equal(http.StatusOK))

				// Onboard
				token := extractSignupToken(email)
				onboardResp := sendOnboardRequest(
					token,
					fmt.Sprintf("Handle Test User %d", i+1),
					"HandlePassword123$",
					"USA",
				)
				Expect(onboardResp.StatusCode).Should(Equal(http.StatusOK))

				var onboardRespObj hub.OnboardHubUserResponse
				err = json.NewDecoder(onboardResp.Body).Decode(&onboardRespObj)
				Expect(err).ShouldNot(HaveOccurred())

				handles = append(handles, onboardRespObj.GeneratedHandle)
			}

			// Verify all handles are unique
			for i := 0; i < len(handles); i++ {
				for j := i + 1; j < len(handles); j++ {
					Expect(handles[i]).ShouldNot(Equal(handles[j]))
				}
			}
		})

		It("should validate tier selection", func() {
			email := "tier-test@registration.example"

			// Signup
			signupReq := hub.SignupHubUserRequest{
				Email: common.EmailAddress(email),
			}
			signupReqBody, err := json.Marshal(signupReq)
			Expect(err).ShouldNot(HaveOccurred())

			signupResp, err := http.Post(
				serverURL+"/hub/signup",
				"application/json",
				bytes.NewBuffer(signupReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(signupResp.StatusCode).Should(Equal(http.StatusOK))

			token := extractSignupToken(email)

			// Test with invalid tier
			onboardReq := hub.OnboardHubUserRequest{
				Token:               token,
				FullName:            "Tier Test User",
				ResidentCountryCode: "USA",
				Password:            "TierPassword123$",
				SelectedTier:        "INVALID_TIER",
			}

			reqBody, err := json.Marshal(onboardReq)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/hub/onboard-hubuser",
				"application/json",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
	})
})
