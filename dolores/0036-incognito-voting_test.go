package dolores

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Incognito Voting API", Ordered, func() {
	var (
		pool *pgxpool.Pool

		// User tokens for dedicated test users
		user0036Token1 string
		user0036Token2 string
		user0036Token3 string
		user0036Token4 string
		user0036Token5 string
		user0036Token6 string
		user0036Token7 string
		user0036Token8 string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0036-incognito-voting-up.pgsql")

		var wg sync.WaitGroup
		wg.Add(8)
		hubSigninAsync(
			"user0036-1@0036-test.com",
			"NewPassword123$",
			&user0036Token1,
			&wg,
		)
		hubSigninAsync(
			"user0036-2@0036-test.com",
			"NewPassword123$",
			&user0036Token2,
			&wg,
		)
		hubSigninAsync(
			"user0036-3@0036-test.com",
			"NewPassword123$",
			&user0036Token3,
			&wg,
		)
		hubSigninAsync(
			"user0036-4@0036-test.com",
			"NewPassword123$",
			&user0036Token4,
			&wg,
		)
		hubSigninAsync(
			"user0036-5@0036-test.com",
			"NewPassword123$",
			&user0036Token5,
			&wg,
		)
		hubSigninAsync(
			"user0036-6@0036-test.com",
			"NewPassword123$",
			&user0036Token6,
			&wg,
		)
		hubSigninAsync(
			"user0036-7@0036-test.com",
			"NewPassword123$",
			&user0036Token7,
			&wg,
		)
		hubSigninAsync(
			"user0036-8@0036-test.com",
			"NewPassword123$",
			&user0036Token8,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(pool, "0036-incognito-voting-down.pgsql")
		pool.Close()
	})

	Describe("Incognito Post Voting", func() {
		It("should upvote a post successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for upvoting.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0036Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			voteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}

			testPOST(
				user0036Token2,
				voteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				user0036Token2,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.MeUpvoted).Should(BeTrue())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(1)))
		})

		It("should downvote a post successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for downvoting.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0036Token3,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			voteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}

			testPOST(
				user0036Token4,
				voteReq,
				"/hub/downvote-incognito-post",
				http.StatusOK,
			)

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				user0036Token4,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.MeDownvoted).Should(BeTrue())
			Expect(getResponse.DownvotesCount).Should(Equal(int32(1)))
		})

		It("should unvote a post successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for unvoting.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0036Token5,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token6,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token6,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusOK,
			)

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				user0036Token6,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.MeUpvoted).Should(BeFalse())
			Expect(getResponse.MeDownvoted).Should(BeFalse())
		})

		It(
			"should return 422 for vote conflict - upvote after downvote",
			func() {
				postReq := hub.AddIncognitoPostRequest{
					Content: "Test post for vote conflict.",
					TagIDs:  []common.VTagID{"startups"},
				}
				postResp := testPOSTGetResp(
					user0036Token7,
					postReq,
					"/hub/add-incognito-post",
					http.StatusOK,
				)
				var postResponse hub.AddIncognitoPostResponse
				err := json.Unmarshal(postResp.([]byte), &postResponse)
				Expect(err).ShouldNot(HaveOccurred())

				downvoteReq := hub.DownvoteIncognitoPostRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
				}
				testPOST(
					user0036Token8,
					downvoteReq,
					"/hub/downvote-incognito-post",
					http.StatusOK,
				)

				upvoteReq := hub.UpvoteIncognitoPostRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
				}
				testPOST(
					user0036Token8,
					upvoteReq,
					"/hub/upvote-incognito-post",
					http.StatusUnprocessableEntity,
				)
			},
		)

		It(
			"should return 422 for vote conflict - downvote after upvote",
			func() {
				postReq := hub.AddIncognitoPostRequest{
					Content: "Test post for downvote conflict.",
					TagIDs:  []common.VTagID{"technology"},
				}
				postResp := testPOSTGetResp(
					user0036Token1,
					postReq,
					"/hub/add-incognito-post",
					http.StatusOK,
				)
				var postResponse hub.AddIncognitoPostResponse
				err := json.Unmarshal(postResp.([]byte), &postResponse)
				Expect(err).ShouldNot(HaveOccurred())

				upvoteReq := hub.UpvoteIncognitoPostRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
				}
				testPOST(
					user0036Token2,
					upvoteReq,
					"/hub/upvote-incognito-post",
					http.StatusOK,
				)

				downvoteReq := hub.DownvoteIncognitoPostRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
				}
				testPOST(
					user0036Token2,
					downvoteReq,
					"/hub/downvote-incognito-post",
					http.StatusUnprocessableEntity,
				)
			},
		)

		It("should return 422 when trying to vote on own post", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for own voting restriction.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0036Token3,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token3,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusUnprocessableEntity,
			)

			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token3,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusUnprocessableEntity,
			)
		})

		It("should return 422 when trying to unvote own post", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for own unvoting restriction.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0036Token4,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token4,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusUnprocessableEntity,
			)
		})

		It("should be idempotent for same vote", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for idempotent voting.",
				TagIDs:  []common.VTagID{"startups"},
			}
			postResp := testPOSTGetResp(
				user0036Token5,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}

			testPOST(
				user0036Token6,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			testPOST(
				user0036Token6,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				user0036Token6,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(1)))
		})

		It("should fail without authentication", func() {
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: "some-post-id",
			}

			testPOST(
				"",
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: "nonexistent-post-id",
			}

			testPOST(
				user0036Token7,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusNotFound,
			)
		})

		It("should return 404 when voting on deleted post", func() {
			// User 1 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post to be deleted",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0036Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 1 deletes their own post
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token1,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusOK,
			)

			// User 2 tries to vote on the deleted post (should return 404)
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token2,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusNotFound,
			)

			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token2,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusNotFound,
			)

			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
			}
			testPOST(
				user0036Token2,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusNotFound,
			)
		})
	})

	Describe("Incognito Comment Voting", func() {
		It("should upvote a comment successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment upvoting.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0036Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for upvoting.",
			}
			commentResp := testPOSTGetResp(
				user0036Token2,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			voteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				user0036Token3,
				voteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should downvote a comment successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment downvoting.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0036Token4,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for downvoting.",
			}
			commentResp := testPOSTGetResp(
				user0036Token5,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			voteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				user0036Token6,
				voteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should unvote a comment successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment unvoting.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0036Token7,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for unvoting.",
			}
			commentResp := testPOSTGetResp(
				user0036Token8,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token1,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token1,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should return 422 for comment vote conflict", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment vote conflict.",
				TagIDs:  []common.VTagID{"startups"},
			}
			postResp := testPOSTGetResp(
				user0036Token2,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for vote conflict.",
			}
			commentResp := testPOSTGetResp(
				user0036Token3,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token4,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusOK,
			)

			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token4,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)
		})

		It("should return 422 when trying to vote on own comment", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for own comment voting restriction.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0036Token5,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for own voting restriction.",
			}
			commentResp := testPOSTGetResp(
				user0036Token6,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token6,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)

			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0036Token6,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)
		})

		It("should fail without authentication", func() {
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: "some-post-id",
				CommentID:       "some-comment-id",
			}

			testPOST(
				"",
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnauthorized,
			)
		})

		It("should return 422 for invalid comment ID", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for invalid comment ID.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0036Token7,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       "nonexistent-comment-id",
			}

			testPOST(
				user0036Token8,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)
		})
	})
})
