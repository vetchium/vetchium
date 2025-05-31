package dolores

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = Describe("Employer Change Password", Ordered, func() {
	var db *pgxpool.Pool

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0030-employer-change-password-up.pgsql")
	})

	AfterAll(func() {
		seedDatabase(db, "0030-employer-change-password-down.pgsql")
		db.Close()
	})

	// Helper function to send change password request
	sendChangePasswordRequest := func(sessionToken, oldPassword, newPassword string) *http.Response {
		reqBody, err := json.Marshal(employer.EmployerChangePasswordRequest{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		})
		Expect(err).ShouldNot(HaveOccurred())

		req, err := http.NewRequest(
			"POST",
			serverURL+"/employer/change-password",
			bytes.NewBuffer(reqBody),
		)
		Expect(err).ShouldNot(HaveOccurred())
		req.Header.Set("Authorization", "Bearer "+sessionToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		Expect(err).ShouldNot(HaveOccurred())
		return resp
	}

	// Helper function to verify login with new password
	verifyNewPassword := func(email, password string) {
		sessionToken, err := employerSignin(
			"0030-changepassword.example",
			email,
			password,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sessionToken).ShouldNot(BeEmpty())
	}

	Describe("Password Change Flow", func() {
		It("should handle various password change scenarios", func() {
			type changePasswordTestCase struct {
				description  string
				email        string
				oldPassword  string
				newPassword  string
				wantStatus   int
				wantErrField string
				verifyLogin  bool
			}

			testCases := []changePasswordTestCase{
				{
					description:  "with invalid old password format",
					email:        "change1@0030-changepassword.example",
					oldPassword:  "short",
					newPassword:  "NewValidPassword123$",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "old_password",
				},
				{
					description:  "with invalid new password format",
					email:        "change2@0030-changepassword.example",
					oldPassword:  "NewPassword123$",
					newPassword:  "short",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "new_password",
				},
				{
					description:  "with incorrect old password",
					email:        "change3@0030-changepassword.example",
					oldPassword:  "WrongPassword123$",
					newPassword:  "NewValidPassword123$",
					wantStatus:   http.StatusUnauthorized,
					wantErrField: "",
				},
				{
					description: "with valid password change",
					email:       "change4@0030-changepassword.example",
					oldPassword: "NewPassword123$",
					newPassword: "NewValidPassword123$",
					wantStatus:  http.StatusOK,
					verifyLogin: true,
				},
				{
					description: "without authorization",
					email:       "change5@0030-changepassword.example",
					oldPassword: "NewPassword123$",
					newPassword: "NewValidPassword123$",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description:  "with empty old password",
					email:        "change6@0030-changepassword.example",
					oldPassword:  "",
					newPassword:  "NewValidPassword123$",
					wantStatus:   http.StatusBadRequest,
					wantErrField: "old_password",
				},
				{
					description: "with same old and new password",
					email:       "change7@0030-changepassword.example",
					oldPassword: "NewPassword123$",
					newPassword: "NewPassword123$",
					wantStatus:  http.StatusOK,
					verifyLogin: true,
				},
				{
					description: "with very long new password",
					email:       "change8@0030-changepassword.example",
					oldPassword: "NewPassword123$",
					newPassword: "VeryLongPassword123$" + string(
						make([]byte, 200),
					),
					wantStatus:   http.StatusBadRequest,
					wantErrField: "new_password",
				},
			}

			for _, tc := range testCases {
				GinkgoWriter.Printf("#### %s\n", tc.description)

				var sessionToken string
				var err error

				// Get session token for the user (except for unauthorized test)
				if tc.wantStatus != http.StatusUnauthorized {
					sessionToken, err = employerSignin(
						"0030-changepassword.example",
						tc.email,
						"NewPassword123$", // Initial password for all test users
					)
					Expect(err).ShouldNot(HaveOccurred())
				}

				// Send change password request
				resp := sendChangePasswordRequest(
					sessionToken,
					tc.oldPassword,
					tc.newPassword,
				)
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				// Check validation errors
				if tc.wantErrField != "" {
					var validationErrors common.ValidationErrors
					err := json.NewDecoder(resp.Body).Decode(&validationErrors)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(
						validationErrors.Errors,
					).Should(ContainElement(tc.wantErrField))
				}

				// Verify login with new password if successful
				if tc.verifyLogin && tc.wantStatus == http.StatusOK {
					verifyNewPassword(tc.email, tc.newPassword)
				}

				// Verify old password no longer works after successful change
				if tc.wantStatus == http.StatusOK &&
					tc.oldPassword != tc.newPassword {
					_, err := employerSignin(
						"0030-changepassword.example",
						tc.email,
						tc.oldPassword,
					)
					Expect(err).Should(HaveOccurred())
				}
			}
		})

		It("should maintain session validity after password change", func() {
			email := "session-test@0030-changepassword.example"
			oldPassword := "NewPassword123$"
			newPassword := "NewSessionPassword123$"

			// Get initial session token
			sessionToken, err := employerSignin(
				"0030-changepassword.example",
				email,
				oldPassword,
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify session works before password change
			req, err := http.NewRequest(
				"GET",
				serverURL+"/employer/get-onboard-status",
				bytes.NewBuffer(
					[]byte(`{"client_id": "0030-changepassword.example"}`),
				),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+sessionToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Change password
			changeResp := sendChangePasswordRequest(
				sessionToken,
				oldPassword,
				newPassword,
			)
			Expect(changeResp.StatusCode).Should(Equal(http.StatusOK))

			// Session should still work after password change
			req2, err := http.NewRequest(
				"GET",
				serverURL+"/employer/get-onboard-status",
				bytes.NewBuffer(
					[]byte(`{"client_id": "0030-changepassword.example"}`),
				),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req2.Header.Set("Authorization", "Bearer "+sessionToken)
			req2.Header.Set("Content-Type", "application/json")

			resp2, err := http.DefaultClient.Do(req2)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp2.StatusCode).Should(Equal(http.StatusOK))

			// Verify new password works for new login
			verifyNewPassword(email, newPassword)

			// Verify old password no longer works
			_, err = employerSignin(
				"0030-changepassword.example",
				email,
				oldPassword,
			)
			Expect(err).Should(HaveOccurred())
		})

		It("should require authentication", func() {
			reqBody, err := json.Marshal(employer.EmployerChangePasswordRequest{
				OldPassword: "OldPassword123$",
				NewPassword: "NewPassword123$",
			})
			Expect(err).ShouldNot(HaveOccurred())

			// Try without authentication
			resp, err := http.Post(
				serverURL+"/employer/change-password",
				"application/json",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))

			// Try with invalid token
			req, err := http.NewRequest(
				"POST",
				serverURL+"/employer/change-password",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer invalid-token")
			req.Header.Set("Content-Type", "application/json")

			resp2, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp2.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})
})
