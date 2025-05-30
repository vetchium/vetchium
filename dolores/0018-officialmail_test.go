package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Official Emails", Ordered, func() {
	var db *pgxpool.Pool
	var addToken, deleteToken, triggerToken, verifyToken, listToken, maxEmailsToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0018-officialmail-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(6)
		hubSigninAsync(
			"addemailuser@0018-hub.example",
			"NewPassword123$",
			&addToken,
			&wg,
		)
		hubSigninAsync(
			"deleteemailuser@0018-hub.example",
			"NewPassword123$",
			&deleteToken,
			&wg,
		)
		hubSigninAsync(
			"triggeruser@0018-hub.example",
			"NewPassword123$",
			&triggerToken,
			&wg,
		)
		hubSigninAsync(
			"verifyuser@0018-hub.example",
			"NewPassword123$",
			&verifyToken,
			&wg,
		)
		hubSigninAsync(
			"listemailsuser@0018-hub.example",
			"NewPassword123$",
			&listToken,
			&wg,
		)
		hubSigninAsync(
			"maxemailuser@0018-hub.example",
			"NewPassword123$",
			&maxEmailsToken,
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

		It(
			"should handle various standard add email test cases correctly",
			func() {
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
						wantStatus: http.StatusBadRequest,
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
						wantStatus: http.StatusConflict,
					},
				}

				for _, tc := range testCases {
					fmt.Fprintf(
						GinkgoWriter,
						"### Testing: %s\n",
						tc.description,
					)
					testPOST(
						tc.token,
						tc.request,
						"/hub/add-official-email",
						tc.wantStatus,
					)
				}
			},
		)

		It(
			"should reject adding email when maximum limit is reached for a user",
			func() {
				// This test uses 'maxemailuser@0018-hub.example' (maxEmailsToken)
				// First add 50 emails for this specific user to reach just under the limit
				for i := 1; i <= 50; i++ {
					email := fmt.Sprintf("max.bulk.%d@officialmail.example", i)
					testPOST(
						maxEmailsToken,
						hub.AddOfficialEmailRequest{
							Email: common.EmailAddress(email),
						},
						"/hub/add-official-email",
						http.StatusOK,
					)
				}

				// Now attempt to add the 51st email, which should be rejected
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing: exceed maximum allowed official emails for dedicated user\n",
				)
				testPOST(
					maxEmailsToken,
					hub.AddOfficialEmailRequest{
						Email: "add.max.final@officialmail.example",
					},
					"/hub/add-official-email",
					http.StatusUnprocessableEntity,
				)
			},
		)
	})

	Describe("Add Official Email - New Domain", Ordered, func() {
		var testUserToken string
		const newDomain = "newlycreated.example.com" // Unique enough for testing
		const newEmailAddress = "contact@" + newDomain

		BeforeAll(func() {
			// Use addToken from parent context
			testUserToken = addToken
		})

		It(
			"should successfully add an official email for a domain not yet in the system",
			func() {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing: Add official email for new domain %s\n",
					newEmailAddress,
				)
				testPOST(
					testUserToken,
					hub.AddOfficialEmailRequest{
						Email: common.EmailAddress(newEmailAddress),
					},
					"/hub/add-official-email",
					http.StatusOK,
				)

				// Verify that a verification email was sent using mailpit
				baseURL, err := url.Parse(mailPitURL + "/api/v1/search")
				Expect(err).ShouldNot(HaveOccurred())
				query := url.Values{}
				query.Add(
					"query",
					fmt.Sprintf(
						"to:%s subject:\"Vetchium - Confirm Email Ownership\"",
						newEmailAddress,
					),
				)
				baseURL.RawQuery = query.Encode()

				var messageID string
				var foundEmail bool
				// Try a few times with delay to allow email to be delivered
				delay := 10 * time.Second
				for i := 0; i < 5; i++ {
					<-time.After(delay)
					delay *= 2
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
						foundEmail = true
						break
					}
				}
				Expect(
					foundEmail,
				).Should(BeTrue(), "Verification email should have been sent")

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

				// Verify the email appears in the user's official emails list
				resp := testPOSTGetResp(
					testUserToken,
					nil,
					"/hub/my-official-emails",
					http.StatusOK,
				).([]byte)

				var emails []hub.OfficialEmail
				err = json.Unmarshal(resp, &emails)
				Expect(err).ShouldNot(HaveOccurred())

				var foundInList bool
				for _, email := range emails {
					if string(email.Email) == newEmailAddress {
						foundInList = true
						Expect(email.LastVerifiedAt).Should(BeNil())
						Expect(email.VerifyInProgress).Should(BeTrue())
						break
					}
				}
				Expect(
					foundInList,
				).Should(BeTrue(), "New email should appear in official emails list")
			},
		)
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
