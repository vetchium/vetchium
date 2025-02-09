package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = Describe("Filter Employers", Ordered, func() {
	var db *pgxpool.Pool
	var hubToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0016-filter-employers-up.pgsql")

		// Login hub user and get token
		var wg sync.WaitGroup
		wg.Add(1)
		hubSigninAsync(
			"user1@filter-employers.example",
			"NewPassword123$",
			&hubToken,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0016-filter-employers-down.pgsql")
		db.Close()
	})

	Describe("Filter Employers Tests", func() {
		type filterEmployersTestCase struct {
			description string
			token       string
			prefix      string
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []filterEmployersTestCase{
				{
					description: "without authentication",
					token:       "",
					prefix:      "acme",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					prefix:      "acme",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "search for existing employer",
					token:       hubToken,
					prefix:      "acme",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var response hub.FilterEmployersResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(response.Employers)).Should(Equal(1))
						Expect(
							response.Employers[0].Name,
						).Should(Equal("Acme Corp"))
						Expect(
							response.Employers[0].Domain,
						).Should(Equal("acme.example"))
					},
				},
				{
					description: "search for non-existent employer",
					token:       hubToken,
					prefix:      "nonexistent",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var response hub.FilterEmployersResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(response.Employers)).Should(BeZero())
					},
				},
				{
					description: "search for domain without employer",
					token:       hubToken,
					prefix:      "domain-without-employer",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var response hub.FilterEmployersResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(response.Employers)).Should(Equal(1))
						Expect(
							response.Employers[0].Name,
						).Should(Equal("domain-without-employer.example"))
						Expect(
							response.Employers[0].Domain,
						).Should(Equal("domain-without-employer.example"))
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)

				// Create request with proper struct
				req := hub.FilterEmployersRequest{
					Prefix: tc.prefix,
				}

				resp := testPOSTGetResp(
					tc.token,
					req,
					"/hub/filter-employers",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
