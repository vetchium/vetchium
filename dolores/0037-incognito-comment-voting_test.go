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

var _ = Describe("Incognito Comment Voting API", Ordered, func() {
	var (
		pool *pgxpool.Pool

		// User tokens for dedicated test users
		user0037Token1  string
		user0037Token2  string
		user0037Token3  string
		user0037Token4  string
		user0037Token5  string
		user0037Token6  string
		user0037Token7  string
		user0037Token8  string
		user0037Token9  string
		user0037Token10 string
		user0037Token11 string
		user0037Token12 string
		user0037Token13 string
		user0037Token14 string
		user0037Token15 string
		user0037Token16 string
		user0037Token17 string
		user0037Token18 string
		user0037Token19 string
		user0037Token20 string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0037-incognito-comment-voting-up.pgsql")

		var wg sync.WaitGroup
		wg.Add(20)
		hubSigninAsync(
			"user0037-1@0037-test.com",
			"NewPassword123$",
			&user0037Token1,
			&wg,
		)
		hubSigninAsync(
			"user0037-2@0037-test.com",
			"NewPassword123$",
			&user0037Token2,
			&wg,
		)
		hubSigninAsync(
			"user0037-3@0037-test.com",
			"NewPassword123$",
			&user0037Token3,
			&wg,
		)
		hubSigninAsync(
			"user0037-4@0037-test.com",
			"NewPassword123$",
			&user0037Token4,
			&wg,
		)
		hubSigninAsync(
			"user0037-5@0037-test.com",
			"NewPassword123$",
			&user0037Token5,
			&wg,
		)
		hubSigninAsync(
			"user0037-6@0037-test.com",
			"NewPassword123$",
			&user0037Token6,
			&wg,
		)
		hubSigninAsync(
			"user0037-7@0037-test.com",
			"NewPassword123$",
			&user0037Token7,
			&wg,
		)
		hubSigninAsync(
			"user0037-8@0037-test.com",
			"NewPassword123$",
			&user0037Token8,
			&wg,
		)
		hubSigninAsync(
			"user0037-9@0037-test.com",
			"NewPassword123$",
			&user0037Token9,
			&wg,
		)
		hubSigninAsync(
			"user0037-10@0037-test.com",
			"NewPassword123$",
			&user0037Token10,
			&wg,
		)
		hubSigninAsync(
			"user0037-11@0037-test.com",
			"NewPassword123$",
			&user0037Token11,
			&wg,
		)
		hubSigninAsync(
			"user0037-12@0037-test.com",
			"NewPassword123$",
			&user0037Token12,
			&wg,
		)
		hubSigninAsync(
			"user0037-13@0037-test.com",
			"NewPassword123$",
			&user0037Token13,
			&wg,
		)
		hubSigninAsync(
			"user0037-14@0037-test.com",
			"NewPassword123$",
			&user0037Token14,
			&wg,
		)
		hubSigninAsync(
			"user0037-15@0037-test.com",
			"NewPassword123$",
			&user0037Token15,
			&wg,
		)
		hubSigninAsync(
			"user0037-16@0037-test.com",
			"NewPassword123$",
			&user0037Token16,
			&wg,
		)
		hubSigninAsync(
			"user0037-17@0037-test.com",
			"NewPassword123$",
			&user0037Token17,
			&wg,
		)
		hubSigninAsync(
			"user0037-18@0037-test.com",
			"NewPassword123$",
			&user0037Token18,
			&wg,
		)
		hubSigninAsync(
			"user0037-19@0037-test.com",
			"NewPassword123$",
			&user0037Token19,
			&wg,
		)
		hubSigninAsync(
			"user0037-20@0037-test.com",
			"NewPassword123$",
			&user0037Token20,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(pool, "0037-incognito-comment-voting-down.pgsql")
		pool.Close()
	})

	Describe("UpvoteIncognitoPostComment", func() {
		It("should upvote a comment successfully", func() {
			// User 1 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment voting",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0037Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 2 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for voting",
			}
			commentResp := testPOSTGetResp(
				user0037Token2,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 3 upvotes the comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token3,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			// Verify vote was recorded
			getReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           25,
			}
			getResp := testPOSTGetResp(
				user0037Token3,
				getReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)
			var getResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Comments)).Should(Equal(1))
			Expect(getResponse.Comments[0].MeUpvoted).Should(BeTrue())
			Expect(getResponse.Comments[0].UpvotesCount).Should(Equal(int32(1)))
			Expect(
				getResponse.Comments[0].DownvotesCount,
			).Should(Equal(int32(0)))
		})
	})

	Describe("DownvoteIncognitoPostComment", func() {
		It("should downvote a comment successfully", func() {
			// User 4 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment downvoting",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0037Token4,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 5 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for downvoting",
			}
			commentResp := testPOSTGetResp(
				user0037Token5,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 6 downvotes the comment
			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token6,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusOK,
			)

			// Verify vote was recorded
			getReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           25,
			}
			getResp := testPOSTGetResp(
				user0037Token6,
				getReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)
			var getResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Comments)).Should(Equal(1))
			Expect(getResponse.Comments[0].MeDownvoted).Should(BeTrue())
			Expect(getResponse.Comments[0].UpvotesCount).Should(Equal(int32(0)))
			Expect(
				getResponse.Comments[0].DownvotesCount,
			).Should(Equal(int32(1)))
		})
	})

	Describe("UnvoteIncognitoPostComment", func() {
		It("should unvote a comment successfully", func() {
			// User 7 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for comment unvoting",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0037Token7,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 8 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for unvoting",
			}
			commentResp := testPOSTGetResp(
				user0037Token8,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 9 upvotes the comment first
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token9,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			// User 9 unvotes the comment
			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token9,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusOK,
			)

			// Verify vote was removed
			getReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           25,
			}
			getResp := testPOSTGetResp(
				user0037Token9,
				getReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)
			var getResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Comments)).Should(Equal(1))
			Expect(getResponse.Comments[0].MeUpvoted).Should(BeFalse())
			Expect(getResponse.Comments[0].MeDownvoted).Should(BeFalse())
			Expect(getResponse.Comments[0].UpvotesCount).Should(Equal(int32(0)))
			Expect(
				getResponse.Comments[0].DownvotesCount,
			).Should(Equal(int32(0)))
		})
	})

	Describe("Vote Own Comment", func() {
		It("should return 422 when trying to vote on own comment", func() {
			// User 10 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for own comment voting",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0037Token10,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Same user adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for own voting",
			}
			commentResp := testPOSTGetResp(
				user0037Token10,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Same user tries to upvote their own comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token10,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)

			// Same user tries to downvote their own comment
			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token10,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)
		})
	})

	Describe("Vote Conflict", func() {
		It(
			"should return 422 when trying to vote in opposite direction",
			func() {
				// User 11 creates a post
				postReq := hub.AddIncognitoPostRequest{
					Content: "Test post for vote conflict",
					TagIDs:  []common.VTagID{"careers"},
				}
				postResp := testPOSTGetResp(
					user0037Token11,
					postReq,
					"/hub/add-incognito-post",
					http.StatusOK,
				)
				var postResponse hub.AddIncognitoPostResponse
				err := json.Unmarshal(postResp.([]byte), &postResponse)
				Expect(err).ShouldNot(HaveOccurred())

				// User 12 adds a comment
				commentReq := hub.AddIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					Content:         "Test comment for vote conflict",
				}
				commentResp := testPOSTGetResp(
					user0037Token12,
					commentReq,
					"/hub/add-incognito-post-comment",
					http.StatusOK,
				)
				var commentResponse hub.AddIncognitoPostCommentResponse
				err = json.Unmarshal(commentResp.([]byte), &commentResponse)
				Expect(err).ShouldNot(HaveOccurred())

				// User 13 downvotes the comment first
				downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					CommentID:       commentResponse.CommentID,
				}
				testPOST(
					user0037Token13,
					downvoteReq,
					"/hub/downvote-incognito-post-comment",
					http.StatusOK,
				)

				// User 13 tries to upvote (should return 422 for vote conflict)
				upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					CommentID:       commentResponse.CommentID,
				}
				testPOST(
					user0037Token13,
					upvoteReq,
					"/hub/upvote-incognito-post-comment",
					http.StatusUnprocessableEntity,
				)
			},
		)
	})

	Describe("Idempotent Voting", func() {
		It("should handle same vote twice idempotently", func() {
			// User 14 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for idempotent voting",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0037Token14,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 15 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for idempotent voting",
			}
			commentResp := testPOSTGetResp(
				user0037Token15,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 16 upvotes the comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token16,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			// User 16 upvotes the same comment again (should be idempotent)
			testPOST(
				user0037Token16,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			// Verify vote count is still 1
			getReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           25,
			}
			getResp := testPOSTGetResp(
				user0037Token16,
				getReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)
			var getResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Comments)).Should(Equal(1))
			Expect(getResponse.Comments[0].MeUpvoted).Should(BeTrue())
			Expect(getResponse.Comments[0].UpvotesCount).Should(Equal(int32(1)))
		})
	})

	Describe("Vote Non-Existent Comment", func() {
		It("should return 422 when voting on non-existent comment", func() {
			// User 17 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for non-existent comment voting",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0037Token17,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to vote on non-existent comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       "non-existent-comment-id",
			}
			testPOST(
				user0037Token17,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusNotFound,
			)

			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       "non-existent-comment-id",
			}
			testPOST(
				user0037Token17,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusNotFound,
			)

			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       "non-existent-comment-id",
			}
			testPOST(
				user0037Token17,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusNotFound,
			)
		})
	})

	Describe("Vote Deleted Comment", func() {
		It("should return 404 when voting on deleted comment", func() {
			// User 18 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for deleted comment voting",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0037Token18,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 19 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment to be deleted",
			}
			commentResp := testPOSTGetResp(
				user0037Token19,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 19 deletes their own comment
			deleteReq := hub.DeleteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token19,
				deleteReq,
				"/hub/delete-incognito-post-comment",
				http.StatusOK,
			)

			// User 20 tries to vote on the deleted comment (should return 404)
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token20,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusNotFound,
			)

			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token20,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusNotFound,
			)

			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				user0037Token20,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusNotFound,
			)
		})
	})

	Describe("Unauthenticated Voting", func() {
		It("should return 401 when voting without authentication", func() {
			// User 18 creates a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for unauthenticated voting",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0037Token18,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// User 19 adds a comment
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for unauthenticated voting",
			}
			commentResp := testPOSTGetResp(
				user0037Token19,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to vote without authentication
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				"",
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnauthorized,
			)

			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				"",
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusUnauthorized,
			)

			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				"",
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusUnauthorized,
			)
		})
	})
})
