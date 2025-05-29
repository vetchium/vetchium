package dolores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var _ = Describe("Comments", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool

		// User tokens
		postAuthorToken      string
		commenterToken       string
		otherUserToken       string
		disableCommentsToken string
		deleteCommentToken   string
		deleteMyCommentToken string

		// Post IDs for different test scenarios
		commentsEnabledPostID  string
		commentsDisabledPostID string
		deleteCommentsPostID   string
		getCommentsPostID      string
		deleteMyCommentPostID  string
		enableCommentsPostID   string

		// Comment IDs for testing
		testCommentIDs []string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0028-comments-up.pgsql")

		// Login hub users and get tokens
		var wg sync.WaitGroup
		wg.Add(6) // 6 hub users to sign in
		hubSigninAsync(
			"post-author@0028-comments.example.com",
			"NewPassword123$",
			&postAuthorToken,
			&wg,
		)
		hubSigninAsync(
			"commenter@0028-comments.example.com",
			"NewPassword123$",
			&commenterToken,
			&wg,
		)
		hubSigninAsync(
			"other-user@0028-comments.example.com",
			"NewPassword123$",
			&otherUserToken,
			&wg,
		)
		hubSigninAsync(
			"disable-comments@0028-comments.example.com",
			"NewPassword123$",
			&disableCommentsToken,
			&wg,
		)
		hubSigninAsync(
			"delete-comment@0028-comments.example.com",
			"NewPassword123$",
			&deleteCommentToken,
			&wg,
		)
		hubSigninAsync(
			"delete-my-comment@0028-comments.example.com",
			"NewPassword123$",
			&deleteMyCommentToken,
			&wg,
		)
		wg.Wait()

		// Create posts for different test scenarios
		// Post 1: For adding comments (comments enabled)
		post1Resp := testPOSTGetResp(
			postAuthorToken,
			hub.AddPostRequest{
				Content: "Post for adding comments test",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post1AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post1Resp, &post1AddResp)).To(Succeed())
		commentsEnabledPostID = post1AddResp.PostID

		// Post 2: For disable comments test
		post2Resp := testPOSTGetResp(
			disableCommentsToken,
			hub.AddPostRequest{
				Content: "Post for disable comments test",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post2AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post2Resp, &post2AddResp)).To(Succeed())
		commentsDisabledPostID = post2AddResp.PostID

		// Post 3: For delete comment test
		post3Resp := testPOSTGetResp(
			deleteCommentToken,
			hub.AddPostRequest{
				Content: "Post for delete comment test",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post3AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post3Resp, &post3AddResp)).To(Succeed())
		deleteCommentsPostID = post3AddResp.PostID

		// Post 4: For get comments test (with multiple comments)
		post4Resp := testPOSTGetResp(
			commenterToken,
			hub.AddPostRequest{
				Content: "Post for get comments test with multiple comments",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post4AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post4Resp, &post4AddResp)).To(Succeed())
		getCommentsPostID = post4AddResp.PostID

		// Post 5: For delete my comment test
		post5Resp := testPOSTGetResp(
			otherUserToken,
			hub.AddPostRequest{
				Content: "Post for delete my comment test",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post5AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post5Resp, &post5AddResp)).To(Succeed())
		deleteMyCommentPostID = post5AddResp.PostID

		// Post 6: For enable comments test (will be disabled first)
		post6Resp := testPOSTGetResp(
			postAuthorToken,
			hub.AddPostRequest{
				Content: "Post for enable comments test",
			},
			"/hub/add-post",
			http.StatusOK,
		).([]byte)
		var post6AddResp hub.AddPostResponse
		Expect(json.Unmarshal(post6Resp, &post6AddResp)).To(Succeed())
		enableCommentsPostID = post6AddResp.PostID

		// Add some comments to getCommentsPostID for pagination testing
		for i := 1; i <= 5; i++ {
			commentResp := testPOSTGetResp(
				commenterToken,
				hub.AddPostCommentRequest{
					PostID:  getCommentsPostID,
					Content: fmt.Sprintf("Test comment %d for pagination", i),
				},
				"/hub/add-post-comment",
				http.StatusOK,
			).([]byte)
			var commentAddResp hub.AddPostCommentResponse
			Expect(json.Unmarshal(commentResp, &commentAddResp)).To(Succeed())
			testCommentIDs = append(testCommentIDs, commentAddResp.CommentID)
		}

		// Add a comment to deleteMyCommentPostID for delete my comment test
		commentResp := testPOSTGetResp(
			deleteMyCommentToken,
			hub.AddPostCommentRequest{
				PostID:  deleteMyCommentPostID,
				Content: "Comment to be deleted by author",
			},
			"/hub/add-post-comment",
			http.StatusOK,
		).([]byte)
		var commentAddResp hub.AddPostCommentResponse
		Expect(json.Unmarshal(commentResp, &commentAddResp)).To(Succeed())
		testCommentIDs = append(testCommentIDs, commentAddResp.CommentID)

		fmt.Fprintf(GinkgoWriter, "### Test Comment IDs: %v\n", testCommentIDs)
	})

	AfterAll(func() {
		seedDatabase(pool, "0028-comments-down.pgsql")
		pool.Close()
	})

	Describe("Add Post Comment", func() {
		type addPostCommentTestCase struct {
			description string
			token       string
			request     hub.AddPostCommentRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various add post comment scenarios", func() {
			testCases := []addPostCommentTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: "This comment should not be added",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: "This comment should not be added",
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "add valid comment",
					token:       commenterToken,
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: "This is a valid comment",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostCommentResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							response.PostID,
						).Should(Equal(commentsEnabledPostID))
						Expect(response.CommentID).ShouldNot(BeEmpty())
					},
				},
				{
					description: "add comment to non-existent post",
					token:       commenterToken,
					request: hub.AddPostCommentRequest{
						PostID:  "non-existent-post-id",
						Content: "Comment on non-existent post",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "add comment with empty content",
					token:       commenterToken,
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: "",
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add comment with content exceeding max length",
					token:       commenterToken,
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: strings.Repeat("x", 4097), // MaxLength is 4096
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "add comment with exactly max content length",
					token:       commenterToken,
					request: hub.AddPostCommentRequest{
						PostID:  commentsEnabledPostID,
						Content: strings.Repeat("y", 4096), // Exactly MaxLength
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var response hub.AddPostCommentResponse
						err := json.Unmarshal(resp, &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.CommentID).ShouldNot(BeEmpty())
					},
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing AddPostComment: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/add-post-comment",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Get Post Comments", func() {
		type getPostCommentsTestCase struct {
			description string
			token       string
			request     hub.GetPostCommentsRequest
			wantStatus  int
			validate    func([]byte)
		}

		It("should handle various get post comments scenarios", func() {
			testCases := []getPostCommentsTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "get comments for valid post - default limit",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []hub.PostComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							len(comments),
						).Should(Equal(5))
						// We added 5 comments
						// Check order (newest first)
						for i := 0; i < len(comments)-1; i++ {
							Expect(
								comments[i].CreatedAt.After(comments[i+1].CreatedAt) ||
									comments[i].CreatedAt.Equal(comments[i+1].CreatedAt),
							).Should(BeTrue())
						}
						// Check comment structure
						Expect(comments[0].ID).ShouldNot(BeEmpty())
						Expect(comments[0].Content).ShouldNot(BeEmpty())
						Expect(
							comments[0].AuthorName,
						).Should(Equal("Commenter User"))
						Expect(
							comments[0].AuthorHandle,
						).Should(Equal(common.Handle("commenter-user")))
					},
				},
				{
					description: "get comments for non-existent post",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "get comments with limit 2",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
						Limit:  2,
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []hub.PostComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(len(comments)).Should(Equal(2))
					},
				},
				{
					description: "get comments with pagination",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID:        getCommentsPostID,
						Limit:         2,
						PaginationKey: testCommentIDs[3], // Use one of the comment IDs
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []hub.PostComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						// Should get remaining comments
						Expect(len(comments)).Should(BeNumerically(">=", 0))
					},
				},
				{
					description: "get comments with invalid pagination key",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID:        getCommentsPostID,
						Limit:         5,
						PaginationKey: "invalid-comment-id",
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []hub.PostComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						// Should return all comments when pagination key is invalid
						Expect(len(comments)).Should(Equal(5))
					},
				},
				{
					description: "get comments with zero limit (should use default)",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
						Limit:  0,
					},
					wantStatus: http.StatusOK,
					validate: func(resp []byte) {
						var comments []hub.PostComment
						err := json.Unmarshal(resp, &comments)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(
							len(comments),
						).Should(Equal(5))
						// Default limit should be applied
					},
				},
				{
					description: "get comments with negative limit (should fail validation)",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
						Limit:  -1,
					},
					wantStatus: http.StatusBadRequest,
				},
				{
					description: "get comments with limit exceeding max (should fail validation)",
					token:       commenterToken,
					request: hub.GetPostCommentsRequest{
						PostID: getCommentsPostID,
						Limit:  41, // Max is 40
					},
					wantStatus: http.StatusBadRequest,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing GetPostComments: %s\n",
					tc.description,
				)
				resp := testPOSTGetResp(
					tc.token,
					tc.request,
					"/hub/get-post-comments",
					tc.wantStatus,
				)
				if tc.validate != nil && tc.wantStatus == http.StatusOK {
					tc.validate(resp.([]byte))
				}
			}
		})
	})

	Describe("Disable Post Comments", func() {
		type disablePostCommentsTestCase struct {
			description string
			token       string
			request     hub.DisablePostCommentsRequest
			wantStatus  int
		}

		It("should handle various disable post comments scenarios", func() {
			testCases := []disablePostCommentsTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.DisablePostCommentsRequest{
						PostID:                 commentsDisabledPostID,
						DeleteExistingComments: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.DisablePostCommentsRequest{
						PostID:                 commentsDisabledPostID,
						DeleteExistingComments: false,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "disable comments for own post without deleting existing",
					token:       disableCommentsToken,
					request: hub.DisablePostCommentsRequest{
						PostID:                 commentsDisabledPostID,
						DeleteExistingComments: false,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "disable comments for non-existent post",
					token:       disableCommentsToken,
					request: hub.DisablePostCommentsRequest{
						PostID:                 "non-existent-post-id",
						DeleteExistingComments: false,
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "disable comments for other user's post",
					token:       commenterToken,
					request: hub.DisablePostCommentsRequest{
						PostID:                 commentsDisabledPostID,
						DeleteExistingComments: false,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "disable comments with delete existing comments",
					token:       disableCommentsToken,
					request: hub.DisablePostCommentsRequest{
						PostID:                 commentsDisabledPostID,
						DeleteExistingComments: true,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing DisablePostComments: %s\n",
					tc.description,
				)
				testPOST(
					tc.token,
					tc.request,
					"/hub/disable-post-comments",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Enable Post Comments", func() {
		type enablePostCommentsTestCase struct {
			description string
			token       string
			request     hub.EnablePostCommentsRequest
			wantStatus  int
		}

		It("should handle various enable post comments scenarios", func() {
			// First disable comments on enableCommentsPostID
			testPOST(
				postAuthorToken,
				hub.DisablePostCommentsRequest{
					PostID:                 enableCommentsPostID,
					DeleteExistingComments: false,
				},
				"/hub/disable-post-comments",
				http.StatusOK,
			)

			testCases := []enablePostCommentsTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.EnablePostCommentsRequest{
						PostID: enableCommentsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.EnablePostCommentsRequest{
						PostID: enableCommentsPostID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "enable comments for own post",
					token:       postAuthorToken,
					request: hub.EnablePostCommentsRequest{
						PostID: enableCommentsPostID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "enable comments for non-existent post",
					token:       postAuthorToken,
					request: hub.EnablePostCommentsRequest{
						PostID: "non-existent-post-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "enable comments for other user's post",
					token:       commenterToken,
					request: hub.EnablePostCommentsRequest{
						PostID: enableCommentsPostID,
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "enable comments for already enabled post",
					token:       postAuthorToken,
					request: hub.EnablePostCommentsRequest{
						PostID: enableCommentsPostID,
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing EnablePostComments: %s\n",
					tc.description,
				)
				testPOST(
					tc.token,
					tc.request,
					"/hub/enable-post-comments",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Delete Post Comment", func() {
		type deletePostCommentTestCase struct {
			description string
			token       string
			request     hub.DeletePostCommentRequest
			wantStatus  int
		}

		It("should handle various delete post comment scenarios", func() {
			// First add a comment to delete
			commentResp := testPOSTGetResp(
				commenterToken,
				hub.AddPostCommentRequest{
					PostID:  deleteCommentsPostID,
					Content: "Comment to be deleted by post author",
				},
				"/hub/add-post-comment",
				http.StatusOK,
			).([]byte)
			var commentAddResp hub.AddPostCommentResponse
			Expect(json.Unmarshal(commentResp, &commentAddResp)).To(Succeed())
			commentToDelete := commentAddResp.CommentID

			testCases := []deletePostCommentTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.DeletePostCommentRequest{
						PostID:    deleteCommentsPostID,
						CommentID: commentToDelete,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.DeletePostCommentRequest{
						PostID:    deleteCommentsPostID,
						CommentID: commentToDelete,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete comment as post author",
					token:       deleteCommentToken,
					request: hub.DeletePostCommentRequest{
						PostID:    deleteCommentsPostID,
						CommentID: commentToDelete,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete comment from non-existent post",
					token:       deleteCommentToken,
					request: hub.DeletePostCommentRequest{
						PostID:    "non-existent-post-id",
						CommentID: "some-comment-id",
					},
					wantStatus: http.StatusNotFound,
				},
				{
					description: "delete comment as non-post-author",
					token:       commenterToken,
					request: hub.DeletePostCommentRequest{
						PostID:    deleteCommentsPostID,
						CommentID: "some-comment-id",
					},
					wantStatus: http.StatusForbidden,
				},
				{
					description: "delete non-existent comment (should succeed)",
					token:       deleteCommentToken,
					request: hub.DeletePostCommentRequest{
						PostID:    deleteCommentsPostID,
						CommentID: "non-existent-comment-id",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing DeletePostComment: %s\n",
					tc.description,
				)
				testPOST(
					tc.token,
					tc.request,
					"/hub/delete-post-comment",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Delete My Comment", func() {
		type deleteMyCommentTestCase struct {
			description string
			token       string
			request     hub.DeleteMyCommentRequest
			wantStatus  int
		}

		It("should handle various delete my comment scenarios", func() {
			// Use the comment we added in BeforeAll
			myCommentID := testCommentIDs[len(testCommentIDs)-1] // Last comment added

			testCases := []deleteMyCommentTestCase{
				{
					description: "without authentication",
					token:       "",
					request: hub.DeleteMyCommentRequest{
						PostID:    deleteMyCommentPostID,
						CommentID: myCommentID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "with invalid token",
					token:       "invalid-token",
					request: hub.DeleteMyCommentRequest{
						PostID:    deleteMyCommentPostID,
						CommentID: myCommentID,
					},
					wantStatus: http.StatusUnauthorized,
				},
				{
					description: "delete my own comment",
					token:       deleteMyCommentToken,
					request: hub.DeleteMyCommentRequest{
						PostID:    deleteMyCommentPostID,
						CommentID: myCommentID,
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete non-existent comment (should succeed)",
					token:       deleteMyCommentToken,
					request: hub.DeleteMyCommentRequest{
						PostID:    deleteMyCommentPostID,
						CommentID: "non-existent-comment-id",
					},
					wantStatus: http.StatusOK,
				},
				{
					description: "delete comment from non-existent post (should succeed)",
					token:       deleteMyCommentToken,
					request: hub.DeleteMyCommentRequest{
						PostID:    "non-existent-post-id",
						CommentID: "some-comment-id",
					},
					wantStatus: http.StatusOK,
				},
			}

			for _, tc := range testCases {
				fmt.Fprintf(
					GinkgoWriter,
					"### Testing DeleteMyComment: %s\n",
					tc.description,
				)
				testPOST(
					tc.token,
					tc.request,
					"/hub/delete-my-comment",
					tc.wantStatus,
				)
			}
		})
	})

	Describe("Comments Disabled Post", func() {
		It("should not allow adding comments to disabled post", func() {
			// First disable comments
			testPOST(
				disableCommentsToken,
				hub.DisablePostCommentsRequest{
					PostID:                 commentsDisabledPostID,
					DeleteExistingComments: false,
				},
				"/hub/disable-post-comments",
				http.StatusOK,
			)

			// Try to add comment to disabled post
			testPOST(
				commenterToken,
				hub.AddPostCommentRequest{
					PostID:  commentsDisabledPostID,
					Content: "This should fail",
				},
				"/hub/add-post-comment",
				http.StatusForbidden,
			)
		})
	})
})
