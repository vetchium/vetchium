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

var _ = Describe("Hub Login", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0007-hub-login-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0007-hub-login-down.pgsql")
		db.Close()
	})

	// Helper functions to reduce code duplication
	getLoginToken := func(email, password string) string {
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
		return loginRespObj.Token
	}

	getTFACode := func(email string) (string, string) {
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

	getSessionToken := func(tfaToken, tfaCode string, rememberMe bool) string {
		tfaReqBody, err := json.Marshal(hub.HubTFARequest{
			TFAToken:   tfaToken,
			TFACode:    tfaCode,
			RememberMe: rememberMe,
		})
		Expect(err).ShouldNot(HaveOccurred())

		resp, err := http.Post(
			serverURL+"/hub/tfa",
			"application/json",
			bytes.NewBuffer(tfaReqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		var tfaResp hub.HubTFAResponse
		err = json.NewDecoder(resp.Body).Decode(&tfaResp)
		Expect(err).ShouldNot(HaveOccurred())
		return tfaResp.SessionToken
	}

	Describe("Hub Login Flow", Ordered, func() {
		type loginTestCase struct {
			description   string
			request       hub.LoginRequest
			wantStatus    int
			wantErrFields []string
		}

		It("should handle login requests correctly", func() {
			testCases := []loginTestCase{
				{
					description: "valid credentials for active user",
					request: hub.LoginRequest{
						Email:    "active@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "invalid password for active user",
					request: hub.LoginRequest{
						Email:    "active@hub.example",
						Password: "WrongPassword123$",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "disabled user",
					request: hub.LoginRequest{
						Email:    "disabled@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "deleted user",
					request: hub.LoginRequest{
						Email:    "deleted@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "non-existent user",
					request: hub.LoginRequest{
						Email:    "nonexistent@hub.example",
						Password: "NewPassword123$",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid email format",
					request: hub.LoginRequest{
						Email:    "invalid-email",
						Password: "NewPassword123$",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"email"},
				},
				{
					description: "empty password",
					request: hub.LoginRequest{
						Email:    "active@hub.example",
						Password: "",
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"password"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				loginReqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/hub/login",
					"application/json",
					bytes.NewBuffer(loginReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())

				if resp.StatusCode != tc.wantStatus {
					body, err := io.ReadAll(resp.Body)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Fprintf(GinkgoWriter, "#### %s\n", string(body))
					Fail(
						fmt.Sprintf(
							"want status %d, got %d",
							tc.wantStatus,
							resp.StatusCode,
						),
					)
					return
				}
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if len(tc.wantErrFields) > 0 {
					var validationErrors common.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
					continue
				}

				if tc.wantStatus == http.StatusOK {
					var loginResp hub.LoginResponse
					err = json.NewDecoder(resp.Body).Decode(&loginResp)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(loginResp.Token).ShouldNot(BeEmpty())
				}
			}
		})

		type tfaTestCase struct {
			description   string
			request       hub.HubTFARequest
			wantStatus    int
			wantErrFields []string
		}

		It("should handle TFA flow correctly", func() {
			// First get a valid TFA token through login
			email := "tfatest@hub.example"
			loginReqBody, err := json.Marshal(hub.LoginRequest{
				Email:    common.EmailAddress(email),
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

			var loginRespObj hub.LoginResponse
			err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			tfaToken := loginRespObj.Token

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

			fmt.Fprintf(GinkgoWriter, "mailURL: %s\n", mailURL)
			// Get the TFA code from mailpit
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

			// Get the email content
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
			tfaCode := matches[1]

			testCases := []tfaTestCase{
				{
					description: "valid TFA token and code",
					request: hub.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "invalid TFA token",
					request: hub.HubTFARequest{
						TFAToken:   "invalid-token",
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid TFA code",
					request: hub.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    "000000",
						RememberMe: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "empty TFA token",
					request: hub.HubTFARequest{
						TFAToken:   "",
						TFACode:    tfaCode,
						RememberMe: false,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"tfa_token"},
				},
				{
					description: "empty TFA code",
					request: hub.HubTFARequest{
						TFAToken:   tfaToken,
						TFACode:    "",
						RememberMe: false,
					},
					wantStatus:    http.StatusBadRequest,
					wantErrFields: []string{"tfa_code"},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				tfaReqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/hub/tfa",
					"application/json",
					bytes.NewBuffer(tfaReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if len(tc.wantErrFields) > 0 {
					var validationErrors common.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElements(tc.wantErrFields))
					continue
				}

				if tc.wantStatus == http.StatusOK {
					var tfaResp hub.HubTFAResponse
					err = json.NewDecoder(resp.Body).Decode(&tfaResp)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(tfaResp.SessionToken).ShouldNot(BeEmpty())
				}
			}

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
		})

		It("should handle remember me flag correctly", func() {
			// First get a valid TFA token through login
			loginReqBody, err := json.Marshal(hub.LoginRequest{
				Email:    "rememberme@hub.example",
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

			var loginRespObj hub.LoginResponse
			err = json.NewDecoder(loginResp.Body).Decode(&loginRespObj)
			Expect(err).ShouldNot(HaveOccurred())
			tfaToken := loginRespObj.Token

			baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
			Expect(err).ShouldNot(HaveOccurred())
			query := url.Values{}
			query.Add(
				"query",
				"to:rememberme@hub.example subject:Vetchium Two Factor Authentication",
			)
			baseURL.RawQuery = query.Encode()
			mailURL := baseURL.String()

			// Get the TFA code from mailpit
			var messageID string
			for i := 0; i < 5; i++ {
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

			// Get the email content
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
			tfaCode := matches[1]

			// Test with remember_me flag
			tfaReqBody, err := json.Marshal(hub.HubTFARequest{
				TFAToken:   tfaToken,
				TFACode:    tfaCode,
				RememberMe: true,
			})
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/hub/tfa",
				"application/json",
				bytes.NewBuffer(tfaReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var tfaResp hub.HubTFAResponse
			err = json.NewDecoder(resp.Body).Decode(&tfaResp)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tfaResp.SessionToken).ShouldNot(BeEmpty())

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
		})

		It("test Hub user logout", func() {
			email := "active@hub.example"
			password := "NewPassword123$"

			// Get login token
			tfaToken := getLoginToken(email, password)

			// Get TFA code from email
			tfaCode, messageID := getTFACode(email)
			defer cleanupEmail(messageID)

			// Get session token
			sessionToken := getSessionToken(tfaToken, tfaCode, false)

			// Test get-my-handle endpoint with valid session
			req, err := http.NewRequest(
				"GET",
				serverURL+"/hub/get-my-handle",
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Test logout
			logoutReq, err := http.NewRequest(
				"POST",
				serverURL+"/hub/logout",
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			logoutReq.Header.Set("Authorization", "Bearer "+sessionToken)

			logoutResp, err := http.DefaultClient.Do(logoutReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(logoutResp.StatusCode).Should(Equal(http.StatusOK))

			// Test get-my-handle endpoint after logout (should fail)
			req2, err := http.NewRequest(
				"GET",
				serverURL+"/hub/get-my-handle",
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			req2.Header.Set("Authorization", "Bearer "+sessionToken)

			resp2, err := http.DefaultClient.Do(req2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp2.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("should handle password change correctly", func() {
			email := "password-change@hub.example"
			oldPassword := "NewPassword123$"
			newPassword := "UpdatedPassword123$"

			// Get initial session token
			tfaToken := getLoginToken(email, oldPassword)
			tfaCode, messageID := getTFACode(email)
			defer cleanupEmail(messageID)
			sessionToken := getSessionToken(tfaToken, tfaCode, false)

			// Test invalid password scenarios
			testCases := []struct {
				description  string
				oldPassword  string
				newPassword  string
				wantStatus   int
				wantErrField string
			}{
				{
					description:  "invalid old password format",
					oldPassword:  "short",
					newPassword:  "NewPassword123$",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "old_password",
				},
				{
					description:  "invalid new password format",
					oldPassword:  "NewPassword123$",
					newPassword:  "weak",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "new_password",
				},
				{
					description: "incorrect old password",
					oldPassword: "WrongPassword123$",
					newPassword: "NewPassword123$",
					wantStatus:  http.StatusUnauthorized,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "Test case: %s\n", tc.description)

				changePasswordReqBody, err := json.Marshal(
					hub.ChangePasswordRequest{
						OldPassword: common.Password(tc.oldPassword),
						NewPassword: common.Password(tc.newPassword),
					},
				)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/hub/change-password",
					bytes.NewBuffer(changePasswordReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+sessionToken)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantErrField != "" {
					var validationErrors common.ValidationErrors
					err = json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElement(tc.wantErrField))
				}
			}

			// Verify original password still works after failed attempts
			tfaToken = getLoginToken(email, oldPassword)
			tfaCode, messageID = getTFACode(email)
			defer cleanupEmail(messageID)
			sessionToken = getSessionToken(tfaToken, tfaCode, false)

			// Test successful password change
			changePasswordReqBody, err := json.Marshal(
				hub.ChangePasswordRequest{
					OldPassword: common.Password(oldPassword),
					NewPassword: common.Password(newPassword),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"POST",
				serverURL+"/hub/change-password",
				bytes.NewBuffer(changePasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Verify old password no longer works
			loginReqBody, err := json.Marshal(hub.LoginRequest{
				Email:    common.EmailAddress(email),
				Password: common.Password(oldPassword),
			})
			Expect(err).ShouldNot(HaveOccurred())

			loginResp, err := http.Post(
				serverURL+"/hub/login",
				"application/json",
				bytes.NewBuffer(loginReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(loginResp.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Verify new password works
			tfaToken = getLoginToken(email, newPassword)
			tfaCode, messageID = getTFACode(email)
			defer cleanupEmail(messageID)
			sessionToken = getSessionToken(tfaToken, tfaCode, false)

			req, err = http.NewRequest(
				"GET",
				serverURL+"/hub/get-my-handle",
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})

		It("Forgot Password and Reset Password", func() {
			type forgotPasswordTestCase struct {
				description string
				request     hub.ForgotPasswordRequest
				wantStatus  int
				checkEmail  bool // whether to check for email
			}

			testCases := []forgotPasswordTestCase{
				{
					description: "with valid email",
					request: hub.ForgotPasswordRequest{
						Email: "password-reset@hub.example",
					},
					wantStatus: http.StatusOK,
					checkEmail: true,
				},
				{
					description: "with non-existent email",
					request: hub.ForgotPasswordRequest{
						Email: "nonexistent@hub.example",
					},
					wantStatus: http.StatusOK,
					checkEmail: false,
				},
				{
					description: "with invalid email format",
					request: hub.ForgotPasswordRequest{
						Email: "invalid-email",
					},
					wantStatus: http.StatusBadRequest,
					checkEmail: false,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)

				// Clear any existing emails before test
				if tc.checkEmail {
					baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
					Expect(err).ShouldNot(HaveOccurred())
					query := url.Values{}
					query.Add(
						"query",
						fmt.Sprintf(
							"to:%s subject:Vetchium Password Reset",
							tc.request.Email,
						),
					)
					baseURL.RawQuery = query.Encode()

					mailPitResp, err := http.Get(baseURL.String())
					Expect(err).ShouldNot(HaveOccurred())
					body, err := io.ReadAll(mailPitResp.Body)
					Expect(err).ShouldNot(HaveOccurred())

					var mailPitRespObj MailPitResponse
					err = json.Unmarshal(body, &mailPitRespObj)
					Expect(err).ShouldNot(HaveOccurred())

					if len(mailPitRespObj.Messages) > 0 {
						deleteReqBody, err := json.Marshal(MailPitDeleteRequest{
							IDs: []string{mailPitRespObj.Messages[0].ID},
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
						Expect(
							deleteResp.StatusCode,
						).Should(Equal(http.StatusOK))
					}
				}

				// Send forgot password request
				forgotPasswordReqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				resp, err := http.Post(
					serverURL+"/hub/forgot-password",
					"application/json",
					bytes.NewBuffer(forgotPasswordReqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if !tc.checkEmail {
					continue
				}

				// Check for password reset email
				var messageID string
				var resetToken string

				baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
				Expect(err).ShouldNot(HaveOccurred())
				query := url.Values{}
				query.Add(
					"query",
					fmt.Sprintf(
						"to:%s subject:Vetchium Password Reset",
						tc.request.Email,
					),
				)
				baseURL.RawQuery = query.Encode()

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

				re := regexp.MustCompile(
					`/reset-password\?token=([a-zA-Z0-9]+)`,
				)
				matches := re.FindStringSubmatch(string(body))
				Expect(len(matches)).Should(BeNumerically(">=", 2))
				resetToken = matches[1]

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

				// Test reset password scenarios
				type resetPasswordTestCase struct {
					description   string
					request       hub.ResetPasswordRequest
					wantStatus    int
					wantErrFields []string
					sleep         time.Duration
				}

				resetTestCases := []resetPasswordTestCase{
					{
						description: "with invalid password format",
						request: hub.ResetPasswordRequest{
							Token:    resetToken,
							Password: "weak",
						},
						wantStatus:    http.StatusBadRequest,
						wantErrFields: []string{"password"},
					},
					{
						description: "with invalid token",
						request: hub.ResetPasswordRequest{
							Token:    "invalid-token",
							Password: "NewPassword123$",
						},
						wantStatus: http.StatusUnauthorized,
					},
					{
						description: "with valid token and password",
						request: hub.ResetPasswordRequest{
							Token:    resetToken,
							Password: "NewPassword123$",
						},
						wantStatus: http.StatusOK,
					},
				}

				for _, rtc := range resetTestCases {
					fmt.Fprintf(
						GinkgoWriter,
						"#### Reset Password: %s\n",
						rtc.description,
					)

					if rtc.sleep > 0 {
						<-time.After(rtc.sleep)
					}

					resetPasswordReqBody, err := json.Marshal(rtc.request)
					Expect(err).ShouldNot(HaveOccurred())

					resetResp, err := http.Post(
						serverURL+"/hub/reset-password",
						"application/json",
						bytes.NewBuffer(resetPasswordReqBody),
					)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resetResp.StatusCode).Should(Equal(rtc.wantStatus))

					if len(rtc.wantErrFields) > 0 {
						var validationErrors common.ValidationErrors
						err = json.NewDecoder(resetResp.Body).
							Decode(&validationErrors)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(validationErrors.Errors).Should(
							ContainElements(rtc.wantErrFields),
						)
					}

					// If password was successfully reset, verify we can login with new password
					if rtc.wantStatus == http.StatusOK {
						// Try logging in with new password
						tfaToken := getLoginToken(
							string(tc.request.Email),
							string(rtc.request.Password),
						)
						tfaCode, messageID := getTFACode(
							string(tc.request.Email),
						)
						sessionToken := getSessionToken(
							tfaToken,
							tfaCode,
							false,
						)

						// Verify session token works
						req, err := http.NewRequest(
							"GET",
							serverURL+"/hub/get-my-handle",
							nil,
						)
						Expect(err).ShouldNot(HaveOccurred())
						req.Header.Set("Authorization", "Bearer "+sessionToken)

						resp, err := http.DefaultClient.Do(req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.StatusCode).Should(Equal(http.StatusOK))

						// Clean up TFA email
						cleanupEmail(messageID)
					}
				}
			}
		})

		It("Password Reset Token Expiry", func() {
			email := "token-expiry@hub.example"

			// Send forgot password request
			forgotPasswordReqBody, err := json.Marshal(
				hub.ForgotPasswordRequest{
					Email: common.EmailAddress(email),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/hub/forgot-password",
				"application/json",
				bytes.NewBuffer(forgotPasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get the reset token from email
			var messageID string
			var resetToken string

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

			// Clean up the email
			cleanupEmail(messageID)

			// Wait for token to expire
			<-time.After(7 * time.Minute)

			// Try to use expired token
			resetPasswordReqBody, err := json.Marshal(
				hub.ResetPasswordRequest{
					Token:    resetToken,
					Password: "NewPassword123$",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			resetResp, err := http.Post(
				serverURL+"/hub/reset-password",
				"application/json",
				bytes.NewBuffer(resetPasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resetResp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})

		It("Password Reset Token Reuse", func() {
			email := "token-reuse@hub.example"

			// Send forgot password request
			forgotPasswordReqBody, err := json.Marshal(
				hub.ForgotPasswordRequest{
					Email: common.EmailAddress(email),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := http.Post(
				serverURL+"/hub/forgot-password",
				"application/json",
				bytes.NewBuffer(forgotPasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Get the reset token from email
			var messageID string
			var resetToken string

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

			// Clean up the email
			cleanupEmail(messageID)

			// First password reset should succeed
			resetPasswordReqBody, err := json.Marshal(
				hub.ResetPasswordRequest{
					Token:    resetToken,
					Password: "NewPassword123$",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			resetResp, err := http.Post(
				serverURL+"/hub/reset-password",
				"application/json",
				bytes.NewBuffer(resetPasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resetResp.StatusCode).Should(Equal(http.StatusOK))

			// Verify login works with new password
			tfaToken := getLoginToken(email, "NewPassword123$")
			tfaCode, messageID := getTFACode(email)
			sessionToken := getSessionToken(tfaToken, tfaCode, false)

			req, err := http.NewRequest(
				"GET",
				serverURL+"/hub/get-my-handle",
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Clean up TFA email
			cleanupEmail(messageID)

			// Try to reuse the same token - should fail
			resetPasswordReqBody, err = json.Marshal(
				hub.ResetPasswordRequest{
					Token:    resetToken,
					Password: "AnotherPassword123$",
				},
			)
			Expect(err).ShouldNot(HaveOccurred())

			resetResp, err = http.Post(
				serverURL+"/hub/reset-password",
				"application/json",
				bytes.NewBuffer(resetPasswordReqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resetResp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})
})
