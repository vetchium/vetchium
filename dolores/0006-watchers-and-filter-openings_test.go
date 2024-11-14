package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

var _ = FDescribe("Openings", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, nonOpeningsToken string
	var recruiterToken, hiringManagerToken string

	BeforeAll(func() {
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
			request     vetchi.FilterOpeningsRequest
			wantStatus  int
			wantCount   int
			wantIDs     []string
		}

		It("should filter openings correctly", func() {
			testCases := []filterOpeningsTestCase{
				{
					description: "with no filters",
					token:       adminToken,
					request:     vetchi.FilterOpeningsRequest{},
					wantStatus:  http.StatusOK,
					wantCount:   4,
					wantIDs: []string{
						"OPENING-001",
						"OPENING-002",
						"OPENING-003",
						"OPENING-004",
					},
				},
				{
					description: "with state filter - draft",
					token:       adminToken,
					request: vetchi.FilterOpeningsRequest{
						State: []vetchi.OpeningState{vetchi.DraftOpening},
					},
					wantStatus: http.StatusOK,
					wantCount:  1,
					wantIDs:    []string{"OPENING-001"},
				},
				{
					description: "with date range filter",
					token:       adminToken,
					request: vetchi.FilterOpeningsRequest{
						FromDate: func() *time.Time {
							t := time.Now().AddDate(0, 0, -15)
							return &t
						}(),
						ToDate: func() *time.Time {
							t := time.Now()
							return &t
						}(),
					},
					wantStatus: http.StatusOK,
					wantCount:  2,
					wantIDs:    []string{"OPENING-003", "OPENING-004"},
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request:     vetchi.FilterOpeningsRequest{},
					wantStatus:  http.StatusUnauthorized,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request:     vetchi.FilterOpeningsRequest{},
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
					var openings []vetchi.OpeningInfo
					err := json.Unmarshal(resp.([]byte), &openings)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(openings).Should(HaveLen(tc.wantCount))

					gotIDs := make([]string, len(openings))
					for i, opening := range openings {
						gotIDs[i] = opening.ID
					}
					Expect(gotIDs).Should(ConsistOf(tc.wantIDs))
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
					openingID:   "OPENING-001",
					wantStatus:  http.StatusOK,
					wantEmails: []string{
						"watcher1@openings0006.example",
						"watcher2@openings0006.example",
					},
				},
				{
					description: "get watchers for opening with single watcher",
					token:       adminToken,
					openingID:   "OPENING-002",
					wantStatus:  http.StatusOK,
					wantEmails:  []string{"watcher1@openings0006.example"},
				},
				{
					description: "get watchers for opening with no watchers",
					token:       adminToken,
					openingID:   "OPENING-003",
					wantStatus:  http.StatusOK,
					wantEmails:  []string{},
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					openingID:   "OPENING-001",
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
					vetchi.GetOpeningWatchersRequest{OpeningID: tc.openingID},
					"/employer/get-opening-watchers",
					tc.wantStatus,
				)

				if tc.wantStatus == http.StatusOK {
					var watchers []vetchi.OrgUserShort
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
				request     vetchi.AddOpeningWatchersRequest
				wantStatus  int
			}

			testCases := []addWatchersTestCase{
				{
					description: "add new watcher to opening",
					token:       adminToken,
					request: vetchi.AddOpeningWatchersRequest{
						OpeningID: "OPENING-003",
						Emails: []vetchi.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "add duplicate watcher",
					token:       adminToken,
					request: vetchi.AddOpeningWatchersRequest{
						OpeningID: "OPENING-001",
						Emails: []vetchi.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request: vetchi.AddOpeningWatchersRequest{
						OpeningID: "OPENING-001",
						Emails: []vetchi.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid opening ID",
					token:       adminToken,
					request: vetchi.AddOpeningWatchersRequest{
						OpeningID: "INVALID-ID",
						Emails: []vetchi.EmailAddress{
							"watcher1@openings0006.example",
						},
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "with non-existent user email",
					token:       adminToken,
					request: vetchi.AddOpeningWatchersRequest{
						OpeningID: "OPENING-001",
						Emails: []vetchi.EmailAddress{
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
				request     vetchi.RemoveOpeningWatcherRequest
				wantStatus  int
			}

			testCases := []removeWatcherTestCase{
				{
					description: "remove existing watcher",
					token:       adminToken,
					request: vetchi.RemoveOpeningWatcherRequest{
						OpeningID: "OPENING-001",
						Email:     "watcher1@openings0006.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "remove non-existent watcher",
					token:       adminToken,
					request: vetchi.RemoveOpeningWatcherRequest{
						OpeningID: "OPENING-001",
						Email:     "nonexistent@openings0006.example",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "with non-openings role",
					token:       nonOpeningsToken,
					request: vetchi.RemoveOpeningWatcherRequest{
						OpeningID: "OPENING-001",
						Email:     "watcher1@openings0006.example",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "with invalid opening ID",
					token:       adminToken,
					request: vetchi.RemoveOpeningWatcherRequest{
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
				vetchi.GetOpeningWatchersRequest{OpeningID: "OPENING-001"},
				"/employer/get-opening-watchers",
				http.StatusOK,
			)

			var watchers []vetchi.OrgUserShort
			err := json.Unmarshal(resp.([]byte), &watchers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(watchers).Should(HaveLen(1))
			Expect(
				string(watchers[0].Email),
			).Should(Equal("watcher2@openings0006.example"))
		})
	})
})
