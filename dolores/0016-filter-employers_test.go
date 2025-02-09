package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
			searchTerm  string
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various test cases correctly", func() {
			testCases := []filterEmployersTestCase{
				{
					description: "without authentication",
					token:       "",
					searchTerm:  "acme",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					searchTerm:  "acme",
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "search for existing employer",
					token:       hubToken,
					searchTerm:  "acme",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var employers []struct {
							Name   string `json:"name"`
							Domain string `json:"domain"`
						}
						err := json.Unmarshal(resp, &employers)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(employers)).Should(Equal(1))
						Expect(employers[0].Name).Should(Equal("Acme Corp"))
						Expect(
							employers[0].Domain,
						).Should(Equal("acme.example"))
					},
				},
				{
					description: "search for non-existent employer",
					token:       hubToken,
					searchTerm:  "nonexistent",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var employers []struct {
							Name   string `json:"name"`
							Domain string `json:"domain"`
						}
						err := json.Unmarshal(resp, &employers)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(employers)).Should(BeZero())
					},
				},
				{
					description: "search for domain without employer",
					token:       hubToken,
					searchTerm:  "domain-without-employer",
					wantStatus:  http.StatusOK,
					validate: func(resp []byte) {
						var employers []struct {
							Name   string `json:"name"`
							Domain string `json:"domain"`
						}
						err := json.Unmarshal(resp, &employers)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(employers)).Should(Equal(1))
						Expect(
							employers[0].Name,
						).Should(Equal("domain-without-employer.example"))
						Expect(
							employers[0].Domain,
						).Should(Equal("domain-without-employer.example"))
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testGETWithQueryGetResp(
					tc.token,
					"/hub/filter-employers",
					map[string]string{"search_term": tc.searchTerm},
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
