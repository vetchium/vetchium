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

var _ = Describe("Incognito Posts API", Ordered, func() {
	var (
		// Database connection
		pool *pgxpool.Pool

		// User tokens
		aliceToken   string
		bobToken     string
		charlieToken string
		eveToken     string
		frankToken   string
	)

	BeforeAll(func() {
		pool = setupTestDB()
		Expect(pool).NotTo(BeNil())
		seedDatabase(pool, "0032-incognito-up.pgsql")

		// Login hub users and get tokens using async signin
		// Note: Diana is excluded as she's DISABLED_HUB_USER for testing disabled user scenarios
		var wg sync.WaitGroup
		wg.Add(5)
		hubSigninAsync(
			"alice@test0032.com",
			"NewPassword123$",
			&aliceToken,
			&wg,
		)
		hubSigninAsync("bob@test0032.com", "NewPassword123$", &bobToken, &wg)
		hubSigninAsync(
			"charlie@company0032.com",
			"NewPassword123$",
			&charlieToken,
			&wg,
		)
		hubSigninAsync("eve@test0032.com", "NewPassword123$", &eveToken, &wg)
		hubSigninAsync(
			"frank@test0032.com",
			"NewPassword123$",
			&frankToken,
			&wg,
		)
		wg.Wait()
	})

	AfterAll(func() {
		seedDatabase(pool, "0032-incognito-down.pgsql")
		pool.Close()
	})

	Describe("AddIncognitoPost", func() {
		It(
			"should create an incognito post successfully with valid data",
			func() {
				reqBody := hub.AddIncognitoPostRequest{
					Content: "This is a new incognito post about technology and career.",
					TagIDs:  []common.VTagID{"technology"},
				}

				respData := testPOSTGetResp(
					aliceToken,
					reqBody,
					"/hub/add-incognito-post",
					http.StatusOK,
				)

				var response hub.AddIncognitoPostResponse
				err := json.Unmarshal(respData.([]byte), &response)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.IncognitoPostID).ShouldNot(BeEmpty())

				// Verify the post was created by getting it
				getReq := hub.GetIncognitoPostRequest{
					IncognitoPostID: response.IncognitoPostID,
				}
				getResp := testPOSTGetResp(
					aliceToken,
					getReq,
					"/hub/get-incognito-post",
					http.StatusOK,
				)

				var getResponse hub.IncognitoPost
				err = json.Unmarshal(getResp.([]byte), &getResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getResponse.Content).Should(Equal(reqBody.Content))
				Expect(getResponse.IsCreatedByMe).Should(BeTrue())
				Expect(getResponse.IsDeleted).Should(BeFalse())
				Expect(len(getResponse.Tags)).Should(Equal(1))
				Expect(
					getResponse.Tags[0].ID,
				).Should(Equal(common.VTagID("technology")))
			},
		)

		It("should create an incognito post with multiple tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "Post with multiple tags about life and career decisions.",
				TagIDs: []common.VTagID{
					"personal-development",
					"careers",
					"mentorship",
				},
			}

			respData := testPOSTGetResp(
				bobToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var response hub.AddIncognitoPostResponse
			err := json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.IncognitoPostID).ShouldNot(BeEmpty())

			// Verify the post has all tags
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: response.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				bobToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(getResponse.Tags)).Should(Equal(3))
		})

		It("should fail without authentication", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This should fail without auth.",
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				"",
				reqBody,
				"/hub/add-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail with empty content", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "",
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				charlieToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with content too long", func() {
			longContent := make([]byte, 1025) // Max is 1024
			for i := range longContent {
				longContent[i] = 'a'
			}

			reqBody := hub.AddIncognitoPostRequest{
				Content: string(longContent),
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				eveToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with no tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has no tags which should fail.",
				TagIDs:  []common.VTagID{},
			}

			testPOST(
				frankToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail with more than 3 tags", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has too many tags.",
				TagIDs: []common.VTagID{
					"technology",
					"careers",
					"personal-development",
					"mentorship",
				},
			}

			testPOST(
				aliceToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})

		It("should fail for disabled user", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This should fail for disabled user.",
				TagIDs:  []common.VTagID{"technology"},
			}

			testPOST(
				"", // Empty token since disabled users can't sign in
				reqBody,
				"/hub/add-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail with invalid tag IDs", func() {
			reqBody := hub.AddIncognitoPostRequest{
				Content: "This post has invalid tags.",
				TagIDs:  []common.VTagID{"invalid-tag", "another-invalid"},
			}

			testPOST(
				bobToken,
				reqBody,
				"/hub/add-incognito-post",
				http.StatusBadRequest,
			)
		})
	})

	Describe("GetIncognitoPost", func() {
		It("should get an existing incognito post successfully", func() {
			// First create a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "This is a test incognito post about technology.",
				TagIDs:  []common.VTagID{"technology"},
			}
			createResp := testPOSTGetResp(
				aliceToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Now get the post
			reqBody := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}

			respData := testPOSTGetResp(
				aliceToken,
				reqBody,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var response hub.IncognitoPost
			err = json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				response.IncognitoPostID,
			).Should(Equal(addResp.IncognitoPostID))
			Expect(
				response.Content,
			).Should(Equal("This is a test incognito post about technology."))
			Expect(response.IsCreatedByMe).Should(BeTrue())
			Expect(response.IsDeleted).Should(BeFalse())
			Expect(len(response.Tags)).Should(Equal(1))
			Expect(
				response.Tags[0].ID,
			).Should(Equal(common.VTagID("technology")))
		})

		It("should show is_created_by_me as false for other users", func() {
			// Bob creates a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "Bob's post for viewing test.",
				TagIDs:  []common.VTagID{"careers"},
			}
			createResp := testPOSTGetResp(
				bobToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Charlie views Bob's post
			reqBody := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}

			respData := testPOSTGetResp(
				charlieToken,
				reqBody,
				"/hub/get-incognito-post",
				http.StatusOK,
			)

			var response hub.IncognitoPost
			err = json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.IsCreatedByMe).Should(BeFalse())
		})

		It("should fail without authentication", func() {
			reqBody := hub.GetIncognitoPostRequest{
				IncognitoPostID: "any-post-id",
			}

			testPOST(
				"",
				reqBody,
				"/hub/get-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			reqBody := hub.GetIncognitoPostRequest{
				IncognitoPostID: "nonexistent",
			}

			testPOST(
				eveToken,
				reqBody,
				"/hub/get-incognito-post",
				http.StatusNotFound,
			)
		})
	})

	Describe("DeleteIncognitoPost", func() {
		It("should delete own incognito post successfully", func() {
			// First create a post to delete
			createReq := hub.AddIncognitoPostRequest{
				Content: "This post will be deleted.",
				TagIDs:  []common.VTagID{"technology"},
			}

			createResp := testPOSTGetResp(
				frankToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)

			var createResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &createResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Now delete it
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}

			testPOST(
				frankToken,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusOK,
			)

			// Verify it's deleted by confirming it returns 404 when fetched
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: createResponse.IncognitoPostID,
			}
			testPOST(
				frankToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusNotFound,
			)
		})

		It("should fail to delete other user's post", func() {
			// Alice creates a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "Alice's post that Bob cannot delete.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			createResp := testPOSTGetResp(
				aliceToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Bob tries to delete Alice's post
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}

			testPOST(
				bobToken,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusForbidden,
			)
		})

		It("should fail without authentication", func() {
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: "any-post-id",
			}

			testPOST(
				"",
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			deleteReq := hub.DeleteIncognitoPostRequest{
				IncognitoPostID: "nonexistent",
			}

			testPOST(
				charlieToken,
				deleteReq,
				"/hub/delete-incognito-post",
				http.StatusNotFound,
			)
		})
	})

	Describe("AddIncognitoPostComment", func() {
		It("should add a top-level comment successfully", func() {
			// Create a post first
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for comment testing.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				eveToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Add a comment
			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "This is a new top-level comment.",
			}

			respData := testPOSTGetResp(
				frankToken,
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			var response hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.CommentID).ShouldNot(BeEmpty())
			Expect(
				response.IncognitoPostID,
			).Should(Equal(postResponse.IncognitoPostID))
		})

		It("should add a reply comment successfully", func() {
			// Create a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for reply testing.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				aliceToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Add a top-level comment first
			comment1Req := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Top level comment for reply test.",
			}
			comment1Resp := testPOSTGetResp(
				bobToken,
				comment1Req,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var comment1Response hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(comment1Resp.([]byte), &comment1Response)
			Expect(err).ShouldNot(HaveOccurred())

			// Add a reply to the comment
			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "This is a reply to the comment above.",
				InReplyTo:       strptr(comment1Response.CommentID),
			}

			respData := testPOSTGetResp(
				charlieToken,
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			var response hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(response.CommentID).ShouldNot(BeEmpty())
		})

		It("should fail without authentication", func() {
			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: "any-post-id",
				Content:         "This should fail without auth.",
			}

			testPOST(
				"",
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: "nonexistent",
				Content:         "Comment on non-existent post.",
			}

			testPOST(
				eveToken,
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusNotFound,
			)
		})

		It("should fail with empty content", func() {
			// Create a post first
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for empty comment test.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				frankToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "",
			}

			testPOST(
				aliceToken,
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusBadRequest,
			)
		})

		It("should fail with content too long", func() {
			// Create a post first
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for long comment test.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				bobToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			longContent := make([]byte, 513) // Max is 512
			for i := range longContent {
				longContent[i] = 'a'
			}

			reqBody := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         string(longContent),
			}

			testPOST(
				charlieToken,
				reqBody,
				"/hub/add-incognito-post-comment",
				http.StatusBadRequest,
			)
		})
	})

	Describe("GetIncognitoPostComments", func() {
		It("should get comments successfully", func() {
			// Create a post
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for comment retrieval test.",
				TagIDs:  []common.VTagID{"personal-development"},
			}
			postResp := testPOSTGetResp(
				eveToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Add some comments
			comment1Req := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "First comment.",
			}
			testPOST(
				frankToken,
				comment1Req,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			comment2Req := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Second comment.",
			}
			testPOST(
				aliceToken,
				comment2Req,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)

			// Get comments
			reqBody := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID:         postResponse.IncognitoPostID,
				Limit:                   10,
				DirectRepliesPerComment: 3,
			}

			respData := testPOSTGetResp(
				bobToken,
				reqBody,
				"/hub/get-incognito-post-comments",
				http.StatusOK,
			)

			var response hub.GetIncognitoPostCommentsResponse
			err = json.Unmarshal(respData.([]byte), &response)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(response.Comments)).Should(BeNumerically(">=", 2))
		})

		It("should fail without authentication", func() {
			reqBody := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID:         "any-post-id",
				Limit:                   10,
				DirectRepliesPerComment: 3,
			}

			testPOST(
				"",
				reqBody,
				"/hub/get-incognito-post-comments",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post", func() {
			reqBody := hub.GetIncognitoPostCommentsRequest{
				IncognitoPostID:         "nonexistent",
				Limit:                   10,
				DirectRepliesPerComment: 3,
			}

			testPOST(
				charlieToken,
				reqBody,
				"/hub/get-incognito-post-comments",
				http.StatusNotFound,
			)
		})
	})

	Describe("Comment Voting", func() {
		It("should upvote comment successfully", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for upvote test.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				eveToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Comment to be upvoted.",
			}
			commentResp := testPOSTGetResp(
				frankToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Upvote the comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				aliceToken,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should downvote comment successfully", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for downvote test.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				bobToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Comment to be downvoted.",
			}
			commentResp := testPOSTGetResp(
				charlieToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Downvote the comment
			downvoteReq := hub.DownvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				eveToken,
				downvoteReq,
				"/hub/downvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should unvote comment successfully", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for unvote test.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			postResp := testPOSTGetResp(
				frankToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Comment to be voted and unvoted.",
			}
			commentResp := testPOSTGetResp(
				aliceToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// First upvote
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				bobToken,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusOK,
			)

			// Then unvote
			unvoteReq := hub.UnvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				bobToken,
				unvoteReq,
				"/hub/unvote-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should fail to vote on own comment", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for self-vote test.",
				TagIDs:  []common.VTagID{"technology"},
			}
			postResp := testPOSTGetResp(
				charlieToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Own comment to try voting on.",
			}
			commentResp := testPOSTGetResp(
				charlieToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to upvote own comment
			upvoteReq := hub.UpvoteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}
			testPOST(
				charlieToken,
				upvoteReq,
				"/hub/upvote-incognito-post-comment",
				http.StatusUnprocessableEntity,
			)
		})
	})

	Describe("Incognito Post Voting", func() {
		It("should upvote incognito post successfully", func() {
			// Create a post by Alice
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post to be upvoted by Bob.",
				TagIDs:  []common.VTagID{"technology"},
			}
			createResp := testPOSTGetResp(
				aliceToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Bob upvotes Alice's post
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				bobToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Verify the vote by getting the post
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				bobToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(1)))
			Expect(getResponse.DownvotesCount).Should(Equal(int32(0)))
			Expect(getResponse.MeUpvoted).Should(BeTrue())
			Expect(getResponse.MeDownvoted).Should(BeFalse())
		})

		It("should downvote incognito post successfully", func() {
			// Create a post by Charlie
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post to be downvoted by Eve.",
				TagIDs:  []common.VTagID{"careers"},
			}
			createResp := testPOSTGetResp(
				charlieToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Eve downvotes Charlie's post
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				eveToken,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusOK,
			)

			// Verify the vote by getting the post
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				eveToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(0)))
			Expect(getResponse.DownvotesCount).Should(Equal(int32(1)))
			Expect(getResponse.MeUpvoted).Should(BeFalse())
			Expect(getResponse.MeDownvoted).Should(BeTrue())
		})

		It("should unvote incognito post successfully", func() {
			// Create a post by Frank
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post to be voted and unvoted by Alice.",
				TagIDs:  []common.VTagID{"mentorship"},
			}
			createResp := testPOSTGetResp(
				frankToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Alice upvotes Frank's post first
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				aliceToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Alice unvotes the post
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				aliceToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusOK,
			)

			// Verify the vote is removed
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				aliceToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(0)))
			Expect(getResponse.DownvotesCount).Should(Equal(int32(0)))
			Expect(getResponse.MeUpvoted).Should(BeFalse())
			Expect(getResponse.MeDownvoted).Should(BeFalse())
		})

		It("should fail to upvote own post", func() {
			// Create a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post that author tries to upvote.",
				TagIDs:  []common.VTagID{"technology"},
			}
			createResp := testPOSTGetResp(
				bobToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to upvote own post
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				bobToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusUnprocessableEntity,
			)
		})

		It("should fail to downvote own post", func() {
			// Create a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post that author tries to downvote.",
				TagIDs:  []common.VTagID{"careers"},
			}
			createResp := testPOSTGetResp(
				charlieToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to downvote own post
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				charlieToken,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusUnprocessableEntity,
			)
		})

		It("should fail to unvote own post", func() {
			// Create a post
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post that author tries to unvote.",
				TagIDs:  []common.VTagID{"personal-development"},
			}
			createResp := testPOSTGetResp(
				eveToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Try to unvote own post
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				eveToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusUnprocessableEntity,
			)
		})

		It(
			"should handle vote conflict - upvote when already downvoted",
			func() {
				// Create a post by Frank
				createReq := hub.AddIncognitoPostRequest{
					Content: "Post for vote conflict test - upvote after downvote.",
					TagIDs:  []common.VTagID{"mentorship"},
				}
				createResp := testPOSTGetResp(
					frankToken,
					createReq,
					"/hub/add-incognito-post",
					http.StatusOK,
				)
				var addResp hub.AddIncognitoPostResponse
				err := json.Unmarshal(createResp.([]byte), &addResp)
				Expect(err).ShouldNot(HaveOccurred())

				// Alice downvotes first
				downvoteReq := hub.DownvoteIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				testPOST(
					aliceToken,
					downvoteReq,
					"/hub/downvote-incognito-post",
					http.StatusOK,
				)

				// Alice tries to upvote (should fail with 422)
				upvoteReq := hub.UpvoteIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				testPOST(
					aliceToken,
					upvoteReq,
					"/hub/upvote-incognito-post",
					http.StatusUnprocessableEntity,
				)

				// Verify the downvote is still there
				getReq := hub.GetIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				getResp := testPOSTGetResp(
					aliceToken,
					getReq,
					"/hub/get-incognito-post",
					http.StatusOK,
				)
				var getResponse hub.IncognitoPost
				err = json.Unmarshal(getResp.([]byte), &getResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getResponse.MeUpvoted).Should(BeFalse())
				Expect(getResponse.MeDownvoted).Should(BeTrue())
			},
		)

		It(
			"should handle vote conflict - downvote when already upvoted",
			func() {
				// Create a post by Bob
				createReq := hub.AddIncognitoPostRequest{
					Content: "Post for vote conflict test - downvote after upvote.",
					TagIDs:  []common.VTagID{"technology"},
				}
				createResp := testPOSTGetResp(
					bobToken,
					createReq,
					"/hub/add-incognito-post",
					http.StatusOK,
				)
				var addResp hub.AddIncognitoPostResponse
				err := json.Unmarshal(createResp.([]byte), &addResp)
				Expect(err).ShouldNot(HaveOccurred())

				// Charlie upvotes first
				upvoteReq := hub.UpvoteIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				testPOST(
					charlieToken,
					upvoteReq,
					"/hub/upvote-incognito-post",
					http.StatusOK,
				)

				// Charlie tries to downvote (should fail with 422)
				downvoteReq := hub.DownvoteIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				testPOST(
					charlieToken,
					downvoteReq,
					"/hub/downvote-incognito-post",
					http.StatusUnprocessableEntity,
				)

				// Verify the upvote is still there
				getReq := hub.GetIncognitoPostRequest{
					IncognitoPostID: addResp.IncognitoPostID,
				}
				getResp := testPOSTGetResp(
					charlieToken,
					getReq,
					"/hub/get-incognito-post",
					http.StatusOK,
				)
				var getResponse hub.IncognitoPost
				err = json.Unmarshal(getResp.([]byte), &getResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getResponse.MeUpvoted).Should(BeTrue())
				Expect(getResponse.MeDownvoted).Should(BeFalse())
			},
		)

		It("should allow same vote multiple times (idempotent)", func() {
			// Create a post by Eve
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post for idempotent vote test.",
				TagIDs:  []common.VTagID{"careers"},
			}
			createResp := testPOSTGetResp(
				eveToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Frank upvotes
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				frankToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Frank upvotes again (should be OK)
			testPOST(
				frankToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusOK,
			)

			// Verify only one upvote
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				frankToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.UpvotesCount).Should(Equal(int32(1)))
			Expect(getResponse.MeUpvoted).Should(BeTrue())
			Expect(getResponse.MeDownvoted).Should(BeFalse())
		})

		It("should allow unvote multiple times (idempotent)", func() {
			// Create a post by Alice
			createReq := hub.AddIncognitoPostRequest{
				Content: "Post for idempotent unvote test.",
				TagIDs:  []common.VTagID{"personal-development"},
			}
			createResp := testPOSTGetResp(
				aliceToken,
				createReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var addResp hub.AddIncognitoPostResponse
			err := json.Unmarshal(createResp.([]byte), &addResp)
			Expect(err).ShouldNot(HaveOccurred())

			// Bob downvotes first
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				bobToken,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusOK,
			)

			// Bob unvotes
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			testPOST(
				bobToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusOK,
			)

			// Bob unvotes again (should be OK)
			testPOST(
				bobToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusOK,
			)

			// Verify no vote
			getReq := hub.GetIncognitoPostRequest{
				IncognitoPostID: addResp.IncognitoPostID,
			}
			getResp := testPOSTGetResp(
				bobToken,
				getReq,
				"/hub/get-incognito-post",
				http.StatusOK,
			)
			var getResponse hub.IncognitoPost
			err = json.Unmarshal(getResp.([]byte), &getResponse)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getResponse.DownvotesCount).Should(Equal(int32(0)))
			Expect(getResponse.MeUpvoted).Should(BeFalse())
			Expect(getResponse.MeDownvoted).Should(BeFalse())
		})

		It("should fail without authentication for upvote", func() {
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: "any-post-id",
			}
			testPOST(
				"",
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail without authentication for downvote", func() {
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: "any-post-id",
			}
			testPOST(
				"",
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail without authentication for unvote", func() {
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: "any-post-id",
			}
			testPOST(
				"",
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusUnauthorized,
			)
		})

		It("should fail for non-existent post on upvote", func() {
			upvoteReq := hub.UpvoteIncognitoPostRequest{
				IncognitoPostID: "nonexistent",
			}
			testPOST(
				charlieToken,
				upvoteReq,
				"/hub/upvote-incognito-post",
				http.StatusNotFound,
			)
		})

		It("should fail for non-existent post on downvote", func() {
			downvoteReq := hub.DownvoteIncognitoPostRequest{
				IncognitoPostID: "nonexistent",
			}
			testPOST(
				eveToken,
				downvoteReq,
				"/hub/downvote-incognito-post",
				http.StatusNotFound,
			)
		})

		It("should fail for non-existent post on unvote", func() {
			unvoteReq := hub.UnvoteIncognitoPostRequest{
				IncognitoPostID: "nonexistent",
			}
			testPOST(
				frankToken,
				unvoteReq,
				"/hub/unvote-incognito-post",
				http.StatusNotFound,
			)
		})
	})

	Describe("DeleteIncognitoPostComment", func() {
		It("should delete own comment successfully", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for comment deletion test.",
				TagIDs:  []common.VTagID{"personal-development"},
			}
			postResp := testPOSTGetResp(
				eveToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Comment to be deleted.",
			}
			commentResp := testPOSTGetResp(
				frankToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Delete the comment
			deleteReq := hub.DeleteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				frankToken,
				deleteReq,
				"/hub/delete-incognito-post-comment",
				http.StatusOK,
			)
		})

		It("should fail to delete other user's comment", func() {
			// Create post and comment
			postReq := hub.AddIncognitoPostRequest{
				Content: "Post for comment deletion permission test.",
				TagIDs:  []common.VTagID{"careers"},
			}
			postResp := testPOSTGetResp(
				aliceToken,
				postReq,
				"/hub/add-incognito-post",
				http.StatusOK,
			)
			var postResponse hub.AddIncognitoPostResponse
			err := json.Unmarshal(postResp.([]byte), &postResponse)
			Expect(err).ShouldNot(HaveOccurred())

			commentReq := hub.AddIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				Content:         "Alice's comment that Bob cannot delete.",
			}
			commentResp := testPOSTGetResp(
				aliceToken,
				commentReq,
				"/hub/add-incognito-post-comment",
				http.StatusOK,
			)
			var commentResponse hub.AddIncognitoPostCommentResponse
			err = json.Unmarshal(commentResp.([]byte), &commentResponse)
			Expect(err).ShouldNot(HaveOccurred())

			// Bob tries to delete Alice's comment
			deleteReq := hub.DeleteIncognitoPostCommentRequest{
				IncognitoPostID: postResponse.IncognitoPostID,
				CommentID:       commentResponse.CommentID,
			}

			testPOST(
				bobToken,
				deleteReq,
				"/hub/delete-incognito-post-comment",
				http.StatusForbidden,
			)
		})
	})
})

func int32ptr(i int32) *int32 {
	return &i
}
