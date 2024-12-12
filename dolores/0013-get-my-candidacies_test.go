package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get My Candidacies", Ordered, func() {
	var db *pgxpool.Pool
	var hubUser1Token, hubUser2Token, hubUser3Token string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0013-get-my-candidacies-up.pgsql")

		// Login hub users and get their tokens
		var wg sync.WaitGroup
		tokens := map[string]*string{
			"hubuser1@my-candidacies.example": &hubUser1Token,
			"hubuser2@my-candidacies.example": &hubUser2Token,
			"hubuser3@my-candidacies.example": &hubUser3Token,
		}

		for email, token := range tokens {
			wg.Add(1)
			hubSigninAsync(email, "NewPassword123$", token, &wg)
		}
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(db, "0013-get-my-candidacies-down.pgsql")
		db.Close()
	})

	Describe("Get My Candidacies Tests", func() {
		type getMyCandidaciesTestCase struct {
			description string
			token       string
			request     hub.MyCandidaciesRequest
			wantStatus  int
			validate    func([]hub.MyCandidacy)
		}

		It("should handle various test cases correctly", func() {
			testCases := []getMyCandidaciesTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.MyCandidaciesRequest{
						Limit: 10,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.MyCandidaciesRequest{
						Limit: 10,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with valid token but no candidacies",
					token:       hubUser3Token,
					request: hub.MyCandidaciesRequest{
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						Expect(candidacies).Should(BeEmpty())
					},
				},
				{
					description: "get all candidacies for user1",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						Limit: 50,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						Expect(candidacies).Should(HaveLen(15))
						// Verify candidacy details
						for _, c := range candidacies {
							Expect(c.OpeningID).ShouldNot(BeEmpty())
							Expect(c.OpeningTitle).ShouldNot(BeEmpty())
							Expect(
								c.CompanyDomain,
							).Should(MatchRegexp(`my-candidacies-\d+\.example`))
						}
					},
				},
				{
					description: "get candidacies with pagination",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						Limit: 5,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						Expect(candidacies).Should(HaveLen(5))
					},
				},
				{
					description: "filter by state - INTERVIEWING",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						CandidacyStates: []common.CandidacyState{
							common.InterviewingCandidacyState,
						},
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						for _, c := range candidacies {
							Expect(
								c.CandidacyState,
							).Should(Equal("INTERVIEWING"))
						}
					},
				},
				{
					description: "filter by state - OFFERED",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						CandidacyStates: []common.CandidacyState{
							common.OfferedCandidacyState,
						},
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						for _, c := range candidacies {
							Expect(c.CandidacyState).Should(Equal("OFFERED"))
						}
					},
				},
				{
					description: "filter by state - CANDIDATE_UNSUITABLE",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						CandidacyStates: []common.CandidacyState{
							common.CandidateUnsuitableCandidacyState,
						},
						Limit: 10,
					},
					wantStatus: http.StatusOK,
					validate: func(candidacies []hub.MyCandidacy) {
						for _, c := range candidacies {
							Expect(
								c.CandidacyState,
							).Should(Equal("CANDIDATE_UNSUITABLE"))
						}
					},
				},
				{
					description: "invalid pagination key",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						PaginationKey: strptr("invalid-key"),
						Limit:         10,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "invalid limit - too high",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						Limit: 1001,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "invalid limit - zero",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						Limit: 0,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "invalid limit - negative",
					token:       hubUser1Token,
					request: hub.MyCandidaciesRequest{
						Limit: -1,
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "Testing: %s\n", tc.description)

				reqBody, err := json.Marshal(tc.request)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/hub/get-my-candidacies",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())

				if tc.token != "" {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(tc.wantStatus))

				if tc.wantStatus == http.StatusOK && tc.validate != nil {
					var candidacies []hub.MyCandidacy
					err = json.NewDecoder(resp.Body).Decode(&candidacies)
					Expect(err).ShouldNot(HaveOccurred())
					tc.validate(candidacies)
				}
			}
		})

		It("should handle pagination correctly", func() {
			allCandidacies := make([]hub.MyCandidacy, 0)
			paginationKey := ""
			pageSize := 3

			for {
				request := hub.MyCandidaciesRequest{
					PaginationKey: &paginationKey,
					Limit:         pageSize,
				}

				reqBody, err := json.Marshal(request)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/hub/get-my-candidacies",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+hubUser1Token)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				var candidacies []hub.MyCandidacy
				err = json.NewDecoder(resp.Body).Decode(&candidacies)
				Expect(err).ShouldNot(HaveOccurred())

				if len(candidacies) == 0 {
					break
				}

				allCandidacies = append(allCandidacies, candidacies...)

				if len(candidacies) < pageSize {
					break
				}

				paginationKey = candidacies[len(candidacies)-1].OpeningID
			}

			// Verify we got all candidacies and they're unique
			Expect(allCandidacies).Should(HaveLen(15))
			seenIDs := make(map[string]bool)
			for _, c := range allCandidacies {
				Expect(seenIDs[c.OpeningID]).Should(BeFalse())
				seenIDs[c.OpeningID] = true
			}
		})
	})
})
