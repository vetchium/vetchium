package dolores

import (
	"encoding/json"
	"fmt"

	// math rand is sufficient as we just need a number; no need for crypto
	"math/rand"

	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

var bachelorEducation_0006 *common.EducationLevel

var _ = Describe("Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOpeningsToken string
	var recruiterToken, hiringManagerToken string

	BeforeAll(func() {
		bachelor := common.BachelorEducation
		bachelorEducation_0006 = &bachelor

		db = setupTestDB()
		seedDatabase(db, "0006-watchers-and-filter-openings-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@openings0006.example":          &adminToken,
			"crud@openings0006.example":           &crudToken,
			"viewer@openings0006.example":         &viewerToken,
			"non-openings@openings0006.example":   &nonOpeningsToken,
			"recruiter@openings0006.example":      &recruiterToken,
			"hiring-manager@openings0006.example": &hiringManagerToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"openings0006.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0006-watchers-and-filter-openings-down.pgsql")
		db.Close()
	})

	Describe("Filter Openings", func() {
		type filterOpeningsTestCase struct {
			description string
			token       string
			request     employer.FilterOpeningsRequest
			wantStatus  int
			wantIDs     []string
		}

		It("should filter openings correctly", func() {
			testCases := []filterOpeningsTestCase{
				{
					description: "with no filters on state (should only return DRAFT, ACTIVE, SUSPENDED states)",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Feb-1")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
					// All DRAFT, ACTIVE, SUSPENDED openings
					wantIDs: []string{
						"2024-Feb-15-001", // DRAFT
						"2024-Feb-25-001", // ACTIVE
						"2024-Mar-01-001", // ACTIVE
						"2024-Mar-01-002", // DRAFT
						"2024-Mar-01-003", // SUSPENDED
						"2024-Mar-01-005", // ACTIVE
						"2024-Mar-06-001", // SUSPENDED
						"2024-Mar-06-002", // ACTIVE
						"2024-Mar-06-003", // DRAFT
						"2024-Mar-06-005", // ACTIVE
						"2024-Mar-06-006", // SUSPENDED
						"2024-Mar-06-007", // DRAFT
						"2024-Mar-06-009", // ACTIVE
						"2024-Mar-06-010", // SUSPENDED
						"2024-Mar-06-011", // ACTIVE
						"2024-Mar-06-012", // DRAFT
						"2024-Mar-06-013", // SUSPENDED
						"2024-Mar-06-014", // ACTIVE
						"2024-Mar-06-015", // DRAFT
					},
				},
				{
					description: "with state filter - draft only",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						State: []common.OpeningState{common.DraftOpening},
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Feb-1")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
					wantIDs: []string{
						"2024-Feb-15-001",
						"2024-Mar-01-002",
						"2024-Mar-06-003",
						"2024-Mar-06-007",
						"2024-Mar-06-012",
						"2024-Mar-06-015",
					},
				},
				{
					description: "with state filter - active only",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						State: []common.OpeningState{common.ActiveOpening},
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Feb-1")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
					wantIDs: []string{
						"2024-Feb-25-001",
						"2024-Mar-01-001",
						"2024-Mar-01-005",
						"2024-Mar-06-002",
						"2024-Mar-06-005",
						"2024-Mar-06-009",
						"2024-Mar-06-011",
						"2024-Mar-06-014",
					},
				},
				{
					description: "with date range filter - March 6 only",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Mar-6")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
						ToDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Mar-7")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
					wantIDs: []string{
						"2024-Mar-06-001",
						"2024-Mar-06-002",
						"2024-Mar-06-003",
						"2024-Mar-06-005",
						"2024-Mar-06-006",
						"2024-Mar-06-007",
						"2024-Mar-06-009",
						"2024-Mar-06-010",
						"2024-Mar-06-011",
						"2024-Mar-06-012",
						"2024-Mar-06-013",
						"2024-Mar-06-014",
						"2024-Mar-06-015",
					},
				},
				{
					description: "with pagination - first page",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Feb-1")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
						Limit: 5,
					},
					wantStatus: http.StatusOK,
					wantIDs: []string{
						"2024-Feb-15-001",
						"2024-Feb-25-001",
						"2024-Mar-01-001",
						"2024-Mar-01-002",
						"2024-Mar-01-003",
					},
				},
				{
					description: "with pagination - using pagination key",
					token:       adminToken,
					request: employer.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t, err := time.Parse("2006-Jan-2", "2024-Feb-1")
							Expect(err).ShouldNot(HaveOccurred())
							return &t
						}(),
						PaginationKey: "2024-Mar-01-003",
						Limit:         5,
					},
					wantStatus: http.StatusOK,
					wantIDs: []string{
						"2024-Mar-01-005",
						"2024-Mar-06-001",
						"2024-Mar-06-002",
						"2024-Mar-06-003",
						"2024-Mar-06-005",
					},
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request:     employer.FilterOpeningsRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request:     employer.FilterOpeningsRequest{},
					wantStatus:  http.StatusForbidden,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/filter-openings",
					tc.wantStatus,
				)

				if tc.wantStatus == http.StatusOK {
					var openings []employer.OpeningInfo
					err := json.Unmarshal(resp.([]byte), &openings)
					Expect(err).ShouldNot(HaveOccurred())

					gotIDs := make([]string, len(openings))
					for i, opening := range openings {
						gotIDs[i] = opening.ID
					}
					Expect(gotIDs).Should(ConsistOf(tc.wantIDs))

					if len(gotIDs) != len(tc.wantIDs) {
						fmt.Fprintf(
							GinkgoWriter,
							"got %d:%v\n",
							len(gotIDs),
							gotIDs,
						)
						fmt.Fprintf(
							GinkgoWriter,
							"want %d:%v\n",
							len(tc.wantIDs),
							tc.wantIDs,
						)
						Fail("got wrong number of openings")
					}
				}
			}
		})
	})

	Describe("Opening Watchers", func() {
		It("should get opening watchers correctly", func() {
			type getWatchersTestCase struct {
				description string
				token       string
				openingID   string
				wantStatus  int
				wantEmails  []string
			}

			testCases := []getWatchersTestCase{
				{
					description: "get watchers for opening with multiple watchers",
					token:       adminToken,
					openingID:   "2024-Feb-15-001",
					wantStatus:  http.StatusOK,
					wantEmails: []string{
						"watcher1@openings0006.example",
						"watcher2@openings0006.example",
					},
				},
				{
					description: "get watchers for opening with single watcher",
					token:       adminToken,
					openingID:   "2024-Feb-25-001",
					wantStatus:  http.StatusOK,
					wantEmails:  []string{"watcher1@openings0006.example"},
				},
				{
					description: "get watchers for opening with no watchers",
					token:       adminToken,
					openingID:   "2024-Mar-06-001",
					wantStatus:  http.StatusOK,
					wantEmails:  []string{},
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					openingID:   "2024-Feb-15-001",
					wantStatus:  http.StatusForbidden,
				},
				{
					description: "with invalid opening ID",
					token:       adminToken,
					openingID:   "INVALID-ID",
					wantStatus:  http.StatusNotFound,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					employer.GetOpeningWatchersRequest{OpeningID: tc.openingID},
					"/employer/get-opening-watchers",
					tc.wantStatus,
				)

				if tc.wantStatus == http.StatusOK {
					var watchers []employer.OrgUserShort
					err := json.Unmarshal(resp.([]byte), &watchers)
					Expect(err).ShouldNot(HaveOccurred())

					gotEmails := make([]string, len(watchers))
					for i, watcher := range watchers {
						gotEmails[i] = string(watcher.Email)
					}
					Expect(gotEmails).Should(ConsistOf(tc.wantEmails))
				}
			}
		})

		It("should add watchers correctly", func() {
			type addWatchersTestCase struct {
				description string
				token       string
				request     employer.AddOpeningWatchersRequest
				wantStatus  int
			}

			testCases := []addWatchersTestCase{
				{
					description: "add new watcher to opening",
					token:       adminToken,
					request: employer.AddOpeningWatchersRequest{
						OpeningID: "2024-Mar-06-001",
						Emails: []common.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "add duplicate watcher",
					token:       adminToken,
					request: employer.AddOpeningWatchersRequest{
						OpeningID: "2024-Feb-15-001",
						Emails: []common.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request: employer.AddOpeningWatchersRequest{
						OpeningID: "2024-Feb-15-001",
						Emails: []common.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid opening ID",
					token:       adminToken,
					request: employer.AddOpeningWatchersRequest{
						OpeningID: "INVALID-ID",
						Emails: []common.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "with non-existent user email",
					token:       adminToken,
					request: employer.AddOpeningWatchersRequest{
						OpeningID: "2024-Feb-15-001",
						Emails: []common.EmailAddress{
							"nonexistent@openings0006.example",
						},
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/add-opening-watchers",
					tc.wantStatus,
				)
			}
		})

		It("should remove watchers correctly", func() {
			type removeWatcherTestCase struct {
				description string
				token       string
				request     employer.RemoveOpeningWatcherRequest
				wantStatus  int
			}

			testCases := []removeWatcherTestCase{
				{
					description: "remove existing watcher",
					token:       adminToken,
					request: employer.RemoveOpeningWatcherRequest{
						OpeningID: "2024-Feb-15-001",
						Email:     "watcher1@openings0006.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "remove non-existent watcher",
					token:       adminToken,
					request: employer.RemoveOpeningWatcherRequest{
						OpeningID: "2024-Feb-15-001",
						Email:     "nonexistent@openings0006.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request: employer.RemoveOpeningWatcherRequest{
						OpeningID: "2024-Feb-15-001",
						Email:     "watcher1@openings0006.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid opening ID",
					token:       adminToken,
					request: employer.RemoveOpeningWatcherRequest{
						OpeningID: "INVALID-ID",
						Email:     "watcher1@openings0006.example",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "#### %s\n", tc.description)
				testPOST(
					tc.token,
					tc.request,
					"/employer/remove-opening-watcher",
					tc.wantStatus,
				)
			}

			// Verify final state
			resp := testPOSTGetResp(
				adminToken,
				employer.GetOpeningWatchersRequest{
					OpeningID: "2024-Feb-15-001",
				},
				"/employer/get-opening-watchers",
				http.StatusOK,
			)

			var watchers []employer.OrgUserShort
			err := json.Unmarshal(resp.([]byte), &watchers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(watchers).Should(HaveLen(1))
			Expect(
				string(watchers[0].Email),
			).Should(Equal("watcher2@openings0006.example"))
		})

		It("should not allow more than 25 watchers", func() {
			const maxWatchersAllowed = 25

			request := employer.CreateOpeningRequest{
				Title:             "Test Opening",
				Positions:         1,
				JD:                "Test Job Description",
				Recruiter:         "admin@openings0006.example",
				HiringManager:     "admin@openings0006.example",
				CostCenterName:    "Engineering",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            0,
				YoeMax:            5,
				MinEducationLevel: common.BachelorEducation,
				Salary: &common.Salary{
					MinAmount: 50000,
					MaxAmount: 100000,
					Currency:  "USD",
				},
			}

			resp := testPOSTGetResp(
				adminToken,
				request,
				"/employer/create-opening",
				http.StatusOK,
			).([]byte)

			var opening employer.CreateOpeningResponse
			err := json.Unmarshal(resp, &opening)
			Expect(err).ShouldNot(HaveOccurred())

			openingID := opening.OpeningID

			for i := 1; i <= maxWatchersAllowed; {
				newWatcherCount := rand.Intn(3) + 1

				var newWatchers []common.EmailAddress
				for j := 0; j < newWatcherCount && i <= maxWatchersAllowed; {
					newWatchers = append(newWatchers, common.EmailAddress(
						fmt.Sprintf("maxwatcher%d@openings0006.example", i),
					))
					j++
					i++
				}
				fmt.Fprintf(GinkgoWriter, "adding %v watchers\n", newWatchers)
				testPOST(
					adminToken,
					employer.AddOpeningWatchersRequest{
						OpeningID: openingID,
						Emails:    newWatchers,
					},
					"/employer/add-opening-watchers",
					http.StatusOK,
				)
			}

			// Try to add one more watcher
			testPOST(
				adminToken,
				employer.AddOpeningWatchersRequest{
					OpeningID: openingID,
					Emails: []common.EmailAddress{
						"maxwatcher26@openings0006.example",
					},
				},
				"/employer/add-opening-watchers",
				http.StatusUnprocessableEntity,
			)
		})
	})
})
