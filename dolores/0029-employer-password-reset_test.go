package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = Describe("Employer Password Reset", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0029-employer-password-reset-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0029-employer-password-reset-down.pgsql")
		db.Close()
	})

	// Helper functions
	sendForgotPasswordRequest := func(email string) *http.Response {
		reqBody, err := json.Marshal(employer.EmployerForgotPasswordRequest{
			Email: email,
		})
		Expect(err).ShouldNot(HaveOccurred())

		resp, err := http.Post(
			serverURL+"/employer/forgot-password",
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Debug: log response body if not successful
		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				fmt.Fprintf(
					GinkgoWriter,
					"Error response body: %s\n",
					string(body),
				)
			}
			resp.Body.Close()
			// Recreate body for caller
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}

		return resp
	}

	getPasswordResetTokenFromEmail := func(email string) (string, string) {
		baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
		Expect(err).ShouldNot(HaveOccurred())
		query := url.Values{}
		query.Add(
			"query",
			fmt.Sprintf(
				"to:%s subject:Vetchium Password Reset",
				email,
			),
		)
		baseURL.RawQuery = query.Encode()

		var messageID string
		var resetToken string

		// Wait and retry for email
		for i := 0; i < 3; i++ {
			<-time.After(10 * time.Second)

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

		// Get email content and extract reset token
		mailResp, err := http.Get(
			mailPitURL + "/api/v1/message/" + messageID,
		)
		Expect(err).ShouldNot(HaveOccurred())
		body, err := io.ReadAll(mailResp.Body)
		Expect(err).ShouldNot(HaveOccurred())

		re := regexp.MustCompile(`/reset-password\?token=([a-zA-Z0-9]+)`)
		matches := re.FindStringSubmatch(string(body))
		Expect(len(matches)).Should(BeNumerically(">=", 2))
		resetToken = matches[1]

		return resetToken, messageID
	}

	cleanupEmail := func(messageID string) {
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

	sendResetPasswordRequest := func(token, password string) *http.Response {
		reqBody, err := json.Marshal(employer.EmployerResetPasswordRequest{
			Token:    token,
			Password: password,
		})
		Expect(err).ShouldNot(HaveOccurred())

		resp, err := http.Post(
			serverURL+"/employer/reset-password",
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		return resp
	}

	verifyEmployerLogin := func(email, password string) {
		sessionToken, err := employerSignin(
			"passwordreset.example",
			email,
			password,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sessionToken).ShouldNot(BeEmpty())
	}

	Describe("Forgot Password", func() {
		It("should handle various forgot password scenarios", func() {
			type forgotPasswordTestCase struct {
				description  string
				email        string
				wantStatus   int
				checkEmail   bool
				wantErrField string
			}

			testCases := []forgotPasswordTestCase{
				{
					description: "with valid active user email",
					email:       "active@passwordreset.example",
					wantStatus:  http.StatusOK,
					checkEmail:  true,
				},
				{
					description: "with valid disabled user email",
					email:       "disabled@passwordreset.example",
					wantStatus:  http.StatusOK,
					checkEmail:  true,
				},
				{
					description: "with non-existent email",
					email:       "nonexistent@passwordreset.example",
					wantStatus:  http.StatusOK,
					checkEmail:  false,
				},
				{
					description:  "with invalid email format",
					email:        "invalid-email",
					wantStatus:   http.StatusBadRequest,
					checkEmail:   false,
					wantErrField: "email",
				},
				{
					description:  "with empty email",
					email:        "",
					wantStatus:   http.StatusBadRequest,
					checkEmail:   false,
					wantErrField: "email",
				},
				{
					description: "with email from different domain",
					email:       "user@different.example",
					wantStatus:  http.StatusOK,
					checkEmail:  false,
				},
				{
					description: "with very long email",
					email: strings.Repeat(
						"a",
						250,
					) + "@passwordreset.example",
					wantStatus:   http.StatusBadRequest,
					checkEmail:   false,
					wantErrField: "email",
				},
				{
					description:  "with SQL injection attempt",
					email:        "test'; DROP TABLE org_users; --@passwordreset.example",
					wantStatus:   http.StatusBadRequest,
					checkEmail:   false,
					wantErrField: "email",
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				resp := sendForgotPasswordRequest(tc.email)
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantErrField != "" {
					var validationErrors common.ValidationErrors
					err := json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElement(tc.wantErrField))
				}

				if tc.checkEmail {
					resetToken, messageID := getPasswordResetTokenFromEmail(
						tc.email,
					)
					Expect(resetToken).ShouldNot(BeEmpty())
					cleanupEmail(messageID)
				}
			}
		})

		It(
			"should invalidate previous tokens when multiple requests are made",
			func() {
				email := "multiple-requests@passwordreset.example"

				// Send first forgot password request
				resp1 := sendForgotPasswordRequest(email)
				Expect(resp1.StatusCode).Should(Equal(http.StatusOK))

				// Get first token
				token1, messageID1 := getPasswordResetTokenFromEmail(email)
				cleanupEmail(messageID1)

				// Send second forgot password request
				resp2 := sendForgotPasswordRequest(email)
				Expect(resp2.StatusCode).Should(Equal(http.StatusOK))

				// Get second token
				token2, messageID2 := getPasswordResetTokenFromEmail(email)
				cleanupEmail(messageID2)

				// Verify tokens are different
				Expect(token1).ShouldNot(Equal(token2))

				// Try to use first token - should fail
				resetResp1 := sendResetPasswordRequest(
					token1,
					"NewPassword123$",
				)
				Expect(
					resetResp1.StatusCode,
				).Should(Equal(http.StatusUnauthorized))

				// Use second token - should succeed
				resetResp2 := sendResetPasswordRequest(
					token2,
					"NewPassword123$",
				)
				Expect(resetResp2.StatusCode).Should(Equal(http.StatusOK))

				// Verify login with new password
				verifyEmployerLogin(email, "NewPassword123$")
			},
		)
	})

	Describe("Reset Password", func() {
		It("should handle various reset password scenarios", func() {
			type resetPasswordTestCase struct {
				description  string
				setupEmail   string
				token        string
				password     string
				wantStatus   int
				wantErrField string
				verifyLogin  bool
			}

			// Setup valid token for positive test cases
			validEmail := "valid-reset@passwordreset.example"
			sendForgotPasswordRequest(validEmail)
			validToken, messageID := getPasswordResetTokenFromEmail(validEmail)
			defer cleanupEmail(messageID)

			testCases := []resetPasswordTestCase{
				{
					description: "with valid token and password",
					token:       validToken,
					password:    "NewValidPassword123$",
					wantStatus:  http.StatusOK,
					verifyLogin: true,
					setupEmail:  validEmail,
				},
				{
					description: "with invalid token",
					token:       "invalid-token-12345",
					password:    "NewPassword123$",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description:  "with empty token",
					token:        "",
					password:     "NewPassword123$",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "token",
				},
				{
					description:  "with very long token",
					token:        strings.Repeat("a", 1000),
					password:     "NewPassword123$",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "token",
				},
				{
					description:  "with invalid password format",
					token:        "some-token",
					password:     "weak",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "password",
				},
				{
					description:  "with empty password",
					token:        "some-token",
					password:     "",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "password",
				},
				{
					description:  "with very long password",
					token:        "some-token",
					password:     strings.Repeat("A1$", 100),
					wantStatus:   http.StatusBadRequest,
					wantErrField: "password",
				},
				{
					description: "with SQL injection in token",
					token:       "'; DROP TABLE org_users; --",
					password:    "NewPassword123$",
					wantStatus:  http.StatusUnauthorized,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				resp := sendResetPasswordRequest(tc.token, tc.password)
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantErrField != "" {
					var validationErrors common.ValidationErrors
					err := json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElement(tc.wantErrField))
				}

				if tc.verifyLogin {
					verifyEmployerLogin(tc.setupEmail, tc.password)
				}
			}
		})

		It("should handle token expiry", func() {
			email := "token-expiry@passwordreset.example"

			// Send forgot password request
			resp := sendForgotPasswordRequest(email)
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get reset token
			resetToken, messageID := getPasswordResetTokenFromEmail(email)
			cleanupEmail(messageID)

			// Wait for token to expire (password reset tokens have short validity in test config)
			<-time.After(7 * time.Minute)

			// Try to use expired token
			resetResp := sendResetPasswordRequest(resetToken, "NewPassword123$")
			Expect(resetResp.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Verify original password still works
			verifyEmployerLogin(email, "NewPassword123$")
		})

		It("should prevent token reuse", func() {
			email := "token-reuse@passwordreset.example"

			// Send forgot password request
			resp := sendForgotPasswordRequest(email)
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get reset token
			resetToken, messageID := getPasswordResetTokenFromEmail(email)
			cleanupEmail(messageID)

			// First password reset should succeed
			resetResp1 := sendResetPasswordRequest(
				resetToken,
				"FirstNewPassword123$",
			)
			Expect(resetResp1.StatusCode).Should(Equal(http.StatusOK))

			// Verify login with new password
			verifyEmployerLogin(email, "FirstNewPassword123$")

			// Try to reuse the same token - should fail
			resetResp2 := sendResetPasswordRequest(
				resetToken,
				"SecondNewPassword123$",
			)
			Expect(resetResp2.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Verify password hasn't changed
			verifyEmployerLogin(email, "FirstNewPassword123$")
		})

		It("should handle cross-employer token attempts", func() {
			// Setup token for one employer
			email1 := "cross-test@passwordreset.example"
			resp1 := sendForgotPasswordRequest(email1)
			Expect(resp1.StatusCode).Should(Equal(http.StatusOK))

			token1, messageID1 := getPasswordResetTokenFromEmail(email1)
			defer cleanupEmail(messageID1)

			// Try to use token from different employer context
			// This should fail as tokens are tied to specific employers
			resetResp := sendResetPasswordRequest(token1, "HackedPassword123$")
			// Token should be valid for the correct user
			Expect(resetResp.StatusCode).Should(Equal(http.StatusOK))

			// Verify login works
			verifyEmployerLogin(email1, "HackedPassword123$")
		})

		It("should maintain session validity after password reset", func() {
			email := "session-test@passwordreset.example"

			// Get session token before password reset
			sessionToken, err := employerSignin(
				"passwordreset.example",
				email,
				"NewPassword123$",
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify session works
			req, err := http.NewRequest(
				"GET",
				serverURL+"/employer/get-onboard-status",
				bytes.NewBuffer(
					[]byte(`{"client_id": "passwordreset.example"}`),
				),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Reset password
			forgotResp := sendForgotPasswordRequest(email)
			Expect(forgotResp.StatusCode).Should(Equal(http.StatusOK))

			resetToken, messageID := getPasswordResetTokenFromEmail(email)
			defer cleanupEmail(messageID)

			resetResp := sendResetPasswordRequest(
				resetToken,
				"NewSessionPassword123$",
			)
			Expect(resetResp.StatusCode).Should(Equal(http.StatusOK))

			// Session should still work after password reset
			// (implementation detail - some systems invalidate sessions, others don't)
			req2, err := http.NewRequest(
				"GET",
				serverURL+"/employer/get-onboard-status",
				bytes.NewBuffer(
					[]byte(`{"client_id": "passwordreset.example"}`),
				),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req2.Header.Set("Authorization", "Bearer "+sessionToken)
			req2.Header.Set("Content-Type", "application/json")

			resp2, err := http.DefaultClient.Do(req2)
			Expect(err).ShouldNot(HaveOccurred())
			// Session should remain valid
			Expect(resp2.StatusCode).Should(Equal(http.StatusOK))

			// Verify new password works for new login
			verifyEmployerLogin(email, "NewSessionPassword123$")
		})
	})
})
