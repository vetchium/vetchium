package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

var _ = Describe("Get Candidacies", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, recruiter1Token, recruiter2Token, viewerToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0012-filter-candidacies-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@filter-candidacy-infos.example":      &adminToken,
			"recruiter1@filter-candidacy-infos.example": &recruiter1Token,
			"recruiter2@filter-candidacy-infos.example": &recruiter2Token,
			"viewer@filter-candidacy-infos.example":     &viewerToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"filter-candidacy-infos.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0012-filter-candidacies-down.pgsql")
		db.Close()
	})

	Describe("Get Candidacy Infos", func() {
		type getCandidaciesTestCase struct {
			description string
			token       string
			request     employer.FilterCandidacyInfosRequest
			wantStatus  int
			validate    func([]employer.Candidacy)
		}

		It("should handle various get candidacies requests correctly", func() {
			testCases := []getCandidaciesTestCase{
				{
					description: "without a session token",
					token:       "",
					request: employer.FilterCandidacyInfosRequest{
						Limit: 10,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid session token",
					token:       "invalid-token",
					request: employer.FilterCandidacyInfosRequest{
						Limit: 10,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with viewer role (should be allowed)",
					token:       viewerToken,
					request: employer.FilterCandidacyInfosRequest{
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						Expect(candidacies).Should(HaveLen(3))
					},
				},
				{
					description: "filter by recruiter1",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						RecruiterEmail: strptr(
							"recruiter1@filter-candidacy-infos.example",
						),
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						Expect(candidacies).Should(HaveLen(2))
						for _, c := range candidacies {
							Expect(c.OpeningID).Should(Or(
								Equal("2024-Mar-01-001"),
							))
						}
					},
				},
				{
					description: "filter by recruiter2",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						RecruiterEmail: strptr(
							"recruiter2@filter-candidacy-infos.example",
						),
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						Expect(candidacies).Should(HaveLen(1))
						Expect(
							candidacies[0].OpeningID,
						).Should(Equal("2024-Mar-01-002"))
					},
				},
				{
					description: "filter by state INTERVIEWING",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						State: strptr(
							string(common.InterviewingCandidacyState),
						),
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						Expect(candidacies).Should(HaveLen(2))
						for _, c := range candidacies {
							Expect(
								c.CandidacyState,
							).Should(Equal(common.InterviewingCandidacyState))
						}
					},
				},
				{
					description: "filter by state OFFERED",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						State: strptr(string(common.OfferedCandidacyState)),
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						Expect(candidacies).Should(HaveLen(1))
						Expect(
							candidacies[0].CandidacyState,
						).Should(Equal(common.OfferedCandidacyState))
					},
				},
				{
					description: "with invalid limit (negative)",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						Limit: -1,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with invalid limit (too large)",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						Limit: 100,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "with pagination",
					token:       adminToken,
					request: employer.FilterCandidacyInfosRequest{
						Limit: 2,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []employer.Candidacy) {
						fmt.Fprintf(GinkgoWriter, "page1: %+v\n", candidacies)
						Expect(candidacies).Should(HaveLen(2))

						paginationKey := &candidacies[1].CandidacyID
						fmt.Fprintf(GinkgoWriter, "pkey: %s\n", *paginationKey)

						// Get next page
						nextPageReq := employer.FilterCandidacyInfosRequest{
							PaginationKey: paginationKey,
							Limit:         2,
						}
						resp := testPOSTGetResp(
							adminToken,
							nextPageReq,
							"/employer/filter-candidacy-infos",
							http.StatusOK,
						).([]byte)

						var nextPage []employer.Candidacy
						err := json.Unmarshal(resp, &nextPage)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(nextPage).Should(HaveLen(1))

						fmt.Fprintf(GinkgoWriter, "page2: %+v\n", nextPage)

						// Verify no duplicates
						firstPageIDs := make(map[string]bool)
						for _, c := range candidacies {
							firstPageIDs[c.CandidacyID] = true
						}
						for _, c := range nextPage {
							Expect(
								firstPageIDs[c.CandidacyID],
							).Should(BeFalse())
						}
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Test case: %s\n", tc.description)

				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/employer/filter-candidacy-infos",
					tc.wantStatus,
				)

				if tc.wantStatus == http.StatusOK {
					var candidacies []employer.Candidacy
					err := json.Unmarshal(resp.([]byte), &candidacies)
					Expect(err).ShouldNot(HaveOccurred())

					if tc.validate != nil {
						tc.validate(candidacies)
					}
				}
			}
		})
	})
})
