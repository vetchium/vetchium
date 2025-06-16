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

var _ = Describe("Incognito Comments API", Ordered, func() {
	var (
		pool *pgxpool.Pool

		// User tokens for dedicated test users
		user0035Token1 string
		user0035Token2 string
		user0035Token3 string
		user0035Token4 string
		user0035Token5 string
		user0035Token6 string
		user0035Token7 string
		user0035Token8 string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0035-incognito-comments-up.pgsql")

		var wg sync.WaitGroup
		wg.Add(8)
		hubSigninAsync(
			"user0035-1@0035-test.com",
			"NewPassword123$",
			&user0035Token1,
			&wg,
		)
		hubSigninAsync(
			"user0035-2@0035-test.com",
			"NewPassword123$",
			&user0035Token2,
			&wg,
		)
		hubSigninAsync(
			"user0035-3@0035-test.com",
			"NewPassword123$",
			&user0035Token3,
			&wg,
		)
		hubSigninAsync(
			"user0035-4@0035-test.com",
			"NewPassword123$",
			&user0035Token4,
			&wg,
		)
		hubSigninAsync(
			"user0035-5@0035-test.com",
			"NewPassword123$",
			&user0035Token5,
			&wg,
		)
		hubSigninAsync(
			"user0035-6@0035-test.com",
			"NewPassword123$",
			&user0035Token6,
			&wg,
		)
		hubSigninAsync(
			"user0035-7@0035-test.com",
			"NewPassword123$",
			&user0035Token7,
			&wg,
		)
		hubSigninAsync(
			"user0035-8@0035-test.com",
			"NewPassword123$",
			&user0035Token8,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(pool, "0035-incognito-comments-down.pgsql")
		pool.Close()
	})

	Describe("AddIncognitoPostComment", func() {
		It("should add a top-level comment successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for top-level comment.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0035Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "This is a top-level comment.",
			}

			commentResp := testPOSTGetResp(
				user0035Token2,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(commentResponse.CommentID).ShouldNot(BeEmpty())
			Expect(
				commentResponse.IncognitoPostID,
			).Should(Equal(postResponse.IncognitoPostID))
		})

		It("should add a reply comment successfully", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for reply comment.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				user0035Token3,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			topCommentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Top-level comment for reply test.",
			}
			topCommentResp := testPOSTGetResp(
				user0035Token4,
				topCommentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var topCommentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(topCommentResp.([]byte), &topCommentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			replyReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "This is a reply to the top-level comment.",
				InReplyTo:       &topCommentResponse.CommentID,
			}

			replyResp := testPOSTGetResp(
				user0035Token5,
				replyReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			var replyResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(replyResp.([]byte), &replyResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(replyResponse.CommentID).ShouldNot(BeEmpty())
		})

		It("should fail with empty content", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for empty comment validation.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0035Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "",
			}

			testPOST(
				user0035Token2,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusBadRequest,
			)
		})

		It("should fail with content too long", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for long comment validation.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0035Token3,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			longContent := make([]byte, 513)
			for i := range longContent {
				longContent[i] = 'a'
			}

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         string(longContent),
			}

			testPOST(
				user0035Token4,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusBadRequest,
			)
		})

		It("should fail for non-existent post", func() {
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: "nonexistent-post-id",
				Content:         "Comment on non-existent post.",
			}

			testPOST(
				user0035Token5,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusNotFound,
			)
		})

		It("should fail for non-existent parent comment", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for invalid parent comment.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0035Token6,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			nonExistentParent := "nonexistent-comment-id"
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Reply to non-existent comment.",
				InReplyTo:       &nonExistentParent,
			}

			testPOST(
				user0035Token7,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusNotFound,
			)
		})

		It("should fail without authentication", func() {
			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: "some-post-id",
				Content:         "Unauthenticated comment.",
			}

			testPOST(
				"",
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusUnauthorized,
			)
		})
	})

	Describe("GetIncognitoPostComments", func() {
		It("should get comments with default parameters", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for getting comments.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0035Token1,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Test comment for retrieval.",
			}
			testPOST(
				user0035Token2,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           25,
			}

			getCommentsResp := testPOSTGetResp(
				user0035Token3,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)

			var getCommentsResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getCommentsResp.([]byte), &getCommentsResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				len(getCommentsResponse.Comments),
			).Should(BeNumerically(">=", 1))
			Expect(
				getCommentsResponse.TotalCommentsCount,
			).Should(BeNumerically(">=", 1))
		})

		It("should respect limit parameter", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for limit parameter.",
				TagIDs:  []common.VTagID{"startups"},
			}
			postResp := testPOSTGetResp(
				user0035Token4,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 5; i++ {
				commentReq := hub.AddIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					Content:         "Comment " + string(rune('A'+i)),
				}
				testPOST(
					user0035Token5,
					commentReq,
					"/hub/add-incognito-post-comment",
					http.StatusOK,
				)
			}

			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Limit:           3,
			}

			getCommentsResp := testPOSTGetResp(
				user0035Token6,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)

			var getCommentsResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getCommentsResp.([]byte), &getCommentsResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				len(getCommentsResponse.Comments),
			).Should(BeNumerically("<=", 3))
		})

		It("should support different sort orders", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Test post for sort order.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				user0035Token7,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 3; i++ {
				commentReq := hub.AddIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					Content:         "Sortable comment " + string(rune('A'+i)),
				}
				testPOST(
					user0035Token8,
					commentReq,
					"/hub/add-incognito-post-comment",
					http.StatusOK,
				)
			}

			for _, sortBy := range []hub.IncognitoPostCommentSortBy{
				hub.IncognitoPostCommentSortByTop,
				hub.IncognitoPostCommentSortByNew,
				hub.IncognitoPostCommentSortByOld,
			} {
				getCommentsReq := hub.GetIncognitoPostCommentsRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					SortBy:          sortBy,
					Limit:           25,
				}

				getCommentsResp := testPOSTGetResp(
					user0035Token1,
					getCommentsReq,
					"/hub/get-incognito-post-comments",
					http.StatusOK,
				)

				var getCommentsResponse hub.GetIncognitoPostCommentsResponse
				err = json.Unmarshal(
					getCommentsResp.([]byte),
					&getCommentsResponse,
				)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(
					len(getCommentsResponse.Comments),
				).Should(BeNumerically(">=", 1))
			}
		})

		It("should handle direct replies preview", func() {
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for direct replies test.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				user0035Token2,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			topCommentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Top level comment.",
			}
			topCommentResp := testPOSTGetResp(
				user0035Token3,
				topCommentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var topCommentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(topCommentResp.([]byte), &topCommentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			for i := 0; i < 3; i++ {
				replyReq := hub.AddIncognitoPostCommentRequest{
					IncognitoPostID: postResponse.IncognitoPostID,
					Content:         "Reply " + string(rune('A'+i)),
					InReplyTo:       &topCommentResponse.CommentID,
				}
				testPOST(
					user0035Token4,
					replyReq,
					"/hub/add-incognito-post-comment",
					http.StatusOK,
				)
			}

			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID:         postResponse.IncognitoPostID,
				DirectRepliesPerComment: 2,
				Limit:                   25,
			}

			getCommentsResp := testPOSTGetResp(
				user0035Token5,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)

			var getCommentsResponse hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(getCommentsResp.([]byte), &getCommentsResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				len(getCommentsResponse.Comments),
			).Should(BeNumerically(">=", 3))
		})

		It("should fail for non-existent post", func() {
			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: "nonexistent-post-id",
				Limit:           25,
			}

			testPOST(
				user0035Token6,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusNotFound,
			)
		})

		It("should fail without authentication", func() {
			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: "some-post-id",
				Limit:           25,
			}

			testPOST(
				"",
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusUnauthorized,
			)
		})

		It("should fail with invalid limit", func() {
			getCommentsReq := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID: "some-post-id",
				Limit:           0,
			}

			testPOST(
				user0035Token7,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusBadRequest,
			)

			getCommentsReq.Limit = 51
			testPOST(
				user0035Token8,
				getCommentsReq,
				"/hub/get-incognito-post-comments",
				http.StatusBadRequest,
			)
		})
	})
})
