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

var _ = Describe("Official Emails", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken1, hubToken2 string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0018-officialmail-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(2)
		hubSigninAsync(
			"officialmailuser1@hub.example",
			"NewPassword123$",
			&hubToken1,
			&wg,
		)
		hubSigninAsync(
			"officialmailuser2@hub.example",
			"NewPassword123$",
			&hubToken2,
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
						Email: "user2@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddOfficialEmailRequest{
						Email: "user2@officialmail.example",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid official email",
					token:       hubToken2,
					request: hub.AddOfficialEmailRequest{
						Email: "user2@officialmail.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "add official email with invalid domain",
					token:       hubToken2,
					request: hub.AddOfficialEmailRequest{
						Email: "user2@invalid-domain.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "add duplicate official email",
					token:       hubToken2,
					request: hub.AddOfficialEmailRequest{
						Email: "user2@officialmail.example",
					},
					wantStatus: http.StatusConflict,
				},
				{
					description: "add email with unverified domain",
					token:       hubToken2,
					request: hub.AddOfficialEmailRequest{
						Email: "user2@unverified.example",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
			}

			// First add 49 more emails to test max limit
			for i := 1; i <= 49; i++ {
				email := fmt.Sprintf("user1.%d@officialmail.example", i)
				testPOST(
					hubToken1,
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
				token:       hubToken1,
				request: hub.AddOfficialEmailRequest{
					Email: "user1.max@officialmail.example",
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
					description: "get official emails for user1",
					token:       hubToken1,
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var emails []hub.OfficialEmail
						err := json.Unmarshal(resp, &emails)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							len(emails),
						).Should(Equal(50))
						// 1 initial + 49 added in test
						Expect(
							string(emails[0].Email),
						).Should(Equal("user1@officialmail.example"))
						Expect(emails[0].LastVerifiedAt).ShouldNot(BeNil())
						Expect(emails[0].VerifyInProgress).Should(BeFalse())
					},
				},
				{
					description: "get official emails for user2",
					token:       hubToken2,
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var emails []hub.OfficialEmail
						err := json.Unmarshal(resp, &emails)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(emails)).Should(Equal(1))
						Expect(
							string(emails[0].Email),
						).Should(Equal("user2@officialmail.example"))
						Expect(emails[0].LastVerifiedAt).Should(BeNil())
						Expect(emails[0].VerifyInProgress).Should(BeTrue())
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
