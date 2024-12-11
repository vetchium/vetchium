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

var _ = FDescribe("Candidacy Comments", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, crudToken, viewerToken, disabledToken string
	var activeHubToken, disabledHubToken, deletedHubToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0011-candidacy-comments-up.pgsql")

		// Get employer tokens
		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@candidacy-comments.example":    &adminToken,
			"crud@candidacy-comments.example":     &crudToken,
			"viewer@candidacy-comments.example":   &viewerToken,
			"disabled@candidacy-comments.example": &disabledToken,
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
		disabledHubToken = hubSignin(
			"0011-disabled@hub.example",
			"NewPassword123$",
		)
		deletedHubToken = hubSignin(
			"0011-deleted@hub.example",
			"NewPassword123$",
		)
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
					description: "crud user can add comment",
					token:       crudToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "CRUD user comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "viewer can add comment",
					token:       viewerToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Viewer comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusOK,
				},
				{
					description: "disabled user cannot add comment",
					token:       disabledToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Disabled user comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid candidacy ID",
					token:       adminToken,
					request: employer.AddEmployerCandidacyCommentRequest{
						CandidacyID: "invalid-id",
						Comment:     "Invalid candidacy comment",
					},
					endpoint:   "/employer/add-candidacy-comment",
					wantStatus: http.StatusBadRequest,
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
				fmt.Fprintf(GinkgoWriter, "Testing: %s\n", tc.description)
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
					description: "disabled hub user cannot add comment",
					token:       disabledHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Disabled hub user comment",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "deleted hub user cannot add comment",
					token:       deletedHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: validCandidacyID,
						Comment:     "Deleted hub user comment",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid candidacy ID",
					token:       activeHubToken,
					request: hub.AddHubCandidacyCommentRequest{
						CandidacyID: "invalid-id",
						Comment:     "Invalid candidacy comment",
					},
					endpoint:   "/hub/add-candidacy-comment",
					wantStatus: http.StatusBadRequest,
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
				fmt.Fprintf(GinkgoWriter, "Testing: %s\n", tc.description)
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
					description: "crud user can get comments",
					token:       crudToken,
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
					description: "viewer can get comments",
					token:       viewerToken,
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
					description: "disabled user cannot get comments",
					token:       disabledToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid candidacy ID",
					token:       adminToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: "invalid-id",
					},
					endpoint:   "/employer/get-candidacy-comments",
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "Testing: %s\n", tc.description)
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
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(BeNumerically(">", 0))
					},
				},
				{
					description: "disabled hub user cannot get comments",
					token:       disabledHubToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/hub/get-candidacy-comments",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "deleted hub user cannot get comments",
					token:       deletedHubToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: validCandidacyID,
					},
					endpoint:   "/hub/get-candidacy-comments",
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "invalid candidacy ID",
					token:       activeHubToken,
					request: common.GetCandidacyCommentsRequest{
						CandidacyID: "invalid-id",
					},
					endpoint:   "/hub/get-candidacy-comments",
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(GinkgoWriter, "Testing: %s\n", tc.description)
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
