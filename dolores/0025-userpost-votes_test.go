package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = FDescribe("User Post Votes", Ordered, func() {
	var db *pgxpool.Pool
	var voter1Token, voter2Token, authorToken string
	var post1ID, post2ID, post3ID string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0025-userpost-votes-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(3)
		hubSigninAsync(
			"voter1@0025-votes.example.com",
			"NewPassword123$",
			&voter1Token,
			&wg,
		)
		hubSigninAsync(
			"voter2@0025-votes.example.com",
			"NewPassword123$",
			&voter2Token,
			&wg,
		)
		hubSigninAsync(
			"author@0025-votes.example.com",
			"NewPassword123$",
			&authorToken,
			&wg,
		)
		wg.Wait()

		// Create test posts using the API
		post1Resp := testPOSTGetResp(
			authorToken,
			hub.AddPostRequest{
				Content: "Test post 1 for voting",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post1AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post1Resp, &post1AddResp)).To(Succeed())
		post1ID = post1AddResp.PostID
		fmt.Fprintf(GinkgoWriter, "Post 1 ID: %s\n", post1ID)

		post2Resp := testPOSTGetResp(
			authorToken,
			hub.AddPostRequest{
				Content: "Test post 2 for voting",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post2AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post2Resp, &post2AddResp)).To(Succeed())
		post2ID = post2AddResp.PostID

		post3Resp := testPOSTGetResp(
			authorToken,
			hub.AddPostRequest{
				Content: "Test post 3 for voting",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post3AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post3Resp, &post3AddResp)).To(Succeed())
		post3ID = post3AddResp.PostID

		// Create initial votes
		testPOSTGetResp(
			voter1Token,
			hub.UpvoteUserPostRequest{PostID: post1ID},
			"/hub/upvote-user-post",
			http.StatusOK,
		)
		testPOSTGetResp(
			voter1Token,
			hub.DownvoteUserPostRequest{PostID: post2ID},
			"/hub/downvote-user-post",
			http.StatusOK,
		)
	})

	AfterAll(func() {
		seedDatabase(db, "0025-userpost-votes-down.pgsql")
		db.Close()
	})

	Describe("Upvote User Post", func() {
		type upvoteTestCase struct {
			description string
			token       string
			request     hub.UpvoteUserPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various upvote scenarios", func() {
			testCases := []upvoteTestCase{
				{
					description: "attempt to upvote own post",
					token:       authorToken,
					request: hub.UpvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "without authentication",
					token:       "",
					request: hub.UpvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.UpvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "upvote non-existent post",
					token:       voter1Token,
					request: hub.UpvoteUserPostRequest{
						PostID: "non-existent-post",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "upvote already upvoted post (should succeed)",
					token:       voter1Token,
					request: hub.UpvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						// Get post details to verify voting fields
						detailsResp := testPOSTGetResp(
							voter1Token,
							hub.GetPostDetailsRequest{PostID: post1ID},
							"/hub/get-post-details",
							http.StatusOK,
						).([]byte)

						var post hub.Post
						err := json.Unmarshal(detailsResp, &post)
						Expect(err).ShouldNot(HaveOccurred())

						// Verify voting fields after upvote
						Expect(
							post.MeUpvoted,
						).Should(BeTrue(), "Post should show as upvoted")
						Expect(
							post.MeDownvoted,
						).Should(BeFalse(), "Post should not show as downvoted")
						Expect(
							post.CanUpvote,
						).Should(BeFalse(), "Should not be able to upvote again")
						Expect(
							post.CanDownvote,
						).Should(BeFalse(), "Should not be able to downvote after upvoting")
					},
				},
				{
					description: "upvote already downvoted post (should fail)",
					token:       voter1Token,
					request: hub.UpvoteUserPostRequest{
						PostID: post2ID,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "new upvote on unvoted post",
					token:       voter2Token,
					request: hub.UpvoteUserPostRequest{
						PostID: post3ID,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/upvote-user-post",
					tc.wantStatus,
				)
				Expect(resp).ToNot(BeNil())
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Downvote User Post", func() {
		type downvoteTestCase struct {
			description string
			token       string
			request     hub.DownvoteUserPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various downvote scenarios", func() {
			testCases := []downvoteTestCase{
				{
					description: "attempt to downvote own post",
					token:       authorToken,
					request: hub.DownvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "without authentication",
					token:       "",
					request: hub.DownvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.DownvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "downvote non-existent post",
					token:       voter1Token,
					request: hub.DownvoteUserPostRequest{
						PostID: "non-existent-post",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "downvote already downvoted post (should succeed)",
					token:       voter1Token,
					request: hub.DownvoteUserPostRequest{
						PostID: post2ID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						// Get post details to verify voting fields
						detailsResp := testPOSTGetResp(
							voter1Token,
							hub.GetPostDetailsRequest{PostID: post2ID},
							"/hub/get-post-details",
							http.StatusOK,
						).([]byte)

						var post hub.Post
						err := json.Unmarshal(detailsResp, &post)
						Expect(err).ShouldNot(HaveOccurred())

						// Verify voting fields after downvote
						Expect(
							post.MeUpvoted,
						).Should(BeFalse(), "Post should not show as upvoted")
						Expect(
							post.MeDownvoted,
						).Should(BeTrue(), "Post should show as downvoted")
						Expect(
							post.CanUpvote,
						).Should(BeFalse(), "Should not be able to upvote after downvoting")
						Expect(
							post.CanDownvote,
						).Should(BeFalse(), "Should not be able to downvote again")
					},
				},
				{
					description: "downvote already upvoted post (should fail)",
					token:       voter1Token,
					request: hub.DownvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "new downvote on unvoted post",
					token:       voter2Token,
					request: hub.DownvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/downvote-user-post",
					tc.wantStatus,
				)
				Expect(resp).ToNot(BeNil())
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Unvote User Post", func() {
		type unvoteTestCase struct {
			description string
			token       string
			request     hub.UnvoteUserPostRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various unvote scenarios", func() {
			testCases := []unvoteTestCase{
				{
					description: "attempt to unvote own post",
					token:       authorToken,
					request: hub.UnvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "without authentication",
					token:       "",
					request: hub.UnvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.UnvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "unvote non-existent post",
					token:       voter1Token,
					request: hub.UnvoteUserPostRequest{
						PostID: "non-existent-post",
					},
					wantStatus: http.StatusUnprocessableEntity,
				},
				{
					description: "unvote upvoted post",
					token:       voter1Token,
					request: hub.UnvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusOK,
					validate: func(respBody []byte) {
						// Get post details to verify voting fields
						detailsResp := testPOSTGetResp(
							voter1Token,
							hub.GetPostDetailsRequest{PostID: post1ID},
							"/hub/get-post-details",
							http.StatusOK,
						).([]byte)

						var post hub.Post
						err := json.Unmarshal(detailsResp, &post)
						Expect(err).ShouldNot(HaveOccurred())

						// Verify voting fields after unvote
						Expect(
							post.MeUpvoted,
						).Should(BeFalse(), "Post should not show as upvoted after unvote")
						Expect(
							post.MeDownvoted,
						).Should(BeFalse(), "Post should not show as downvoted")
						Expect(
							post.CanUpvote,
						).Should(BeTrue(), "Should be able to upvote after unvote")
						Expect(
							post.CanDownvote,
						).Should(BeTrue(), "Should be able to downvote after unvote")
					},
				},
				{
					description: "unvote downvoted post",
					token:       voter1Token,
					request: hub.UnvoteUserPostRequest{
						PostID: post2ID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "unvote already unvoted post",
					token:       voter2Token,
					request: hub.UnvoteUserPostRequest{
						PostID: post1ID,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/unvote-user-post",
					tc.wantStatus,
				)
				Expect(resp).ToNot(BeNil())
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})
})
