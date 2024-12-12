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
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = Describe("Candidacy Comments", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, hiringManagerToken, recruiterToken, watcherToken, regularUserToken string
	var activeHubToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0011-candidacy-comments-up.pgsql")

		// Get employer tokens
		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@candidacy-comments.example":         &adminToken,
			"hiringmanager@candidacy-comments.example": &hiringManagerToken,
			"recruiter@candidacy-comments.example":     &recruiterToken,
			"watcher@candidacy-comments.example":       &watcherToken,
			"regular@candidacy-comments.example":       &regularUserToken,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"candidacy-comments.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()

		// Get hub user tokens
		activeHubToken = hubSignin("0011-active@hub.example", "NewPassword123$")
	})

	AfterAll(func() {
		seedDatabase(db, "0011-candidacy-comments-down.pgsql")
		db.Close()
	})

	var _ = Describe("Add Candidacy Comments", func() {
		type addCommentTestCase struct {
			description string
			token       string
			request     interface{}
			endpoint    string
			wantStatus  int
		}

		It("should handle employer comments correctly", func() {
			validCandidacyID := "12345678-0011-0011-0011-000000060001"

			testCases := []addCommentTestCase{
				{
					description: "admin can add comment",
					token:       adminToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Admin comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "hiring manager can add comment",
					token:       hiringManagerToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Hiring manager comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "recruiter can add comment",
					token:       recruiterToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Recruiter comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "watcher can add comment",
					token:       watcherToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Watcher comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "regular user cannot add comment",
					token:       regularUserToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Regular user comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusForbidden,
				},
				{
					description: "invalid candidacy ID",
					token:       adminToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: "invalid-id",
						Comment:     "Invalid candidacy comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusForbidden,
				},
				{
					description: "empty comment",
					token:       adminToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(tc.token, tc.request, tc.endpoint, tc.wantStatus)
			}
		})

		It("should handle hub user comments correctly", func() {
			validCandidacyID := "12345678-0011-0011-0011-000000060001"

			testCases := []addCommentTestCase{
				{
					description: "active hub user can add comment",
					token:       activeHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Active hub user comment",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "invalid candidacy ID",
					token:       activeHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: "invalid-id",
						Comment:     "Invalid candidacy comment",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusForbidden,
				},
				{
					description: "empty comment",
					token:       activeHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				testPOST(tc.token, tc.request, tc.endpoint, tc.wantStatus)
			}
		})
	})

	var _ = Describe("Get Candidacy Comments", func() {
		type getCommentsTestCase struct {
			description string
			token       string
			request     interface{}
			endpoint    string
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle employer comment retrieval correctly", func() {
			validCandidacyID := "12345678-0011-0011-0011-000000060001"

			testCases := []getCommentsTestCase{
				{
					description: "admin can get comments",
					token:       adminToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "hiring manager can get comments",
					token:       hiringManagerToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "regular user can get comments",
					token:       regularUserToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "invalid candidacy ID",
					token:       adminToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: "invalid-id",
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeZero())
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					tc.endpoint,
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})

		It("should handle hub user comment retrieval correctly", func() {
			validCandidacyID := "12345678-0011-0011-0011-000000060001"

			testCases := []getCommentsTestCase{
				{
					description: "active hub user can get comments",
					token:       activeHubToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/hub/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						fmt.Fprintf(
							GinkgoWriter,
							"Response Body: %s\n",
							string(resp),
						)
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "invalid candidacy ID",
					token:       activeHubToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: "invalid-id",
					},
					endpoint:   "/hub/get-candidacy-comments",
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []common.CandidacyComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeZero())
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "### Testing: %s\n", tc.description)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					tc.endpoint,
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
