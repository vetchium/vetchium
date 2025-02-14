package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = FDescribe("Official Emails", Ordered, func() {
	var db *pgxpool.Pool
	var addToken, deleteToken, triggerToken, verifyToken, listToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0018-officialmail-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(5)
		hubSigninAsync(
			"addemailuser@hub.example",
			"NewPassword123$",
			&addToken,
			&wg,
		)
		hubSigninAsync(
			"deleteemailuser@hub.example",
			"NewPassword123$",
			&deleteToken,
			&wg,
		)
		hubSigninAsync(
			"triggeruser@hub.example",
			"NewPassword123$",
			&triggerToken,
			&wg,
		)
		hubSigninAsync(
			"verifyuser@hub.example",
			"NewPassword123$",
			&verifyToken,
			&wg,
		)
		hubSigninAsync(
			"listemailsuser@hub.example",
			"NewPassword123$",
			&listToken,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0018-officialmail-down.pgsql")
		db.Close()
	})

	Describe("Add Official Email", func() {
		type addOfficialEmailTestCase struct {
			description string
			token       string
			request     hub.AddOfficialEmailRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []addOfficialEmailTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddOfficialEmailRequest{
						Email: "add.new@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddOfficialEmailRequest{
						Email: "add.new@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add email with invalid domain format",
					token:       addToken,
					request: hub.AddOfficialEmailRequest{
						Email: "add.new@invalid-domain",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "add valid official email",
					token:       addToken,
					request: hub.AddOfficialEmailRequest{
						Email: "add.new@officialmail.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "add duplicate official email",
					token:       addToken,
					request: hub.AddOfficialEmailRequest{
						Email: "add.new@officialmail.example",
					},
					wantStatus: http.StatusPreconditionFailed,
				},
			}

			// First add 49 more emails to test max limit
			for i := 1; i <= 49; i++ {
				email := fmt.Sprintf("add.bulk.%d@officialmail.example", i)
				testPOST(
					addToken,
					hub.AddOfficialEmailRequest{
						Email: common.EmailAddress(email),
					},
					"/hub/add-official-email",
					http.StatusOK,
				)
			}

			// Now add test case for max emails reached
			testCases = append(testCases, addOfficialEmailTestCase{
				description: "exceed maximum allowed official emails",
				token:       addToken,
				request: hub.AddOfficialEmailRequest{
					Email: "add.max@officialmail.example",
				},
				wantStatus: http.StatusPreconditionFailed,
			})

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/add-official-email",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Delete Official Email", func() {
		type deleteOfficialEmailTestCase struct {
			description string
			token       string
			request     hub.DeleteOfficialEmailRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []deleteOfficialEmailTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.DeleteOfficialEmailRequest{
						Email: "delete.verified@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.DeleteOfficialEmailRequest{
						Email: "delete.verified@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete non-existent email",
					token:       deleteToken,
					request: hub.DeleteOfficialEmailRequest{
						Email: "nonexistent@officialmail.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "delete email belonging to another user",
					token:       addToken,
					request: hub.DeleteOfficialEmailRequest{
						Email: "delete.verified@officialmail.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "delete pending verification email",
					token:       deleteToken,
					request: hub.DeleteOfficialEmailRequest{
						Email: "delete.pending@officialmail.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete verified email",
					token:       deleteToken,
					request: hub.DeleteOfficialEmailRequest{
						Email: "delete.verified@officialmail.example",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/delete-official-email",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Trigger Verification", func() {
		type triggerVerificationTestCase struct {
			description string
			token       string
			request     hub.TriggerVerificationRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []triggerVerificationTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.TriggerVerificationRequest{
						Email: "trigger.old@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.TriggerVerificationRequest{
						Email: "trigger.old@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "trigger verification for non-existent email",
					token:       triggerToken,
					request: hub.TriggerVerificationRequest{
						Email: "nonexistent@officialmail.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "trigger verification for recently verified email",
					token:       triggerToken,
					request: hub.TriggerVerificationRequest{
						Email: "trigger.recent@officialmail.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "trigger verification for old verified email",
					token:       triggerToken,
					request: hub.TriggerVerificationRequest{
						Email: "trigger.old@officialmail.example",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/trigger-verification",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Verify Official Email", func() {
		type verifyOfficialEmailTestCase struct {
			description string
			token       string
			request     hub.VerifyOfficialEmailRequest
			wantStatus  int
		}

		It("should handle various test cases correctly", func() {
			testCases := []verifyOfficialEmailTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.pending@officialmail.example",
						Code:  "VERIFY123",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.pending@officialmail.example",
						Code:  "VERIFY123",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "verify non-existent email",
					token:       verifyToken,
					request: hub.VerifyOfficialEmailRequest{
						Email: "nonexistent@officialmail.example",
						Code:  "VERIFY123",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "verify with incorrect code",
					token:       verifyToken,
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.pending@officialmail.example",
						Code:  "WRONG123",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "verify with expired code",
					token:       verifyToken,
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.expired@officialmail.example",
						Code:  "EXPIRED",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "verify with correct code",
					token:       verifyToken,
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.pending@officialmail.example",
						Code:  "VERIFY123",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "verify already verified email",
					token:       verifyToken,
					request: hub.VerifyOfficialEmailRequest{
						Email: "verify.pending@officialmail.example",
						Code:  "VERIFY123",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/hub/verify-official-email",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("My Official Emails", func() {
		type myOfficialEmailsTestCase struct {
			description string
			token       string
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []myOfficialEmailsTestCase{
				{
					description: "without authentication",
					token:       "",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "get official emails for list user",
					token:       listToken,
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var emails []hub.OfficialEmail
						err := json.Unmarshal(resp, &emails)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(emails)).Should(Equal(3))

						// Verify the emails are present
						var emailAddresses []string
						for _, email := range emails {
							emailAddresses = append(
								emailAddresses,
								string(email.Email),
							)
						}
						Expect(emailAddresses).Should(ContainElements(
							"list.verified@officialmail.example",
							"list.pending@officialmail.example",
							"list.expired@officialmail.example",
						))

						// Verify the states
						for _, email := range emails {
							switch email.Email {
							case "list.verified@officialmail.example":
								Expect(email.LastVerifiedAt).ShouldNot(BeNil())
								Expect(email.VerifyInProgress).Should(BeFalse())
							case "list.pending@officialmail.example":
								Expect(email.LastVerifiedAt).Should(BeNil())
								Expect(email.VerifyInProgress).Should(BeTrue())
							case "list.expired@officialmail.example":
								Expect(email.LastVerifiedAt).Should(BeNil())
								Expect(email.VerifyInProgress).Should(BeTrue())
							}
						}
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					nil,
					"/hub/my-official-emails",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
